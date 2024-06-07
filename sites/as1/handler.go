package main

import (
	"crawlkit/crawler"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/playwright-community/playwright-go"
	"log/slog"
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
				slog.Info(fullUrl)
			} else {
				slog.Error("Category URL not found.")
			}
		})
	})

	return categoryUrls
}

func productCodeHandler(document *goquery.Document, url string) []string {
	fmt.Println("productCodeHandler Url", url)
	urlParts := strings.Split(strings.Trim(url, "/"), "/")
	return []string{urlParts[len(urlParts)-1]}
}

func productNameHandler(document *goquery.Document, url string) string {
	return strings.Trim(document.Find(".details .intro h2").Text(), " \n")
}

func getUrlHandler(document *goquery.Document, urlCollection crawler.UrlCollection) string {
	return urlCollection.Url
}
