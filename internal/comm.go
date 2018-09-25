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
	regAddress       = regexp.MustCompile(`^(?:(?:https?://)|(?:git@))(?:[\w-]{1,61}\.)+[A-Za-z]{2,6}`)
)

func checkVendorDir(vendorDir string) error {
	var err error
	if _, err = os.Stat(vendorDir); err != nil {
		return os.Mkdir(vendorDir, 0755)
	}
	return nil
}

func getSubDeps(di *defInfo) (*mod, error) {
	p := newMod()
	err := p.Read(path.Join(di.storePath, PackageFile))
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
	repoAddr, version, storePath, branch string
}

func (di *defInfo) CurrentCommitId() (string, error) {
	var (
		result []byte
		err    error
	)
	cmdGetTags := fmt.Sprintf(`cd %s && git log -1`, di.storePath)
	if result, err = exec.Command("sh", "-c", cmdGetTags).CombinedOutput(); err != nil {
		return "", fmt.Errorf("get commit id err: %s", err)
	}

	return strings.TrimPrefix(regCommitId.FindString(string(result)), "commit "), nil
}

func (di *defInfo) String() string {
	return fmt.Sprintf("%s -> %s(%s)", di.repoAddr, di.storePath, di.version)
}
