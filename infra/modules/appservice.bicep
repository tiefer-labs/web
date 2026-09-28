// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

// Linux App Service running the site container. It accepts traffic only
// from this Front Door profile (service tag plus the X-Azure-FDID header),
// has no FTP or basic-auth deployment, and reads its two secrets from Key
// Vault with its managed identity.

param planName string
param appName string
param location string
@description('App Service plan SKU, for example B1 or P0v3.')
param planSku string
param appIdentityId string
param appIdentityClientId string
param registryLoginServer string
param imageTag string
param frontDoorId string
param keyVaultName string
param workspaceId string
@description('Application settings without secrets (see main.bicep).')
param settings object
@description('True when SMTP is configured, so SMTP_PASS is read from Key Vault.')
param smtpEnabled bool

resource plan 'Microsoft.Web/serverfarms@2024-11-01' = {
  name: planName
  location: location
  kind: 'linux'
  sku: {
    name: planSku
  }
  properties: {
    reserved: true
  }
}

var secretSettings = union({
  CSRF_SECRET: '@Microsoft.KeyVault(VaultName=${keyVaultName};SecretName=csrf-secret)'
}, smtpEnabled ? {
  SMTP_PASS: '@Microsoft.KeyVault(VaultName=${keyVaultName};SecretName=smtp-pass)'
} : {})

var fixedSettings = {
  ENV: 'production'
  PORT: '8080'
  WEBSITES_PORT: '8080'
  WEBSITES_ENABLE_APP_SERVICE_STORAGE: 'false'
  BEHIND_FRONT_DOOR: 'true'
  FRONT_DOOR_ID: frontDoorId
}

resource app 'Microsoft.Web/sites@2024-11-01' = {
  name: appName
  location: location
  kind: 'app,linux,container'
  identity: {
    type: 'UserAssigned'
    userAssignedIdentities: {
      '${appIdentityId}': {}
    }
  }
  properties: {
    serverFarmId: plan.id
    httpsOnly: true
    clientAffinityEnabled: false
    publicNetworkAccess: 'Enabled'
    keyVaultReferenceIdentity: appIdentityId
    siteConfig: {
      linuxFxVersion: 'DOCKER|${registryLoginServer}/tiefer-web:${imageTag}'
      acrUseManagedIdentityCreds: true
      acrUserManagedIdentityID: appIdentityClientId
      alwaysOn: true
      http20Enabled: true
      minTlsVersion: '1.2'
      scmMinTlsVersion: '1.2'
      ftpsState: 'Disabled'
      remoteDebuggingEnabled: false
      healthCheckPath: '/healthz'
      appSettings: [for s in items(union(settings, fixedSettings, secretSettings)): {
        name: s.key
        value: s.value
      }]
      ipSecurityRestrictionsDefaultAction: 'Deny'
      ipSecurityRestrictions: [
        {
          name: 'front-door'
          priority: 100
          action: 'Allow'
          tag: 'ServiceTag'
          ipAddress: 'AzureFrontDoor.Backend'
          headers: {
            'x-azure-fdid': [
              frontDoorId
            ]
          }
        }
      ]
      // The deployment site is not needed: images are pushed to the
      // registry and the site pulls them.
      scmIpSecurityRestrictionsUseMain: false
      scmIpSecurityRestrictionsDefaultAction: 'Deny'
    }
  }
}

resource ftp 'Microsoft.Web/sites/basicPublishingCredentialsPolicies@2024-11-01' = {
  parent: app
  name: 'ftp'
  properties: {
    allow: false
  }
}

resource scm 'Microsoft.Web/sites/basicPublishingCredentialsPolicies@2024-11-01' = {
  parent: app
  name: 'scm'
  properties: {
    allow: false
  }
}

// The application's own log lines (no IP addresses, no form content).
// The HTTP log of App Service, which would hold IP addresses, stays off.
resource diagnostics 'Microsoft.Insights/diagnosticSettings@2021-05-01-preview' = {
  scope: app
  name: 'console-log'
  properties: {
    workspaceId: workspaceId
    logs: [
      {
        category: 'AppServiceConsoleLogs'
        enabled: true
      }
    ]
  }
}

output name string = app.name
output defaultHostName string = app.properties.defaultHostName
