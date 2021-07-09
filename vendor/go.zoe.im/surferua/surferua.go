package surferua

import (
	"math/rand"
	"strings"
	"time"
)

//go:generate go run ./hack/gen.go config.yml
//go:generate go fmt -w -s

func init() {
	Seed(time.Now().UnixNano())
}

func Seed(seed int64) {
	rand.Seed(seed)
}

const Mozilla = "Mozilla"
const MozillaWithVersion = "Mozilla/5.0"

type UserAgent struct {
	// Common version: 5.0
	version string

	// Browser information
	browser *Browser

	// Platform information
	platform *Platform
}

func (ua *UserAgent) String() (us string) {

	// firefox has different tpl
	// 1. use `.` to connect semver instead of `_`
	// 2. ends with `rv:`
	if strings.Contains(ua.browser.Name, "irefox") {
		return ua.version + " (" + strings.Replace(ua.platform.String(), "_", ".", -1) + "; rv:" + ua.browser.Semver.String() + ") " + ua.browser.String()
	}
	return ua.version + " (" + ua.platform.String() + ") " + ua.browser.String()
}

func New(keys ...string) (ua *UserAgent) {

	ua = &UserAgent{
		version: MozillaWithVersion,
	}

	// if we need specific platform or browser
	// we just need to set the value again with the specific function.
	// this is a easy way to attach our global

	// random platform with version

	// random browser with version
	if ua.browser == nil {
		ua.browser = browserDB[rand.Intn(browserDBSize)].Random()
	}

	if ua.platform == nil {
		x := rand.Intn(platformTypeSize)
		ua.platform = platformDB[x][rand.Intn(platformSize[x])].Random()
	}

	return
}
