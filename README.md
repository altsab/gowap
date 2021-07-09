# Gowap [[Wappalyzer](https://github.com/AliasIO/Wappalyzer) implementation in Go]

[![Build Status](https://github.com/dranih/gowap/workflows/Build%20and%20test/badge.svg)](https://github.com/dranih/gowap/actions?workflow=Build%20and%20test)
[![coverage](https://codecov.io/gh/dranih/gowap/branch/master/graph/badge.svg)](https://codecov.io/gh/dranih/gowap)
[![report card](https://goreportcard.com/badge/github.com/dranih/gowap)](https://goreportcard.com/report/github.com/dranih/gowap)

## Notes

* This is a fork of [unstppbl/gowap](https://github.com/unstppbl/gowap). The main goal here is for me to improve my skills in Go. **Therefore, this is not production ready nor bug free**. Any comment or contribution are welcome.

* A [pull request](https://github.com/unstppbl/gowap/pull/2) is open with the original project

* This implementation adds the following to the original GoWap project :
  - JS analysing (using [Rod](https://github.com/go-rod/rod))
  - DNS scraping
  - Confidence rate
  - Recursive crawling
  - [Rod](https://github.com/go-rod/rod) browser integration ([Colly](https://github.com/gocolly/colly) can still be used - faster but not loading JS)
  - Can be used with as a cmd (technologies.json file embeded)
  - Test coverage 100%
  - robots.txt compliance

## Usage
### Using the package
`go get github.com/dranih/gowap`

Call `Init()` function with a `Config` object created with the `NewConfig()` function. It will return `Wappalyzer` object on which you can call Analyze method with URL string as argument.

```golang
    //Create a Config object and customize it
	config := gowap.NewConfig()
    //Path to override default technologies.json file
	config.AppsJSONPath = "path/to/my/technologies.json"
    //Timeout in seconds for fetching the url
	config.TimeoutSeconds = 5
    //Timeout in seconds for loading the page
	config.LoadingTimeoutSeconds = 5
    //Don't analyze page when depth superior to this number. Default (0) means no recursivity (only first page will be analyzed)
	config.MaxDepth = 2
    //Max number of pages to visit. Exit when reached
	config.MaxVisitedLinks = 10
    //Delay in ms between requests
	config.MsDelayBetweenRequests = 200
    //Choose scraper between rod (default) and colly
	config.Scraper = "colly"
    //Override the user-agent string
	config.UserAgent = "GoWap"
    //Output as a JSON string
    config.JSON = true

    //Initialisation
	wapp, err := gowap.Init(config)
    //Scraping 
    url := "https://scrapethissite.com/"
	res, err := wapp.Analyze(url)

```
### Using the cmd
You can build the cmd using the commande :
`go build -o gowap cmd/gowap/main.go`

Then using the compiled binary :
```
You must specify a url to analyse
Usage : gowap [options] <url>
  -delay int
    	Delay in ms between requests (default 100)
  -depth int
    	Don't analyze page when depth superior to this number. Default (0) means no recursivity (only first page will be analyzed)
  -file string
    	Path to override default technologies.json file
  -h	Help
  -loadtimeout int
    	Timeout in seconds for loading the page (default 3)
  -maxlinks int
    	Max number of pages to visit. Exit when reached (default 5)
  -pretty
    	Pretty print json output
  -scraper string
    	Choose scraper between rod (default) and colly (default "rod")
  -timeout int
    	Timeout in seconds for fetching the url (default 3)
  -useragent string
    	Override the user-agent string
```

## To Do
List of some ideas  :
- [ ] analyse robots (field certIssuer)
- [X] analyse certificates (field certIssuer)
- [ ] anayse css (field css)
- [ ] anayse xhr requests (field xhr)
- [ ] scrape an url list from a file in args
- [ ] ability to choose what is scraped (DNS, cookies, HTML, scripts, etc...)
- [ ] more tests in "real life"
- [ ] perf ? regex html seems long
- [X] should output be the same as original wappalizer ? + ordering