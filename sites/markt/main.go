package main

import (
	"crawlkit/constant"
	"crawlkit/crawler"
)

const siteName = "markt"
const siteUrl = "https://markt-mall.jp/"

func main() {
	app := crawler.NewCrawler(siteName, siteUrl, crawler.Engine{
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

	app.Collection(constant.Categories).CrawlUrls(app.GetBaseCollection(), crawler.UrlSelector{
		Selector:       ".l-category-button-list__in",
		SingleResult:   false,
		FindSelector:   "a.c-category-button",
		Attr:           "href",
		ToCollection:   constant.Categories,
		FromCollection: app.GetBaseCollection(),
	})
	app.Collection(constant.Products).CrawlUrls(constant.Categories, handleProducts)

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
	app.Collection(constant.ProductDetails).CrawlPageDetail(constant.Products)
}
