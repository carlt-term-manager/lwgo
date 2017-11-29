package internal

import (
	"flag"
	"fmt"
	"os"
)

func cmdAddDepsUsage() {
	fmt.Println(`  Golang package manager for Carlt
------------------------------------------------
  lwgo add src-addr [version [dst-addr]]
    - src-addr  package repo address.

	- dst-addr  alias address, it's used to create paths and import, but src-addr's content
	            instead of the source.

    - version
	  depends of package version:
	  -git branch
	  -git commitId
	  -git tag
	  -\* or empty  match default branch at git repo
	`)

	os.Exit(0)
}

func CmdAddDeps() {
	args := flag.Args()[1:]

	var src, ver, dst string
	switch len(args) {
	default:
		cmdAddDepsUsage()
	case 2:
		src, ver = args[0], args[1]
	case 3:
		src, ver, dst = args[0], args[1], args[2]
	}

	p := newEntry()
	handleError(p.Read(PackageFile))
	p.Deps = append(p.Deps, Dep{Src: src, Ver: ver, Dst: dst})
	handleError(p.Validate())
	p.Save(PackageFile)

	handleError(doCmdGets(p))
}
