// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package content

// English is the copy of the site in English, the default locale.
var English = Site{
	Locale: Locale{Code: "en", Prefix: ""},

	Meta: Meta{
		SiteName:    "Tiefer",
		Title:       "Tiefer | Edge AI in orbit",
		Description: "Tiefer builds AI software that runs on Earth observation satellites: it filters useless pixels in orbit, detects events on board and sends kilobyte-sized alerts instead of gigabytes of raw imagery. Built in Baku.",
		OGImageAlt:  "Tiefer. See deeper.",
		ThemeColor:  "#0C003D",
	},

	UI: UI{
		SkipLink:   "Skip to content",
		NavLabel:   "Main",
		MenuOpen:   "Menu",
		MenuClose:  "Close",
		FooterNav:  "Legal and links",
		LegalLabel: "Legal",
		BackHome:   "Tiefer home",
	},

	Nav: []Link{
		{Label: "Product", Href: "#product"},
		{Label: "How it works", Href: "#how"},
		{Label: "Use cases", Href: "#use-cases"},
		{Label: "Principles", Href: "#principles"},
		{Label: "Contact", Href: "#contact"},
	},
	CTA: Link{Label: "Talk to us", Href: "#contact"},

	Hero: Hero{
		Position:  "BAKU 40.41 N 49.87 E",
		AlertTag:  "ALERT 48 KB",
		Stored:    "KEPT ON BOARD",
		Eyebrow:   "Edge AI in orbit",
		Title:     "See deeper.",
		Lead:      "Tiefer builds AI software that runs on Earth observation satellites. It filters out useless pixels in orbit, detects events on board, and sends kilobyte-sized alerts to the ground instead of gigabytes of raw imagery.",
		Primary:   Link{Label: "Talk to us", Href: "#contact"},
		Secondary: Link{Label: "See how it works", Href: "#how"},
	},

	Problem: Problem{
		Label: "01 / PROBLEM",
		Title: "Satellites see more than they can send.",
		Intro: "An Earth observation satellite stores what it captures and sends it down only when it passes over a ground station. By the time the image is processed, a fire has spread or a spill has drifted. And about two thirds of Earth's surface is under cloud at any moment, so much of what is sent is never useful.",
		Items: []Item{
			{Title: "Hours of delay", Text: "Images wait on board for the next ground contact, then wait again for processing."},
			{Title: "Wasted downlink", Text: "Cloudy and empty frames use the same scarce bandwidth as the ones that matter."},
			{Title: "Passive cameras", Text: "The satellite captures everything and understands nothing until humans look."},
		},
	},

	Product: Product{
		Label: "02 / PRODUCT",
		Title: "Four stages. One onboard system.",
		Body:  "The four lines in our logo are the four stages Tiefer runs on the satellite. Each can be tested on its own. Together they turn a camera into an analyst.",
		Stages: []Stage{
			{Number: "01", Name: "FILTER", Text: "Checks every frame for cloud and quality as it is captured. Useless frames are compressed and kept, not sent."},
			{Number: "02", Name: "DETECT", Text: "Runs computer vision models for the events that matter: fires, oil spills, vessels, floods."},
			{Number: "03", Name: "ALERT", Text: "Sends a small packet with location, event type, confidence and an image chip through the first available link."},
			{Number: "04", Name: "UPDATE", Text: "New or better models reach orbit as small, signed updates. No new satellite needed."},
		},
		Note: "Tiefer is software. It runs on the onboard compute operators already fly: GPU, VPU or FPGA.",
	},

	How: How{
		Label: "03 / HOW IT WORKS",
		Title: "From capture to alert, on board",
		Steps: []Item{
			{Title: "Capture", Text: "The sensor takes a frame."},
			{Title: "Filter", Text: "Cloudy or empty? Compress and keep on board."},
			{Title: "Detect", Text: "Useful frame? Run the event models."},
			{Title: "Alert", Text: "Event found? Build a small alert packet."},
			{Title: "Downlink", Text: "Send it through the first available channel to the operator."},
		},
		Branch:      "KEPT ON BOARD",
		Note:        "How fast an alert reaches the ground depends on the communication link: the next ground station pass, a relay, or a wider ground network. Tiefer makes the alert small enough to use whichever comes first.",
		GroundTitle: "On the ground too",
		Ground: []Item{
			{Title: "Tiefer Lab", Text: "Train, quantise and test models in a virtual orbit simulator and on flight-like hardware before anything flies."},
			{Title: "Tiefer Ground", Text: "Run the same models at the ground station today. When your satellites carry Tiefer, the same code moves to orbit."},
		},
	},

	Demo: Demo{
		Label:     "04 / ALERT PACKET",
		Title:     "An alert, not an archive.",
		Body:      "When Tiefer finds an event, the satellite does not send the whole scene. It sends what an operator needs to act, and keeps the rest on board until someone asks for it. Every packet says what was observed and what was inferred.",
		CardLabel: "Illustrative example. Not real data.",
		CardTitle: "Alert packet",
		Fields: []Field{
			{Name: "Event", Value: "Wildfire"},
			{Name: "Location", Value: "Sector 12"},
			{Name: "Detected", Value: "On board, seconds after capture"},
			{Name: "Confidence", Value: "High"},
			{Name: "Packet size", Value: "48 KB"},
			{Name: "Full scene", Value: "Kept on board. Downlink on request."},
			{Name: "Cloud in frame", Value: "12%"},
			{Name: "Observed", Value: "Active fire front, about 1.2 km long.", Kind: "observed"},
			{Name: "Inferred", Value: "Likely spread to the north-east. Confidence: medium.", Kind: "inferred"},
		},
		SceneLabel: "Full scene",
		AlertLabel: "Alert",
	},

	UseCases: UseCases{
		Label: "05 / USE CASES",
		Title: "Built for events where hours matter",
		Cases: []UseCase{
			{Label: "EVENT 01", Title: "Wildfires", Text: "Spot active fires in forests and grasslands and alert emergency services while the fire is still small."},
			{Label: "EVENT 02", Title: "Oil spills", Text: "Flag slicks around offshore platforms and shipping lanes for rapid response."},
			{Label: "EVENT 03", Title: "Maritime awareness", Text: "Detect vessels at sea; match them with AIS on the ground to find ships that are not broadcasting."},
			{Label: "EVENT 04", Title: "Floods and reservoirs", Text: "Track rising water and reservoir levels during flood season."},
		},
		Note: "We start with operators in Azerbaijan and the wider region, and work in English, Azerbaijani, Turkish and Russian.",
	},

	Principles: Principles{
		Label: "06 / PRINCIPLES",
		Title: "How we build flight software",
		Items: []Principle{
			{Number: "P1", Statement: "Nothing is deleted blindly.", Text: "Filtered frames are compressed and kept. The operator sets the rules."},
			{Number: "P2", Statement: "Observed is not inferred.", Text: "Every alert labels the two separately, with a confidence value."},
			{Number: "P3", Statement: "No invented pixels.", Text: "Generated imagery is never sent as evidence."},
			{Number: "P4", Statement: "Signed updates only.", Text: "Every model in orbit is signed, versioned and can be rolled back."},
			{Number: "P5", Statement: "Measured, not claimed.", Text: "We publish accuracy, latency and power figures with the method behind them."},
			{Number: "P6", Statement: "People are not targets.", Text: "Tiefer is not used to track individuals. Customers are screened."},
		},
	},

	Roadmap: Roadmap{
		Label:    "07 / ROADMAP",
		Title:    "Ground first. Then orbit.",
		NowLabel: "Now",
		Items: []Milestone{
			{Title: "Benchmark", Text: "Cloud filter and three event models, measured on flight-like hardware in a virtual orbit.", Now: true},
			{Title: "Ground pilot", Text: "Tiefer Ground at an operator's ground segment."},
			{Title: "First flight", Text: "Tiefer models running in orbit on a hosted compute platform."},
			{Title: "National satellites", Text: "Tiefer on board operator satellites, agreed with the manufacturer."},
			{Title: "Region", Text: "More operators, more events, radar data on board."},
		},
	},

	Name: Name{
		Label: "08 / NAME",
		Line:  "Tiefer is German for deeper.",
		Body:  "A camera records the surface. An analyst looks deeper. We are putting the analyst in orbit.",
	},

	Contact: Contact{
		Label:        "09 / CONTACT",
		Title:        "We are early, and we are listening.",
		Body:         "We are looking for our first partners: satellite operators, space agencies, integrators and hosted-compute platforms. Tell us which events you need to catch first.",
		FormLabel:    "Contact form",
		NameLabel:    "Name",
		EmailLabel:   "Work email",
		OrgLabel:     "Organisation",
		RoleLabel:    "Your role",
		RoleEmpty:    "Choose one",
		MessageLabel: "What would you want your satellite to detect first?",
		MessageHint:  "Up to 2,000 characters.",
		Optional:     "optional",
		Roles: []Option{
			{Value: "operator", Label: "Satellite operator"},
			{Value: "agency", Label: "Space agency"},
			{Value: "manufacturer", Label: "Manufacturer or integrator"},
			{Value: "hosted-compute", Label: "Hosted compute platform"},
			{Value: "investor", Label: "Investor"},
			{Value: "other", Label: "Other"},
		},
		ConsentBefore:  "I agree that Tiefer may use these details to reply to my message. See the ",
		ConsentLink:    "privacy notice",
		ConsentAfter:   ".",
		HoneypotLabel:  "Leave this field empty",
		Submit:         "Send message",
		Sending:        "Sending",
		Success:        "Thank you. We will reply within two working days.",
		ErrorBefore:    "Something went wrong. Please email us at",
		ErrorAfter:     ".",
		Invalid:        "Please check the fields marked below.",
		Expired:        "This form has expired. Please send it again.",
		Limited:        "Too many messages from this connection. Please try again later or email us.",
		EmailUs:        "Email us",
		PreferEmail:    "Prefer email?",
		LinkedIn:       "Follow Tiefer on LinkedIn",
		NameMissing:    "Please enter your name.",
		NameTooLong:    "Please use at most 100 characters.",
		EmailMissing:   "Please enter your work email.",
		EmailInvalid:   "Please enter a valid email address, for example name@company.com.",
		OrgTooLong:     "Please use at most 150 characters.",
		RoleInvalid:    "Please choose a role from the list.",
		MessageMissing: "Please tell us what you would want to detect first.",
		MessageTooLong: "Please use at most 2,000 characters.",
		ConsentMissing: "Please confirm that we may use these details to reply.",
	},

	Footer: Footer{
		Copyright:     "© 2026 Tiefer",
		LegalNotice:   "Legal notice",
		Privacy:       "Privacy",
		AcceptableUse: "Acceptable use",
		LinkedIn:      "LinkedIn",
		SourceCode:    "Source code",
		Tagline:       "See deeper. Built in Baku.",
	},

	Legal: Legal{
		ReviewNote: "Placeholder text. To be reviewed before launch.",
		Contents:   "On this page",

		Notice: LegalPage{
			Title:       "Legal notice",
			Description: "Company information for Tiefer, a space technology company in Baku, Azerbaijan.",
			Intro:       "Information about the company that operates this website.",
			Sections: []LegalSection{
				{
					Title: "Company",
					Facts: []Fact{
						{Label: "Legal name", Value: "{legal_name}"},
						{Label: "Legal form", Value: "{legal_form}"},
						{Label: "Registered address", Value: "{address}"},
						{Label: "Tax identification number (VÖEN)", Value: "{tax_id}"},
						{Label: "State registration", Value: "{registration}"},
						{Label: "Managing director", Value: "{director}"},
						{Label: "Email", Value: "{contact_email}"},
						{Label: "Website", Value: "{site_url}"},
					},
				},
				{
					Title: "Responsibility for content",
					Paras: []string{
						"{legal_name} is responsible for the content of this website. [CONTENT_RESPONSIBILITY: confirm with counsel whether a named person must be stated here.]",
					},
				},
				{
					Title: "Source code, name and logo",
					Paras: []string{
						"The source code of this website is published under the Mozilla Public License 2.0. The Tiefer name and logo are not covered by that licence and may not be used to suggest endorsement or to brand a fork.",
						"The fonts Mozilla Headline and Mozilla Text are used under the SIL Open Font License 1.1.",
					},
				},
				{
					Title: "Links to other websites",
					Paras: []string{
						"This website links to LinkedIn and to its source code repository. Those websites are run by other companies under their own terms and privacy policies.",
					},
				},
			},
			Updated: "Last updated: [LAST_UPDATED]",
		},

		Privacy: LegalPage{
			Title:       "Privacy notice",
			Description: "How the Tiefer website handles personal data: no cookies, no analytics, no requests to other servers from your browser.",
			Intro:       "This notice explains what personal data this website processes, why, and what rights you have. It covers this website and its contact form, and is structured to address the Law of the Republic of Azerbaijan on Personal Data of 11 May 2010 and, for visitors in the European Union, the General Data Protection Regulation (GDPR).",
			Sections: []LegalSection{
				{
					Title: "Who is responsible",
					Facts: []Fact{
						{Label: "Controller", Value: "{legal_name}"},
						{Label: "Address", Value: "{address}"},
						{Label: "Email", Value: "{contact_email}"},
						{Label: "EU representative", Value: "[EU_REPRESENTATIVE: name and address, or confirm that none is required]"},
					},
				},
				{
					Title: "What this website does not do",
					List: []string{
						"It does not set cookies and stores nothing in your browser.",
						"It does not use analytics, tracking pixels or advertising.",
						"Your browser loads nothing from other servers. The pages are built on our server, and the fonts, the stylesheet and the script come from the same server.",
						"It does not embed social media. The LinkedIn link is a plain link: LinkedIn receives data only if you follow it.",
					},
				},
				{
					Title: "Hosting and protection from attacks",
					Paras: []string{
						"The website runs on Microsoft Azure: Azure Front Door and its web application firewall receive every request and pass it to Azure App Service, which runs the website. To deliver the pages and to protect the website from attacks, these services process the IP address of every visitor together with technical data of the request, such as the time and the page requested.",
						"The web application firewall keeps logs of the requests it inspects, including IP addresses, for {log_retention_days} days, and then deletes them. Legal basis: [LEGAL_BASIS_HOSTING].",
						"Hosting provider: {hosting_provider} ({hosting_country}).",
					},
				},
				{
					Title: "Our own server logs",
					Paras: []string{
						"For each request, our server records the time, the page requested, the response status, the response size and how long the response took. It does not record IP addresses, and it never records what you write in the contact form.",
						"To protect the contact form from abuse, the server keeps the IP address of each sender (for IPv6, only its network part) in memory for up to one hour, to limit the number of messages per connection. It is never written to disk or to a log, and it is gone when the server restarts.",
					},
				},
				{
					Title: "The contact form",
					Paras: []string{
						"If you use the contact form, we receive your name, your work email address, your organisation and role if you give them, your message, and your confirmation that we may reply.",
						"We use these details to reply to you and to discuss working together.",
						"The website does not store your message in a database. It sends it by email, through our email provider {smtp_provider} ({smtp_country}), to our mailbox, and keeps no copy on the web server.",
						"If you email us directly, we receive your email address and what you write.",
					},
				},
				{
					Title: "Legal basis",
					List: []string{
						"Delivering the website and protecting it from attacks: [LEGAL_BASIS_HOSTING].",
						"Replying to your message: [LEGAL_BASIS_CONTACT].",
					},
					Paras: []string{
						"[LEGAL_BASIS: to be determined with counsel, under the Law of the Republic of Azerbaijan on Personal Data and under Article 6 GDPR.]",
					},
				},
				{
					Title: "Who receives your data",
					Paras: []string{
						"Our hosting provider {hosting_provider} and our email provider {smtp_provider} process data on our behalf. We do not sell personal data and do not share it with advertisers. [RECIPIENTS: confirm the list of recipients and the processing agreements.]",
					},
				},
				{
					Title: "International transfers",
					Paras: []string{
						"Tiefer is based in Azerbaijan. Your data may be processed in {hosting_country} (hosting) and {smtp_country} (email). [INTERNATIONAL_TRANSFERS: describe the transfer safeguards under the GDPR and the cross-border transfer rules under Azerbaijani law.]",
					},
				},
				{
					Title: "How long we keep it",
					List: []string{
						"Web application firewall logs: {log_retention_days} days.",
						"Our own server logs: [APP_LOG_RETENTION].",
						"Contact messages: [RETENTION_PERIOD].",
						"The in-memory record for the contact form limit: up to one hour.",
					},
				},
				{
					Title: "Your rights",
					Paras: []string{
						"You can ask for access to your personal data, for its correction or deletion, and for its processing to be restricted. You can object to processing and ask for your data in a portable format. Where processing is based on consent, you can withdraw it at any time; this does not affect processing that happened before.",
						"To use these rights, write to {contact_email}. [DATA_SUBJECT_RIGHTS: confirm the rights and response times under Azerbaijani law and the GDPR.]",
					},
				},
				{
					Title: "How to complain",
					Paras: []string{
						"You can complain to a data protection supervisory authority, in particular in the EU member state where you live or work. In Azerbaijan: [SUPERVISORY_AUTHORITY_AZ].",
					},
				},
			},
			Updated: "Last updated: [LAST_UPDATED]",
		},

		AcceptableUse: LegalPage{
			Title:       "Acceptable use policy",
			Description: "How Tiefer may and may not be used.",
			Intro:       "Tiefer detects events in satellite imagery to support human decisions. This policy sets out what it may not be used for. It applies to every pilot, account and deployment. [AUP_SCOPE: confirm how this policy becomes part of customer contracts.]",
			Sections: []LegalSection{
				{
					Title: "What we commit to",
					List: []string{
						"Every alert labels what was observed separately from what was inferred, with a confidence value.",
						"Generated imagery is never sent as evidence.",
						"Every model in orbit is signed, versioned and can be rolled back.",
						"Alerts support human decisions; they do not replace them.",
					},
				},
				{
					Title: "What is not allowed",
					List: []string{
						"Tracking, identifying or profiling individual people.",
						"Any use that violates human rights.",
						"Any use that breaches applicable sanctions or export control rules. [SANCTIONS_REGIMES: list the sanctions and export control regimes Tiefer applies.]",
						"Selecting targets for attack or planning violence.",
						"Presenting inferred findings as observed facts.",
						"Using satellite imagery against the licence terms of its provider.",
					},
				},
				{
					Title: "Customer screening",
					Paras: []string{
						"We screen every customer before granting access, including against sanctions lists, and we may ask how Tiefer will be used. We may decline any request. [SCREENING_PROCESS: describe the screening steps.]",
					},
				},
				{
					Title: "Suspension",
					Paras: []string{
						"We may suspend or end access if we believe this policy has been broken. [ENFORCEMENT_TERMS: align with customer contracts.]",
					},
				},
				{
					Title: "Reporting misuse",
					Paras: []string{
						"If you believe Tiefer is being misused, write to {contact_email}.",
					},
				},
			},
			Updated: "Last updated: [LAST_UPDATED]",
		},
	},

	NotFound: NotFound{
		Code:  "404 / NOT FOUND",
		Title: "This page does not exist.",
		Body:  "The link may be old, or the address may contain a typo.",
		Home:  "Back to the start page",
	},
}
