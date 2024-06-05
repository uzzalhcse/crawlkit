package main

import (
	"crawlkit/crawler"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/playwright-community/playwright-go"
	"strings"
	"time"
)

func handleProducts(document *goquery.Document, collection *crawler.UrlCollection, page playwright.Page) []crawler.UrlCollection {
	var urls []crawler.UrlCollection
	productLinkSelector := "a.c-text-link.u-color-text--link.c-text-link--underline"
	clickAndWaitButton(".u-hidden-sp li button", page)

	items, err := page.Locator("ul.p-card-list-no-scroll li.p-product-card.p-product-card--large").All()
	if err != nil {
		fmt.Println("Error fetching items:", err)
		return urls
	}

	for i, item := range items {
		err := item.Click(playwright.LocatorClickOptions{Timeout: playwright.Float(10000)})
		if err != nil {
			fmt.Println("Failed to click on Product Card:", err)
			continue
		}

		// Wait for the modal to open and the link to be available
		_, err = page.WaitForSelector(productLinkSelector, playwright.PageWaitForSelectorOptions{
			Timeout: playwright.Float(10000),
		})
		if err != nil {
			fmt.Println("Timeout waiting for product link:", err)
			crawler.WritePageContentToFile(page)
			continue
		}

		doc, err := crawler.GetPageDom(page)
		if err != nil {
			fmt.Println("Error getting page DOM:", err)
			continue
		}

		productLink, exist := doc.Find(productLinkSelector).First().Attr("href")

		fullUrl := crawler.GetFullUrl(productLink)
		if !exist {
			fmt.Println("Failed to find product link")
		} else {
			fmt.Println("Product Link:", fullUrl)
		}

		// Close the modal
		closeModal := page.Locator("#__next > div.l-background__wrap > div.l-background__in > div > button")
		if closeModal != nil {
			err = closeModal.Click(playwright.LocatorClickOptions{Timeout: playwright.Float(5000)})
			if err != nil {
				fmt.Println("Failed to close modal:", err)
				crawler.WritePageContentToFile(page)
			}
		} else {
			fmt.Println("Modal close button not found.")
		}

		if exist {
			urls = append(urls, crawler.UrlCollection{Url: fullUrl}) // Assuming you want to collect URLs
		}

		// Add a delay after every 15 items
		if (i+1)%5 == 0 {
			fmt.Println("Sleeping for 3 seconds...")
			time.Sleep(3 * time.Second)
		}
	}

	return urls
}
func clickAndWaitButton(selector string, page playwright.Page) {
	for {
		button := page.Locator(selector)
		err := button.Click()
		page.WaitForSelector(selector)
		if err != nil {
			fmt.Println("No more button available")
			break
		}
	}
}
func productCodeHandler(document *goquery.Document, url string) []string {
	fmt.Println("productCodeHandler Url", url)
	urlParts := strings.Split(strings.Trim(url, "/"), "/")
	return []string{urlParts[len(urlParts)-1]}
}

func productNameHandler(document *goquery.Document, url string) string {
	return strings.Trim(document.Find("h2.ProductInfo_Head_Main_ProductName").Text(), " \n")
}

func getUrlHandler(document *goquery.Document, urlCollection crawler.UrlCollection) string {
	return urlCollection.Url
}
