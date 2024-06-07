package main

import (
	"crawlkit/constant"
	"crawlkit/crawler"
	"github.com/PuerkitoBio/goquery"
	"strings"
)

const siteName = "kyocera"
const siteUrl = "https://www.kyocera.co.jp/prdct/tool/category/product"

func main() {
	app := crawler.NewCrawler(siteName, siteUrl, crawler.Engine{
		BrowserType:     "chromium",
		ConcurrentLimit: 10,
		IsDynamic:       false,
		BlockResources:  true,
		BlockedURLs:     []string{"syncsearch.jp"},
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
	app.Collection(constant.Categories).CrawlUrls(app.GetBaseCollection(), categorySelector)
	app.Collection(constant.Products).CrawlUrls(app.GetBaseCollection(), categoryProductSelector)
	app.Collection(constant.Other).CrawlUrls(app.GetBaseCollection(), categoryOtherSelector)
	app.Collection(constant.Products).CrawlUrls(constant.Categories, productSelector)
	app.Collection(constant.Products).CrawlUrls(constant.Other, productSelector)

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
	app.Collection(constant.ProductDetails).CrawlPageDetail(constant.Products)
}
