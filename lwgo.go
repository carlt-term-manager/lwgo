package main

import (
	"flag"
	"fmt"
	lib "github.com/carltd/lwgo/internal"
	"os"
)

const (
	Version = "2.00.00"
)

func usage() {
	fmt.Printf(`  Golang package manager for Carlt v%s
------------------------------------------------
  init    init a project
  add     add deps to project's vendor directory
  gets    gets all deps
  help    this help page
`, Version)
	os.Exit(0)
}

func main() {
	flag.Usage = usage

	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		usage()
	}

	switch args[0] {
	case "init":
		lib.CmdInit()
	case "add":
		lib.CmdAddDeps()
	case "gets":
		lib.CmdGets()
	default:
		usage()
	}
}
