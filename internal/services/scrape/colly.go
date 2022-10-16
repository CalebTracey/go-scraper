package scrape

import (
	"crypto/tls"
	"github.com/calebtracey/go-scraper/internal/models"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
	log "github.com/sirupsen/logrus"
	"go.zoe.im/surferua"
	"net"
	"net/http"
	"strings"
	"time"
)

const (
	baseUrl               = "https://www.yellowpages.com"
	baseUrlSearch         = "https://www.yellowpages.com/search?"
	searchTermsField      = "search_terms="
	geoLocationTermsField = "&geo_location_terms="
	urlComma              = "%2C+"
)

const (
	proxy1 = "socks5://184.181.217.204:4145"
	proxy2 = "socks5://98.185.94.94:4145"
	proxy3 = "socks5://184.178.172.18:15280"
)

type ScraperI interface {
	Init() (*colly.Collector, error)
	SetDepth(depth int)
}

type Config struct {
	AppsJSONPath           string
	TimeoutSeconds         int
	LoadingTimeoutSeconds  int
	JSON                   bool
	MaxDepth               int
	visitedLinks           int
	MaxVisitedLinks        int
	MsDelayBetweenRequests int
	UserAgent              string
}

type CollyScraper struct {
	Collector             *colly.Collector
	Transport             *http.Transport
	Response              *http.Response
	TimeoutSeconds        int
	LoadingTimeoutSeconds int
	UserAgent             string
	depth                 int
}

func NewConfig() *Config {
	return &Config{
		AppsJSONPath:           "",
		TimeoutSeconds:         3,
		LoadingTimeoutSeconds:  3,
		JSON:                   true,
		MaxDepth:               0,
		visitedLinks:           0,
		MaxVisitedLinks:        100,
		MsDelayBetweenRequests: 10,
		UserAgent:              surferua.New().Desktop().Chrome().String(),
	}
}

func (c *CollyScraper) CanRenderPage() bool {
	return false
}

func (c *CollyScraper) SetDepth(depth int) {
	c.depth = depth
}

type GoWapTransport struct {
	trans        *http.Transport
	respCallBack func(resp *http.Response)
}

func NewGoWapTransport(t *http.Transport, f func(resp *http.Response)) *GoWapTransport {
	return &GoWapTransport{trans: t, respCallBack: f}
}

func (t *GoWapTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req = req.WithContext(req.Context())
	return t.trans.RoundTrip(req)
}

func (c *CollyScraper) Init() (*colly.Collector, error) {
	log.Infoln("Colly initialization")
	c.Transport = &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   time.Second * time.Duration(c.TimeoutSeconds),
			KeepAlive: 180 * time.Second,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   2 * time.Second,
		ExpectContinueTimeout: time.Duration(c.TimeoutSeconds) * time.Second,
		TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
	}
	c.Collector = colly.NewCollector(
		colly.Async(true),
		colly.MaxDepth(2),
	)
	err := c.Collector.Limit(&colly.LimitRule{DomainGlob: "*", RandomDelay: 1 * time.Second, Parallelism: 6})
	if err != nil {
		return nil, err
	}
	c.Collector.UserAgent = c.UserAgent
	c.Collector.WithTransport(c.Transport)
	c.Collector.OnRequest(func(r *colly.Request) {
		log.Println("Visiting", r.URL)
	})
	c.Collector.OnResponse(func(r *colly.Response) {
		log.Println("response received", r.StatusCode)
	})
	c.Collector.OnError(func(r *colly.Response, err error) {
		log.Println("error:", r.StatusCode, err)
	})
	setResp := func(r *http.Response) {
		c.Response = r
	}
	//proxySwitcher, err := proxy.RoundRobinProxySwitcher(proxy1, proxy2, proxy3)
	//if err != nil {
	//	return nil, err
	//}
	//c.Collector.SetProxyFunc(proxySwitcher)
	c.Collector.WithTransport(NewGoWapTransport(c.Transport, setResp))
	extensions.Referer(c.Collector)

	return c.Collector, nil
}

func BuildScrapeUrl(req models.ScrapeRequest) string {
	terms := strings.ReplaceAll(strings.TrimSpace(req.Terms), " ", "+")
	city := strings.ReplaceAll(strings.TrimSpace(req.City), " ", "+")
	return strings.Join([]string{baseUrlSearch, searchTermsField, terms, geoLocationTermsField, city, urlComma, strings.TrimSpace(req.State)}, "")
}
