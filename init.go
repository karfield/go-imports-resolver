package main

import (
	"os"
	"strings"

	"github.com/codegangsta/cli"
)

func (app *ResolverApp) init() {
	app.goRoot = os.Getenv("GOROOT")
	app.goPaths = strings.Split(os.Getenv("GOPATH"), ":")

	app.cli.Name = "go-imports-resolver"
	app.cli.Usage = "resolve go imports automatically"
	app.cli.Flags = []cli.Flag{}
}
