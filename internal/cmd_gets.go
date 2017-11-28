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

type defInfo struct {
	fullURI    string
	branch     string
	tagValue   string
	tagMatcher *versionMatcher
	ustPath    string
}

func (di *defInfo) CurrentCommitId() (string, error) {
	var (
		result []byte
		err    error
	)
	cmdGetTags := fmt.Sprintf(`cd %s && git log -1`, di.ustPath)
	if result, err = exec.Command("sh", "-c", cmdGetTags).CombinedOutput(); err != nil {
		return "", fmt.Errorf("get commit id err: %s", err)
	}

	return strings.TrimPrefix(regCommitId.FindString(string(result)), "commit "), nil
}

func parseAddress(addr string) (*defInfo, error) {
	di := &defInfo{branch: defBranch}
	switch {
	case strings.HasPrefix(addr, "git@"): // git@xxx.xxx:/xx/xx.git

		di.fullURI = addr
		di.ustPath = strings.Replace(strings.Replace(addr, "git@", "", 1), ":", "/", 1)
	case regHttpAddress.MatchString(addr): // http(s)://xxx.xxx/xx/xx.git
		di.fullURI = addr
		di.ustPath = strings.TrimLeft(addr, "https://")
		di.ustPath = strings.TrimLeft(addr, "http://")
	case regAddress.MatchString(addr): // ddd.ddd.ddd.ddd or xxx.xxx or xxx.xxx.xxx
		di.fullURI = "http://" + addr
		di.ustPath = addr
	default:
		return nil, fmt.Errorf("can not resolve git address: %s", addr)
	}

	di.ustPath = path.Join(Vendor, strings.TrimRight(di.ustPath, ".git"))

	return di, nil
}

func parseVersionLimit(v string, dst *defInfo) {
	dst.tagValue = v
	var strMatcher string
	switch v[0] {
	case '~': // match sub version
		vs := strings.SplitN(v[1:], ".", 3)
		strMatcher = fmt.Sprintf("%s.%s.*", vs[0], vs[1])
	case '^': // match main version
		vs := strings.SplitN(v[1:], ".", 2)
		strMatcher = fmt.Sprintf("%s.*.*", vs[0])
	case '*': // match any version
		strMatcher = "*.*.*"
	default:
		strMatcher = v
	}
	dst.tagMatcher = getVersionRange(strMatcher)
}

func parseVersion(v string, dst *defInfo) error {
	ss := strings.Split(v, "#")
	if len(ss) > 2 {
		return fmt.Errorf("invalid version: %s", v)
	}

	if len(ss) < 2 {
		dst.branch = defBranch
		parseVersionLimit(v, dst)
	} else {
		dst.branch = ss[0]
		if len(ss[1]) > 0 {
			parseVersionLimit(ss[1], dst)
		} else {
			parseVersionLimit("*.*.*", dst)
		}
	}
	return nil
}

func cloneDep(di *defInfo) error {
	cmdArgs := []string{"clone", di.fullURI, di.ustPath}
	if _, err := exec.Command("git", cmdArgs...).CombinedOutput(); err != nil {
		return fmt.Errorf("repo(%s) clone err:%s", di.fullURI, err)
	}
	return nil
}

func updateDep(di *defInfo) error {
	cmdLookupBranch := fmt.Sprintf(`cd %s && git fetch`, di.ustPath)
	if _, err := exec.Command("sh", "-c", cmdLookupBranch).CombinedOutput(); err != nil {
		return fmt.Errorf("repo(%s) update err: %s", di.fullURI, err)
	}
	return nil
}

func switch2Branch(di *defInfo) error {
	cmdCheckoutBranch := fmt.Sprintf(`cd %s && git checkout %s`, di.ustPath, di.branch)
	if _, err := exec.Command("sh", "-c", cmdCheckoutBranch).CombinedOutput(); err != nil {
		return fmt.Errorf("repo(%s) switch branch(%s) err: (%s)", di.fullURI, di.branch, err)
	}

	return nil
}

func updateTag(di *defInfo) (err error) {
	// if match value is *, it's means match any version.
	if di.tagValue == "*" {
		return nil
	}

	var result []byte
	cmdGetTags := fmt.Sprintf(`cd %s && git tag -l`, di.ustPath)
	if result, err = exec.Command("sh", "-c", cmdGetTags).CombinedOutput(); err != nil {
		return fmt.Errorf("repo(%s) gets all tag error at %s: %s", di.fullURI, di.branch, err)
	}

	commitId, err := di.CurrentCommitId()
	if err != nil {
		return err
	}

	tags := strings.Split(string(result), "\n")
	if len(tags) == 0 {
		return fmt.Errorf("repo(%s) no any tag at %s", di.fullURI, di.branch)
	}

	matched := di.tagMatcher.MatchLastVersion(tags)
	if matched == "0.0.0" {
		return fmt.Errorf("repo(%s) not contains version(%s)", di.fullURI, di.tagValue)
	}

	cmdIsTagCommitId := fmt.Sprintf(`cd %s && git tag -l %s --points-at %s`, di.ustPath, matched, commitId)
	if result, err = exec.Command("sh", "-c", cmdIsTagCommitId).CombinedOutput(); err != nil {
		return fmt.Errorf("repo(%s) verify %s error: %s", di.fullURI, matched, err)
	}

	if string(result) != "" {
		return nil
	}

	cmdCheckoutTag := fmt.Sprintf(`cd %s && git checkout %s`, di.ustPath, matched)
	if _, err = exec.Command("sh", "-c", cmdCheckoutTag).CombinedOutput(); err != nil {
		return fmt.Errorf("repo(%s) switch branch(%s) error:%s", di.fullURI, matched, err)
	}

	return nil
}

func doCmdGets(p *Entry) error {
	var (
		di  *defInfo
		err error
	)
	processedDeps := map[string]bool{}

	// collocate all deps
	subDeps := make(map[string]string)

	for address, version := range p.Deps {
		if processedDeps[address] {
			continue
		}
		processedDeps[address] = true
		if di, err = parseAddress(address); err != nil {
			return err
		}
		handleError(parseVersion(version, di))

		// if repo does not exists, clone it
		if _, err = os.Stat(di.ustPath); err != nil {
			if !os.IsNotExist(err) {
				handleError(fmt.Errorf("proc deps(%s) error :%s", address, err))
			}
			handleError(cloneDep(di))
		}

		// checkout version
		handleError(updateDep(di))
		handleError(switch2Branch(di))
		handleError(updateTag(di))

		deps := getSubDeps(di)
		for repo, ver := range deps.Deps {
			subDeps[repo] = ver
		}
	}

	for len(subDeps) > 0 {
		tmp := subDeps
		subDeps = make(map[string]string)
		for address, version := range tmp {
			if processedDeps[address] {
				continue
			}
			processedDeps[address] = true
			if di, err = parseAddress(address); err != nil {
				return err
			}
			handleError(parseVersion(version, di))

			if _, err = os.Stat(di.ustPath); err != nil {
				if !os.IsNotExist(err) {
					handleError(fmt.Errorf("proc deps(%s) error: %s", address, err))
				}
				handleError(cloneDep(di))
			}

			handleError(updateDep(di))
			handleError(switch2Branch(di))
			handleError(updateTag(di))

			deps := getSubDeps(di)
			for repo, ver := range deps.Deps {
				subDeps[repo] = ver
			}
		}
	}
	return nil
}
