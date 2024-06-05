package main

import (
	"crawlkit/crawler"
	"crawlkit/crawler/constant"
	"github.com/PuerkitoBio/goquery"
	"strings"
)

func main() {
	app := crawler.NewCrawler(crawler.Engine{
		BrowserType:     "chromium",
		ConcurrentLimit: 10,
	})
	app.Start()
	defer app.Stop()
	handleDynamicCrawl(app)

}
func handleDynamicCrawl(app *crawler.Crawler) {
	categorySelector := crawler.UrlSelector{
		Selector:     "div.index.clearfix ul.clearfix li",
		SingleResult: false,
		FindSelector: "a",
		Attr:         "href",
		Handler:      customHandler,
	}
	categoryProductSelector := crawler.UrlSelector{
		Selector:     "div.index.clearfix ul.clearfix li",
		SingleResult: false,
		FindSelector: "a",
		Attr:         "href",
		Handler: func(collection crawler.UrlCollection, fullUrl string, a *goquery.Selection) (string, map[string]interface{}) {
			if strings.Contains(fullUrl, "/tool/product/") {
				return fullUrl, nil
			}
			return "", nil
		},
	}
	categoryOtherSelector := crawler.UrlSelector{
		Selector:     "div.index.clearfix ul.clearfix li",
		SingleResult: false,
		FindSelector: "a",
		Attr:         "href",
		Handler: func(collection crawler.UrlCollection, fullUrl string, a *goquery.Selection) (string, map[string]interface{}) {
			if strings.Contains(fullUrl, "/tool/sgs/") {
				return fullUrl, nil
			}
			return "", nil
		},
	}
	productSelector := crawler.UrlSelector{
		Selector:     "ul.product-list li.product-item,ul.heightLineParent.clearfix li",
		SingleResult: false,
		FindSelector: "a,div dl dt a",
		Attr:         "href",
	}
	app.Collection(constant.CATEGORIES).CrawlUrls(constant.SITES, categorySelector)
	app.Collection(constant.PRODUCTS).CrawlUrls(constant.SITES, categoryProductSelector)
	app.Collection(constant.OTHER).CrawlUrls(constant.SITES, categoryOtherSelector)
	app.Collection(constant.PRODUCTS).CrawlUrls(constant.CATEGORIES, productSelector)
	app.Collection(constant.PRODUCTS).CrawlUrls(constant.OTHER, productSelector)

	app.ProductDetailSelector = crawler.ProductDetailSelector{
		Jan: "",
		PageTitle: &crawler.SingleSelector{
			Selector: "title",
		},
		Url: getUrlHandler,
		Images: &crawler.MultiSelectors{
			Selectors: []crawler.Selector{
				{Query: ".details .intro .image img", Attr: "src"},
			},
		},
		ProductCodes: "",
		Maker:        "",
		Brand:        "",
		ProductName:  productNameHandler,
		Category:     "",
		Description:  "",
	}
	app.Collection(constant.PRODUCT_DETAILS).CrawlPageDetail(constant.PRODUCTS)
}
