package internal

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	Vendor      = "vendor"
	PackageFile = "package.json"
)

var (
	packageDir       = "./"
	regTagMatcherExp = regexp.MustCompile(`^(?:\*)|master|[\w\-.~]+$`)
	regCommitId      = regexp.MustCompile("^commit ([a-f0-9]+)")
	regHttpReplace   = regexp.MustCompile(`^https?://`)
	regAddress       = regexp.MustCompile(`^(?:(?:https?://)|(?:git@))(?:[\w-]{1,61}\.)+[A-Za-z]{2,6}`)
)

func checkVendorDir(vendorDir string) error {
	var err error
	if _, err = os.Stat(vendorDir); err != nil {
		return os.Mkdir(vendorDir, 0755)
	}
	return nil
}

func getSubDeps(di *defInfo) (*Entry, error) {
	p := newEntry()
	err := p.Read(path.Join(di.usedPath, PackageFile))
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	p.Fill()
	return p, nil
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

type defInfo struct {
	origin, version, usedPath, branch string
}

func (di *defInfo) CurrentCommitId() (string, error) {
	var (
		result []byte
		err    error
	)
	cmdGetTags := fmt.Sprintf(`cd %s && git log -1`, di.usedPath)
	if result, err = exec.Command("sh", "-c", cmdGetTags).CombinedOutput(); err != nil {
		return "", fmt.Errorf("get commit id err: %s", err)
	}

	return strings.TrimPrefix(regCommitId.FindString(string(result)), "commit "), nil
}
