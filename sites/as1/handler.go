package main

import (
	"crawlkit/crawler"
	"github.com/PuerkitoBio/goquery"
	"github.com/playwright-community/playwright-go"
	"strings"
)

func categoryHandler(document *goquery.Document, urlCollection *crawler.UrlCollection, page playwright.Page) []crawler.UrlCollection {

	categoryUrls := []crawler.UrlCollection{}
	categoryDiv := document.Find("#menuList > div.submenu > div.accordion > dl")

	// Get the total number of items in the selection
	totalCats := categoryDiv.Length()

	// Iterate over all items except the last three
	categoryDiv.Slice(0, totalCats-3).Each(func(i int, cat *goquery.Selection) {
		cat.Find("dd > div > ul > li ul > li ul > li").Each(func(j int, lMain *goquery.Selection) {
			if href, ok := lMain.Find("a").Attr("href"); ok {
				fullUrl := crawler.App.BaseUrl + href
				categoryUrls = append(categoryUrls, crawler.UrlCollection{
					Url:      fullUrl,
					MetaData: nil,
				})
			} else {
				crawler.App.Logger.Error("Category URL not found.")
			}
		})
	})

	return categoryUrls
}

func productCodeHandler(document *goquery.Document, urlCollection crawler.UrlCollection) []string {
	urlParts := strings.Split(strings.Trim(urlCollection.Url, "/"), "/")
	return []string{urlParts[len(urlParts)-1]}
}

func productNameHandler(document *goquery.Document, urlCollection crawler.UrlCollection) string {
	return strings.Trim(document.Find(".details .intro h2").Text(), " \n")
}

func getUrlHandler(document *goquery.Document, urlCollection crawler.UrlCollection) string {
	return urlCollection.Url
}
