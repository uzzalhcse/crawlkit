package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log/slog"
	"strings"
)

func main() {

	OpenPage("chromium")

	defer ClosePage()

	pageData, _, err := GetAsyncPageData("https://www.kyocera.co.jp/prdct/tool/category/product", false)
	if err != nil {
		return
	}
	productCategoryUrls := []string{}
	productDetailUrls := []string{}
	otherUrls := []string{}

	pageData.Find("div.index.clearfix ul.clearfix li").Each(func(i int, s *goquery.Selection) {
		s.Find("a").Each(func(j int, h *goquery.Selection) {
			href, ok := h.Attr("href")
			if ok {
				fullUrl := href
				fmt.Println("fullUrl:", fullUrl)
				if strings.Contains(href, "/category/") {
					productCategoryUrls = append(productCategoryUrls, fullUrl)
				} else if strings.Contains(href, "/product/") {
					productDetailUrls = append(productDetailUrls, fullUrl)
				} else {
					otherUrls = append(otherUrls, fullUrl)
				}
			} else {
				slog.Error("URL not found.")
			}
		})

	})
}
