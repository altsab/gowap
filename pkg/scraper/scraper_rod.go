package scraper

import (
	"crypto/tls"
	"errors"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
	"github.com/temoto/robotstxt"

	log "github.com/sirupsen/logrus"
)

type RodScraper struct {
	Browser               *rod.Browser
	Page                  *rod.Page
	TimeoutSeconds        int
	LoadingTimeoutSeconds int
	UserAgent             string
	protoUserAgent        *proto.NetworkSetUserAgentOverride
	lock                  *sync.RWMutex
	robotsMap             map[string]*robotstxt.RobotsData
	depth                 int
}

func (s *RodScraper) CanRenderPage() bool {
	return true
}

func (s *RodScraper) SetDepth(depth int) {
	s.depth = depth
}

func (s *RodScraper) Init() error {
	log.Infoln("Rod initialization")
	return rod.Try(func() {
		s.lock = &sync.RWMutex{}
		s.robotsMap = make(map[string]*robotstxt.RobotsData)
		s.protoUserAgent = &proto.NetworkSetUserAgentOverride{UserAgent: s.UserAgent}
		s.Browser = rod.
			New().
			MustConnect().
			MustIgnoreCertErrors(true)
	})
}

func (s *RodScraper) Scrape(paramURL string) (*ScrapedData, error) {

	scraped := &ScrapedData{}

	parsedURL, err := url.Parse(paramURL)
	if err != nil {
		return scraped, err
	}
	if s.depth > 0 {
		if err := s.checkRobots(parsedURL); err != nil {
			return scraped, err
		}
	}

	var e proto.NetworkResponseReceived
	s.Page = s.Browser.MustPage("")
	wait := s.Page.WaitEvent(&e)
	go s.Page.MustHandleDialog()

	errRod := rod.Try(func() {
		s.Page.
			Timeout(time.Duration(s.TimeoutSeconds) * time.Second).
			MustSetUserAgent(s.protoUserAgent).
			MustNavigate(paramURL)
	})
	if errRod != nil {
		log.Errorf("Error while visiting %s : %s", paramURL, errRod.Error())
		return scraped, errRod
	}

	wait()
	if e.Response.SecurityDetails != nil && len(e.Response.SecurityDetails.Issuer) > 0 {
		scraped.CertIssuer = append(scraped.CertIssuer, e.Response.SecurityDetails.Issuer)
	}
	scraped.URLs = ScrapedURL{e.Response.URL, e.Response.Status}
	scraped.Headers = make(map[string][]string)
	for header, value := range e.Response.Headers {
		lowerCaseKey := strings.ToLower(header)
		scraped.Headers[lowerCaseKey] = append(scraped.Headers[lowerCaseKey], value.String())
	}

	scraped.DNS = scrapeDNS(paramURL)

	//TODO : headers and cookies could be parsed before load completed
	errRod = rod.Try(func() {
		s.Page.
			Timeout(time.Duration(s.LoadingTimeoutSeconds) * time.Second).
			MustWaitLoad()
	})
	if errRod != nil {
		log.Errorf("Error while loading %s : %s", paramURL, errRod.Error())
		return scraped, errRod
	}

	scraped.HTML = s.Page.MustHTML()

	scripts, _ := s.Page.Elements("script")
	for _, script := range scripts {
		if src, _ := script.Property("src"); src.Val() != nil {
			scraped.Scripts = append(scraped.Scripts, src.String())
		}
	}

	metas, _ := s.Page.Elements("meta")
	scraped.Meta = make(map[string][]string)
	for _, meta := range metas {
		name, _ := meta.Attribute("name")
		if name == nil {
			name, _ = meta.Attribute("property")
		}
		if name != nil {
			if content, _ := meta.Attribute("content"); content != nil {
				nameLower := strings.ToLower(*name)
				scraped.Meta[nameLower] = append(scraped.Meta[nameLower], *content)
			}
		}
	}

	scraped.Cookies = make(map[string]string)
	str := []string{}
	cookies, _ := s.Page.Cookies(str)
	for _, cookie := range cookies {
		scraped.Cookies[cookie.Name] = cookie.Value
	}

	return scraped, nil
}

func (s *RodScraper) EvalJS(jsProp string) (*string, error) {
	res, err := s.Page.Eval(jsProp)
	if err == nil && res != nil && res.Value.Val() != nil {
		value := ""
		if res.Type == "string" || res.Type == "number" {
			value = res.Value.String()
		}
		return &value, err
	} else {
		return nil, err
	}
}

// checkRobots function implements the robots.txt file checking for rod scraper
// Borrowed from Colly : https://github.com/gocolly/colly/blob/e664321b4e5b94ed568999d37a7cbdef81d61bda/colly.go#L777
// Return nil if no robot.txt or cannot be parsed
func (s *RodScraper) checkRobots(u *url.URL) error {
	s.lock.RLock()
	robot, ok := s.robotsMap[u.Host]
	s.lock.RUnlock()
	if !ok {
		// no robots file cached
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client := &http.Client{Transport: tr}
		resp, err := client.Get(u.Scheme + "://" + u.Host + "/robots.txt")
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		robot, err = robotstxt.FromResponse(resp)
		if err != nil {
			return err
		}
		s.lock.Lock()
		s.robotsMap[u.Host] = robot
		s.lock.Unlock()
	}

	uaGroup := robot.FindGroup(s.UserAgent)

	eu := u.EscapedPath()
	if u.RawQuery != "" {
		eu += "?" + u.Query().Encode()
	}
	if !uaGroup.Test(eu) {
		return errors.New("ErrRobotsTxtBlocked")
	}
	return nil
}
