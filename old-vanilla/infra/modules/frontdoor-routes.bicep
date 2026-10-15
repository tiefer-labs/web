// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

// Endpoint, custom domains, origin and routes of the Front Door profile.
// www and plain HTTP redirect to https://tiefer.space.

param profileName string
param endpointName string
param wafPolicyId string
@description('Default host name of the web app, the only origin.')
param originHostName string
@description('The site domain, for example tiefer.space.')
param siteHost string
@description('Also serve the site on the endpoint default domain (*.azurefd.net). Useful before DNS is set up; switch off afterwards.')
param linkDefaultDomain bool

resource profile 'Microsoft.Cdn/profiles@2024-09-01' existing = {
  name: profileName
}

resource endpoint 'Microsoft.Cdn/profiles/afdEndpoints@2024-09-01' = {
  parent: profile
  name: endpointName
  location: 'global'
  properties: {
    enabledState: 'Enabled'
  }
}

resource apex 'Microsoft.Cdn/profiles/customDomains@2024-09-01' = {
  parent: profile
  name: replace(siteHost, '.', '-')
  properties: {
    hostName: siteHost
    tlsSettings: {
      certificateType: 'ManagedCertificate'
      minimumTlsVersion: 'TLS12'
    }
  }
}

resource www 'Microsoft.Cdn/profiles/customDomains@2024-09-01' = {
  parent: profile
  name: 'www-${replace(siteHost, '.', '-')}'
  properties: {
    hostName: 'www.${siteHost}'
    tlsSettings: {
      certificateType: 'ManagedCertificate'
      minimumTlsVersion: 'TLS12'
    }
  }
}

resource originGroup 'Microsoft.Cdn/profiles/originGroups@2024-09-01' = {
  parent: profile
  name: 'app'
  properties: {
    loadBalancingSettings: {
      sampleSize: 4
      successfulSamplesRequired: 3
      additionalLatencyInMilliseconds: 50
    }
    healthProbeSettings: {
      probePath: '/healthz'
      probeRequestType: 'HEAD'
      probeProtocol: 'Https'
      probeIntervalInSeconds: 100
    }
    sessionAffinityState: 'Disabled'
  }
}

resource origin 'Microsoft.Cdn/profiles/originGroups/origins@2024-09-01' = {
  parent: originGroup
  name: 'app-service'
  properties: {
    hostName: originHostName
    originHostHeader: originHostName
    httpPort: 80
    httpsPort: 443
    priority: 1
    weight: 1000
    enabledState: 'Enabled'
    enforceCertificateNameCheck: true
  }
}

resource ruleSet 'Microsoft.Cdn/profiles/ruleSets@2024-09-01' = {
  parent: profile
  name: 'canonicalhost'
}

resource redirectWww 'Microsoft.Cdn/profiles/ruleSets/rules@2024-09-01' = {
  parent: ruleSet
  name: 'wwwtoapex'
  properties: {
    order: 1
    matchProcessingBehavior: 'Stop'
    conditions: [
      {
        name: 'HostName'
        parameters: {
          typeName: 'DeliveryRuleHostNameConditionParameters'
          operator: 'Equal'
          negateCondition: false
          matchValues: [
            'www.${siteHost}'
          ]
          transforms: [
            'Lowercase'
          ]
        }
      }
    ]
    actions: [
      {
        name: 'UrlRedirect'
        parameters: {
          typeName: 'DeliveryRuleUrlRedirectActionParameters'
          redirectType: 'PermanentRedirect'
          destinationProtocol: 'Https'
          customHostname: siteHost
        }
      }
    ]
  }
}

var domains = [
  {
    id: apex.id
  }
  {
    id: www.id
  }
]

// Pages: never cached at the edge (they carry a single-use form token).
resource pages 'Microsoft.Cdn/profiles/afdEndpoints/routes@2024-09-01' = {
  parent: endpoint
  name: 'pages'
  dependsOn: [
    origin
    redirectWww
  ]
  properties: {
    customDomains: domains
    originGroup: {
      id: originGroup.id
    }
    ruleSets: [
      {
        id: ruleSet.id
      }
    ]
    supportedProtocols: [
      'Http'
      'Https'
    ]
    patternsToMatch: [
      '/*'
    ]
    forwardingProtocol: 'HttpsOnly'
    httpsRedirect: 'Enabled'
    linkToDefaultDomain: linkDefaultDomain ? 'Enabled' : 'Disabled'
    enabledState: 'Enabled'
  }
}

// Hashed static files never change, so the edge may cache them.
resource assets 'Microsoft.Cdn/profiles/afdEndpoints/routes@2024-09-01' = {
  parent: endpoint
  name: 'static'
  dependsOn: [
    origin
    pages
  ]
  properties: {
    customDomains: domains
    originGroup: {
      id: originGroup.id
    }
    ruleSets: [
      {
        id: ruleSet.id
      }
    ]
    supportedProtocols: [
      'Http'
      'Https'
    ]
    patternsToMatch: [
      '/static/*'
    ]
    forwardingProtocol: 'HttpsOnly'
    httpsRedirect: 'Enabled'
    linkToDefaultDomain: linkDefaultDomain ? 'Enabled' : 'Disabled'
    enabledState: 'Enabled'
    cacheConfiguration: {
      queryStringCachingBehavior: 'IgnoreQueryString'
      compressionSettings: {
        isCompressionEnabled: false
      }
    }
  }
}

resource security 'Microsoft.Cdn/profiles/securityPolicies@2024-09-01' = {
  parent: profile
  name: 'waf'
  properties: {
    parameters: {
      type: 'WebApplicationFirewall'
      wafPolicy: {
        id: wafPolicyId
      }
      associations: [
        {
          domains: concat(domains, [
            {
              id: endpoint.id
            }
          ])
          patternsToMatch: [
            '/*'
          ]
        }
      ]
    }
  }
}

output endpointHostName string = endpoint.properties.hostName
output apexValidationToken string = apex.properties.validationProperties.validationToken
output wwwValidationToken string = www.properties.validationProperties.validationToken
