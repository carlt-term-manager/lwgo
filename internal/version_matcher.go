package internal

import (
	"math"
	"strconv"
	"strings"
)

type versionMatcher struct {
	v1value, v2value, v3value int
	invalid                   bool
}

func (vr *versionMatcher) MatchLastVersion(vns []string) string {
	var last = &VersionStruct{}
	if !vr.invalid {

		vss := ToVersionStructArray(vns).Sort()

		if vr.v1value == math.MaxInt32 { // *.*.*
			for _, v := range vss {
				if !v.invalid {
					if last.Less(v) {
						last = v
					}
				}
			}
		} else if vr.v2value == math.MaxInt32 { // n.*.*
			for _, v := range vss {
				if !v.invalid && v.v1 == vr.v1value {
					if last.Less(v) {
						last = v
					}
				}
			}
		} else if vr.v3value == math.MaxInt32 { // n.n.*
			for _, v := range vss {
				if !v.invalid && v.v1 == vr.v1value && v.v2 == vr.v2value {
					if last.Less(v) {
						last = v
					}
				}
			}
		} else { // n.n.n
			for _, v := range vss {
				if !v.invalid && v.v1 == vr.v1value && v.v2 == vr.v2value && v.v3 == vr.v3value {
					last = &VersionStruct{v1: vr.v1value, v2: vr.v2value, v3: vr.v3value}
					break
				}
			}
		}
	}
	return last.String()
}

func _atoi(s string) int {
	if s != "*" {
		v, _ := strconv.Atoi(s)
		return v
	}
	return math.MaxInt32
}

func getVersionRange(v string) *versionMatcher {
	r := &versionMatcher{}
	vs := strings.SplitN(v, ".", 3)

	if len(vs) >= 3 {
		r.v1value = _atoi(vs[0])
		r.v2value = _atoi(vs[1])
		r.v3value = _atoi(vs[2])
	} else {
		r.invalid = true
	}
	return r
}
