package internal

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
)

const (
	Vendor      = "vendor"
	PackageFile = "package.json"
	defBranch   = "master"
)

var (
	packageDir       = "./"
	regTagMatcherExp = regexp.MustCompile(`^(?:\*)|master|[\w\-.]+|(?:[\^~]?\d+.\d+.\d+)$`)
	regCommitId      = regexp.MustCompile("^commit ([a-f0-9]+)")
	regHttpAddress   = regexp.MustCompile("^https*:\\\\")
	regAddress       = regexp.MustCompile(`^(?:[\w-]{1,61}\.)+[A-Za-z]{2,6}`)
)

func checkVendorDir(vendorDir string) error {
	var err error
	if _, err = os.Stat(vendorDir); err != nil {
		return os.Mkdir(vendorDir, 0755)
	}
	return nil
}

func getSubDeps(di *defInfo) *Entry {
	p := newEntry()
	p.Read(path.Join(di.ustPath, PackageFile))
	return p
}

func handleError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	packageDir, _ = filepath.Abs(packageDir)
}
