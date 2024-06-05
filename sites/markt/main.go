package main

import (
	"crawlkit/crawler"
	"crawlkit/crawler/constant"
)

func main() {
	app := crawler.NewCrawler(crawler.Engine{
		BrowserType:     "chromium",
		ConcurrentLimit: 1,
		DevCrawlLimit:   1,
		IsDynamic:       true,
		BlockResources:  true,
	})
	app.Start()
	defer app.Stop()
	handleDynamicCrawl(app)

}
func handleDynamicCrawl(app *crawler.Crawler) {

	//app.Collection(constant.CATEGORIES).CrawlUrls(constant.SITES, crawler.UrlSelector{
	//	Selector:       ".l-category-button-list__in",
	//	SingleResult:   false,
	//	FindSelector:   "a.c-category-button",
	//	Attr:           "href",
	//	ToCollection:   constant.CATEGORIES,
	//	FromCollection: constant.SITES,
	//})
	app.Collection(constant.PRODUCTS).CrawlUrls(constant.CATEGORIES, handleProducts)

	app.ProductDetailSelector = crawler.ProductDetailSelector{
		Jan: "",
		PageTitle: &crawler.SingleSelector{
			Selector: "title",
		},
		Url: getUrlHandler,
		Images: &crawler.MultiSelectors{
			Selectors: []crawler.Selector{
				{Query: "img#image-item", Attr: "src"},
				{Query: "section.ProductDetail_Section_Function a", Attr: "href"},
				{Query: "section.ProductDetail_Section_Spec img", Attr: "src"},
			},
		},
		ProductCodes: productCodeHandler,
		Maker:        "",
		Brand:        "",
		ProductName:  productNameHandler,
		Category:     "",
		Description:  "",
	}
	//app.Collection(constant.PRODUCT_DETAILS).CrawlPageDetail(constant.PRODUCTS)
}
