package main

import (
	"crawlkit/crawler"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log/slog"
	"strings"
)

func categoryHandler(document *goquery.Document, urlCollection crawler.UrlCollection) ([]string, []interface{}) {
	items := []string{}

	// Iterate over the desired elements in the document
	document.Find("div.index.clearfix ul.clearfix li").Each(func(i int, s *goquery.Selection) {
		s.Find("a").Each(func(j int, h *goquery.Selection) {
			href, ok := h.Attr("href")
			if ok {
				fullUrl := crawler.App.Config.Site.BaseUrl + href

				items = append(items, fullUrl)
			} else {
				slog.Error("URL not found in anchor tag.")
			}
		})
	})

	return items, nil
}

func productCodeHandler(document *goquery.Document, urlCollection crawler.UrlCollection) []string {
	fmt.Println("productCodeHandler Url", urlCollection.Url)
	urlParts := strings.Split(strings.Trim(urlCollection.Url, "/"), "/")
	return []string{urlParts[len(urlParts)-1]}
}

func getJanService(doc *goquery.Document, urlCollection crawler.UrlCollection) string {
	var janCode string
	doc.Find("dl#product_detail_standard").Find("span.product_detail_item").Each(func(i int, s *goquery.Selection) {
		dt := s.Find("dt").Text()
		if dt == "JANコード" {
			janCode = s.Find("dd").Text()
		}
	})
	return janCode
}

func getUrlHandler(document *goquery.Document, urlCollection crawler.UrlCollection) string {
	fmt.Println("getUrlHandler Url", urlCollection.Url)
	return urlCollection.Url
}

func getBrandService(document *goquery.Document, urlCollection crawler.UrlCollection) string {
	fmt.Println("getUrlHandler Url", urlCollection.Url)
	return urlCollection.Url
}
