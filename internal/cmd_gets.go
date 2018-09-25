package internal

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
)

func CmdGets() {
	handleError(
		doCmdGets(newMod()),
	)
}

func (u *updater) cloneRepo(di *defInfo) error {
	cmdArgs := []string{"clone"}
	if di.branch != "" {
		cmdArgs = append(cmdArgs, "--branch", di.branch)
	}

	cmdArgs = append(cmdArgs, di.repoAddr, di.storePath)
	if _, err := exec.Command("git", cmdArgs...).CombinedOutput(); err != nil {
		return fmt.Errorf("repo(%s) clone err:%s", di.repoAddr, err)
	}
	return nil
}

func (u *updater) updateBranch(di *defInfo) error {
	cmdLookupBranch := fmt.Sprintf(`cd %s && git fetch && git pull origin %s`, di.storePath, di.branch)
	if _, err := exec.Command("sh", "-c", cmdLookupBranch).CombinedOutput(); err != nil {
		return fmt.Errorf("repo(%s) update err: %s", di.repoAddr, err)
	}
	return nil
}

func (u *updater) switchToTag(di *defInfo) error {
	if di.version != "" {
		cmdCheckoutBranch := fmt.Sprintf(`cd %s && git checkout %s`, di.storePath, di.version)
		if _, err := exec.Command("sh", "-c", cmdCheckoutBranch).CombinedOutput(); err != nil {
			return fmt.Errorf("repo(%s) checkout(%s) err: (%s)", di.repoAddr, di.version, err)
		}
	}
	return nil
}

func (u *updater) isProcessed(name string) bool {
	if u.processed[name] {
		return true
	}
	u.processed[name] = true
	return false
}

// 更新依赖
func (u *updater) Run(p *mod) error {
	col := &mod{}
	for _, d := range p.Deps {
		fmt.Printf("[*] %s (%s)...", d.Name, d.Version)

		// 处理过的不再处理
		if u.isProcessed(d.Replace) {
			fmt.Println("skipped.")
			continue
		}

		di, err := u.parseAddress(d)
		if err != nil {
			return err
		}

		// if repo does not exists, clone it
		if _, err := os.Stat(di.storePath); err != nil {
			if !os.IsNotExist(err) {
				handleError(fmt.Errorf("repo(%s) proc error :%s", di.repoAddr, err))
			}
			handleError(u.cloneRepo(di))
		}

		// 优先版本， 没版本按分支
		if di.version != "" {
			handleError(u.switchToTag(di))
		} else {
			handleError(u.updateBranch(di))
		}

		deps, err := getSubDeps(di)
		handleError(err)

		for _, c := range deps.Deps {
			col.Deps = append(col.Deps, c)
		}
		fmt.Println("succeed.")
	}
	if len(col.Deps) == 0 {
		return nil
	}

	return u.Run(col)
}

func (u *updater) parseAddress(d ModItem) (*defInfo, error) {
	di := &defInfo{repoAddr: d.Replace, branch: d.Branch, version: d.Version}
	switch {

	// git@xxx.xxx:xx/xx.git
	case strings.HasPrefix(d.Name, "git@"):
		di.storePath = strings.Replace(strings.TrimPrefix(d.Name, "git@"), ":", "/", 1)

	// http://xxx.xxx/xx/xx.git
	case strings.HasPrefix(d.Name, "http://"):
		di.storePath = strings.TrimPrefix(d.Name, "http://")

	// https://xxx.xxx/xx/xx.git
	case strings.HasPrefix(d.Name, "https://"):
		di.storePath = strings.TrimPrefix(d.Name, "https://")

	default:
		return nil, fmt.Errorf("repo(%s) can not resolve", d.Replace)
	}

	di.storePath = path.Join(Vendor, strings.TrimSuffix(di.storePath, ".git"))

	return di, nil
}

type updater struct {
	processed map[string]bool
}

func doCmdGets(p *mod) error {
	handleError(p.Read(PackageFile))
	p.Fill()
	updaterNew := &updater{processed: make(map[string]bool, 0)}
	return updaterNew.Run(p)
}
