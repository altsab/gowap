package surferua

import "strings"

var browserDBSize = 0
var browserDB []BrowserInfo

// Browser is br
type Browser struct {
	// The name of the browser.
	Name string

	// The name of the browser's engine.
	Engine Engine

	// The version of the browser.
	Semver Semver
}

func (b *Browser) String() string {
	// chrome has different format tpl
	if b == nil {
		return ""
	}

	if strings.Contains(b.Name, "hrome") {
		// make safari random
		return b.Engine.String() + " " + b.Name + "/" + b.Semver.String() + " Safari/537.36"
	}
	return b.Engine.String() + " " + b.Name + "/" + b.Semver.String()
}

type BrowserInfo struct {
	Name        string
	EngineInfo  EngineInfo
	VersionInfo VersionInfo
}

func (bi *BrowserInfo) Random() *Browser {
	return &Browser{Name: bi.Name, Engine: *bi.EngineInfo.Random(), Semver: *bi.VersionInfo.Random()}
}

func NewBrowserInfo(name string, m interface{}) (bi *BrowserInfo) {
	if mMap, ok := m.(map[interface{}]interface{}); ok {
		bi = &BrowserInfo{
			Name:        name,
			EngineInfo:  NewEngineInfo(mMap["engine"]),
			VersionInfo: NewVersionInfo(mMap["version"]),
		}
	}
	return
}

func NewBrowserInfoList(m map[interface{}]interface{}) (l []*BrowserInfo) {
	for name, mMpa := range m {
		bi := NewBrowserInfo(name.(string), mMpa)
		if bi != nil {
			l = append(l, bi)
		}
	}
	return
}
