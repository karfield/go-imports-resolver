package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func parseRecursively(dir string, imports map[string]struct{}) error {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}

	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, dir, nil, parser.ImportsOnly)
	if err != nil {
		return err
	}

	//ast.Print(fset, pkgs)
	for _, pkgAst := range pkgs {
		for _, fileAst := range pkgAst.Files {
			for _, decl := range fileAst.Decls {
				if genDecl, ok := decl.(*ast.GenDecl); ok {
					if genDecl.Tok == token.IMPORT {
						for _, spec := range genDecl.Specs {
							if importSpec, ok := spec.(*ast.ImportSpec); ok {
								if importSpec.Path.Kind == token.STRING {
									//fmt.Printf("package: %s\n", importSpec.Path.Value)
									importPkg, err := strconv.Unquote(importSpec.Path.Value)
									if err != nil {
										fmt.Printf("package: %s\n", importSpec.Path.Value)
									}
									imports[importPkg] = struct{}{}
								}
							}
						}
					}
				}
			}
		}
	}

	for _, f := range files {
		if f.IsDir() {
			if err := parseRecursively(filepath.Join(dir, f.Name()), imports); err != nil {
				return err
			}
		}
	}

	return nil
}

func dirExists(dir ...string) bool {
	stat, err := os.Stat(filepath.Join(dir...))
	if err != nil {
		return false
	}
	return stat.IsDir()
}

func getPkgName(s string) string {
	segs := strings.Split(s, "/")
	if len(segs) > 1 {
		domain := segs[0]
		switch domain {
		case "github.com", "gopkg.in", "golang.org", "bitbucket.org", "hub.jazz.net", "git.apache.org", "9fans.net", "sourcegraph.com":
			if len(segs) > 2 {
				return filepath.Join(domain, segs[1], segs[2])
			} else {
				return ""
			}
		case "karfield.com", "apporture.com":
			return filepath.Join(domain, segs[1])
		}
	}
	return ""
}

func main() {

	cwd, err := os.Getwd()
	if err != nil {
		return
	}

	goroot := os.Getenv("GOROOT")
	gopaths := strings.Split(os.Getenv("GOPATH"), string(os.PathListSeparator))

	currpkg := ""
	for _, gopath := range gopaths {
		p := filepath.Join(gopath, "src")
		if strings.HasPrefix(cwd, p+"/") {
			currpkg = cwd[len(p)+1:]
		}
	}
	currpkg = getPkgName(currpkg)

	imports := map[string]struct{}{}

	err = parseRecursively(cwd, imports)
	if err != nil {
		panic(err)
	}

	import_filtered := map[string]struct{}{}
	for imp, _ := range imports {
		if dirExists(filepath.Join(goroot, "src", imp)) {
			continue
		}

		if strings.HasPrefix(imp, cwd) {
			// inner package, ignore
			continue
		}

		n := getPkgName(imp)
		if n != "" {
			if n != currpkg {
				import_filtered[n] = struct{}{}
			}
		}
	}

	if len(import_filtered) == 0 {
		fmt.Printf("Nothing extra packages need to import!\n")
	}

	output := ""
	for _, o := range os.Args[1:] {
		if dirExists(o) {
			output = filepath.Join(o, "get-go-deps.sh")
			break
		} else if strings.HasSuffix(o, ".sh") {
			if dirExists(filepath.Dir(o)) {
				output = o
				break
			}
		}
	}

	if output != "" {
		f, err := os.OpenFile(output, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.FileMode(0644))
		if err != nil {
			fmt.Fprintf(os.Stderr, "cannot write to %s: %+v\n", output, err)
		}
		defer f.Close()
		for p, _ := range import_filtered {
			fmt.Fprintf(f, "go get -u %s\n", p)
		}
		fmt.Printf("%s has generated\n", output)
		f.Chmod(os.FileMode(0755))
	} else {
		for p, _ := range import_filtered {
			fmt.Printf("go get -u %s\n", p)
		}
	}
}
