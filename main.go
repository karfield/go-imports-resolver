package main

import "github.com/codegangsta/cli"

type ResolverApp struct {
	cli        *cli.App
	resolveMap map[string][]string
	resolved   map[string]bool
	goPaths    []string
	goRoot     string
}

func newResolver() *ResolverApp {
	app := new(ResolverApp)
	app.cli = cli.NewApp()
	app.resolveMap = make(map[string][]string)
	app.resolved = make(map[string]bool)
	app.goPaths = []string{}
	return app
}

func main() {
	resolver := newResolver()
	resolver.init()
	resolver.resolve()
}

func (app *ResolverApp) addSrcfileByImport(impt, srcfile string) {
	if app.resolveMap[impt] != nil {
		if len(app.resolveMap[impt]) > 10 {
			return
		}
		app.resolveMap[impt] = append(app.resolveMap[impt], srcfile)
	} else {
		app.resolveMap[impt] = []string{srcfile}
	}
}
