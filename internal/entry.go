package internal

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Entry struct {
	Name    string            `json:"name"`
	Version string            `json:"version"`
	Deps    map[string]string `json:"deps"`
}

func newEntry() *Entry {
	return &Entry{
		Deps: make(map[string]string, 0),
	}
}

func (p *Entry) Read(f string) error {
	buf, err := ioutil.ReadFile(f)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(buf, p); err == nil {
		for addr, v := range p.Deps {
			if !regTagMatcherExp.Match([]byte(v)) {
				return fmt.Errorf("invalid version matcher(%s) for %s", v, addr)
			}
		}
	}
	return err
}

func (p *Entry) Save(f string) error {
	result, err := json.MarshalIndent(p, "", "\t")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(f, result, 0755)
}
