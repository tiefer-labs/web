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
				Question: "How does Tiefer deal with AI regulation?",
				Answer:   "We follow the EU Artificial Intelligence Act and the laws of Azerbaijan. Tiefer tells users they are working with AI, marks every brief as generated by AI, labels what was observed separately from what was inferred, and leaves decisions to people. It is not built for any use the AI Act prohibits.",
				Link:     Link{Label: "Read the AI policy", Href: "/ai-policy"},
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
		AIPolicy:      "AI policy",
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
			Description: "How the Tiefer website handles personal data under the laws of Azerbaijan and the GDPR: no cookies, no analytics, no requests to third parties.",
			Intro:       "This notice explains what personal data this website processes, why, on what legal basis, and what rights you have. It covers this website and its contact form. It follows the laws of the Republic of Azerbaijan and, for people in the European Union, the General Data Protection Regulation (GDPR).",
			Sections: []LegalSection{
				{
					Heading: "Laws this notice follows",
					List: []string{
						"Azerbaijan: the Constitution of the Republic of Azerbaijan (Article 32, the right to privacy), the Law on Personal Data of 11 May 2010, and the Law on Information, Informatisation and Protection of Information of 3 April 1998.",
						"The Council of Europe Convention for the Protection of Individuals with regard to Automatic Processing of Personal Data (Convention 108), to which Azerbaijan is a party.",
						"European Union: the General Data Protection Regulation (Regulation (EU) 2016/679). It applies to us where we offer our services to people in the EU (Article 3(2) GDPR).",
						"European Union: the ePrivacy Directive (Directive 2002/58/EC) and the national laws that implement it, on storing and reading information on your device.",
						"The national laws of EU member states that add to the GDPR, where they apply.",
						"[APPLICABLE_LAW: confirm this list with counsel in Azerbaijan and in the EU, including amendments made after this notice was written.]",
					},
				},
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
					Paras: []string{
						"Under Article 5(3) of the ePrivacy Directive, storing or reading information on your device needs your consent unless it is strictly necessary. This website stores and reads nothing, so it asks for no consent and shows no cookie banner.",
					},
					List: []string{
						"It does not set cookies and does not use local storage, session storage or any other storage in your browser.",
						"It does not use analytics, tracking pixels, fingerprinting or advertising.",
						"It does not load fonts, scripts or images from other servers. Everything, including the fonts, is served from our own server.",
						"It does not embed social media plugins. The LinkedIn link is a plain link: LinkedIn receives data only if you follow it.",
						"It does not build profiles of visitors.",
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
						"If you use the contact form, we receive your name, your work email address, your organisation and role if you give them, your question, and your consent that we may reply.",
						"Name, email address, question and consent are needed to reply to you; organisation and role are optional. You do not have to use the form, and you can write to us by email instead.",
						"We use these details only to reply to you and to discuss access to Tiefer. We do not use them for marketing.",
						"The website does not store your message in a database. It sends it by email, through our email provider {smtp_provider} ({smtp_country}), to the Tiefer mailbox, and keeps no copy on the web server. The content of your message is not written to the server logs.",
						"If you email us directly, we receive your email address and what you write.",
					},
				},
				{
					Heading: "Purposes and legal bases",
					List: []string{
						"Running the website and keeping it secure, including the in-memory rate limit record: our legitimate interest in a working and secure website (Article 6(1)(f) GDPR).",
						"Replying to a message about access to Tiefer: steps you ask for before a possible contract (Article 6(1)(b) GDPR). For other questions: our legitimate interest in answering them (Article 6(1)(f) GDPR).",
						"Under the Law on Personal Data of Azerbaijan: your consent, which you give with the checkbox of the contact form or by writing to us, and which you can withdraw at any time. [AZ_LEGAL_BASIS: confirm the basis for the rate limit record under Azerbaijani law.]",
						"[LEGAL_BASIS: confirm these legal bases with counsel.]",
					},
				},
				{
					Heading: "Who receives your data",
					Paras: []string{
						"Our hosting provider {hosting_provider} and our email provider {smtp_provider} process data on our behalf, under agreements that bind them to our instructions (Article 28 GDPR). We do not sell personal data and do not share it with advertisers. We disclose it to authorities only where the law requires it. [RECIPIENTS: confirm the list of recipients and processing agreements.]",
					},
				},
				{
					Heading: "Transfers to other countries",
					Paras: []string{
						"Tiefer is based in Azerbaijan. The European Commission has not adopted an adequacy decision for Azerbaijan. When you write to us from the EU, you send your message directly to us in Azerbaijan.",
						"Where our providers process your data in another country, {hosting_country} (hosting) and {smtp_country} (email), we use the safeguards required by Chapter V of the GDPR, such as the standard contractual clauses of the European Commission. Under the Law on Personal Data, personal data may leave Azerbaijan only where the receiving country protects it at a comparable level or another condition of that law is met. [INTERNATIONAL_TRANSFERS: confirm the safeguards under the GDPR and under Azerbaijani law.]",
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
					Heading: "Security",
					Paras: []string{
						"The website is served only over HTTPS, and browsers are told never to connect to it without encryption. Messages from the contact form travel to our email provider over an encrypted connection. The website keeps no database of messages, so there is none to lose. [MAILBOX_SECURITY: describe how the mailbox is protected, for example with two-factor authentication.]",
					},
				},
				{
					Heading: "Automated decisions and AI",
					Paras: []string{
						"This website makes no decisions about you based solely on automated processing, including profiling (Article 22 GDPR). It does not pass your message to an AI system: the message is sent by email and read by a person at Tiefer. [AI_IN_MAILBOX: confirm that no AI tool processes messages in the mailbox, or name it here.]",
						"How the Tiefer product itself uses AI is explained in the AI policy.",
					},
				},
				{
					Heading: "Your rights",
					Paras: []string{
						"Under the GDPR and the Law on Personal Data you have the rights below. To use them, write to {contact_email}. We answer within one month; for complex requests this can be extended by two further months, and we will tell you if it is (Article 12(3) GDPR). [DATA_SUBJECT_RIGHTS: confirm the rights and response times under the Law on Personal Data.]",
					},
					List: []string{
						"Access: to know whether we process your data and to receive a copy (Article 15 GDPR).",
						"Correction of inaccurate or incomplete data (Article 16 GDPR).",
						"Deletion (Article 17 GDPR).",
						"Restriction of processing, or blocking under the Law on Personal Data (Article 18 GDPR).",
						"Portability: to receive your data in a machine-readable format, where processing is based on consent or a contract (Article 20 GDPR).",
						"Objection to processing based on our legitimate interests (Article 21 GDPR).",
						"Withdrawal of consent at any time, without affecting processing that happened before (Article 7(3) GDPR).",
					},
				},
				{
					Heading: "Complaints",
					Paras: []string{
						"You can complain to a data protection supervisory authority, in particular in the EU member state where you live or work or where you believe the infringement took place (Article 77 GDPR). In Azerbaijan: [SUPERVISORY_AUTHORITY_AZ]. You can also go to court.",
					},
				},
				{
					Heading: "Children",
					Paras: []string{
						"This website is meant for professionals and is not directed at children. We do not knowingly collect personal data from children.",
					},
				},
				{
					Heading: "Changes to this notice",
					Paras: []string{
						"We update this notice when the website, our providers or the law change. The date below shows the last change.",
					},
				},
			},
			Updated: "Last updated: [LAST_UPDATED]",
		},

		AI: LegalPage{
			Title:       "AI policy",
			Description: "How Tiefer develops and offers AI under the EU Artificial Intelligence Act and the laws of Azerbaijan.",
			Intro:       "Tiefer is an AI system that analyses satellite data and writes intelligence briefs. This policy explains which laws on artificial intelligence we follow in the European Union and in Azerbaijan, and how we meet them. It applies to the Tiefer product, its demo and its pilots.",
			Sections: []LegalSection{
				{
					Heading: "Laws this policy follows",
					List: []string{
						"European Union: the Artificial Intelligence Act (Regulation (EU) 2024/1689), in force since 1 August 2024, with obligations that apply in stages. It applies to providers outside the EU whose AI systems are offered in the EU or whose output is used there (Article 2(1)).",
						"European Union: the GDPR, wherever an AI system processes personal data, including its rules on automated decisions (Article 22 GDPR).",
						"European Union: the Product Liability Directive (Directive (EU) 2024/2853), which treats software, including AI systems, as a product, for products placed on the market from 9 December 2026.",
						"Azerbaijan: there is no law specific to artificial intelligence yet. AI systems are subject to the general laws, in particular the Law on Personal Data and the Law on Information, Informatisation and Protection of Information, and to the national Artificial Intelligence Strategy for 2025 to 2028.",
						"Export control and sanctions law in both jurisdictions, as set out in the acceptable use policy.",
						"[AI_LAW_REVIEW: confirm this list with counsel, including changes to the application dates of the AI Act and any new Azerbaijani AI law.]",
					},
				},
				{
					Heading: "Our role",
					Paras: []string{
						"Under the AI Act, we are the provider of the Tiefer AI system: we develop it and offer it under our own name. Our customers are deployers: they use Tiefer under their own authority and decide what to do with its briefs.",
						"Where Tiefer is built on general-purpose AI models from other companies, the providers of those models carry the obligations for them, and we build on the documentation they give us. [AI_ACT_ROLE: confirm the roles and list the general-purpose models in use.]",
					},
				},
				{
					Heading: "Risk classification",
					Paras: []string{
						"Tiefer analyses places and objects in satellite data. It is not designed to identify, recognise or track people, and it does not carry out any practice the AI Act prohibits (Article 5), such as social scoring, biometric identification or categorisation, or emotion recognition.",
						"Tiefer is not designed or offered for the high-risk uses listed in Annex III of the AI Act, such as biometrics, law enforcement, migration, asylum and border control, or the administration of justice. A customer who wants to use it in one of these areas needs our written agreement first. We then assess whether the rules for high-risk AI systems apply, and meet them before any use. [AI_ACT_CLASSIFICATION: confirm the classification, including for public sector and infrastructure use cases.]",
					},
				},
				{
					Heading: "Transparency",
					List: []string{
						"People who use Tiefer are told that they are working with an AI system (Article 50(1) of the AI Act).",
						"Every brief is marked as generated by an AI system, in a form people can read and in a machine-readable form (Article 50(2) of the AI Act).",
						"Every brief labels what was observed separately from what was inferred.",
						"Every finding links to its source image, sensor, date and processing steps.",
						"Generated imagery is never shown as evidence, and Tiefer does not create images that could be taken for real satellite observations.",
						"Every brief states what could not be seen.",
					},
				},
				{
					Heading: "Human oversight",
					Paras: []string{
						"Tiefer supports decisions; people make them. A brief is an analysis for a person to review, not a decision about anyone.",
						"Tiefer orders new imagery automatically only within the budget rules a customer sets. Anything beyond them waits for approval by a person.",
						"No decision with legal or similarly significant effects on a person is made by Tiefer alone (Article 22 GDPR).",
					},
				},
				{
					Heading: "Data",
					Paras: []string{
						"Tiefer works with satellite data: public data such as the Copernicus Sentinel missions, Landsat and VIIRS, and commercial imagery used under the licence terms of its provider. It does not use satellite data to identify people.",
						"Questions and briefs of customers are not used to train AI models without the customer's written agreement. [CUSTOMER_DATA: confirm with counsel and align with customer contracts.]",
						"Where a question or its data contains personal data, the GDPR and the Law on Personal Data apply, and we process only what the question needs.",
					},
				},
				{
					Heading: "Accuracy and limits",
					Paras: []string{
						"Satellite data has hard limits: cloud, resolution, how often a satellite passes, and licence terms. General language models can misread satellite images, so Tiefer leaves the image analysis to specialised geospatial models; the language model plans and explains.",
						"We test Tiefer on known cases before each release and record its known limits in the documentation for customers. [AI_TESTING: describe the testing and documentation process.]",
					},
				},
				{
					Heading: "AI literacy",
					Paras: []string{
						"Everyone who builds, sells or supports Tiefer learns how it works, where it fails and what the law requires (Article 4 of the AI Act). Customers receive documentation on what Tiefer can and cannot do.",
					},
				},
				{
					Heading: "Reporting problems",
					Paras: []string{
						"If a brief is wrong, if you believe Tiefer is being misused, or if you have a question about this policy, write to {contact_email}. We record every report, look into it and correct Tiefer where needed. Where the law requires it, we report serious incidents to the competent authorities.",
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
