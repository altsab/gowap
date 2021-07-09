package core

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/stretchr/testify/assert"
)

func TestBadUrl(t *testing.T) {
	url := "https://badurlformat"
	wapp, err := Init(NewConfig())
	if assert.NoError(t, err, "GoWap Init error") {
		_, err = wapp.Analyze(url)
		assert.Error(t, err, "Bad formatted URL should throw error")
	}

	url = "https://thiswebsitedoesnot.exists"
	_, err = wapp.Analyze(url)
	assert.Error(t, err, "Bad URL should throw error")
}

func TestColly(t *testing.T) {
	ts := MockHTTP(`<html><head><script src="jquery-3.5.1.min.js"></script></head></html>`)
	defer ts.Close()
	config := NewConfig()
	config.Scraper = "colly"
	wapp, err := Init(config)
	if assert.NoError(t, err, "GoWap Init error") {
		res, err := wapp.Analyze(ts.URL)
		if assert.NoError(t, err, "GoWap Analyze error") {
			var output output
			err = json.UnmarshalFromString(res.(string), &output)
			if assert.NoError(t, err, "Unmarshal error") {
				//We should have jquery in the output
				var expected technology
				for _, v := range output.Technologies {
					if v.Name == "jQuery" {
						expected = v
					}
				}
				assert.Equal(t, "3.5.1", expected.Version, "We should find jQuery version 3.5.1")
			}
		}
	}
}

func TestJSEval(t *testing.T) {
	ts := MockHTTP(`<html><head></head><script>jQuery=[];jQuery.fn=[];jQuery.fn.jquery="1.11.3"</script></html>`)
	defer ts.Close()
	config := NewConfig()
	wapp, err := Init(config)
	if assert.NoError(t, err, "GoWap Init error") {
		res, err := wapp.Analyze(ts.URL)
		if assert.NoError(t, err, "GoWap Analyze error") {
			var output output
			err = json.UnmarshalFromString(res.(string), &output)
			if assert.NoError(t, err, "Unmarshal error") {
				//We should have jquery in the output
				var expected technology
				for _, v := range output.Technologies {
					if v.Name == "jQuery" {
						expected = v
					}
				}
				assert.Equal(t, "1.11.3", expected.Version, "We should find jQuery version 1.11.3")
			}
		}
	}
}

func TestDomSearch(t *testing.T) {
	ts := MockHTTP(`<html><head></head><body><a href='https://amzn.to'>Link<a/><div id='jira'></div></body></html>`)
	defer ts.Close()
	config := NewConfig()
	wapp, err := Init(config)
	if assert.NoError(t, err, "GoWap Init error") {
		res, err := wapp.Analyze(ts.URL)
		if assert.NoError(t, err, "GoWap Analyze error") {
			var output output
			err = json.UnmarshalFromString(res.(string), &output)
			if assert.NoError(t, err, "Unmarshal error") {
				var found bool
				for _, v := range output.Technologies {
					if v.Name == "Atlassian Jira" {
						found = true
					}
				}
				assert.True(t, found, "Atlassian Jira should be found in DOM")
			}
		}
	}
}

func TestLoadExternalTechnologiesJSON(t *testing.T) {
	config := NewConfig()
	config.AppsJSONPath = "assets/technologies.json"
	_, err := Init(config)
	assert.NoError(t, err, "External JSON technologies file should open")

	config.AppsJSONPath = "assets/nofile.json"
	_, err = Init(config)
	assert.NoError(t, err, "Should load internal JSON if file not present")
}

func TestTechnologiesFileParsing(t *testing.T) {
	//Bad file format
	config := NewConfig()
	wapp, err := Init(config)
	if assert.NoError(t, err, "GoWap Init error") {
		unmarshalError := []byte(`{"this":"is","notgood"}`)
		err = parseTechnologiesFile(&unmarshalError, wapp)
		assert.Error(t, err, "Unmarshalling should throw an error")
		unmarshalCategoriesError := []byte(`{"categories":{"is":"notgood"}}`)
		err = parseTechnologiesFile(&unmarshalCategoriesError, wapp)
		assert.Error(t, err, "Unmarshalling Categories should throw an error")
		unmarshalAppsError := []byte(`{"categories":{"1":{"name":"CMS","priority":1}},"technologies":{"is":"notgood"}}`)
		err = parseTechnologiesFile(&unmarshalAppsError, wapp)
		assert.Error(t, err, "Unmarshalling Apps should throw an error")
		noCategoryFound := []byte(`{"this":"isgood"}`)
		err = parseTechnologiesFile(&noCategoryFound, wapp)
		assert.Error(t, err, "Should throw an error NoCategoryFound")
		noTechnologyFound := []byte(`{"categories":{"1":{"name":"CMS","priority":1}},"this":"isgood"}`)
		err = parseTechnologiesFile(&noTechnologyFound, wapp)
		assert.Error(t, err, "Should throw an error NoTechnologyFound")
	}

	//Error loading included asset
	embedPath = "does/not/exist"
	_, err = Init(config)
	assert.Error(t, err, "GoWap Init should throw an error trying to open non existing embed file")
	embedPath = "assets/technologies.json"

}

func TestImpliesExcludes(t *testing.T) {
	ts := MockHTTP(`<html><head></head><body><script>Drupal="test"; Backdrop="test";</script><div></div></body></html>`)
	defer ts.Close()
	config := NewConfig()
	wapp, err := Init(config)
	if assert.NoError(t, err, "GoWap Init error") {
		res, err := wapp.Analyze(ts.URL)
		if assert.NoError(t, err, "GoWap Analyze error") {
			var output output
			err = json.UnmarshalFromString(res.(string), &output)
			if assert.NoError(t, err, "Unmarshal error") {
				var found bool
				for _, v := range output.Technologies {
					assert.NotEqual(t, "AngularDart", v.Name, "Backdrop should exclude Drupal")
					if v.Name == "PHP" {
						found = true
					}
				}
				assert.True(t, found, "Backdrop should imply PHP")
			}
		}
	}
}

func TestMeta(t *testing.T) {
	ts := MockHTTP(`<html><head><meta name="generator" content="TiddlyWiki" /></head><body><div></div></body></html>`)
	defer ts.Close()
	config := NewConfig()
	wapp, err := Init(config)
	if assert.NoError(t, err, "GoWap Init error") {
		res, err := wapp.Analyze(ts.URL)
		if assert.NoError(t, err, "GoWap Analyze error") {
			var output output
			err = json.UnmarshalFromString(res.(string), &output)
			if assert.NoError(t, err, "Unmarshal error") {
				var found bool
				for _, v := range output.Technologies {
					if v.Name == "TiddlyWiki" {
						found = true
					}
				}
				assert.True(t, found, "TiddlyWiki should be find in meta")
			}
		}
	}
}

func TestUrl(t *testing.T) {
	config := NewConfig()
	wapp, err := Init(config)
	wapp.Config.TimeoutSeconds = 5
	wapp.Config.LoadingTimeoutSeconds = 5
	if assert.NoError(t, err, "GoWap Init error") {
		res, err := wapp.Analyze("https://twitter.github.io/")
		if assert.NoError(t, err, "GoWap Analyze error") {
			var output output
			err = json.UnmarshalFromString(res.(string), &output)
			if assert.NoError(t, err, "Unmarshal error") {
				var found, foundCert bool
				for _, v := range output.Technologies {
					if v.Name == "GitHub Pages" {
						found = true
					}
					if v.Name == "DigiCert" {
						foundCert = true
					}
				}
				assert.True(t, found, "GitHub Pages should be find in URL")
				assert.True(t, foundCert, "Digicert should be found in certs")
			}
		}
	}
}

func TestHTML(t *testing.T) {
	ts := MockHTTP(`<html><head><title>RoundCube</title></head><body><div></div></body></html>`)
	defer ts.Close()
	config := NewConfig()
	wapp, err := Init(config)
	if assert.NoError(t, err, "GoWap Init error") {
		res, err := wapp.Analyze(ts.URL)
		if assert.NoError(t, err, "GoWap Analyze error") {
			var output output
			err = json.UnmarshalFromString(res.(string), &output)
			if assert.NoError(t, err, "Unmarshal error") {
				var found bool
				for _, v := range output.Technologies {
					if v.Name == "RoundCube" {
						found = true
					}
				}
				assert.True(t, found, "RoundCube should be find in HTML")
			}
		}
	}
	//Testing raw output
	wapp.Config.JSON = false
	res, err := wapp.Analyze(ts.URL)
	if assert.NoError(t, err, "GoWap Analyze error") {
		var found bool
		for _, v := range res.(*output).Technologies {
			if v.Name == "RoundCube" {
				found = true
			}
		}
		assert.True(t, found, "RoundCube should be find in HTML")
	}
}

func TestUnkownScraper(t *testing.T) {
	config := NewConfig()
	config.Scraper = "Unknown"
	_, err := Init(config)
	assert.Error(t, err, "Should throw an error")
}

func TestConfidence(t *testing.T) {
	ts := MockHTTP(`<html><head><script src="alpine.min.js"></script></head><body><div x-data="dropdown()"></div></body></html>`)
	defer ts.Close()
	config := NewConfig()
	wapp, err := Init(config)
	if assert.NoError(t, err, "GoWap Init error") {
		res, err := wapp.Analyze(ts.URL)
		if assert.NoError(t, err, "GoWap Analyze error") {
			var output output
			err = json.UnmarshalFromString(res.(string), &output)
			if assert.NoError(t, err, "Unmarshal error") {
				var found bool
				for _, v := range output.Technologies {
					if v.Name == "Alpine.js" {
						assert.Equal(t, 100, v.Confidence, "Alpine.js confidence should be 100")
						found = true
					}
				}
				assert.True(t, found, "AOS should be found")
			}
		}
	}
}

func TestVersion(t *testing.T) {
	ts := MockHTTP(`<html><head><script src="4.5.6/modernizr.1.2.3.js"></script></head><body><div></div></body></html>`)
	defer ts.Close()
	config := NewConfig()
	wapp, err := Init(config)
	if assert.NoError(t, err, "GoWap Init error") {
		res, err := wapp.Analyze(ts.URL)
		if assert.NoError(t, err, "GoWap Analyze error") {
			var output output
			err = json.UnmarshalFromString(res.(string), &output)
			if assert.NoError(t, err, "Unmarshal error") {
				var found bool
				for _, v := range output.Technologies {
					if v.Name == "Modernizr" {
						assert.Equal(t, "4.5.6", v.Version, "Modernizr version should be 4.5.6")
						found = true
					}
				}
				assert.True(t, found, "Modernizr should be found")
			}
		}
	}

	ts2 := MockHTTP(`<html><head><script src="abc/modernizr.1.2.3.js"></script></head><body><div></div></body></html>`)
	defer ts2.Close()
	res, err := wapp.Analyze(ts2.URL)
	if assert.NoError(t, err, "GoWap Analyze error") {
		var output output
		err = json.UnmarshalFromString(res.(string), &output)
		if assert.NoError(t, err, "Unmarshal error") {
			var found bool
			for _, v := range output.Technologies {
				if v.Name == "Modernizr" {
					assert.Equal(t, "1.2.3", v.Version, "Modernizr version should be 1.2.3")
					found = true
				}
			}
			assert.True(t, found, "Modernizr should be found")
		}
	}
}

func TestParsePattern(t *testing.T) {
	patterns := make(map[string]int)
	//Logging output should be tested here
	parsePatterns(patterns)
	patterns2 := make(map[string]interface{})
	patterns2["test"] = patterns
	parsePatterns(patterns2)
}

func TestAnalyseDom(t *testing.T) {
	app := &application{}
	godoc := &goquery.Document{}
	detectedApp := &detected{}
	app.Dom = false
	//Logging output should be tested here
	analyzeDom(app, godoc, detectedApp)
}

func TestRecursivity(t *testing.T) {
	url := "https://scrapethissite.com/"
	//url := "https://quotes.toscrape.com/"
	config := NewConfig()
	config.MaxDepth = 1
	config.MaxVisitedLinks = 3
	config.Scraper = "colly"
	wapp, err := Init(config)
	if assert.NoError(t, err, "GoWap Init error") {
		res, err := wapp.Analyze(url)
		if assert.NoError(t, err, "GoWap Analyze error") {
			var output output
			err = json.UnmarshalFromString(res.(string), &output)
			if assert.NoError(t, err, "Unmarshal error") {
				assert.Equal(t, 3, len(output.URLs), "Should have parsed 3 URL")
			}
		}
	}
}

func MockHTTP(content string) *httptest.Server {
	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintln(w, content)
			}))
	return ts
}
