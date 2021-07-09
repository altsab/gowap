package surferua

var platformTypeSize int
var platformSize []int
var platformDB [][]PlatformInfo

// TODO: use tpl
// _tpl: "{{.MozillaWithVersion}} (iPhone; CPU iPhone OS {{.Version.Major}}_{{.Version.Minor}}_{{.Version.Patch}} like Mac OS X) {{.BrowserS}} {{.Comment}}"
// comment: "Version/{{.Version.Major}}.0 Mobile/11A466 Safari/602.1"

// Represents full information on the operating system extracted from the user agent.
type Platform struct {
	// Name of the operating system. This is sometimes a shorter version of the
	// operating system name, e.g. "Mac OS X" instead of "Intel Mac OS X"
	Name string

	// Operating system version, e.g. 7 for Windows 7 or 10.8 for Max OS X Mountain Lion
	Semver Semver

	// Comment
	Comment string
}

func (p *Platform) String() string {
	if p == nil {
		return ""
	}

	switch p.Name {
	case "iOS":
		return "iPhone; CPU iPhone OS " + p.Semver.String("_") + " like Mac OS X"
	case "Android":
		// Miss manufacturer
		return "Linux; Android " + p.Semver.String("_") + "; Build/" + RandStringBytesMaskImpr(5)
	case "Windows":
		// where is x86?
		return "Windows NT " + p.Semver.String("_") + "; Win64; x64"
	case "Linux":
		// Ubuntu ...?
		return "X11; Linux x86_64"
	case "MacOS":
		return "Macintosh; Intel Mac OS X " + p.Semver.String("_")
	default:
		return p.Name + "; " + p.Name + p.Semver.String("_")
	}
}

type PlatformInfo struct {
	Name string

	VersionInfo VersionInfo

	// this should be random
	Comment string

	// store phone or pc?
}

func (p *PlatformInfo) Random() *Platform {
	return &Platform{Name: p.Name, Semver: *p.VersionInfo.Random(), Comment: p.Comment}
}

func NewPlatformInfo(name string, m interface{}) (pi *PlatformInfo) {
	if mMap, ok := m.(map[interface{}]interface{}); ok {
		pi = &PlatformInfo{
			Name:        name,
			VersionInfo: NewVersionInfo(mMap["version"]),
			// Comment: ...
		}
	}
	return
}

func NewPlatformInfoList(m map[interface{}]interface{}) (types []string, l [][]*PlatformInfo) {
	for type0, mMap := range m {
		subList := []*PlatformInfo{}
		if mmMap, ok := mMap.(map[interface{}]interface{}); ok {
			for name, mmMapMap := range mmMap {
				pi := NewPlatformInfo(name.(string), mmMapMap)
				if pi != nil {
					subList = append(subList, pi)
				}
			}
		}

		// append sub list to the father list
		if len(subList) > 0 {
			l = append(l, subList)
			types = append(types, type0.(string))
		}
	}
	return
}
