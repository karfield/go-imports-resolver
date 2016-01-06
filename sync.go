package main

import (
	"fmt"
	"os"
	"os/exec"
)

func (app *ResolverApp) syncPack(packname string) {
	fmt.Println("Sync package: " + packname)
	cmd := exec.Command("go", "get", packname)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()
	cmd.Run()
}
