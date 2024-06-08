package main

import (
	"crawlkit/crawler"
	"github.com/PuerkitoBio/goquery"
	"strings"
)

func productCodeHandler(document *goquery.Document, urlCollection crawler.UrlCollection) []string {
	crawler.App.Logger.Info("Um,here...")
	urlParts := strings.Split(strings.Trim(urlCollection.Url, "/"), "/")
	return []string{urlParts[len(urlParts)-1]}
}

func productNameHandler(document *goquery.Document, urlCollection crawler.UrlCollection) string {
	return strings.Trim(document.Find("h2.ProductInfo_Head_Main_ProductName").Text(), " \n")
}

func getUrlHandler(document *goquery.Document, urlCollection crawler.UrlCollection) string {
	return urlCollection.Url
}
