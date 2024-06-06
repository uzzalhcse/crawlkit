package test

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/playwright-community/playwright-go"
	"log/slog"
	"os"
	"strings"
)

var page playwright.Page

func OpenPage(browserType string) {
	//err := playwright.Install()
	//if err != nil {
	//	slog.Error("Failed to install playwright", "error", err)
	//}
	pw, err := playwright.Run()
	if err != nil {
		slog.Info(fmt.Sprintf("Failed to launch playwright: %v", err))
	}

	var browser playwright.Browser
	launchOptions := playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(true),
	}

	switch browserType {
	case "chromium":
		browser, err = pw.Chromium.Launch(launchOptions)
	case "firefox":
		browser, err = pw.Firefox.Launch(launchOptions)
	case "webkit":
		browser, err = pw.WebKit.Launch(launchOptions)
	default:
		slog.Error(fmt.Sprintf("unsupported browser type: %s", browserType))
	}
	if err != nil {
		slog.Info(fmt.Sprintf("failed to launch browser: %w", err))
	}
	// Check for USER_AGENT environment variable
	userAgent := os.Getenv("USER_AGENT")
	if userAgent == "" {
		userAgent = pw.Devices["Desktop Edge"].UserAgent
	}

	page, err = browser.NewPage(playwright.BrowserNewPageOptions{
		UserAgent: playwright.String(userAgent),
	})
	if err != nil {
		slog.Info(fmt.Sprintf("failed to create page: %w", err))
	}

}
func ClosePage() {
	err := page.Close()
	if err != nil {
		slog.Error(err.Error())
	}
}
func GetAsyncPageData(url string, isDynamicSite bool) (*goquery.Document, playwright.Page, error) {
	document, err := NavigateToURL(page, url, isDynamicSite)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to navigate: %v\n", err)
	}
	return document, page, nil
}

func NavigateToURL(page playwright.Page, url string, isDynamicSite bool) (*goquery.Document, error) {
	waitUntil := playwright.WaitUntilStateDomcontentloaded
	if isDynamicSite {
		waitUntil = playwright.WaitUntilStateNetworkidle
	}

	_, err := page.Goto(url, playwright.PageGotoOptions{
		WaitUntil: waitUntil,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to navigate to URL: %w", err)
	}
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
