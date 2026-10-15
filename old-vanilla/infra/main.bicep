// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

// Tiefer website on Azure: Front Door with WAF in front of a Linux App
// Service container, images in Container Registry, secrets in Key Vault,
// logs in Log Analytics, and a GitHub OIDC identity for deployments.
// Deploy into an existing resource group; see infra/README.md.

targetScope = 'resourceGroup'

@description('Azure region for the regional resources.')
param location string = 'germanywestcentral'

@description('Prefix of every resource name.')
param prefix string = 'tiefer-web'

@description('The site domain.')
param siteHost string = 'tiefer.space'

@description('GitHub repository allowed to deploy.')
param githubRepository string = 'tiefer-labs/web'

@description('GitHub deployment environment allowed to deploy.')
param githubEnvironment string = 'production'

@description('Front Door tier. Premium adds the managed WAF rule sets (default rule set and bot protection); Standard has the custom rules only.')
@allowed([
  'Premium_AzureFrontDoor'
  'Standard_AzureFrontDoor'
])
param frontDoorSku string = 'Premium_AzureFrontDoor'

@description('App Service plan SKU.')
param planSku string = 'B1'

@description('Image tag to run. The deploy workflow sets it to the commit SHA.')
param imageTag string = 'latest'

@description('Days to keep the WAF and application logs. Must equal LOG_RETENTION_DAYS, which the privacy notice shows.')
@minValue(30)
param logRetentionDays int = 30

@description('Serve the site on the Front Door default domain too (before DNS is set up).')
param linkDefaultDomain bool = true

@description('Add preload to HSTS. Founder decision; hard to undo.')
param hstsPreload bool = false

@description('Public contact address.')
param contactEmail string = 'hello@tiefer.space'

@description('SMTP host for the contact form. Empty hides the form.')
param smtpHost string = ''
param smtpPort int = 587
param smtpUser string = ''
param contactTo string = ''
param contactFrom string = ''

@description('LEGAL_* values for the legal pages, for example { LEGAL_NAME: \'...\' }. Unset values show placeholders.')
param legal object = {}

var suffix = take(uniqueString(resourceGroup().id), 6)
var names = {
  workspace: '${prefix}-logs'
  appIdentity: '${prefix}-app-id'
  deployIdentity: '${prefix}-deploy-id'
  registry: replace('${prefix}${suffix}', '-', '')
  vault: '${prefix}-kv-${suffix}'
  plan: '${prefix}-plan'
  app: '${prefix}-${suffix}'
  profile: '${prefix}-fd'
  waf: '${replace(prefix, '-', '')}waf'
  endpoint: '${prefix}-${suffix}'
}

module monitoring 'modules/monitoring.bicep' = {
  name: 'monitoring'
  params: {
    name: names.workspace
    location: location
    retentionDays: logRetentionDays
  }
}

module identity 'modules/identity.bicep' = {
  name: 'identity'
  params: {
    appIdentityName: names.appIdentity
    deployIdentityName: names.deployIdentity
    location: location
    githubRepository: githubRepository
    githubEnvironment: githubEnvironment
  }
}

module registry 'modules/registry.bicep' = {
  name: 'registry'
  params: {
    name: names.registry
    location: location
    appPrincipalId: identity.outputs.appPrincipalId
    deployPrincipalId: identity.outputs.deployPrincipalId
  }
}

module vault 'modules/keyvault.bicep' = {
  name: 'keyvault'
  params: {
    name: names.vault
    location: location
    appPrincipalId: identity.outputs.appPrincipalId
  }
}

module frontDoor 'modules/frontdoor-profile.bicep' = {
  name: 'frontdoor-profile'
  params: {
    profileName: names.profile
    wafPolicyName: names.waf
    sku: frontDoorSku
    workspaceId: monitoring.outputs.id
  }
}

module app 'modules/appservice.bicep' = {
  name: 'appservice'
  params: {
    planName: names.plan
    appName: names.app
    location: location
    planSku: planSku
    appIdentityId: identity.outputs.appId
    appIdentityClientId: identity.outputs.appClientId
    registryLoginServer: registry.outputs.loginServer
    imageTag: imageTag
    frontDoorId: frontDoor.outputs.frontDoorId
    keyVaultName: vault.outputs.name
    workspaceId: monitoring.outputs.id
    smtpEnabled: !empty(smtpHost)
    settings: union(legal, {
      SITE_URL: 'https://${siteHost}'
      CONTACT_EMAIL: contactEmail
      LINKEDIN_URL: 'https://www.linkedin.com/company/tiefer/'
      REPO_URL: 'https://github.com/${githubRepository}'
      LOG_RETENTION_DAYS: string(logRetentionDays)
      HSTS_PRELOAD: string(hstsPreload)
    }, empty(smtpHost) ? {} : {
      SMTP_HOST: smtpHost
      SMTP_PORT: string(smtpPort)
      SMTP_USER: smtpUser
      CONTACT_TO: contactTo
      CONTACT_FROM: contactFrom
    })
  }
}

module routes 'modules/frontdoor-routes.bicep' = {
  name: 'frontdoor-routes'
  params: {
    profileName: frontDoor.outputs.profileName
    endpointName: names.endpoint
    wafPolicyId: frontDoor.outputs.wafPolicyId
    originHostName: app.outputs.defaultHostName
    siteHost: siteHost
    linkDefaultDomain: linkDefaultDomain
  }
}

// Deployment identity may update this web app only.
var websiteContributor = 'de139f84-1756-47ae-9be6-808fbbe84772'

resource site 'Microsoft.Web/sites@2024-11-01' existing = {
  name: names.app
}

resource deployRole 'Microsoft.Authorization/roleAssignments@2022-04-01' = {
  scope: site
  name: guid(resourceGroup().id, names.app, 'deploy', websiteContributor)
  dependsOn: [
    app
  ]
  properties: {
    principalId: identity.outputs.deployPrincipalId
    principalType: 'ServicePrincipal'
    roleDefinitionId: subscriptionResourceId('Microsoft.Authorization/roleDefinitions', websiteContributor)
  }
}

output appName string = app.outputs.name
output registryLoginServer string = registry.outputs.loginServer
output keyVaultName string = vault.outputs.name
output deployClientId string = identity.outputs.deployClientId
output frontDoorEndpoint string = routes.outputs.endpointHostName
output apexValidationToken string = routes.outputs.apexValidationToken
output wwwValidationToken string = routes.outputs.wwwValidationToken
