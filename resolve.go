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

	for packname, _ := range app.resolveMap {
		if app.check(packname) {
			app.resolveMap[packname] = true
		}
	}

	for packname, resolved := range app.resolveMap {
		if !resolved {
			app.syncPack(packname)
		}
	}
}

func (app *ResolverApp) check(packname string) bool {
	if app.checkInPath(packname, app.goRoot) {
		return true
	}
	for _, gopath := range app.goPaths {
		if app.checkInPath(packname, gopath) {
			return true
		}
	}
	return false
}

func (app *ResolverApp) checkInPath(packname, _path string) bool {
	return dry.FileIsDir(path.Join(_path, "src", packname))
}
