package main

import (
	"go/parser"
	"go/token"
	"strconv"
	"strings"
	"unicode"
)

const buildMatch = "+build "

func (app *ResolverApp) scanImports(srcfile string) bool {
	pf, err := parser.ParseFile(token.NewFileSet(), srcfile, nil, parser.ParseComments)
	if err != nil {
		return false
	}

	if len(pf.Comments) > 0 {
		for _, c := range pf.Comments {
			ct := c.Text()
			if i := strings.Index(ct, buildMatch); i != -1 {
				for _, b := range strings.FieldsFunc(ct[i+len(buildMatch):], func(r rune) bool {
					return unicode.IsSpace(r) || r == ','
				}) {
					//TODO: appengine is a special case for now: https://github.com/tools/godep/issues/353
					if b == "ignore" || b == "appengine" {
						return true
					}
				}
			}
		}
	}

	for _, impt := range pf.Imports {
		name, err := strconv.Unquote(impt.Path.Value)
		if err != nil {
			continue
		}
		app.resolveMap[name] = false
	}

	return true
}
