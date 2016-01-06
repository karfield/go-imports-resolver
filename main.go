package main

import "github.com/codegangsta/cli"

type ResolverApp struct {
	cli        *cli.App
	resolveMap map[string]bool
	goPaths    []string
	goRoot     string
}

func newResolver() *ResolverApp {
	app := new(ResolverApp)
	app.cli = cli.NewApp()
	app.resolveMap = make(map[string]bool)
	app.goPaths = []string{}
	return app
}

func main() {
	resolver := newResolver()
	resolver.init()
	resolver.resolve()
}
