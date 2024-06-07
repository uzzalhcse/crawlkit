package crawler

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/playwright-community/playwright-go"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// processDocument processes the document based on the UrlSelector type
func processDocument(doc *goquery.Document, selector UrlSelector, urlCollection UrlCollection) []UrlCollection {
	if selector.SingleResult {
		// Process a single result
		return processSingleResult(doc, selector, urlCollection)
	} else {
		// Process multiple results
		var items []UrlCollection

		doc.Find(selector.Selector).Each(func(i int, selection *goquery.Selection) {
			item := processSelection(selection, selector, urlCollection)
			items = append(items, item...)
		})

		return items
	}
}

// processSingleResult processes a single result based on the selector
func processSingleResult(doc *goquery.Document, selector UrlSelector, urlCollection UrlCollection) []UrlCollection {
	selection := doc.Find(selector.Selector).First()
	return processSelection(selection, selector, urlCollection)
}

// processSelection processes each selection and extracts attribute values
func processSelection(selection *goquery.Selection, selector UrlSelector, urlCollection UrlCollection) []UrlCollection {
	items := []UrlCollection{}

	selection.Find(selector.FindSelector).Each(func(j int, s *goquery.Selection) {
		attrValue, ok := s.Attr(selector.Attr)
		if !ok {
			log.Println("Attribute not found.")
		} else {
			fullUrl := GetFullUrl(attrValue)
			if selector.Handler != nil {
				url, meta := selector.Handler(urlCollection, fullUrl, s)
				if url != "" {
					items = append(items, UrlCollection{
						Url:      url,
						MetaData: meta,
					})
				}
			} else {
				items = append(items, UrlCollection{
					Url:      fullUrl,
					MetaData: nil,
				})
			}
		}
	})

	return items
}

func getItemsFromAttrOrText(selection *goquery.Selection, selector *CategorySelector) []string {
	var items []string
	selection.Each(func(i int, s *goquery.Selection) {
		var value string
		var ok bool // Declare ok here to avoid shadowing
		s.Find("span.gt").Remove()
		if selector.Attr != "" {
			if value, ok = s.Attr(selector.Attr); ok {
				items = append(items, value)
			}
		} else {
			value = strings.TrimSpace(s.Text())
			items = append(items, value)
		}
	})
	return items
}

func GetPageDom(page playwright.Page) (*goquery.Document, error) {
	html, err := page.Content()
	if err != nil {
		return nil, err
	}
	document, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, err
	}
	return document, nil
}
func WritePageContentToFile(page playwright.Page) error {
	content, err := page.Content()
	if err != nil {
		return err
	}
	content = fmt.Sprintf("<!-- Page Url: %s -->\n%s", page.URL(), content)
	filename := GenerateFilename(page.URL())
	directory := filepath.Join("storage", "logs", App.Name)
	err = os.MkdirAll(directory, 0755)
	if err != nil {
		return err
	}
	filePath := filepath.Join(directory, filename)
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		return err
	}

	return nil
}

// GenerateFilename generates a filename based on URL and current date
func GenerateFilename(rawURL string) string {

	// Replace remaining characters not allowed in file names
	invalidChars := []string{"/", "\\", ":", "*", "?", "\"", "<", ">", "|"}
	for _, char := range invalidChars {
		rawURL = strings.ReplaceAll(rawURL, char, "_")
	}

	// Combine the encoded path with current date and a suitable extension
	currentDate := time.Now().Format("2006-01-02")
	return currentDate + "_" + rawURL + ".html"
}

func GetFullUrl(url string) string {
	fullUrl := ""
	if strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://") {
		// If href is already a full URL, don't concatenate with baseUrl
		fullUrl = url
	} else {
		fullUrl = App.BaseUrl + url
	}
	return fullUrl
}

// shouldBlockResource checks if a resource should be blocked based on its type and URL.
func shouldBlockResource(resourceType string, url string) bool {
	if resourceType == "image" || resourceType == "font" {
		return true
	}

	for _, blockedURL := range App.engine.BlockedURLs {
		if strings.Contains(url, blockedURL) {
			return true
		}
	}

	return false
}
func getBaseUrl(urlString string) string {
	parsedURL, err := url.Parse(urlString)
	if err != nil {
		fmt.Println("failed to parse Url:", "Error", err)
		return ""
	}

	baseURL := parsedURL.Scheme + "://" + parsedURL.Host
	return baseURL
}
