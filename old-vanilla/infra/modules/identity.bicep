// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

// Two user-assigned identities: one the web app runs as (pulls its image,
// reads its secrets), one GitHub Actions deploys as through OpenID Connect,
// without any stored password.

param appIdentityName string
param deployIdentityName string
param location string
@description('GitHub repository, for example tiefer-labs/web.')
param githubRepository string
@description('GitHub deployment environment allowed to deploy.')
param githubEnvironment string

resource app 'Microsoft.ManagedIdentity/userAssignedIdentities@2023-01-31' = {
  name: appIdentityName
  location: location
}

resource deploy 'Microsoft.ManagedIdentity/userAssignedIdentities@2023-01-31' = {
  name: deployIdentityName
  location: location
}

// Only workflow runs in the given GitHub environment can use the deploy
// identity: not other branches, pull requests or forks.
resource github 'Microsoft.ManagedIdentity/userAssignedIdentities/federatedIdentityCredentials@2023-01-31' = {
  parent: deploy
  name: 'github-${githubEnvironment}'
  properties: {
    issuer: 'https://token.actions.githubusercontent.com'
    subject: 'repo:${githubRepository}:environment:${githubEnvironment}'
    audiences: [
      'api://AzureADTokenExchange'
    ]
  }
}

output appId string = app.id
output appPrincipalId string = app.properties.principalId
output appClientId string = app.properties.clientId
output deployPrincipalId string = deploy.properties.principalId
output deployClientId string = deploy.properties.clientId
