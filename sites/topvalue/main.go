package main

import (
	"crawlkit/constant"
	"crawlkit/crawler"
)

const siteName = "topvalu"
const siteUrl = "https://www.topvalu.net/"

func main() {
	app := crawler.NewCrawler(siteName, siteUrl, crawler.Engine{
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
	app.Collection(constant.Categories).
		CrawlUrls(app.GetBaseCollection(), categorySelector)
}
