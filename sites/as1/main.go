package main

import (
	"crawlkit/constant"
	"crawlkit/crawler"
)

/*
This is under development
*/
const siteName = "as1"
const siteUrl = "https://axel.as-1.co.jp/"

func main() {
	app := crawler.NewCrawler(siteName, siteUrl, crawler.Engine{
		BrowserType:     "chromium",
		ConcurrentLimit: 10,
		IsDynamic:       false,
	})
	app.Start()
	defer app.Stop()
	handleDynamicCrawl(app)

}
func handleDynamicCrawl(app *crawler.Crawler) {
	subCategorySelector := crawler.UrlSelector{
		Selector:     "#af-categories-list > li",
		SingleResult: false,
		FindSelector: "a",
		Attr:         "href",
	}
	productSelector := crawler.UrlSelector{
		Selector:     "ul.product-list li.product-item,ul.heightLineParent.clearfix li",
		SingleResult: false,
		FindSelector: "a,div dl dt a",
		Attr:         "href",
	}
	app.Collection(constant.Categories).CrawlUrls(app.GetBaseCollection(), categoryHandler)
	app.Collection(constant.SubCategories).CrawlUrls(constant.Categories, subCategorySelector)
	app.Collection(constant.Products).CrawlUrls(constant.SubCategories, productSelector)

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
