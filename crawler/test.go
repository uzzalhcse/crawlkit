package crawler

import (
	"crawlkit/crawler/constant"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
)

func (e *Engine) TestCrawl(collection string, processor TestSelector) {
	urlCollections := App.GetUrlCollections(collection)
	// Clear the urls slice to avoid appending to existing data
	var items []UrlCollection
	isLocalEnv := App.Config.Site.SiteEnv == constant.LOCAL
	// Function to check if we should continue adding URLs
	shouldContinue := func() bool {
		return !isLocalEnv || len(items) < App.engine.DevCrawlLimit
	}
	total := len(items)
	for _, urlCollection := range urlCollections {
		if !shouldContinue() {
			total = App.engine.DevCrawlLimit
			log.Printf("Crawl limit reached in local environment: %d", App.engine.DevCrawlLimit)
			break
		}
		log.Printf("Crawling %s", urlCollection.Url)
		page, err := GetPage(App.browser)
		if err != nil {
			log.Fatalf("failed to create page: %v\n", err)
		}
		// Navigate to the Url
		doc, err := NavigateToURL(page, urlCollection.Url)
		if err != nil {
			log.Println("Error navigating to Url:", err)
			continue // Continue to the next Url on error
		}

		selector := doc.Find(processor.Selector)
		if processor.UseFirst {
			selector = selector.First()
		}
		selector.Each(func(i int, s *goquery.Selection) {
			selector2 := s.Find(processor.Selector2)
			if processor.UseFirst2 {
				selector2 = selector2.First()
			}
			selector2.Each(func(i int, s *goquery.Selection) {
				selector3 := s.Find(processor.Selector3)
				if processor.UseFirst3 {
					selector3 = selector3.First()
				}
				selector3.Each(func(i int, s *goquery.Selection) {
					var attrValue string
					var ok bool // Declare ok here to avoid shadowing

					if processor.Attr != "" {
						attrValue, ok = s.Attr(processor.Attr)
						if !ok {
							log.Println("Attribute not found.")
							return // Skip processing if attribute not found
						}
					} else {
						attrValue = s.Text()
					}

					if processor.Handler != nil {
						data := processor.Handler(urlCollection, attrValue, s)
						fmt.Println("Handler data", data)
					} else {
						fmt.Println("attrValue", attrValue)
					}
				})
			})

		})

		//doc.Find(processor.Selector).Each(func(i int, selection *goquery.Selection) {
		//	selection.Find(processor.FindSelector).Each(func(j int, s *goquery.Selection) {
		//		attrValue, ok := s.Attr(processor.Attr)
		//		if !ok {
		//			log.Println("Attribute not found.")
		//		} else {
		//			if processor.Handler != nil {
		//				data := processor.Handler(urlCollection, attrValue, s)
		//				fmt.Println("Handler data", data)
		//
		//			} else {
		//				fmt.Println("attrValue", attrValue)
		//			}
		//		}
		//	})
		//})
	}

	log.Printf("Total %v urls: %v", App.collection, total)
}
