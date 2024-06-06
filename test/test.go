package test

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log/slog"
	"strings"
)

func Kyocera() {

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
func Suntory() {

	OpenPage("chromium")

	defer ClosePage()
	categories := []string{}

	doc, _, err := GetAsyncPageData("https://products.suntory.co.jp?ke=hd", true)
	if err != nil {
		slog.Error(err.Error())
		return
	}

	html, err := doc.Html()
	fmt.Println("HTML:", html)

	// Extract data from the HTML
	doc.Find("ul#drink_list").Each(func(i int, s *goquery.Selection) {
		s.Find("li a").Each(func(j int, a *goquery.Selection) {
			href, exists := a.Attr("href")
			if exists {
				if !strings.HasPrefix(href, "https") {
					href = "https://products.suntory.co.jp" + href
				}
				fmt.Println("link:", href)
				categories = append(categories, href)
			}
		})
	})
	fmt.Println("Total categories:", len(categories))
}
