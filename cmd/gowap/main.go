package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"os"

	gowap "github.com/dranih/gowap/pkg/core"
)

func main() {

	var url, appsJSONPath, scraper, userAgent string
	var help, pretty bool
	var timeoutSeconds, loadingTimeoutSeconds, maxDepth, maxVisitedLinks, msDelayBetweenRequests int
	flag.StringVar(&appsJSONPath, "file", "", "Path to override default technologies.json file")
	flag.StringVar(&scraper, "scraper", "rod", "Choose scraper between rod (default) and colly")
	flag.StringVar(&userAgent, "useragent", "", "Override the user-agent string")
	flag.IntVar(&timeoutSeconds, "timeout", 3, "Timeout in seconds for fetching the url")
	flag.IntVar(&loadingTimeoutSeconds, "loadtimeout", 3, "Timeout in seconds for loading the page")
	flag.IntVar(&maxDepth, "depth", 0, "Don't analyze page when depth superior to this number. Default (0) means no recursivity (only first page will be analyzed)")
	flag.IntVar(&maxVisitedLinks, "maxlinks", 5, "Max number of pages to visit. Exit when reached")
	flag.IntVar(&msDelayBetweenRequests, "delay", 100, "Delay in ms between requests")
	flag.BoolVar(&pretty, "pretty", false, "Pretty print json output")
	flag.BoolVar(&help, "h", false, "Help")
	flag.Parse()

	var Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage : gowap [options] <url>")
		flag.PrintDefaults()
	}

	if help {
		Usage()
		os.Exit(1)
	}
	if flag.NArg() == 0 {
		fmt.Fprintln(os.Stderr, "You must specify a url to analyse")
		Usage()
		os.Exit(1)
	} else if flag.NArg() > 1 {
		fmt.Fprintf(os.Stderr, "Too many arguments %s", flag.Args())
		Usage()
		os.Exit(1)
	} else {
		url = flag.Arg(0)
	}
	if scraper != "rod" && scraper != "colly" {
		fmt.Fprintf(os.Stderr, "Unknown scraper %s : only supporting rod and colly", scraper)
		Usage()
		os.Exit(1)
	}

	config := gowap.NewConfig()
	config.AppsJSONPath = appsJSONPath
	config.TimeoutSeconds = timeoutSeconds
	config.LoadingTimeoutSeconds = loadingTimeoutSeconds
	config.MaxDepth = maxDepth
	config.MaxVisitedLinks = maxVisitedLinks
	config.MsDelayBetweenRequests = msDelayBetweenRequests
	config.Scraper = scraper
	if userAgent != "" {
		config.UserAgent = userAgent
	}

	wapp, err := gowap.Init(config)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	res, err := wapp.Analyze(url)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if pretty {
		var prettyJSON bytes.Buffer
		err = json.Indent(&prettyJSON, []byte(res.(string)), "", "  ")
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		fmt.Println(&prettyJSON)
	} else {
		fmt.Println(res)

	}
}
