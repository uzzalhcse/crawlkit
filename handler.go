package main

import (
	"crawlkit/crawler"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"strings"
)

func customHandler(collection crawler.UrlCollection, fullUrl string, a *goquery.Selection) (string, map[string]interface{}) {
	if strings.Contains(fullUrl, "/tool/category/") {
		return fullUrl, nil
	}
	return "", nil
}
func productCodeHandler(document *goquery.Document, url string) []string {
	fmt.Println("productCodeHandler Url", url)
	urlParts := strings.Split(strings.Trim(url, "/"), "/")
	return []string{urlParts[len(urlParts)-1]}
}

func productNameHandler(document *goquery.Document, url string) string {
	return strings.Trim(document.Find(".details .intro h2").Text(), " \n")
}

func getUrlHandler(document *goquery.Document, url string) string {
	fmt.Println("getUrlHandler Url", url)
	return url
}
