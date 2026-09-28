// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package content

// Site is the complete copy of the site in one locale.
type Site struct {
	Locale     Locale
	Meta       Meta
	UI         UI
	Nav        []Link
	CTA        Link
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

// Locale identifies a language version of the site. The default locale
// has an empty Prefix and is served at the root; others live under their
// prefix, for example /az.
type Locale struct {
	Code   string // BCP 47 tag for lang and hreflang, for example "en"
	Prefix string // URL prefix, "" for the default locale
}

// Meta holds page metadata.
type Meta struct {
	SiteName    string
	Title       string // title of the index page
	Description string
	OGImageAlt  string
	ThemeColor  string
}

// UI holds interface strings that are not part of one section.
type UI struct {
	SkipLink   string
	NavLabel   string
	MenuOpen   string
	MenuClose  string
	FooterNav  string
	LegalLabel string // label of the navigation on legal pages
	BackHome   string
}

// Link is a label and a target: a fragment ("#contact"), a site path
// ("/privacy") or an absolute URL.
type Link struct {
	Label string
	Href  string
}

// Item is a title with a sentence.
type Item struct {
	Title string
	Text  string
}

// Hero is the dark opening section.
type Hero struct {
	Position  string // technical label, top right on desktop
	AlertTag  string // telemetry label on the downlink line
	Stored    string // telemetry label at the ground point
	Eyebrow   string
	Title     string
	Lead      string
	Primary   Link
	Secondary Link
}

// Problem is the statement of the problem.
type Problem struct {
	Label string
	Title string
	Intro string
	Items []Item
}

// Stage is one onboard stage in the product table.
type Stage struct {
	Number string
	Name   string
	Text   string
}

// Product lists the four onboard stages.
type Product struct {
	Label  string
	Title  string
	Body   string
	Stages []Stage
	Note   string
}

// How is the capture-to-alert flow and the ground products.
type How struct {
	Label       string
	Title       string
	Steps       []Item
	Note        string
	GroundTitle string
	Ground      []Item
}

// Field is one row of the illustrative alert packet.
type Field struct {
	Name  string
	Value string
	Kind  string // "", "observed" or "inferred"
}

// Demo is the illustrative alert packet.
type Demo struct {
	Label      string
	Title      string
	Body       string
	CardLabel  string // "Illustrative example. Not real data."
	CardTitle  string // caption of the packet table
	Fields     []Field
	SceneLabel string
	AlertLabel string
}

// UseCase is one event type.
type UseCase struct {
	Label string
	Title string
	Text  string
}

// UseCases lists the event types.
type UseCases struct {
	Label string
	Title string
	Cases []UseCase
	Note  string
}

// Principle is one clause of the principles.
type Principle struct {
	Number    string
	Statement string
	Text      string
}

// Principles lists how Tiefer builds flight software.
type Principles struct {
	Label string
	Title string
	Items []Principle
}

// Milestone is one node of the roadmap.
type Milestone struct {
	Title string
	Text  string
	Now   bool
}

// Roadmap is the timeline.
type Roadmap struct {
	Label    string
	Title    string
	NowLabel string
	Items    []Milestone
}

// Name explains the name.
type Name struct {
	Label string
	Line  string
	Body  string
}

// Option is one choice of a select field.
type Option struct {
	Value string
	Label string
}

// Contact is the contact section and its form.
type Contact struct {
	Label          string
	Title          string
	Body           string
	FormLabel      string // accessible name of the form
	NameLabel      string
	EmailLabel     string
	OrgLabel       string
	RoleLabel      string
	RoleEmpty      string
	Roles          []Option
	MessageLabel   string
	MessageHint    string
	Optional       string
	ConsentBefore  string
	ConsentLink    string
	ConsentAfter   string
	HoneypotLabel  string
	Submit         string
	Sending        string
	Success        string
	ErrorBefore    string // followed by the contact email as a link
	ErrorAfter     string
	Invalid        string // summary when fields need attention
	Expired        string // the form token expired or was used
	Limited        string // too many messages
	EmailUs        string
	PreferEmail    string
	LinkedIn       string
	NameMissing    string
	NameTooLong    string
	EmailMissing   string
	EmailInvalid   string
	OrgTooLong     string
	RoleInvalid    string
	MessageMissing string
	MessageTooLong string
	ConsentMissing string
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

// Fact is a label and value on a legal page; Value may contain {tokens}.
type Fact struct {
	Label string
	Value string
}

// LegalSection is one titled part of a legal page. Paragraphs, list
// items and fact values may contain {tokens}.
type LegalSection struct {
	Title string
	Facts []Fact
	Paras []string
	List  []string
}

// LegalPage is one legal page.
type LegalPage struct {
	Title       string
	Description string
	Intro       string
	Sections    []LegalSection
	Updated     string
}

// Legal holds the legal pages and their shared strings. Words in curly
// braces, such as {legal_name}, are filled from configuration; words in
// square brackets, such as [RETENTION_PERIOD], are placeholders for a
// lawyer to replace here before launch.
type Legal struct {
	ReviewNote    string
	Contents      string
	Notice        LegalPage
	Privacy       LegalPage
	AcceptableUse LegalPage
}

// NotFound is the 404 page.
type NotFound struct {
	Code  string
	Title string
	Body  string
	Home  string
}

// Locales lists every locale the site serves. The first is the default
// and is served without a URL prefix.
var Locales = []*Site{&English}

// Default returns the default locale.
func Default() *Site { return Locales[0] }
