package main

import (
	"crawlkit/constant"
	"crawlkit/crawler"
)

const siteName = "aqua"
const siteUrl = "https://aqua-has.com"

func main() {
	app := crawler.NewCrawler(siteName, siteUrl, crawler.Engine{
		BrowserType:     "chromium",
		ConcurrentLimit: 1,
		DevCrawlLimit:   100,
		BlockResources:  true,
		BlockedURLs:     []string{},
		BoostCrawling:   false,
		//ProxyServers:    []crawler.Proxy{},
	})
	app.Start()
	defer app.Stop()
	handleDynamicCrawl(app)

}
func handleDynamicCrawl(app *crawler.Crawler) {
	app.PageSelector(
		crawler.UrlSelector{
			Selector:       "ul.Header_Navigation_List_Item_Sub_Group_Inner",
			SingleResult:   true,
			FindSelector:   "a",
			Attr:           "href",
			ToCollection:   constant.Categories,
			FromCollection: app.GetBaseCollection(),
		})
	app.PageSelector(
		crawler.UrlSelector{
			Selector:       "div.CategoryTop_Series_Item_Content_List",
			SingleResult:   false,
			FindSelector:   "a",
			Attr:           "href",
			ToCollection:   constant.Products,
			FromCollection: constant.Categories,
		})
	app.StartUrlCrawling()

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
