// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package content

// English is the English copy of the site. Edit the text here; the
// templates in web/templates only decide where it appears.
var English = Site{
	Locale: Locale{
		Code:     "en",
		Prefix:   "",
		Name:     "English",
		OGLocale: "en_GB",
	},

	Meta: Meta{
		SiteName:    "Tiefer",
		Title:       "Tiefer | Space intelligence you can ask",
		Description: "Tiefer is an AI analyst for satellite data. Ask a question in plain language and get a sourced intelligence brief from optical, radar, thermal and night-light imagery. Built in Baku.",
		OGImageAlt:  "The Tiefer logo in white on dark blue",
	},

	UI: UI{
		SkipLink:  "Skip to content",
		MenuOpen:  "Menu",
		MenuClose: "Close",
		LogoAlt:   "Tiefer",
		NavLabel:  "Main",
		FooterNav: "Footer",
	},

	Nav: Nav{
		Links: []Link{
			{Label: "Product", Href: "product"},
			{Label: "How it works", Href: "how"},
			{Label: "Use cases", Href: "use-cases"},
			{Label: "Principles", Href: "principles"},
			{Label: "FAQ", Href: "faq"},
			{Label: "Contact", Href: "contact"},
		},
		CTA: Link{Label: "Request access", Href: "contact"},
	},

	Hero: Hero{
		Coordinates:     "BAKU 40.41 N 49.87 E",
		ArtLabels:       []string{"PASS 034", "SAR", "CLOUD 82%"},
		Eyebrow:         "Space intelligence you can ask",
		Title:           "See deeper.",
		Lead:            "Ask a question about any place on Earth in plain language. Tiefer finds and fuses optical, radar, thermal and night-light satellite data, tasks new imagery when the archive is not enough, and returns a sourced intelligence brief.",
		Question:        "Which ports on the eastern Caspian coast saw more vessel activity in the last 14 days?",
		QuestionCaption: "Example question",
		Primary:         Link{Label: "Request access", Href: "contact"},
		Secondary:       Link{Label: "See how it works", Href: "how"},
	},

	Problem: Problem{
		Title: "Satellite data is everywhere. Answers are not.",
		Intro: "Getting one answer from space still takes days: find the right images, buy them, clean them, align optical, radar and thermal layers by hand, look for change, write it up. The people who need the answer rarely have the specialists to do it.",
		Items: []Item{
			{Title: "Scattered catalogues", Text: "Every satellite has its own archive, format, resolution and revisit."},
			{Title: "Clouds and darkness", Text: "Optical images fail under cloud and at night. Radar and thermal fill the gap, but few teams can combine them."},
			{Title: "Specialist work", Text: "Cleaning, aligning and analysing imagery needs remote sensing experts most teams do not have."},
			{Title: "Slow tasking", Text: "When no recent image exists, ordering a new one is a manual process that takes days."},
		},
	},

	Product: Product{
		Title: "Four layers. One answer.",
		Body:  "The four lines in our logo are the four layers Tiefer fuses. Each keeps its own source. Together they answer the question.",
		Layers: []Layer{
			{Icon: "optical", Label: "OPTICAL", Text: "What the surface looks like. Sentinel-2, Landsat and commercial imagery.", Specs: []Fact{
				{Label: "Open data", Value: "Sentinel-2 (10 m), Landsat 8 and 9 (30 m)"},
				{Label: "On demand", Value: "Commercial optical imagery with finer detail"},
				{Label: "Limit", Value: "Cannot see through cloud or in the dark"},
			}},
			{Icon: "radar", Label: "RADAR", Text: "Sees through cloud and at night. Sentinel-1 and commercial SAR.", Specs: []Fact{
				{Label: "Open data", Value: "Sentinel-1 (about 10 m), day and night"},
				{Label: "On demand", Value: "Commercial SAR"},
				{Label: "Limit", Value: "Shows structure and texture, not colour"},
			}},
			{Icon: "thermal", Label: "THERMAL", Text: "Heat: fires, industrial activity, surface temperature.", Specs: []Fact{
				{Label: "Open data", Value: "Landsat thermal bands (100 m), ECOSTRESS (about 70 m)"},
				{Label: "On demand", Value: "Commercial thermal imagery"},
				{Label: "Limit", Value: "Coarse pixels: small heat sources can be missed"},
			}},
			{Icon: "night", Label: "NIGHT LIGHTS", Text: "Where light appears or disappears after dark.", Specs: []Fact{
				{Label: "Open data", Value: "VIIRS night lights (about 500 to 750 m), VIIRS active fires (375 m)"},
				{Label: "On demand", Value: "Commercial night imagery"},
				{Label: "Limit", Value: "Small or weak lights often do not show"},
			}},
		},
	},

	How: How{
		Title: "From question to brief in five steps",
		Steps: []Item{
			{Title: "Ask", Text: "Tiefer turns your question into a place, a time window and the signs to look for. If something is unclear, it asks."},
			{Title: "Plan", Text: "It picks the sensors that can answer. Cloudy area? Radar first."},
			{Title: "Fuse", Text: "It finds, cleans and aligns the imagery, then runs change and object detection."},
			{Title: "Task", Text: "No recent image? It prepares an order for a commercial satellite, within your budget rules. Anything beyond them waits for your approval."},
			{Title: "Brief", Text: "You get a map, the findings, a confidence level for each, the source images, and what could not be seen."},
		},
		Note: Item{
			Title: "Under the hood",
			Text:  "Two kinds of models share the work. A language model turns your question into a plan and explains the result. Specialised geospatial models analyse the images: they find change and detect objects such as vessels, tanks and new buildings. The language model does not make up findings.",
		},
	},

	Demo: Demo{
		Title:         "Every finding shows its source.",
		Body:          "An answer you cannot check is not intelligence. Tiefer labels what a satellite observed separately from what a model inferred, links every finding to its source image, and never presents generated pixels as evidence.",
		CardLabel:     "Illustrative example. Not real data.",
		ContentsTitle: "Every brief contains",
		Contents: []string{
			"A map of the area and the time window you asked about.",
			"Each finding, labelled observed or inferred, with a confidence level.",
			"For every finding: the sensor, the satellite, the date of the image, the processing steps and the model version.",
			"What could not be seen, and why.",
			"A proposed next step when the data is not enough, such as a radar pass or new imagery.",
		},
		Rows: []BriefRow{
			{Label: "Question", Value: "Vessel activity at Port A, last 14 days"},
			{Label: "Sensors", Value: "Radar (SAR), optical"},
			{Label: "Passes used", Value: "9 radar, 3 optical (4 optical rejected: cloud)"},
			{Label: "Finding 1", Tag: "observed", TagLabel: "Observed", Value: "Vessels at berth rose from 6 to 11 between day 1 and day 14."},
			{Label: "Finding 2", Tag: "inferred", TagLabel: "Inferred", Value: "Likely increase in loading activity. Confidence: medium."},
			{Label: "Not seen", Value: "No optical view on days 5 to 9 (cloud). Radar only."},
		},
		// Day 1 to day 14. Matches the rows above: 9 radar passes,
		// 3 usable optical passes, 4 optical passes rejected on days 5 to 9.
		Strip: Strip{
			Days: []PassDay{
				{Radar: 1},
				{Optical: 1},
				{Radar: 1},
				{},
				{Radar: 1, Rejected: 1},
				{Radar: 1, Rejected: 1},
				{Radar: 1},
				{Radar: 1, Rejected: 1},
				{Radar: 1, Rejected: 1},
				{},
				{Radar: 1},
				{Optical: 1},
				{Radar: 1},
				{Optical: 1},
			},
			FirstDay:       "DAY 1",
			LastDay:        "DAY 14",
			LegendRadar:    "Radar pass",
			LegendOptical:  "Optical pass",
			LegendRejected: "Rejected: cloud",
		},
	},

	UseCases: UseCases{
		Title:         "Built for people who decide on what happens on the ground",
		QuestionLabel: "Example questions",
		Cards: []UseCase{
			{
				Title: "Commodity traders",
				Text:  "Ports, terminals, storage and crops: see supply move before the statistics do.",
				Questions: []string{
					"How full are the tanks at this terminal compared with last month?",
					"How has the number of vessels at this port changed over the last two weeks?",
					"How are the grain fields in this region developing this season?",
				},
			},
			{
				Title: "Insurers",
				Text:  "Floods, fires and damage mapped in hours, with the images to back each claim.",
				Questions: []string{
					"Which insured sites in this district were under water on Tuesday?",
					"Which properties lie inside the area burned by this fire?",
				},
			},
			{
				Title: "Governments",
				Text:  "Sovereign deployments that keep sensitive questions and data inside the country.",
				Questions: []string{
					"What changed at this infrastructure site since January?",
					"Where has new construction appeared in this area since last year?",
				},
			},
		},
		Note: "We start in the Caspian, Caucasus, Black Sea and Central Asia region. Tiefer answers in English, Azerbaijani, Turkish and Russian.",
	},

	Principles: Principles{
		Title: "How we work",
		Items: []Item{
			{Title: "Observed is not inferred.", Text: "Every brief labels the two separately."},
			{Title: "No invented pixels.", Text: "Generated imagery is never shown as evidence. If it was cloudy, we say so."},
			{Title: "Every finding is traceable.", Text: "Sensor, date, source image and processing steps travel with it."},
			{Title: "You control the spend.", Text: "Automatic tasking runs only within your rules. Anything else waits for approval."},
			{Title: "Limits are stated.", Text: "Every brief says what could not be seen."},
			{Title: "People are not targets.", Text: "Tiefer is not used to track individuals. Customers are screened."},
		},
	},

	Limits: Limits{
		Title:         "What satellites cannot see",
		Intro:         "Space data has hard limits. We state them up front, and every brief says what could not be seen.",
		LimitLabel:    "The limit",
		ResponseLabel: "What Tiefer does",
		Items: []Limit{
			{
				Limit:    "Optical satellites cannot see through cloud or in the dark. Filling the gap with generated imagery would be a guess, not an observation.",
				Response: "Switches to radar, which sees through cloud and at night, and says in the brief which days had no optical view.",
			},
			{
				Limit:    "Free night-light data has pixels of about 500 to 750 metres. A small or weak light often does not show.",
				Response: "Uses free data for strong lights and open fires, and proposes commercial night or thermal imagery for small objects.",
			},
			{
				Limit:    "What can be detected depends on the resolution of the sensor, on how often a satellite passes and on the weather.",
				Response: "Lists the passes it used and the ones it rejected, and states what could not be seen.",
			},
			{
				Limit:    "Commercial imagery costs money, comes with licence terms and can be restricted for some areas.",
				Response: "Orders only within your budget rules, checks the licence, and asks a person to approve anything beyond the rules.",
			},
			{
				Limit:    "General language models can misread satellite images and make things up.",
				Response: "Leaves the image analysis to specialised geospatial models. The language model plans and explains; it does not invent findings.",
			},
		},
	},

	Roadmap: Roadmap{
		Title:    "Where we are",
		NowLabel: "Now",
		Items: []Milestone{
			{Title: "Demo", When: "Oct to Dec 2026", Text: "Three question types on free satellite data: port activity, flood mapping, land change.", Now: true},
			{Title: "Pilots", When: "First half of 2027", Text: "First paying teams, watchlists and weekly briefs."},
			{Title: "Tasking", When: "Second half of 2027", Text: "Commercial imagery ordered by Tiefer, with human approval first."},
			{Title: "Sovereign", When: "2028", Text: "On-premise deployments for governments and large companies."},
			{Title: "More regions", When: "From 2028", Text: "The same analyst, new places and new sensors."},
		},
	},

	Name: Name{
		Line: "Tiefer is German for deeper.",
		Body: "We look past the surface of an image to what is actually happening on the ground, and we show you how we know.",
	},

	FAQ: FAQ{
		Title: "Questions and answers",
		Items: []Question{
			{
				Question: "Does Tiefer operate its own satellites?",
				Answer:   "No. Tiefer works with public satellite data, such as the Copernicus Sentinel missions, Landsat and VIIRS, and orders imagery from commercial satellite operators when the archive is not enough. It is not tied to any one provider.",
			},
			{
				Question: "Which places can I ask about?",
				Answer:   "Any place on Earth. We start sales and example briefs in the Caspian, Caucasus, Black Sea and Central Asia region, the region we know best.",
			},
			{
				Question: "How detailed and how recent is the data?",
				Answer:   "It depends on the sensor. Free optical and radar data show details down to about 10 metres; free thermal and night-light data are much coarser. When a question needs finer detail or a more recent image, Tiefer prepares an order for commercial imagery.",
			},
			{
				Question: "What happens when it is cloudy?",
				Answer:   "Optical satellites cannot see through cloud, so Tiefer switches to radar. It never fills the gap with generated pixels, and the brief says which days had no optical view.",
			},
			{
				Question: "Who decides when new imagery is bought?",
				Answer:   "You do. Your organisation sets budget rules. Orders within the rules can go ahead automatically; anything beyond them waits for approval by a person.",
			},
			{
				Question: "Can Tiefer run on our own servers?",
				Answer:   "Sovereign deployments that run on your own servers, inside your country, are on our roadmap for governments and large companies.",
			},
			{
				Question: "Can Tiefer be used to follow people?",
				Answer:   "No. Tiefer answers questions about places, not about people. Customers are screened, and access can be suspended when the acceptable use policy is broken.",
				Link:     Link{Label: "Read the acceptable use policy", Href: "/acceptable-use"},
			},
			{
				Question: "What does it cost?",
				Answer:   "We are setting prices together with our first pilot teams. Commercial imagery is priced by its provider and is only ordered within your budget rules.",
			},
			{
				Question: "Which languages does Tiefer answer in?",
				Answer:   "English, Azerbaijani, Turkish and Russian.",
			},
			{
				Question: "What do SAR, tasking and revisit mean?",
				Answer:   "SAR (synthetic aperture radar) is a radar that images the ground through cloud and at night. Tasking means ordering a satellite to image a specific place at a specific time. Revisit is how often a satellite passes over the same place.",
			},
		},
	},

	Contact: Contact{
		Title:         "We are early, and we are listening.",
		Body:          "We are looking for our first pilot teams: traders, insurers, analysts and public institutions with questions about the physical world. Tell us the question you would ask first.",
		FormLabel:     "Request access",
		NameLabel:     "Name",
		EmailLabel:    "Work email",
		OrgLabel:      "Organisation",
		RoleLabel:     "Your role",
		RoleEmpty:     "Choose one",
		QuestionLabel: "The first question you would ask Tiefer",
		QuestionHint:  "Name the place, the time window and what you want to know. Up to 2,000 characters.",
		OptionalLabel: "optional",
		Roles: []Option{
			{Value: "trader-analyst", Label: "Trader or analyst"},
			{Value: "insurer", Label: "Insurer"},
			{Value: "public-sector", Label: "Public sector"},
			{Value: "investor", Label: "Investor"},
			{Value: "other", Label: "Other"},
		},
		ConsentBefore: "I agree that Tiefer may use these details to reply to my message. See the ",
		ConsentLink:   "privacy notice",
		ConsentAfter:  ".",
		HoneypotLabel: "Leave this field empty",
		Submit:        "Request access",
		Sending:       "Sending",
		Success:       "Thank you. We will reply within two working days.",
		Errors: FormErrors{
			NameRequired:     "Please enter your name.",
			NameTooLong:      "Please shorten your name to 100 characters.",
			EmailRequired:    "Please enter your work email.",
			EmailInvalid:     "Please enter a valid email address, for example name@company.com.",
			OrgTooLong:       "Please shorten the organisation name to 150 characters.",
			RoleInvalid:      "Please choose a role from the list.",
			QuestionRequired: "Please tell us the first question you would ask.",
			QuestionTooLong:  "Please keep your question to 2,000 characters or fewer.",
			ConsentRequired:  "Please confirm that we may use these details to reply.",
			Invalid:          "This field contains characters that are not allowed.",
			Summary:          "Please check the highlighted fields.",
			Expired:          "This form has expired. Please send it again.",
			ErrorBefore:      "Something went wrong. Please email us at",
			ErrorAfter:       ".",
		},
		EmailUs:     "Email us",
		PreferEmail: "Prefer email?",
		LinkedIn:    "Follow Tiefer on LinkedIn",
	},

	Footer: Footer{
		CompanyTitle:  "Company",
		ContactTitle:  "Contact",
		Security:      "Security",
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
			Description: "Who operates the Tiefer website: a space technology startup in Baku, Azerbaijan.",
			Intro:       "Information about who operates this website.",
			Sections: []LegalSection{
				{
					Heading: "Status",
					Stage:   Founding,
					Paras: []string{
						"Tiefer is a startup in its founding stage and is not yet registered as a company. Registration as a limited liability company (MMC) in Azerbaijan is planned. Until then, this website is operated by the founder of Tiefer as a private individual.",
					},
				},
				{
					Heading: "Operator",
					Stage:   Founding,
					Facts: []Fact{
						{Label: "Name", Value: "{founder}"},
						{Label: "Role", Value: "Founder of Tiefer"},
						{Label: "Location", Value: "Baku, Azerbaijan [POSTAL_ADDRESS: confirm with counsel whether a full postal address must be stated here.]"},
						{Label: "Email", Value: "{contact_email}"},
						{Label: "Website", Value: "{site_url}"},
					},
				},
				{
					Heading: "Company",
					Stage:   Company,
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
					Heading: "Responsibility for content",
					Paras: []string{
						"{operator} is responsible for the content of this website. [CONTENT_RESPONSIBILITY: confirm with counsel whether a named person must be stated here.]",
					},
				},
				{
					Heading: "Source code and brand",
					Paras: []string{
						"The source code of this website is published under the Mozilla Public License 2.0. The Tiefer name and logo are not covered by that licence and may not be used to suggest endorsement or to brand a fork.",
						"The fonts Mozilla Headline and Mozilla Text are used under the SIL Open Font License 1.1.",
					},
				},
				{
					Heading: "External links",
					Paras: []string{
						"This website links to LinkedIn and, where shown, to our source code repository. Those sites are operated by other companies under their own terms and privacy policies.",
					},
				},
			},
			Updated: "Last updated: [LAST_UPDATED]",
		},

		Privacy: LegalPage{
			Title:       "Privacy notice",
			Description: "How the Tiefer website handles personal data: no cookies, no analytics, no requests to third parties.",
			Intro:       "This notice explains what personal data this website processes, why, and what rights you have. It covers this website and its contact form. It is structured to address the Law of the Republic of Azerbaijan on Personal Data of 11 May 2010 and, for visitors in the European Union, the General Data Protection Regulation (GDPR).",
			Sections: []LegalSection{
				{
					Heading: "Who is responsible",
					Stage:   Founding,
					Paras: []string{
						"Tiefer is not yet registered as a company. Until it is, the founder is responsible for the personal data this website processes. Once the company is registered, it takes over this role and this notice will be updated.",
					},
					Facts: []Fact{
						{Label: "Controller", Value: "{founder}, founder of Tiefer"},
						{Label: "Location", Value: "Baku, Azerbaijan"},
						{Label: "Email", Value: "{contact_email}"},
						{Label: "EU representative (Article 27 GDPR)", Value: "[EU_REPRESENTATIVE: name and address, or state that none is required]"},
					},
				},
				{
					Heading: "Who is responsible",
					Stage:   Company,
					Facts: []Fact{
						{Label: "Controller", Value: "{legal_name}"},
						{Label: "Legal form", Value: "{legal_form}"},
						{Label: "Address", Value: "{address}"},
						{Label: "Email", Value: "{contact_email}"},
						{Label: "EU representative (Article 27 GDPR)", Value: "[EU_REPRESENTATIVE: name and address, or state that none is required]"},
					},
				},
				{
					Heading: "What this website does not do",
					List: []string{
						"It does not set cookies.",
						"It does not use analytics, tracking pixels or advertising.",
						"It does not load fonts, scripts or images from other servers. Everything, including the fonts, is served from our own server.",
						"It does not embed social media plugins. The LinkedIn link is a plain link: LinkedIn receives data only if you follow it.",
					},
				},
				{
					Heading: "Server logs",
					Paras: []string{
						"To run the website and fix errors, our server records each request: the time, the page requested, the response status, the response size and how long the response took. It does not write IP addresses to its logs.",
						"To protect the contact form from abuse, the server keeps the IP address of each sender (for IPv6, only its network part) in memory for up to one hour, to limit the number of messages per connection. It is never written to disk or to the logs, and it is gone when the server restarts.",
						"Our hosting provider, {hosting_provider} ({hosting_country}), may process connection data such as IP addresses to operate its network. [HOSTING_PROVIDER_LOGS: describe what the provider records and for how long.]",
					},
				},
				{
					Heading: "The contact form",
					Paras: []string{
						"If you use the contact form, we receive your name, your work email address, your organisation and role if you give them, your question, and your confirmation that we may reply.",
						"We use these details only to reply to you and to discuss access to Tiefer.",
						"The website does not store your message in a database. It sends it by email, through our email provider {smtp_provider} ({smtp_country}), to the Tiefer mailbox, and keeps no copy on the web server. The content of your message is not written to the server logs.",
						"If you email us directly, we receive your email address and what you write.",
					},
				},
				{
					Heading: "Legal basis",
					Paras: []string{
						"[LEGAL_BASIS: to be determined with counsel, under the Law of the Republic of Azerbaijan on Personal Data and under Article 6 GDPR.]",
					},
				},
				{
					Heading: "Who receives your data",
					Paras: []string{
						"Our hosting provider {hosting_provider} and our email provider {smtp_provider} process data on our behalf. We do not sell personal data and do not share it with advertisers. [RECIPIENTS: confirm the list of recipients and processing agreements.]",
					},
				},
				{
					Heading: "Transfers to other countries",
					Paras: []string{
						"Tiefer is based in Azerbaijan. Your data may be processed in {hosting_country} (hosting) and {smtp_country} (email). [INTERNATIONAL_TRANSFERS: describe the transfer safeguards under the GDPR and the cross-border transfer rules under Azerbaijani law.]",
					},
				},
				{
					Heading: "How long we keep it",
					List: []string{
						"Contact messages: [RETENTION_PERIOD].",
						"Server logs: [LOG_RETENTION_PERIOD].",
						"The in-memory rate limit record: up to one hour.",
					},
				},
				{
					Heading: "Your rights",
					Paras: []string{
						"You can ask us for access to your personal data and for its correction or deletion, and you can ask us to restrict its processing. You can object to processing and ask for your data in a portable format. Where processing is based on consent, you can withdraw it at any time; this does not affect processing that happened before.",
						"To use these rights, write to {contact_email}. [DATA_SUBJECT_RIGHTS: confirm the rights and response times under Azerbaijani law and the GDPR.]",
					},
				},
				{
					Heading: "Complaints",
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
			Intro:       "Tiefer answers questions about places, not about people. This policy sets out how Tiefer may be used. It applies to every pilot, account and brief. [AUP_SCOPE: confirm how this policy becomes part of customer contracts.]",
			Sections: []LegalSection{
				{
					Heading: "What we commit to",
					List: []string{
						"Every brief labels what was observed separately from what was inferred.",
						"Generated imagery is never shown as evidence.",
						"Every finding links to its source image, sensor, date and processing steps.",
						"Spending on new imagery stays under human control.",
						"Every brief states what could not be seen.",
					},
				},
				{
					Heading: "What is not allowed",
					List: []string{
						"Tracking, identifying or profiling individual people, including following a specific person or household over time.",
						"Any use that violates human rights, including surveillance of people because of their ethnicity, religion, political opinion or other protected characteristics.",
						"Planning or supporting violence, or selecting targets for attack.",
						"Any use that breaches applicable sanctions or export control laws. [SANCTIONS_REGIMES: list the sanctions and export control regimes Tiefer applies.]",
						"Presenting inferred findings as observed facts, or presenting generated imagery as evidence.",
						"Using satellite imagery against the licence terms of its provider.",
						"Getting around budget rules or approval steps for tasking.",
					},
				},
				{
					Heading: "Customer screening",
					Paras: []string{
						"We screen every customer before granting access, including against sanctions lists, and we may ask how Tiefer will be used. We may decline any request. [SCREENING_PROCESS: describe the screening steps.]",
					},
				},
				{
					Heading: "Suspension",
					Paras: []string{
						"We may suspend or end access if we believe this policy has been breached, and we may refuse or cancel imagery orders. [ENFORCEMENT_TERMS: align with customer contracts.]",
					},
				},
				{
					Heading: "Reporting misuse",
					Paras: []string{
						"If you believe Tiefer is being misused, write to {contact_email}.",
					},
				},
			},
			Updated: "Last updated: [LAST_UPDATED]",
		},
	},

	NotFound: NotFound{
		Label: "ERROR 404",
		Title: "This page does not exist.",
		Body:  "The link may be old, or the address may contain a typo.",
		Home:  "Back to the start page",
	},
}
