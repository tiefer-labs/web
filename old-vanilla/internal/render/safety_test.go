// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package render

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"text/template/parse"

	"github.com/tiefer-labs/web/web"
)

// allowedFuncs are the functions templates may call: the comparison and
// logic builtins and the functions of Funcs. Escaping builtins (html, js,
// urlquery), print functions and call are not allowed: html/template
// escapes by context, and building markup or URLs with printf would work
// around it.
var allowedFuncs = map[string]bool{
	"and": true, "or": true, "not": true, "eq": true, "ne": true,
	"lt": true, "le": true, "gt": true, "ge": true, "len": true, "index": true,
}

// TestTemplateSafety walks the parse tree of every template and fails on
// functions outside the allowed list. It also rejects markup patterns
// that would weaken the CSP or open the page to injection.
func TestTemplateSafety(t *testing.T) {
	for name := range Funcs(&Assets{}) {
		allowedFuncs[name] = true
	}
	root := web.Templates()
	forbidden := map[string]*regexp.Regexp{
		"inline event handler":         regexp.MustCompile(`(?i)\son[a-z]+\s*=`),
		"inline style attribute":       regexp.MustCompile(`(?i)\sstyle\s*=`),
		"style element":                regexp.MustCompile(`(?i)<style\b`),
		"javascript: URL":              regexp.MustCompile(`(?i)javascript:`),
		"inline script":                regexp.MustCompile(`(?i)<script\b(?:[^>]*\bsrc=)?[^>]*>\s*[^<\s]`),
		"frame or object":              regexp.MustCompile(`(?i)<(iframe|object|embed|frame)\b`),
		"target attribute":             regexp.MustCompile(`(?i)\starget\s*=`),
		"base element":                 regexp.MustCompile(`(?i)<base\b`),
		"third-party http(s) resource": regexp.MustCompile(`(?i)\b(src|href)="https?://`),
	}
	err := fs.WalkDir(root, ".", func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		b, err := fs.ReadFile(root, p)
		if err != nil {
			return err
		}
		src := string(b)
		for what, re := range forbidden {
			if loc := re.FindStringIndex(src); loc != nil {
				t.Errorf("%s:%d: %s", p, strings.Count(src[:loc[0]], "\n")+1, what)
			}
		}
		trees, err := parse.Parse(p, src, "{{", "}}", builtinsForParse())
		if err != nil {
			t.Errorf("%s: %v", p, err)
			return nil
		}
		for _, tree := range trees {
			walk(tree.Root, func(n parse.Node) {
				if id, ok := n.(*parse.IdentifierNode); ok && !allowedFuncs[id.Ident] {
					loc, _ := tree.ErrorContext(n)
					t.Errorf("%s: function %q is not allowed in templates", loc, id.Ident)
				}
			})
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}

// builtinsForParse lists every name the parser must accept, allowed or
// not, so that disallowed calls reach the walker and are reported.
func builtinsForParse() map[string]any {
	m := map[string]any{}
	for n := range allowedFuncs {
		m[n] = true
	}
	for _, n := range []string{"html", "js", "urlquery", "print", "printf", "println", "call", "slice"} {
		m[n] = true
	}
	return m
}

func walk(n parse.Node, fn func(parse.Node)) {
	if n == nil {
		return
	}
	fn(n)
	switch n := n.(type) {
	case *parse.ListNode:
		if n == nil {
			return
		}
		for _, c := range n.Nodes {
			walk(c, fn)
		}
	case *parse.ActionNode:
		walk(n.Pipe, fn)
	case *parse.PipeNode:
		if n == nil {
			return
		}
		for _, c := range n.Cmds {
			walk(c, fn)
		}
	case *parse.CommandNode:
		for _, a := range n.Args {
			walk(a, fn)
		}
	case *parse.IfNode:
		walk(n.Pipe, fn)
		walk(n.List, fn)
		walk(n.ElseList, fn)
	case *parse.RangeNode:
		walk(n.Pipe, fn)
		walk(n.List, fn)
		walk(n.ElseList, fn)
	case *parse.WithNode:
		walk(n.Pipe, fn)
		walk(n.List, fn)
		walk(n.ElseList, fn)
	case *parse.TemplateNode:
		walk(n.Pipe, fn)
	}
}

// unsafeTypes are the html/template types that mark a value as already
// safe. Converting to them skips escaping.
var unsafeTypes = map[string]bool{"HTML": true, "HTMLAttr": true, "JS": true, "JSStr": true, "CSS": true, "URL": true, "Srcset": true}

// allowedConversions names the only functions that may convert a value
// to one of those types, with the reason.
var allowedConversions = map[string]string{
	"internal/server/seo.go:buildJSONLD": "JSON-LD from json.Marshal, which escapes <, > and &; allowed by its CSP hash",
}

// TestNoUnsafeConversions scans the Go code for conversions to the
// html/template types that disable escaping.
func TestNoUnsafeConversions(t *testing.T) {
	root := filepath.Join("..", "..")
	fset := token.NewFileSet()
	for _, dir := range []string{"cmd", "internal", "web"} {
		err := filepath.WalkDir(filepath.Join(root, dir), func(p string, d fs.DirEntry, err error) error {
			if err != nil || d.IsDir() || !strings.HasSuffix(p, ".go") || strings.HasSuffix(p, "_test.go") {
				return err
			}
			src, err := os.ReadFile(p) // #nosec G304 -- Go files of this repository
			if err != nil {
				return err
			}
			f, err := parser.ParseFile(fset, p, src, 0)
			if err != nil {
				return err
			}
			rel, _ := filepath.Rel(root, p)
			rel = filepath.ToSlash(rel)
			for _, decl := range f.Decls {
				fn, ok := decl.(*ast.FuncDecl)
				if !ok {
					continue
				}
				ast.Inspect(fn, func(n ast.Node) bool {
					call, ok := n.(*ast.CallExpr)
					if !ok {
						return true
					}
					sel, ok := call.Fun.(*ast.SelectorExpr)
					if !ok {
						return true
					}
					pkg, ok := sel.X.(*ast.Ident)
					if ok && pkg.Name == "template" && unsafeTypes[sel.Sel.Name] {
						key := rel + ":" + fn.Name.Name
						if _, allowed := allowedConversions[key]; !allowed {
							t.Errorf("%s: template.%s conversion in %s is not allowed", fset.Position(call.Pos()), sel.Sel.Name, fn.Name.Name)
						}
					}
					return true
				})
			}
			return nil
		})
		if err != nil {
			t.Fatal(err)
		}
	}
}

// TestSafetyWalkerCatches makes sure the walker sees calls nested in
// pipelines and blocks, so the check above cannot pass by accident.
func TestSafetyWalkerCatches(t *testing.T) {
	src := `{{define "x"}}{{if .A}}<a href="{{printf "%s" .B | html}}">{{end}}{{range .C}}{{call .D}}{{end}}{{end}}`
	trees, err := parse.Parse("x", src, "{{", "}}", builtinsForParse())
	if err != nil {
		t.Fatal(err)
	}
	var found []string
	for _, tree := range trees {
		walk(tree.Root, func(n parse.Node) {
			if id, ok := n.(*parse.IdentifierNode); ok && !allowedFuncs[id.Ident] {
				found = append(found, id.Ident)
			}
		})
	}
	if strings.Join(found, ",") != "printf,html,call" {
		t.Errorf("walker found %v", found)
	}
}
