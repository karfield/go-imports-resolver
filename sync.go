package main

import (
	"fmt"
	"os"
	"os/exec"
)

func (app *ResolverApp) syncPack(packname string, users []string) {
	fmt.Printf("Sync package: %s (used by %+v)\n", packname, users)
	cmd := exec.Command("go", "get", packname)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()
	cmd.Run()
}
