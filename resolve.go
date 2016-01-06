package main

import (
	"path"
	"path/filepath"

	"github.com/ungerik/go-dry"
)

func (app *ResolverApp) resolve() {
	srcs, _ := filepath.Glob("*.go")
	srcs2, _ := filepath.Glob("**/*.go")
	srcfiles := append(srcs2, srcs...)

	for _, srcfile := range srcfiles {
		app.scanImports(srcfile)
	}

	for packname, users := range app.resolveMap {
		app.syncPack(packname, users)
	}
}

func (app *ResolverApp) checkInPath(packname, _path string) bool {
	return dry.FileIsDir(path.Join(_path, "src", packname))
}

func (app *ResolverApp) check(packname string) bool {
	if app.resolved[packname] {
		return true
	}
	if app.checkInPath(packname, app.goRoot) {
		app.resolved[packname] = true
		return true
	}
	for _, gopath := range app.goPaths {
		if app.checkInPath(packname, gopath) {
			app.resolved[packname] = true
			return true
		}
	}
	return false
}
