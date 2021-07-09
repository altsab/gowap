package scraper

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDnsScraping(t *testing.T) {
	scraperTest := &CollyScraper{}
	err := scraperTest.Init()
	assert.NoError(t, err, "Scraper Init error")
	res, err := scraperTest.Scrape("https://scrapethissite.com/")
	assert.NoError(t, err, "Colly scraping error")
	assert.NotEmpty(t, res.DNS, "There should be some DNS results")
}

func TestCollyScraper(t *testing.T) {
	scraperTest := &CollyScraper{}

	assert.False(t, scraperTest.CanRenderPage(), "Colly cannot render JS")
	_, err := scraperTest.EvalJS("jQuery")
	assert.Error(t, err, "Colly cannot render JS")

	err = scraperTest.Init()
	assert.NoError(t, err, "Scraper Init error")

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		c := &http.Cookie{Name: "test", Value: "testv", HttpOnly: false}
		http.SetCookie(w, c)
		w.WriteHeader(200)
		//nolint:errcheck
		w.Write([]byte(`<html><head><meta property="generator" content="TiddlyWiki" /></head><script scr="jquery.js"/><body><div></div></body></html>`))
	})
	ts := httptest.NewServer(mux)
	defer ts.Close()
	res, err := scraperTest.Scrape(ts.URL)
	if assert.NoError(t, err, "Scrap should work") {
		assert.NotEmpty(t, res.HTML, "There should be some HTML content")
	}
}

func TestRodScraper(t *testing.T) {
	scraperTest := &RodScraper{TimeoutSeconds: 2, LoadingTimeoutSeconds: 2}

	assert.True(t, scraperTest.CanRenderPage(), "Rod can render JS")

	err := scraperTest.Init()
	assert.NoError(t, err, "Scraper Init error")

	ts := MockHTTP("<html><script>var now = Date.now();var end = now + 2000;while (now < end) { now = Date.now(); }</script></html>")
	defer ts.Close()
	err = scraperTest.Init()
	scraperTest.LoadingTimeoutSeconds = 0
	assert.NoError(t, err, "GoWap Init error")
	_, err = scraperTest.Scrape(ts.URL)
	assert.Error(t, err, "Timeout should throw error")
	scraperTest.LoadingTimeoutSeconds = 2

	url := "https://doesnotexist"
	err = scraperTest.Init()
	assert.NoError(t, err, "GoWap Init error")
	_, err = scraperTest.Scrape(url)
	assert.Error(t, err, "Bad URL should throw error")

	url = ":foo"
	err = scraperTest.Init()
	assert.NoError(t, err, "GoWap Init error")
	_, err = scraperTest.Scrape(url)
	assert.Error(t, err, "Bad URL should throw error")

	ts = MockHTTP(`<html><head><meta property="generator" content="TiddlyWiki" /></head><body><div></div></body></html>`)
	defer ts.Close()
	res, err := scraperTest.Scrape(ts.URL)
	if assert.NoError(t, err, "Scrap should work") {
		assert.Equal(t, "TiddlyWiki", res.Meta["generator"][0], "Scrap meta should work")
	}
	scraperTest.TimeoutSeconds = 0
	_, err = scraperTest.Scrape(ts.URL + "/doesnotexists")
	assert.Error(t, err, "Timeout should throw error")
	scraperTest.TimeoutSeconds = 2

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		c := &http.Cookie{Name: "test", Value: "testv", HttpOnly: false}
		http.SetCookie(w, c)
		w.WriteHeader(200)
		//nolint:errcheck
		w.Write([]byte(`<html><head><meta property="generator" content="TiddlyWiki" /></head><script scr="jquery.js"/><body><div></div></body></html>`))
	})
	ts = httptest.NewServer(mux)
	defer ts.Close()

	res, err = scraperTest.Scrape(ts.URL)
	assert.NoError(t, err, "Colly scraping error")
	assert.NotEmpty(t, res.HTML, "There should be some HTML content")
	resJS, err := scraperTest.EvalJS(`"test"`)
	assert.Equal(t, "test", *resJS, "Test string should eval as test string...")
	assert.NoError(t, err, "Rod should render JS")
	resJS, err = scraperTest.EvalJS("this.should.throw.error")
	assert.Nil(t, resJS, "Should return nil")
	assert.Error(t, err, "Rod should throw error on rendering bad JS")

	url = "https://twitter.github.io/"
	err = scraperTest.Init()
	assert.NoError(t, err, "GoWap Init error")
	res, err = scraperTest.Scrape(url)
	if assert.NoError(t, err, "GoWap Analyze error") {
		assert.Contains(t, res.CertIssuer[0], "DigiCert", "Scrap cert manager should work")
	}
}

func TestRobot(t *testing.T) {

	var robotsFile = `
		User-agent: GoWap
		Allow: /allowed
		Disallow: /disallowed
		Disallow: /allowed*q=
		`

	mux := http.NewServeMux()
	mux.HandleFunc("/robots.txt", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		//nolint:errcheck
		w.Write([]byte(robotsFile))
	})

	mux.HandleFunc("/allowed", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		//nolint:errcheck
		w.Write([]byte("allowed"))
	})

	mux.HandleFunc("/disallowed", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		//nolint:errcheck
		w.Write([]byte("disallowed"))
	})
	ts := httptest.NewServer(mux)
	defer ts.Close()

	collyScraperTest := &CollyScraper{UserAgent: "GoWap"}
	err := collyScraperTest.Init()
	collyScraperTest.SetDepth(1)
	if assert.NoError(t, err, "Scraper Init error") {
		_, err := collyScraperTest.Scrape(ts.URL + "/allowed")
		assert.NoError(t, err, "Robot should allowed this url")
		_, err = collyScraperTest.Scrape(ts.URL + "/disallowed")
		assert.Error(t, err, "Robot should block this url")
	}

	rodScraperTest := &RodScraper{TimeoutSeconds: 2, LoadingTimeoutSeconds: 2, UserAgent: "GoWap"}
	err = rodScraperTest.Init()
	rodScraperTest.SetDepth(1)
	if assert.NoError(t, err, "Scraper Init error") {
		_, err := rodScraperTest.Scrape(ts.URL + "/allowed")
		assert.NoError(t, err, "Robot should allowed this url")
		_, err = rodScraperTest.Scrape(ts.URL + "/disallowed")
		assert.Error(t, err, "Robot should block this url")
		_, err = rodScraperTest.Scrape(ts.URL + "/allowed?q=1")
		assert.Error(t, err, "Robot should block this url")
	}

	rodScraperTest.UserAgent = "NotListed"
	_, err = rodScraperTest.Scrape(ts.URL + "/disallowed")
	assert.NoError(t, err, "Robot should not block this url (user agent not listed)")

	robotsFile = `Disallow: /`
	mux2 := http.NewServeMux()
	mux2.HandleFunc("/robots.txt", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		//nolint:errcheck
		w.Write([]byte(robotsFile))
	})

	ts2 := httptest.NewServer(mux2)
	defer ts2.Close()

	_, err = rodScraperTest.Scrape(ts2.URL)
	assert.Error(t, err, "Bad robot format should throw an error")

	_, err = rodScraperTest.Scrape("https://doesnotexist")
	assert.Error(t, err, "Navigation should fail")

}

func MockHTTP(content string) *httptest.Server {
	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintln(w, content)
			}))
	return ts
}
