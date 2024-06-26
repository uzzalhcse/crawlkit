package main

import (
	"crawlkit/crawler"
	"github.com/PuerkitoBio/goquery"
	"strings"
)

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
	return urlCollection.Url
}

func getBrandService(document *goquery.Document, urlCollection crawler.UrlCollection) string {
	return urlCollection.Url
}
func getListPriceService(document *goquery.Document, urlCollection crawler.UrlCollection) string {
	var listPrice string
	document.Find("dl#product_detail_standard").Find("span.product_detail_item").Each(func(
		i int, s *goquery.Selection) {
		dt := s.Find("dt").Text()
		if dt == "希望小売価格" {
			listPrice = strings.TrimSuffix(s.Find("dd").Text(), "円")
		}
	})

	return listPrice
}

func breadCrumbHandler(doc *goquery.Document, urlCollection crawler.UrlCollection) string {
	var items []string
	doc.Find(".topicpath.pc_only > *:not(span.gt)").Each(func(i int, s *goquery.Selection) {
		s.Find("span.gt").Remove()
		items = append(items, strings.TrimSpace(s.Text()))
	})
	return strings.Join(items, " > ")
}

func getDescriptionService(doc *goquery.Document, urlCollection crawler.UrlCollection) string {
	// Replace <br> tags with newline characters
	doc.Find("p#product_detail_exp br").Each(func(i int, s *goquery.Selection) {
		s.ReplaceWithHtml("\n")
	})

	// Extract the text content
	descriptionText := doc.Find("p#product_detail_exp").Text()

	// Trim leading and trailing whitespace
	return strings.TrimSpace(descriptionText)
}
