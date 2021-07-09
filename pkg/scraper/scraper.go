package scraper

import (
	"net"
	"net/url"
	"strings"
)

type ScrapedURL struct {
	URL    string `json:"url,omitempty"`
	Status int    `json:"status,omitempty"`
}

type ScrapedData struct {
	URLs       ScrapedURL
	HTML       string
	Headers    map[string][]string
	Scripts    []string
	Cookies    map[string]string
	Meta       map[string][]string
	DNS        map[string][]string
	CertIssuer []string
}

// Scraper is an interface for different scrapping brower (colly, rod)
type Scraper interface {
	Init() error
	CanRenderPage() bool
	Scrape(paramURL string) (*ScrapedData, error)
	EvalJS(jsProp string) (*string, error)
	SetDepth(depth int)
}

func scrapeDNS(paramURL string) map[string][]string {
	scrapedDNS := make(map[string][]string)
	u, _ := url.Parse(paramURL)
	parts := strings.Split(u.Hostname(), ".")
	domain := parts[len(parts)-2] + "." + parts[len(parts)-1]
	nsSlice, _ := net.LookupNS(domain)
	for _, ns := range nsSlice {
		scrapedDNS["NS"] = append(scrapedDNS["NS"], string(ns.Host))
	}
	mxSlice, _ := net.LookupMX(domain)
	for _, mx := range mxSlice {
		scrapedDNS["MX"] = append(scrapedDNS["MX"], string(mx.Host))
	}
	txtSlice, _ := net.LookupTXT(domain)
	scrapedDNS["TXT"] = append(scrapedDNS["TXT"], txtSlice...)
	cname, _ := net.LookupCNAME(domain)
	scrapedDNS["CNAME"] = append(scrapedDNS["CNAME"], cname)

	return scrapedDNS
}
