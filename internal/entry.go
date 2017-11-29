package internal

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Dep struct {
	Src string `json:"src"`
	Dst string `json:"dst,omitempty"`
	Ver string `json:"ver,omitempty"`
}

type Entry struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Deps    []Dep  `json:"deps"`
}

func newEntry() *Entry {
	return &Entry{
		Deps: make([]Dep, 0),
	}
}

func (p *Entry) Read(f string) error {
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

func (p *Entry) Save(f string) error {
	p.Merge()
	result, err := json.MarshalIndent(p, "", "\t")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(f, result, 0755)
}

func (p *Entry) Validate() error {
	for _, v := range p.Deps {
		if !regAddress.MatchString(v.Src) {
			return fmt.Errorf("repo(%s) invalid src address", v.Src)
		}
		if v.Dst != "" && !regAddress.MatchString(v.Dst) {
			return fmt.Errorf("repo(%s) invalid dst address", v.Dst)
		}
		if v.Ver != "" && !regTagMatcherExp.Match([]byte(v.Ver)) {
			return fmt.Errorf("repo(%s) invalid version matcher(%s)", v.Src, v.Ver)
		}
	}
	return nil
}

func (p *Entry) Merge() {
	fmt.Printf("合并前: %#v\n", p.Deps)
	s := make(map[string]Dep)
	for _, d := range p.Deps {
		s[d.Src] = d
	}
	p.Deps = make([]Dep, 0)
	for _, d := range s {
		p.Deps = append(p.Deps, d)
	}
	fmt.Printf("合并后: %#v\n", p.Deps)

}

func (p *Entry) Fill() {
	for idx, d := range p.Deps {
		if d.Dst == "" {
			d.Dst = d.Src
		}
		p.Deps[idx] = d
	}
}
