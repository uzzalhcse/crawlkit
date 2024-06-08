package main

import (
	"crawlkit/crawler"
	"github.com/PuerkitoBio/goquery"
	"strings"
)

func customHandler(collection crawler.UrlCollection, fullUrl string, a *goquery.Selection) (string, map[string]interface{}) {
	if strings.Contains(fullUrl, "/tool/category/") {
		return fullUrl, nil
	}
	return "", nil
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
