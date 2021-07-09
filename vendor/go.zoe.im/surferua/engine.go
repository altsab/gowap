package surferua

import "strings"

type Engine struct {
	Name   string
	Semver *Semver
}

func (e *Engine) String() string {
	if e == nil {
		return ""
	}
	if strings.Contains(e.Name, "ecko") {
		// this can be modify?
		return e.Name + "/20100101"
	}
	return e.Name + "/" + e.Semver.String() + " (KHTML, like Gecko)"

}

type EngineInfo struct {
	Name        string
	VersionInfo VersionInfo
}

func (ei *EngineInfo) Random() (e *Engine) {
	if ei == nil {
		return nil
	}
	return &Engine{Name: ei.Name, Semver: ei.VersionInfo.Random()}
}

func NewEngineInfo(m interface{}) (ei EngineInfo) {
	if mMap, ok := m.(map[interface{}]interface{}); ok {
		if name, ok := mMap["name"].(string); ok {
			return EngineInfo{Name: name, VersionInfo: NewVersionInfo(mMap["version"])}
		}
	}
	return
}
