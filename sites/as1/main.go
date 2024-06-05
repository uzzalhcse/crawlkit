package main

import (
	"crawlkit/crawler"
	"crawlkit/crawler/constant"
)

/*
This is under development
*/
func main() {
	app := crawler.NewCrawler(crawler.Engine{
		BrowserType:     "chromium",
		ConcurrentLimit: 10,
		IsDynamic:       false,
	})
	app.Start()
	defer app.Stop()
	handleDynamicCrawl(app)

}
func handleDynamicCrawl(app *crawler.Crawler) {
	//subCategorySelector := crawler.UrlSelector{
	//	Selector:     "#af-categories-list > li",
	//	SingleResult: false,
	//	FindSelector: "a",
	//	Attr:         "href",
	//}
	productSelector := crawler.UrlSelector{
		Selector:     "ul.product-list li.product-item,ul.heightLineParent.clearfix li",
		SingleResult: false,
		FindSelector: "a,div dl dt a",
		Attr:         "href",
	}
	//app.Collection(constant.CATEGORIES).CrawlUrls(constant.SITES, categoryHandler)
	//app.Collection(constant.SUB_CATEGORIES).CrawlUrls(constant.CATEGORIES, subCategorySelector)
	app.Collection(constant.PRODUCTS).CrawlUrls(constant.SUB_CATEGORIES, productSelector)

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
	//app.Collection(constant.PRODUCT_DETAILS).CrawlPageDetail(app.GetUrlsFromCollection(constant.PRODUCTS))
}
