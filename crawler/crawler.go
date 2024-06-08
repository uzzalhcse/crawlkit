package crawler

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/playwright-community/playwright-go"
	"log"
	"sync"
)

// crawlWorker is a worker function that handles crawling URLs and processing results.
// It uses the specified proxy to navigate to URLs and process the results using the provided processor function.
func (e *Engine) crawlWorker(urlChan <-chan UrlCollection, resultChan chan<- interface{}, proxy Proxy, processor interface{}, isLocalEnv bool) {
	browser, page, err := GetBrowserPage(App.pw, App.engine.BrowserType, proxy)
	if err != nil {
		log.Fatalf("failed to initialize browser with proxy: %v\n", err)
	}
	defer browser.Close()
	defer page.Close()

	for {
		urlCollection, more := <-urlChan
		if !more {
			break
		}
		if isLocalEnv && len(resultChan) >= App.engine.DevCrawlLimit {
			return
		}
		log.Printf("Crawling %s using proxy %s", urlCollection.Url, proxy.Server)

		doc, err := NavigateToURL(page, urlCollection.Url)
		if err != nil {
			log.Println("Error navigating to URL:", err)
			continue
		}

		var results interface{}
		switch v := processor.(type) {
		case func(*goquery.Document, *UrlCollection, playwright.Page) []UrlCollection:
			results = v(doc, &urlCollection, page)

		case UrlSelector:
			results = processDocument(doc, v, urlCollection)

		case ProductDetailSelector:
			results = handleProductDetail(doc, urlCollection)

		default:
			log.Fatalf("Unsupported processor type: %T", processor)
		}

		select {
		case resultChan <- results:
		default:
			log.Println("Result channel is full, dropping result")
		}
	}
}

// CrawlUrls initiates the crawling process for the URLs from the specified collection.
// It distributes the work among multiple goroutines and uses proxies if available.
func (e *Engine) CrawlUrls(collection string, processor interface{}) {
	urlCollections := App.GetUrlCollections(collection)
	var items []UrlCollection
	isLocalEnv := App.Config.Site.SiteEnv == Local

	var wg sync.WaitGroup
	urlChan := make(chan UrlCollection, len(urlCollections))
	resultChan := make(chan interface{}, len(urlCollections))

	for _, urlCollection := range urlCollections {
		urlChan <- urlCollection
	}
	close(urlChan)

	proxyCount := len(e.ProxyServers)
	batchSize := App.engine.ConcurrentLimit
	totalUrls := len(urlCollections)
	goroutineCount := min(max(proxyCount, 1)*batchSize, totalUrls) // Determine the required number of goroutines

	for i := 0; i < goroutineCount; i++ {
		proxy := Proxy{}
		if proxyCount > 0 {
			proxy = e.ProxyServers[i%proxyCount]
		}
		wg.Add(1)
		go func(proxy Proxy) {
			defer wg.Done()
			e.crawlWorker(urlChan, resultChan, proxy, processor, isLocalEnv)
		}(proxy)
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	for results := range resultChan {
		switch v := results.(type) {
		case []UrlCollection:
			items = append(items, v...)
			for _, item := range v {
				App.Insert(items, item.Url)
			}
		}
	}

	log.Printf("Total %v urls: %v", App.collection, len(items))
}

// CrawlPageDetail initiates the crawling process for detailed page information from the specified collection.
// It distributes the work among multiple goroutines and uses proxies if available.
func (e *Engine) CrawlPageDetail(collection string) {
	urlCollections := App.GetUrlCollections(collection)
	isLocalEnv := App.Config.Site.SiteEnv == Local

	var wg sync.WaitGroup
	urlChan := make(chan UrlCollection, len(urlCollections))
	resultChan := make(chan interface{}, len(urlCollections))

	for _, urlCollection := range urlCollections {
		urlChan <- urlCollection
	}
	close(urlChan)

	proxyCount := len(e.ProxyServers)
	batchSize := App.engine.ConcurrentLimit
	totalUrls := len(urlCollections)
	goroutineCount := min(max(proxyCount, 1)*batchSize, totalUrls) // Determine the required number of goroutines

	for i := 0; i < goroutineCount; i++ {
		proxy := Proxy{}
		if proxyCount > 0 {
			proxy = e.ProxyServers[i%proxyCount]
		}
		wg.Add(1)
		go func(proxy Proxy) {
			defer wg.Done()
			e.crawlWorker(urlChan, resultChan, proxy, App.ProductDetailSelector, isLocalEnv)
		}(proxy)
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	total := 0
	for results := range resultChan {
		switch v := results.(type) {
		case *ProductDetail:
			fmt.Println("Saving Url", v.Url)
			App.SaveProductDetail(v)
			total++
			if isLocalEnv && total >= App.engine.DevCrawlLimit {
				break
			}
		}
	}

	log.Printf("Total %v %v Inserted ", total, App.collection)
}

// PageSelector adds a new URL selector to the crawler.
func (a *Crawler) PageSelector(selector UrlSelector) *Crawler {
	a.UrlSelectors = append(a.UrlSelectors, selector)
	return a
}

// StartUrlCrawling initiates the URL crawling process for all added selectors.
func (a *Crawler) StartUrlCrawling() *Crawler {
	for _, selector := range a.UrlSelectors {
		a.Collection(selector.ToCollection).
			CrawlUrls(selector.FromCollection, selector)
	}
	return a
}
