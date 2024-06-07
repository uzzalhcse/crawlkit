package crawler

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/playwright-community/playwright-go"
	"log"
	"log/slog"
	"reflect"
	"sync"
)

// CrawlUrls handles both dynamic crawling and URL extraction based on the provided function or selector
func (e *Engine) CrawlUrls(collection string, processor interface{}) {
	urlCollections := App.GetUrlCollections(collection)
	var items []UrlCollection
	isLocalEnv := App.Config.Site.SiteEnv == Local
	shouldContinue := func() bool {
		return !isLocalEnv || len(items) < App.engine.DevCrawlLimit
	}

	var wg sync.WaitGroup
	urlChan := make(chan UrlCollection, len(urlCollections))
	resultChan := make(chan []UrlCollection, len(urlCollections))

	for _, urlCollection := range urlCollections {
		urlChan <- urlCollection
	}
	close(urlChan)

	for i := 0; i < App.engine.ConcurrentLimit; i++ { // Number of concurrent workers
		wg.Add(1)
		go func() {
			defer wg.Done()

			page, err := GetPage(App.browser)
			if err != nil {
				log.Fatalf("failed to create page: %v\n", err)
			}
			defer page.Close()

			for urlCollection := range urlChan {
				if !shouldContinue() {
					break
				}
				log.Printf("Crawling %s", urlCollection.Url)

				doc, err := NavigateToURL(page, urlCollection.Url)
				if err != nil {
					log.Println("Error navigating to URL:", err)
					continue
				}

				var results []UrlCollection
				switch v := processor.(type) {
				case func(*goquery.Document, *UrlCollection, playwright.Page) []UrlCollection:
					results = v(doc, &urlCollection, page)

				case UrlSelector:
					results = processDocument(doc, v, urlCollection)

				default:
					funcValue := reflect.ValueOf(processor)
					funcType := funcValue.Type()
					if funcType.Kind() == reflect.Func {
						log.Fatalf("Invalid function signature: expected func(*goquery.Document, *UrlCollection, playwright.Page) []UrlCollection, got %v", funcType)
					} else {
						log.Fatalf("Unsupported type: %T", processor)
					}
				}

				select {
				case resultChan <- results:
				default:
					log.Println("Result channel is full, dropping result")
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	for results := range resultChan {
		items = append(items, results...)
		for _, item := range results {
			App.Insert(items, item.Url)
		}
	}

	log.Printf("Total %v urls: %v", App.collection, len(items))
}

// CrawlPageDetail handles crawling of page details with concurrency
func (e *Engine) CrawlPageDetail(collection string) {
	urlCollections := App.GetUrlCollections(collection)
	isLocalEnv := App.Config.Site.SiteEnv == Local

	var wg sync.WaitGroup
	urlChan := make(chan UrlCollection, len(urlCollections))
	resultChan := make(chan *ProductDetail, len(urlCollections))

	for _, urlCollection := range urlCollections {
		urlChan <- urlCollection
	}
	close(urlChan)

	for i := 0; i < App.engine.ConcurrentLimit; i++ { // Number of concurrent workers
		wg.Add(1)
		go func() {
			defer wg.Done()
			page, err := GetPage(App.browser)
			if err != nil {
				log.Fatalf("failed to create page: %v\n", err)
			}
			defer page.Close()

			for urlCollection := range urlChan {
				if isLocalEnv && len(resultChan) >= App.engine.DevCrawlLimit {
					return
				}

				document, err := NavigateToURL(page, urlCollection.Url)
				productDetail := handleProductDetail(document, urlCollection, err)
				select {
				case resultChan <- productDetail:
				default:
					log.Println("Result channel is full, dropping product detail")
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	total := 0
	for productDetail := range resultChan {
		fmt.Println("Saving Url", productDetail.Url)
		App.SaveProductDetail(productDetail)
		total++
		if isLocalEnv && total >= App.engine.DevCrawlLimit {
			break
		}
	}

	slog.Info(fmt.Sprintf("Total %v %v Inserted ", total, App.collection))
}

func (a *Crawler) PageSelector(selector UrlSelector) *Crawler {
	a.UrlSelectors = append(a.UrlSelectors, selector)
	return a
}

func (a *Crawler) StartUrlCrawling() *Crawler {
	for _, selector := range a.UrlSelectors {
		a.Collection(selector.ToCollection).
			CrawlUrls(selector.FromCollection, selector)
	}
	return a
}
