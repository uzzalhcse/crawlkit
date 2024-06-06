package main

import (
	"crawlkit/crawler"
	"crawlkit/crawler/constant"
)

func main() {
	app := crawler.NewCrawler(crawler.Engine{
		BrowserType:     "chromium",
		ConcurrentLimit: 1,
		BlockResources:  true,
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
	app.Collection(constant.CATEGORIES).
		CrawlUrls(constant.SITES, categorySelector)
}
