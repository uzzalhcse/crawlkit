package crawler

import (
	"crawlkit/config"
	"fmt"
	"github.com/playwright-community/playwright-go"
	"log"
	"log/slog"
	"sync"
	"time"
)

var Once sync.Once
var App *Crawler
var startTime time.Time

const Local = "local"
const baseCollection = "sites"

type Crawler struct {
	*Client
	Name                  string
	Url                   string
	BaseUrl               string
	Config                *config.Config
	pw                    *playwright.Playwright
	browser               playwright.Browser
	collection            string
	url                   string
	UrlSelectors          []UrlSelector
	ProductDetailSelector ProductDetailSelector
	engine                *Engine
}

func NewCrawler(name, url string, engines ...Engine) *Crawler {
	Once.Do(func() {
		startTime = time.Now()
		slog.Info("Program started! 🚀")

		// Create default engine configuration
		defaultEngine := Engine{
			BrowserType:     "chromium",
			ConcurrentLimit: 10,
			IsDynamic:       false,
			DevCrawlLimit:   10,
			BlockResources:  false,
			BlockedURLs: []string{
				"www.googletagmanager.com",
				"google.com",
				"googleapis.com",
				"gstatic.com",
			},
		}

		// Override defaults with provided engine configuration if available
		if len(engines) > 0 {
			eng := engines[0]
			if eng.BrowserType != "" {
				defaultEngine.BrowserType = eng.BrowserType
			}
			if eng.ConcurrentLimit > 0 {
				defaultEngine.ConcurrentLimit = eng.ConcurrentLimit
			}
			if eng.IsDynamic {
				defaultEngine.IsDynamic = eng.IsDynamic
			}
			if eng.DevCrawlLimit > 0 {
				defaultEngine.DevCrawlLimit = eng.DevCrawlLimit
			}
			if eng.BlockResources {
				defaultEngine.BlockResources = eng.BlockResources
			}
			defaultEngine.BlockedURLs = append(defaultEngine.BlockedURLs, eng.BlockedURLs...)
		}

		App = &Crawler{
			Name:       name,
			Url:        url,
			BaseUrl:    getBaseUrl(url),
			Config:     config.NewConfig(),
			collection: App.GetBaseCollection(),
			engine:     &defaultEngine,
		}
	})

	return App
}

func (a *Crawler) Start() {
	client := ConnectDB()
	client.NewSite()
	pw, err := GetPlaywright()
	if err != nil {
		log.Fatalf("failed to initialize playwright: %v\n", err)
	}

	browser, err := GetBrowser(pw, a.engine.BrowserType)
	if err != nil {
		log.Fatalf("failed to launch browser: %v\n", err)
	}

	a.Client = client
	a.pw = pw
	a.browser = browser
}

func (a *Crawler) Stop() {
	if a.browser != nil {
		a.browser.Close()
	}
	if a.pw != nil {
		a.pw.Stop()
	}
	duration := time.Since(startTime)
	slog.Info(fmt.Sprintf("Program stopped in ⚡ %v", duration))
}

func (a *Crawler) Collection(collection string) *Engine {
	a.collection = collection
	return a.engine
}

func (a *Crawler) GetUrl() string {
	return a.url
}

func (a *Crawler) GetCollection() string {
	return a.collection
}

func (a *Crawler) GetBaseCollection() string {
	return baseCollection
}
