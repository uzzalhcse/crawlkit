package main

import (
	"crawlkit/crawler"
	"crawlkit/crawler/constant"
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
		Selector:     "div.pulldown__list--head div.pulldown__inner--parent",
		SingleResult: false,
		FindSelector: "a.pulldown__ttl",
		Attr:         "href",
	}
	//categoryProductSelector := crawler.UrlSelector{
	//	Selector:     "div.index.clearfix ul.clearfix li",
	//	SingleResult: false,
	//	FindSelector: "a",
	//	Attr:         "href",
	//	Handler: func(url string) string {
	//		if strings.Contains(url, "/tool/product/") {
	//			return url
	//		}
	//		return ""
	//	},
	//}
	//categoryOtherSelector := crawler.UrlSelector{
	//	Selector:     "div.index.clearfix ul.clearfix li",
	//	SingleResult: false,
	//	FindSelector: "a",
	//	Attr:         "href",
	//	Handler: func(url string) string {
	//		if strings.Contains(url, "/tool/sgs/") {
	//			return url
	//		}
	//		return ""
	//	},
	//}
	//productSelector := crawler.UrlSelector{
	//	Selector:     "ul.product-list li.product-item,ul.heightLineParent.clearfix li",
	//	SingleResult: false,
	//	FindSelector: "a,div dl dt a",
	//	Attr:         "href",
	//}
	app.Collection(constant.CATEGORIES).
		CrawlUrls(app.GetUrlsFromCollection(constant.SITES), categorySelector)
	//app.Collection(constant.PRODUCTS).
	//	CrawlUrls(app.GetUrlsFromCollection(constant.SITES), categoryProductSelector)
	//app.Collection(constant.OTHER).
	//	CrawlUrls(app.GetUrlsFromCollection(constant.SITES), categoryOtherSelector)
	//app.Collection(constant.PRODUCTS).
	//	CrawlUrls(app.GetUrlsFromCollection(constant.CATEGORIES), productSelector)
	//app.Collection(constant.PRODUCTS).
	//	CrawlUrls(app.GetUrlsFromCollection(constant.OTHER), productSelector)

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
	app.Collection(constant.PRODUCT_DETAILS).CrawlPageDetail(app.GetUrlsFromCollection(constant.PRODUCTS))
}
