package crawler

import (
	"crawlkit/crawler/constant"
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
	isLocalEnv := App.Config.Site.SiteEnv == constant.LOCAL
	shouldContinue := func() bool {
		return !isLocalEnv || len(items) < App.engine.DevCrawlLimit
	}

	var wg sync.WaitGroup
	urlChan := make(chan UrlCollection, len(urlCollections))

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

				switch v := processor.(type) {
				case func(*goquery.Document, *UrlCollection, playwright.Page) []UrlCollection:
					results := v(doc, &urlCollection, page)
					items = append(items, results...)
					App.Insert(items, urlCollection.Url)

				case UrlSelector:
					results := processDocument(doc, v, urlCollection)
					items = append(items, results...)
					App.Insert(items, urlCollection.Url)

				default:
					funcValue := reflect.ValueOf(processor)
					funcType := funcValue.Type()
					if funcType.Kind() == reflect.Func {
						log.Fatalf("Invalid function signature: expected func(*goquery.Document, *UrlCollection, playwright.Page) []UrlCollection, got %v", funcType)
					} else {
						log.Fatalf("Unsupported type: %T", processor)
					}
				}

			}

			page.Close()
		}()
	}

	wg.Wait()
	log.Printf("Total %v urls: %v", App.collection, len(items))
}

// CrawlPageDetail handles crawling of page details with concurrency
func (e *Engine) CrawlPageDetail(collection string) {
	urlCollections := App.GetUrlCollections(collection)
	isLocalEnv := App.Config.Site.SiteEnv == constant.LOCAL

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
			for urlCollection := range urlChan {
				if isLocalEnv && len(resultChan) >= App.engine.DevCrawlLimit {
					return
				}

				document, err := NavigateToURL(page, urlCollection.Url)
				productDetail := handleProductDetail(document, urlCollection, err)
				resultChan <- productDetail
			}
			page.Close()
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
