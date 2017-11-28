package internal

import (
	"fmt"
	"strconv"
	"strings"
)

type VersionStruct struct {
	v1, v2, v3 int
	invalid    bool
}

func (s *VersionStruct) String() string {
	if s.invalid {
		return ""
	}
	return fmt.Sprintf("%d.%d.%d", s.v1, s.v2, s.v3)
}

func (s *VersionStruct) Less(vs *VersionStruct) bool {
	if s.invalid {
		return true
	}
	if vs.invalid {
		return false
	}
	if s.v1 == vs.v1 {
		if s.v2 == vs.v2 {
			return s.v3 < vs.v3
		}
		return s.v2 < vs.v2
	}

	return s.v1 < vs.v1
}

func (s *VersionStruct) parse(v string) {
	s.invalid = true
	vs := strings.SplitN(v, ".", 3)
	if len(vs) >= 3 {
		for {
			var err error
			if s.v1, err = strconv.Atoi(strings.TrimPrefix(vs[0], "v")); err != nil {
				break
			}
			if s.v2, _ = strconv.Atoi(vs[1]); err != nil {
				break
			}
			if s.v3, _ = strconv.Atoi(vs[2]); err != nil {
				break
			}
			s.invalid = false
			break
		}
	}
}

type VersionStructArray []*VersionStruct

func (vsa VersionStructArray) Sort() VersionStructArray {
	return vsa
}

func ToVersionStructArray(vss []string) VersionStructArray {
	vsa := make(VersionStructArray, len(vss))
	for i, v := range vss {
		vsa[i] = &VersionStruct{}
		vsa[i].parse(v)
	}
	return vsa
}
