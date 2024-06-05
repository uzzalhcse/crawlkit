package crawler

import (
	"crawlkit/config"
	"crawlkit/crawler/constant"
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

type Crawler struct {
	*Client
	Config                *config.Config
	pw                    *playwright.Playwright
	browser               playwright.Browser
	Page                  playwright.Page
	collection            string
	url                   string
	UrlSelectors          []UrlSelector
	ProductDetailSelector ProductDetailSelector
	engine                *Engine
}

func NewCrawler(engines ...Engine) *Crawler {
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
		}

		App = &Crawler{
			Config:     config.NewConfig(),
			collection: constant.SITES,
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
	page, err := GetPage(browser)
	if err != nil {
		log.Fatalf("failed to create page: %v\n", err)
	}
	a.Client = client
	a.pw = pw
	a.browser = browser
	a.Page = page
}
func (a *Crawler) Stop() {
	defer a.pw.Stop()
	defer a.browser.Close()
	defer a.Page.Close()
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
