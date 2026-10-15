// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

using 'main.bicep'

param location = 'germanywestcentral'
param prefix = 'tiefer-web'
param siteHost = 'tiefer.space'
param githubRepository = 'tiefer-labs/web'
param githubEnvironment = 'production'
param frontDoorSku = 'Premium_AzureFrontDoor'
param planSku = 'B1'
param logRetentionDays = 30
param linkDefaultDomain = true
param hstsPreload = false
param contactEmail = 'hello@tiefer.space'

// Fill these once the SMTP provider is chosen (SMTP_PASS goes to Key Vault).
param smtpHost = ''
param smtpPort = 587
param smtpUser = ''
param contactTo = ''
param contactFrom = ''

// Company details for the legal pages, once known.
param legal = {}
