package surferua

import (
	"math/rand"
	"strconv"
)

// Semantic versioning 2.0.0
// http://semver.org/

var defaultVersionRange = &Endpoint{0, 0}

type Semver struct {
	Major      int
	Minor      int
	Patch      int
	PreRelease PreRelease
	Metadata   string
}

type PreRelease string

func (s *Semver) String(seps ...string) string {
	var sep = "."
	if len(seps) >= 1 {
		sep = seps[0]
	}
	// default use '.' as sep
	// remove the last 0
	if s.Patch == 0 {
		return strconv.Itoa(s.Major) + sep + strconv.Itoa(s.Minor)
	}
	return strconv.Itoa(s.Major) + sep + strconv.Itoa(s.Minor) + sep + strconv.Itoa(s.Patch)
}

type VersionInfo struct {
	Major Endpoint
	Minor Endpoint
	Patch Endpoint
}

type Endpoint struct {
	Start int
	End   int
}

func (vi *VersionInfo) Random() (s *Semver) {
	s = &Semver{
		Major: randomRange(vi.Major.Start, vi.Major.End),
		Minor: randomRange(vi.Minor.Start, vi.Minor.End),
		Patch: randomRange(vi.Patch.Start, vi.Patch.End),
	}
	return s
}

func randomRange(start, end int) int {
	if start == end {
		return start
	} else {
		return start + rand.Intn(end-start)
	}
}

func mustInt(v interface{}) int {
	switch value := v.(type) {
	case int:
		return value
	case int64:
		return int(value)
	default:
		return 0
	}
}

func getVersionRange(v interface{}) Endpoint {
	// check type of v
	switch value := v.(type) {
	case int:
		// if value if int64 return []
		return Endpoint{value, value}
	case int64:
		// if value if int64 return []
		intValue := int(value)
		return Endpoint{intValue, intValue}
	case []int:
		// if value if list of int, can be int64?
		switch len(value) {
		case 0:
			return *defaultVersionRange
		case 1:
			return Endpoint{value[0], value[0]}
		default:
			if value[0] > value[1] {
				return Endpoint{value[1], value[0]}
			}
			return Endpoint{value[0], value[1]}
		}
	case []int64:
		// if value if list of int, can be int64?
		switch len(value) {
		case 0:
			return *defaultVersionRange
		case 1:
			intValue := int(value[0])
			return Endpoint{intValue, intValue}
		default:
			if value[0] > value[1] {
				return Endpoint{int(value[1]), int(value[0])}
			}
			return Endpoint{int(value[0]), int(value[1])}
		}
	case []interface{}:
		// if value if list of int, can be int64?
		switch len(value) {
		case 0:
			return *defaultVersionRange
		case 1:
			intValue := mustInt(value[0])
			return Endpoint{intValue, intValue}
		default:
			v0 := mustInt(value[0])
			v1 := mustInt(value[1])
			if v0 > v1 {
				return Endpoint{v1, v0}
			}
			return Endpoint{v0, v1}
		}
	default:
		// default return default 0.0.0
		return *defaultVersionRange
	}
}

func NewVersionInfo(m interface{}) (vi VersionInfo) {
	// we need to get data from map[string]interface{}
	if mMap, ok := m.(map[interface{}]interface{}); ok {
		return VersionInfo{
			Major: getVersionRange(mMap["major"]),
			Minor: getVersionRange(mMap["minor"]),
			Patch: getVersionRange(mMap["patch"]),
		}
	}
	return
}
