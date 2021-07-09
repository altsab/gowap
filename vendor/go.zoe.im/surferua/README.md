<p align="center">
    <img src="https://cloud.githubusercontent.com/assets/597902/16172506/9debc136-357a-11e6-90fb-c7c46f50dff0.png" alt="surferua">
</p>
<p align="center">
    <i>Logo from <a href="https://github.com/avct/uasurfer">uasurfer</a></i>
</p>
<h3 align="center">Surfer UA</h3>
<p align="center">High performance User-Agent generator in Golang.</p>
<p align="center">
    <a href="https://travis-ci.org/jiusanzhou/surferua"><img src="https://img.shields.io/travis/jiusanzhou/surferua.svg?label=build"></a>
	<a href="https://godoc.org/github.com/jiusanzhou/surferua"><img src="https://img.shields.io/badge/godoc-reference-blue.svg"></a>
	<a href="https://goreportcard.com/report/jiusanzhou/surferua"><img src="https://goreportcard.com/badge/github.com/jiusanzhou/surferua"></a>
	<a href="https://twitter.com/jiusanzhou" title="@Zoe on Twitter"><img src="https://img.shields.io/badge/twitter-@jiusanzhou-55acee.svg" alt="@Zoe on Twitter"></a>
	<a href="https://sourcegraph.com/github.com/jiusanzhou/surferua?badge" title="SuerferUA on Sourcegraph"><img src="https://sourcegraph.com/github.com/jiusanzhou/surferua/-/badge.svg" alt="SurferUA on Sourcegraph"></a>
</p>

---

Surfer User-Agent (surferua) is a lightweight Golang package that generate HTTP User-Agent strings with particular attention to device type.


## Start


### Basic usage

- Install this package to your `$GOPATH` with command: `go get go.zoe.im/surferua`

- **Enjoy it!**


Some example code below:

```golang
package main

import (
	"go.zoe.im/surferua"

	"fmt"
)

func main() {
	fmt.Println(surferua.New().String())

	// functions depends on your generated inputting data.

	fmt.Println(surferua.NewBot())
	fmt.Println(surferua.NewBotGoogle())
	fmt.Println(surferua.New().Phone().String())
	fmt.Println(surferua.New().Android().String())
	fmt.Println(surferua.New().Desktop().Chrome().String())
}
```

### Add more UA

- Open the `config.yml` with any editor you like.
- Edit the c`yaml` file as you can image what it do

```yaml
# You can get all User-Agent from: http://www.webapps-online.com/online-tools/user-agent-strings
# If you want to add UA of some version or change them,
# you just need to edit this file and generate code with `go generate` again.
browsers:
  Firefox:
    engine:
      name: Gecko
    version:
      major: [35, 56]
      minor: [35, 56]
      patch: [0, 3]
  Chrome:
    engine:
      name: AppleWebKit
      version:
        major: [534, 603]
        minor: [35, 56]
    version:
      major: [39, 64]
      minor: 0
      patch: [0, 3000]
  Safari:
    engine:
      name: WebKit
      version:
        major: [534, 603]
        minor: [0, 21]
        patch: [0, 10]
    version:
      major: [5, 11]
      minor: [0, 2]
      patch: [0, 10]
platforms:
  Desktop:
    Linux:
    MacOS:
    Windows:
  Phone:
    iOS:
      version:
        major: [6, 11]
        minor: [0, 3]
        patch: [0, 3]
    Android:
      version:
        major: [4, 8]
        minor: [0, 4]
        patch: [0, 4]
bots: # Auto upper the first letter
  google: Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)
  bing: Mozilla/5.0 (compatible; bingbot/2.0; +http://www.bing.com/bingbot.htm)
  yahoo: Mozilla/5.0 (compatible; Yahoo! Slurp; http://help.yahoo.com/help/us/ysearch/slurp)
  mj12: Mozilla/5.0 (compatible; MJ12bot/v1.4.5; http://www.majestic12.co.uk/bot.php)
  simplePie: SimplePie/1.3.1 (Feed Parser; http://simplepie.org; Allow like Gecko)
  blex: Mozilla/5.0 (compatible; BLEXBot/1.0; +http://webmeup-crawler.com/)
  yandex: Mozilla/5.0 (compatible; YandexBot/3.0; +http://yandex.com/bots)
  scoutJet: Mozilla/5.0 (compatible; ScoutJet; +http://www.scoutjet.com/)
  duck: DuckDuckBot/1.1; (+http://duckduckgo.com/duckduckbot.html)
  baidu: Mozilla/5.0 (compatible; Baiduspider/2.0; +http://www.baidu.com/search/spider.html)
  sogou: Sogou web spider/4.0(+http://www.sogou.com/docs/help/webmasters.htm#07)
  sogouOrin: Sogou Orion spider/3.0( http://www.sogou.com/docs/help/webmasters.htm#07)
  exa: Mozilla/5.0 (compatible; Konqueror/3.5; Linux) KHTML/3.5.5 (like Gecko) (Exabot-Thumbnails)
  facebook: facebookexternalhit/1.1 (+http://www.facebook.com/externalhit_uatext.php)
  alexa: ia_archiver (+http://www.alexa.com/site/help/webmasters; crawler@alexa.com)
  msn: msnbot/2.0b (+http://search.msn.com/msnbot.htm)
```

- `go generate` or `make gen` in `*nix`

## TODO

- [ ] Add functional way to create UA factory
- [ ] Add testing
