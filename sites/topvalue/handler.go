package main

import (
	"crawlkit/crawler"
	"crawlkit/crawler/constant"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log/slog"
	"strings"
)

func customHandler(url string) string {
	if strings.Contains(url, "/tool/category/") {
		return url
	}
	return ""
}
func categoryHandler(document *goquery.Document, url string) []string {
	productCategoryUrls := []string{}
	productDetailUrls := []string{}
	otherUrls := []string{}

	// Iterate over the desired elements in the document
	document.Find("div.index.clearfix ul.clearfix li").Each(func(i int, s *goquery.Selection) {
		s.Find("a").Each(func(j int, h *goquery.Selection) {
			href, ok := h.Attr("href")
			if ok {
				fullUrl := crawler.App.Config.Site.BaseUrl + href
				// Categorize URLs based on their paths
				if strings.Contains(href, "/category/") {
					productCategoryUrls = append(productCategoryUrls, fullUrl)
				} else if strings.Contains(href, "/product/") {
					productDetailUrls = append(productDetailUrls, fullUrl)
				} else {
					otherUrls = append(otherUrls, fullUrl)
				}
			} else {
				slog.Error("URL not found in anchor tag.")
			}
		})
	})

	// Insert categorized URLs into respective collections
	crawler.App.Collection(constant.OTHER).Insert(otherUrls, url, nil)
	crawler.App.Collection(constant.CATEGORIES).Insert(productCategoryUrls, url, nil)
	crawler.App.Collection(constant.PRODUCTS).Insert(productDetailUrls, url, nil)

	return productCategoryUrls
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
