package internal

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func newMod() *mod {
	return &mod{
		Deps: make([]ModItem, 0),
	}
}

type (
	ModItem struct {
		Name    string `json:"src"`
		Replace string `json:"dst,omitempty"`
		Version string `json:"ver,omitempty"`
		Branch  string `json:"branch,omitempty"`
	}

	mod struct {
		Name    string    `json:"repoAddr"`
		Version string    `json:"version"`
		Deps    []ModItem `json:"deps"`
	}
)

func (p *mod) Read(f string) error {
	buf, err := ioutil.ReadFile(f)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(buf, p); err == nil {
		p.Merge()
		return p.Validate()
	}
	return fmt.Errorf("proc file %s, err:%s", f, err)
}

func (p *mod) Save(f string) error {
	p.Merge()
	result, err := json.MarshalIndent(p, "", "\t")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(f, result, 0755)
}

func (p *mod) Validate() error {
	for _, v := range p.Deps {
		if !regAddress.MatchString(v.Name) {
			return fmt.Errorf("repo(%s) invalid src address", v.Name)
		}
		if v.Replace != "" && !regAddress.MatchString(v.Replace) {
			return fmt.Errorf("repo(%s) invalid dst address", v.Replace)
		}
		if v.Version != "" && !regTagMatcherExp.Match([]byte(v.Version)) {
			return fmt.Errorf("repo(%s) invalid version matcher(%s)", v.Name, v.Version)
		}
	}
	return nil
}

// 去掉重复名称的依赖
func (p *mod) Merge() {
	s := make(map[string]ModItem)
	for _, d := range p.Deps {
		s[d.Name] = d
	}
	p.Deps = make([]ModItem, 0)
	for _, d := range s {
		p.Deps = append(p.Deps, d)
	}
}

func (p *mod) Fill() {
	for idx, d := range p.Deps {
		if d.Replace == "" {
			d.Replace = d.Name
		}
		p.Deps[idx] = d
	}
}
