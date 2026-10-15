# Azure runbook

How to create the hosting for https://tiefer.space from `infra/` and deploy
the site. Nothing here has been deployed yet. Every step that touches
Azure, DNS or GitHub settings is done by a person with the right access.

## What gets created

| Resource | Name (prefix `tiefer-web`) | Purpose |
|---|---|---|
| Front Door profile, Premium | `tiefer-web-fd` | TLS, `www` and HTTP redirects, edge cache for `/static/*` |
| WAF policy | `tieferwebwaf` | Prevention mode: allowed methods, rate limits, Microsoft default rule set 2.1, bot manager 1.1 |
| App Service plan (Linux, B1) and web app | `tiefer-web-plan`, `tiefer-web-<suffix>` | Runs the container; only reachable from this Front Door profile |
| Container Registry (Basic) | `tieferweb<suffix>` | Site images; no admin user |
| Key Vault | `tiefer-web-kv-<suffix>` | `csrf-secret` and `smtp-pass`; RBAC only, purge protection |
| Log Analytics workspace | `tiefer-web-logs` | WAF log and the application log, kept `logRetentionDays` (30) days |
| Managed identities | `tiefer-web-app-id`, `tiefer-web-deploy-id` | The app pulls images and reads secrets; GitHub deploys with OpenID Connect |

`<suffix>` is derived from the resource group, so names are stable.

Access rules:

- The web app accepts traffic only from the Front Door service tag with
  this profile's `X-Azure-FDID`, and the application checks the header
  again (`BEHIND_FRONT_DOOR=true`, `FRONT_DOOR_ID`).
- FTP and basic-auth deployment are off; the deployment site (SCM) denies
  all traffic.
- The deploy identity can push images and update this one web app. It
  trusts only workflow runs of `tiefer-labs/web` in the GitHub environment
  `production`.
- Only the WAF log and the application's own log are collected. The
  Front Door access log and the App Service HTTP log (both with IP
  addresses) are not enabled, which matches the privacy notice.

## Before you start

- **Legal**: confirm with a lawyer in Azerbaijan that personal data from
  the contact form may be processed in the chosen region (default
  `germanywestcentral`, Frankfurt), and by the SMTP provider. Record the
  result in the privacy notice (`LEGAL_HOSTING_*`, `LEGAL_SMTP_*`).
- **Cost**: Front Door Premium is needed for the managed WAF rule sets.
  Standard is much cheaper but keeps only the custom rules
  (`frontDoorSku = 'Standard_AzureFrontDoor'`).
- Tools: Azure CLI with Bicep (`az bicep install`), and Owner (or User
  Access Administrator plus Contributor) on the resource group, because
  the templates create role assignments.

## 1. Create the resources

```sh
az login
az account set --subscription "<subscription id>"
az group create --name tiefer-web-rg --location germanywestcentral

# Review infra/main.bicepparam first (SMTP, legal details).
az deployment group what-if --resource-group tiefer-web-rg \
  --template-file infra/main.bicep --parameters infra/main.bicepparam
az deployment group create --resource-group tiefer-web-rg \
  --template-file infra/main.bicep --parameters infra/main.bicepparam
```

Note the outputs: `appName`, `registryLoginServer`, `keyVaultName`,
`deployClientId`, `frontDoorEndpoint`, `apexValidationToken`,
`wwwValidationToken`. The web app cannot start until the first image is
pushed (step 4).

## 2. Put the secrets into Key Vault

Give yourself the role "Key Vault Secrets Officer" on the vault, set the
secrets, then remove the role again:

```sh
az keyvault secret set --vault-name <keyVaultName> --name csrf-secret \
  --value "$(openssl rand -hex 32)" --output none
# Only when SMTP is configured:
az keyvault secret set --vault-name <keyVaultName> --name smtp-pass \
  --value "<SMTP password>" --output none
```

The values never go into the repository, a parameter file or a log. If a
reference cannot be resolved, the site refuses to start rather than run
with the literal reference text as its key.

## 3. Set up GitHub

In the repository settings (organisation owner):

1. Environment `production`: required reviewer (the founder), deployment
   branches limited to `main`.
2. Environment variables (not secrets; they are identifiers):
   `AZURE_CLIENT_ID` (output `deployClientId`), `AZURE_TENANT_ID`,
   `AZURE_SUBSCRIPTION_ID`, `AZURE_RESOURCE_GROUP` (`tiefer-web-rg`),
   `AZURE_WEBAPP_NAME` (output `appName`), `AZURE_REGISTRY_LOGIN_SERVER`
   (output `registryLoginServer`).
3. Branch protection on `main`: pull requests, the CI and CodeQL checks
   required, no force pushes.

## 4. First deployment

Run the workflow **Deploy** by hand (Actions, Deploy, Run workflow). It
runs CI, builds the image, pushes it as `tiefer-web:<commit>`, points the
web app at it and checks `https://tiefer.space/healthz`. Until DNS is set
up (step 5) the smoke test fails; check the site on the Front Door
endpoint instead (`https://<frontDoorEndpoint>/`). From then on every push
to `main` deploys after CI passes and the reviewer approves.

**Roll back** by pointing the web app at an earlier image:

```sh
az webapp config set --resource-group tiefer-web-rg --name <appName> \
  --linux-fx-version "DOCKER|<registryLoginServer>/tiefer-web:<earlier commit>"
```

## 5. DNS for tiefer.space

| Name | Type | Value |
|---|---|---|
| `_dnsauth.tiefer.space` | TXT | `apexValidationToken` |
| `_dnsauth.www.tiefer.space` | TXT | `wwwValidationToken` |
| `tiefer.space` | ALIAS, ANAME or flattened CNAME | `frontDoorEndpoint` |
| `www.tiefer.space` | CNAME | `frontDoorEndpoint` |
| `tiefer.space` | CAA | `0 issue "digicert.com"` (Front Door managed certificates) |
| `tiefer.space` | CAA | `0 issuewild ";"` |
| `tiefer.space` | CAA | `0 iodef "mailto:hello@tiefer.space"` |

If the registrar cannot point the bare domain at a host name (no ALIAS
or flattening), move the zone to Azure DNS, which can. Turn on DNSSEC and
the transfer lock at the registrar, with two-factor authentication on the
registrar account.

Front Door renews its managed certificates itself; for the bare domain,
check the certificate state in the portal every few months, because
renewal depends on the domain still pointing at the endpoint.

When both domains show "Approved" and their certificates are active, set
`linkDefaultDomain = false` in `main.bicepparam` and deploy the templates
again, so the site answers only on tiefer.space.

## 6. Mail for @tiefer.space

The contact form sends from `CONTACT_FROM` through the SMTP provider.
Set up, with the provider's values:

| Name | Type | Value |
|---|---|---|
| `tiefer.space` | TXT | `v=spf1 include:<provider> -all` |
| `<selector>._domainkey.tiefer.space` | TXT or CNAME | DKIM key from the provider |
| `_dmarc.tiefer.space` | TXT | `v=DMARC1; p=none; rua=mailto:<reports>` for two weeks, then `p=reject` |
| `_smtp._tls.tiefer.space` | TXT | `v=TLSRPTv1; rua=mailto:<reports>` |

## 7. HSTS preload (decision)

The site sends `Strict-Transport-Security: max-age=63072000;
includeSubDomains`. Setting `hstsPreload = true` adds `preload`; after that
the domain can be submitted at https://hstspreload.org. Browsers then never
use plain HTTP for tiefer.space **or any subdomain**, and removal takes
months. Decide only when every current and future subdomain will serve
HTTPS.

## 8. Checks after go-live

```sh
curl -sI https://tiefer.space/ | head -30        # security headers, no Server header
curl -sI http://tiefer.space/                    # redirect to https
curl -sI https://www.tiefer.space/               # redirect to https://tiefer.space/
curl -s https://tiefer.space/.well-known/security.txt
curl -sI https://<appName>.azurewebsites.net/    # 403: the origin is closed
```

Then run SSL Labs, the Mozilla HTTP Observatory and internet.nl once.

## Operating

- **Logs**: Log Analytics, table `AzureDiagnostics` for the WAF
  (`Category == "FrontDoorWebApplicationFirewallLog"`) and
  `AppServiceConsoleLogs` for the application.
- **Secret rotation**: set a new `csrf-secret` version; open forms become
  invalid, nothing else. Restart the web app to load it. Rotate
  `smtp-pass` at the provider first.
- **Updates**: Dependabot proposes updates to Go tools, actions and base
  images; CI runs govulncheck weekly.
- **Scaling out**: used form tokens are remembered per instance (see
  `docs/SECURITY-DECISIONS.md`); keep one instance or move that state
  before adding more.
