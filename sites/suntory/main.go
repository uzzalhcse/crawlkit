package main

import (
	"crawlkit/constant"
	"crawlkit/crawler"
)

const siteName = "suntory"
const siteUrl = "https://products.suntory.co.jp"

func main() {
	app := crawler.NewCrawler(siteName, siteUrl, crawler.Engine{
		BrowserType:     "chromium",
		ConcurrentLimit: 10,
		IsDynamic:       true,
	})
	app.Start()
	defer app.Stop()
	handleDynamicCrawl(app)

}
func handleDynamicCrawl(app *crawler.Crawler) {
	//categorySelector := crawler.UrlSelector{
	//	Selector:     "ul#drink_list, ul#liquor_list",
	//	SingleResult: false,
	//	FindSelector: "li a",
	//	Attr:         "href",
	//}
	//subCategorySelector := crawler.UrlSelector{
	//	Selector:     "ul.category_list li",
	//	SingleResult: false,
	//	FindSelector: "div.category_order h4 a",
	//	Attr:         "href",
	//	Handler: func(collection crawler.UrlCollection, fullUrl string, a *goquery.Selection) (string, map[string]interface{}) {
	//		brand1 := a.Text()
	//		return fullUrl, map[string]any{
	//			"brand1": brand1,
	//		}
	//	},
	//}
	//productSelector := crawler.UrlSelector{
	//	Selector:     "ul.category_list li",
	//	SingleResult: false,
	//	FindSelector: "div.category_order h4 a",
	//	Attr:         "href",
	//	Handler: func(urlCollection crawler.UrlCollection, fullUrl string, a *goquery.Selection) (string, map[string]interface{}) {
	//		brand2 := a.Text()
	//		return fullUrl, map[string]any{
	//			"brand1": urlCollection.MetaData["brand1"],
	//			"brand2": brand2,
	//		}
	//	},
	//}
	//app.Collection(constant.Categories).SetCrawlLimit(3).CrawlUrls(app.GetBaseCollection(), categorySelector)
	//app.Collection(constant.SUB_CATEGORIES).SetCrawlLimit(5).CrawlUrls(constant.Categories, subCategorySelector)
	//app.Collection(constant.Products).CrawlUrls(constant.SUB_CATEGORIES, productSelector)

	app.ProductDetailSelector = crawler.ProductDetailSelector{
		Jan: getJanService,
		PageTitle: &crawler.SingleSelector{
			Selector: "title",
		},
		Url: getUrlHandler,
		Images: &crawler.MultiSelectors{
			Selectors: []crawler.Selector{
				{Query: "p#product_img img", Attr: "src"},
			},
		},
		ProductCodes: []string{},
		Maker:        "",
		Brand:        "",
		ProductName: &crawler.SingleSelector{
			Selector: "div#product h2",
		},
		Category:    breadCrumbHandler,
		Description: getDescriptionService,
		ListPrice:   getListPriceService,
	}
	app.Collection(constant.ProductDetails).IsDynamicPage(false).SetCrawlLimit(2).CrawlPageDetail(constant.Products)
}
