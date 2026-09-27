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
			{Icon: "optical", Label: "OPTICAL", Text: "What the surface looks like. Sentinel-2, Landsat and commercial imagery."},
			{Icon: "radar", Label: "RADAR", Text: "Sees through cloud and at night. Sentinel-1 and commercial SAR."},
			{Icon: "thermal", Label: "THERMAL", Text: "Heat: fires, industrial activity, surface temperature."},
			{Icon: "night", Label: "NIGHT LIGHTS", Text: "Where light appears or disappears after dark."},
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
	},

	Demo: Demo{
		Title:     "Every finding shows its source.",
		Body:      "An answer you cannot check is not intelligence. Tiefer labels what a satellite observed separately from what a model inferred, links every finding to its source image, and never presents generated pixels as evidence.",
		CardLabel: "Illustrative example. Not real data.",
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
		QuestionLabel: "Example question",
		Cards: []UseCase{
			{
				Title:    "Commodity traders",
				Text:     "Ports, terminals, storage and crops: see supply move before the statistics do.",
				Question: "How full are the tanks at this terminal compared with last month?",
			},
			{
				Title:    "Insurers",
				Text:     "Floods, fires and damage mapped in hours, with the images to back each claim.",
				Question: "Which insured sites in this district were under water on Tuesday?",
			},
			{
				Title:    "Governments",
				Text:     "Sovereign deployments that keep sensitive questions and data inside the country.",
				Question: "What changed at this infrastructure site since January?",
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

	Roadmap: Roadmap{
		Title:    "Where we are",
		NowLabel: "Now",
		Items: []Milestone{
			{Title: "Demo", Text: "Three question types on free satellite data: port activity, flood mapping, land change.", Now: true},
			{Title: "Pilots", Text: "First paying teams, watchlists and weekly briefs."},
			{Title: "Tasking", Text: "Commercial imagery ordered by Tiefer, with human approval first."},
			{Title: "Sovereign", Text: "On-premise deployments for governments and large companies."},
			{Title: "More regions", Text: "The same analyst, new places and new sensors."},
		},
	},

	Name: Name{
		Line: "Tiefer is German for deeper.",
		Body: "We look past the surface of an image to what is actually happening on the ground, and we show you how we know.",
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
		QuestionHint:  "Up to 2,000 characters.",
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
		Copyright:     "© 2026 Tiefer",
		LegalNotice:   "Legal notice",
		Privacy:       "Privacy",
		AcceptableUse: "Acceptable use",
		LinkedIn:      "LinkedIn",
		SourceCode:    "Source code",
		Tagline:       "See deeper. Built in Baku.",
	},

	NotFound: NotFound{
		Label: "ERROR 404",
		Title: "This page does not exist.",
		Body:  "The link may be old, or the address may contain a typo.",
		Home:  "Back to the start page",
	},
}
