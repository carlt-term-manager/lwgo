package internal

import (
	"flag"
	"fmt"
	"os"
)

func cmdAddDepsUsage() {
	fmt.Println(`  Golang package manager for Carlt
------------------------------------------------
  lwgo add address version
    - address
      golang package's git address

    - version
	  depends of package
	  -git branch
	  -git commitId
	  -git tag
	  -\*            match default branch at git repo
	`)

	os.Exit(0)
}

func CmdAddDeps() {
	args := flag.Args()[1:]
	if len(args) != 2 {
		cmdAddDepsUsage()
	}

	p := newEntry()
	handleError(p.Read(PackageFile))

	p.Deps[args[0]] = args[1]
	p.Save(PackageFile)

	handleError(doCmdGets(p))
}
