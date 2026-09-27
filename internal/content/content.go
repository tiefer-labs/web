// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

// Package content holds every piece of website copy as typed Go values,
// one file per locale (en.go, later az.go, de.go, ...). Templates only
// read from these structs, so changing a sentence here changes the site.
package content

// Locale identifies one language version of the site.
type Locale struct {
	Code     string // BCP 47 language tag, used in lang and hreflang
	Prefix   string // URL prefix: "" for the default locale, "/az" and so on for others
	Name     string // Name of the language in that language
	OGLocale string // Open Graph locale, for example en_GB
}

// Link is a label with a target.
type Link struct {
	Label string
	Href  string
}

// Item is a titled piece of text, used for grids and lists.
type Item struct {
	Title string
	Text  string
}

// Site is the complete copy of the website in one locale.
type Site struct {
	Locale     Locale
	Meta       Meta
	UI         UI
	Nav        Nav
	Hero       Hero
	Problem    Problem
	Product    Product
	How        How
	Demo       Demo
	UseCases   UseCases
	Principles Principles
	Roadmap    Roadmap
	Name       Name
	Contact    Contact
	Footer     Footer
	Legal      Legal
	NotFound   NotFound
}

// Meta holds the document metadata of the index page.
type Meta struct {
	SiteName    string
	Title       string
	Description string
	OGImageAlt  string
}

// UI holds small interface strings.
type UI struct {
	SkipLink  string
	MenuOpen  string
	MenuClose string
	LogoAlt   string
	NavLabel  string
	FooterNav string
}

// Nav is the main navigation.
type Nav struct {
	Links []Link // Href is an anchor on the index page, for example "product"
	CTA   Link
}

// Hero is the dark opening section.
type Hero struct {
	Coordinates     string // technical label, desktop only
	ArtLabels       []string
	Eyebrow         string
	Title           string
	Lead            string
	Question        string
	QuestionCaption string
	Primary         Link
	Secondary       Link
}

// Problem describes why answers from satellite data are hard to get.
type Problem struct {
	Title string
	Intro string
	Items []Item
}

// Layer is one of the four data layers. Icon selects the inline SVG icon.
type Layer struct {
	Icon  string // optical, radar, thermal or night
	Label string
	Text  string
	Specs []Fact // data sources and limits, shown under the text
}

// Product explains the four layers.
type Product struct {
	Title  string
	Body   string
	Layers []Layer
}

// How lists the steps from question to brief.
type How struct {
	Title string
	Steps []Item
	Note  Item // how the work is split between models
}

// BriefRow is one row of the illustrative brief. Tag is "observed",
// "inferred" or empty; TagLabel is its visible text.
type BriefRow struct {
	Label    string
	Tag      string
	TagLabel string
	Value    string
}

// PassDay lists the satellite passes of one day in the demo strip.
type PassDay struct {
	Radar    int
	Optical  int
	Rejected int
}

// Strip is the decorative 14 day pass strip under the brief.
type Strip struct {
	Days           []PassDay
	FirstDay       string
	LastDay        string
	LegendRadar    string
	LegendOptical  string
	LegendRejected string
}

// Demo is the illustrative brief section.
type Demo struct {
	Title     string
	Body      string
	CardLabel string // caption of the brief table
	// ContentsTitle and Contents list what every brief contains.
	ContentsTitle string
	Contents      []string
	Rows          []BriefRow
	Strip         Strip
}

// UseCase is one customer group.
type UseCase struct {
	Title    string
	Text     string
	Question string
}

// UseCases lists who Tiefer is for.
type UseCases struct {
	Title         string
	QuestionLabel string
	Cards         []UseCase
	Note          string
}

// Principles lists how Tiefer works.
type Principles struct {
	Title string
	Items []Item
}

// Milestone is one node on the roadmap.
type Milestone struct {
	Title string
	Text  string
	Now   bool
}

// Roadmap shows where Tiefer is.
type Roadmap struct {
	Title    string
	NowLabel string
	Items    []Milestone
}

// Name explains the name.
type Name struct {
	Line string
	Body string
}

// Option is one choice of a select field.
type Option struct {
	Value string
	Label string
}

// FormErrors are the messages shown next to invalid form fields and in
// the form status area.
type FormErrors struct {
	NameRequired     string
	NameTooLong      string
	EmailRequired    string
	EmailInvalid     string
	OrgTooLong       string
	RoleInvalid      string
	QuestionRequired string
	QuestionTooLong  string
	ConsentRequired  string
	Invalid          string // general characters not allowed
	Summary          string // shown above the form when fields are invalid
	Expired          string
	ErrorBefore      string // followed by the contact email as a link
	ErrorAfter       string
}

// Contact is the contact section and form.
type Contact struct {
	Title         string
	Body          string
	FormLabel     string
	NameLabel     string
	EmailLabel    string
	OrgLabel      string
	RoleLabel     string
	RoleEmpty     string
	Roles         []Option
	QuestionLabel string
	QuestionHint  string
	OptionalLabel string
	ConsentBefore string
	ConsentLink   string
	ConsentAfter  string
	HoneypotLabel string
	Submit        string
	Sending       string
	Success       string
	Errors        FormErrors
	EmailUs       string
	PreferEmail   string
	LinkedIn      string
}

// Footer is the site footer.
type Footer struct {
	Copyright     string
	LegalNotice   string
	Privacy       string
	AcceptableUse string
	LinkedIn      string
	SourceCode    string
	Tagline       string
}

// Fact is a label and value pair on a legal page. Value may contain
// {tokens} that are filled from configuration (see Legal.Tokens).
type Fact struct {
	Label string
	Value string
}

// LegalSection is one titled block of a legal page. Paragraphs, list
// items and fact values may contain {tokens}.
type LegalSection struct {
	Heading string
	Paras   []string
	List    []string
	Facts   []Fact
}

// LegalPage is the complete text of one legal page.
type LegalPage struct {
	Title       string
	Description string
	Intro       string
	Sections    []LegalSection
	Updated     string
}

// Legal holds the three legal pages and their shared strings.
//
// Tokens in curly braces are replaced with configuration values:
// {legal_name}, {legal_form}, {address}, {tax_id}, {registration},
// {director}, {contact_email}, {site_url}, {hosting_provider},
// {hosting_country}, {smtp_provider}, {smtp_country}.
// Tokens in square brackets, such as [RETENTION_PERIOD], are placeholders
// that a lawyer must replace in this file before launch.
type Legal struct {
	ReviewNote    string
	Contents      string // heading of the list of sections
	Notice        LegalPage
	Privacy       LegalPage
	AcceptableUse LegalPage
}

// NotFound is the 404 page.
type NotFound struct {
	Label string
	Title string
	Body  string
	Home  string
}

// Locales lists every locale the site serves. The first entry is the
// default locale and is served without a URL prefix.
var Locales = []*Site{&English}

// Default returns the default locale.
func Default() *Site { return Locales[0] }
