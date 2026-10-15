// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

// Azure Front Door profile and its web application firewall policy. The
// routes are added in frontdoor-routes.bicep once the web app exists,
// because the web app needs this profile's ID first.

param profileName string
@description('Letters and digits only.')
param wafPolicyName string
@allowed([
  'Premium_AzureFrontDoor'
  'Standard_AzureFrontDoor'
])
param sku string
param workspaceId string

var premium = sku == 'Premium_AzureFrontDoor'

resource profile 'Microsoft.Cdn/profiles@2024-09-01' = {
  name: profileName
  location: 'global'
  sku: {
    name: sku
  }
  properties: {
    originResponseTimeoutSeconds: 60
  }
}

resource waf 'Microsoft.Network/FrontDoorWebApplicationFirewallPolicies@2024-02-01' = {
  name: wafPolicyName
  location: 'global'
  sku: {
    name: sku
  }
  properties: {
    policySettings: {
      enabledState: 'Enabled'
      mode: 'Prevention'
      requestBodyCheck: 'Enabled'
      customBlockResponseStatusCode: 403
    }
    customRules: {
      rules: [
        {
          // The site answers GET, HEAD and POST only.
          name: 'AllowedMethods'
          priority: 10
          enabledState: 'Enabled'
          ruleType: 'MatchRule'
          action: 'Block'
          matchConditions: [
            {
              matchVariable: 'RequestMethod'
              operator: 'Equal'
              negateCondition: true
              matchValue: [
                'GET'
                'HEAD'
                'POST'
              ]
            }
          ]
        }
        {
          // Contact form posts per client address, before they reach the
          // application's own limit.
          name: 'ContactPostRate'
          priority: 20
          enabledState: 'Enabled'
          ruleType: 'RateLimitRule'
          rateLimitDurationInMinutes: 5
          rateLimitThreshold: 10
          action: 'Block'
          matchConditions: [
            {
              matchVariable: 'RequestMethod'
              operator: 'Equal'
              matchValue: [
                'POST'
              ]
            }
          ]
        }
        {
          // All requests per client address.
          name: 'RequestRate'
          priority: 30
          enabledState: 'Enabled'
          ruleType: 'RateLimitRule'
          rateLimitDurationInMinutes: 1
          rateLimitThreshold: 300
          action: 'Block'
          matchConditions: [
            {
              matchVariable: 'RequestUri'
              operator: 'Any'
              matchValue: []
            }
          ]
        }
      ]
    }
    // Managed rule sets need the Premium tier.
    managedRules: premium ? {
      managedRuleSets: [
        {
          ruleSetType: 'Microsoft_DefaultRuleSet'
          ruleSetVersion: '2.1'
          ruleSetAction: 'Block'
        }
        {
          ruleSetType: 'Microsoft_BotManagerRuleSet'
          ruleSetVersion: '1.1'
        }
      ]
    } : {
      managedRuleSets: []
    }
  }
}

// Only the firewall log is kept: it is the log the privacy notice names.
// The access log is not enabled.
resource diagnostics 'Microsoft.Insights/diagnosticSettings@2021-05-01-preview' = {
  scope: profile
  name: 'waf-log'
  properties: {
    workspaceId: workspaceId
    logs: [
      {
        category: 'FrontDoorWebApplicationFirewallLog'
        enabled: true
      }
    ]
  }
}

output profileName string = profile.name
output frontDoorId string = profile.properties.frontDoorId
output wafPolicyId string = waf.id
