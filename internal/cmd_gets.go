package internal

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
)

func CmdGets() {
	p := newEntry()
	handleError(p.Read(PackageFile))
	handleError(doCmdGets(p))
}

func parseAddress(d Dep) (*defInfo, error) {
	di := &defInfo{}
	switch {
	case strings.HasPrefix(d.Src, "git@"): // git@xxx.xxx:/xx/xx.git
		di.origin = d.Dst
		di.usedPath = strings.Replace(strings.TrimPrefix(d.Src, "git@"), ":", "/", 1)
	case regAddress.MatchString(d.Src): // http(s)://xxx.xxx/xx/xx.git
		di.origin = d.Dst
		di.usedPath = regHttpReplace.ReplaceAllString(d.Src, "")
	default:
		return nil, fmt.Errorf("repo(%s) can not resolve", d.Dst)
	}

	di.usedPath = path.Join(Vendor, strings.TrimSuffix(di.usedPath, ".git"))
	di.branch = d.Branch

	return di, nil
}

func cloneDep(di *defInfo) error {
	cmdArgs := []string{"clone"}
	if di.branch != "" {
		cmdArgs = append(cmdArgs, "--branch", di.branch)
	}
	cmdArgs = append(cmdArgs, di.origin, di.usedPath)
	if _, err := exec.Command("git", cmdArgs...).CombinedOutput(); err != nil {
		return fmt.Errorf("repo(%s) clone err:%s", di.origin, err)
	}
	return nil
}

func updateDep(di *defInfo) error {
	cmdLookupBranch := fmt.Sprintf(`cd %s && git fetch && git pull origin %s`, di.usedPath, di.branch)
	if _, err := exec.Command("sh", "-c", cmdLookupBranch).CombinedOutput(); err != nil {
		return fmt.Errorf("repo(%s) update err: %s", di.origin, err)
	}
	return nil
}

func co(di *defInfo) error {
	if di.version != "" {
		cmdCheckoutBranch := fmt.Sprintf(`cd %s && git checkout %s`, di.usedPath, di.version)
		if _, err := exec.Command("sh", "-c", cmdCheckoutBranch).CombinedOutput(); err != nil {
			return fmt.Errorf("repo(%s) checkout(%s) err: (%s)", di.origin, di.version, err)
		}
	}
	return nil
}

var processedDeps = map[string]bool{}

func updateDeps(p *Entry) error {
	col := &Entry{}
	for _, d := range p.Deps {
		if processedDeps[d.Dst] {
			continue
		}
		processedDeps[d.Dst] = true
		di, err := parseAddress(d)
		if err != nil {
			return err
		}
		// if repo does not exists, clone it
		if _, err := os.Stat(di.usedPath); err != nil {
			if !os.IsNotExist(err) {
				handleError(fmt.Errorf("repo(%s) proc error :%s", di.origin, err))
			}
			handleError(cloneDep(di))
		}

		// checkout version
		handleError(updateDep(di))
		handleError(co(di))

		deps, err := getSubDeps(di)
		handleError(err)

		for _, c := range deps.Deps {
			col.Deps = append(col.Deps, c)
		}
	}
	if len(col.Deps) == 0 {
		return nil
	}

	return updateDeps(col)
}

func doCmdGets(p *Entry) error {
	p.Fill()
	processedDeps = map[string]bool{}
	return updateDeps(p)
}
