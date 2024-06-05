package crawler

import (
	"crawlkit/crawler/constant"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/playwright-community/playwright-go"
	"os"
)

func GetPlaywright() (*playwright.Playwright, error) {
	err := playwright.Install()
	if err != nil {
		return nil, err
	}
	pw, err := playwright.Run()
	if err != nil {
		fmt.Printf("failed to initialize playwright: %v\n", err)
		return nil, err
	}
	return pw, nil
}

func GetBrowser(pw *playwright.Playwright, browserType string) (playwright.Browser, error) {
	var browser playwright.Browser
	var err error

	isLocalEnv := App.Config.Site.SiteEnv == constant.LOCAL
	var browserTypeLaunchOptions playwright.BrowserTypeLaunchOptions
	browserTypeLaunchOptions.Headless = playwright.Bool(!isLocalEnv)
	browserTypeLaunchOptions.Devtools = playwright.Bool(isLocalEnv)

	switch browserType {
	case "chromium":
		browser, err = pw.Chromium.Launch(browserTypeLaunchOptions)
	case "firefox":
		browser, err = pw.Firefox.Launch(browserTypeLaunchOptions)
	case "webkit":
		browser, err = pw.WebKit.Launch(browserTypeLaunchOptions)
	default:
		return nil, fmt.Errorf("unsupported browser type: %s", browserType)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to launch browser: %w", err)
	}
	return browser, nil
}

func GetPage(browser playwright.Browser) (playwright.Page, error) {
	// Check for USER_AGENT environment variable
	userAgent := os.Getenv("SITE_USER_AGENT")

	page, err := browser.NewPage(playwright.BrowserNewPageOptions{
		UserAgent: playwright.String(userAgent),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create page: %w", err)
	}

	// Conditionally intercept and block images, CSS, and fonts based on configuration
	if App.engine.BlockResources {
		err = page.Route("**/*", func(route playwright.Route) {
			if route.Request().ResourceType() == "image" || route.Request().ResourceType() == "font" {
				//fmt.Println("Abort request", route.Request().ResourceType())
				route.Abort()
			} else {
				route.Continue()
			}
		})
		if err != nil {
			return nil, fmt.Errorf("failed to set up request interception: %w", err)
		}
	}

	return page, nil
}

func NavigateToURL(page playwright.Page, url string) (*goquery.Document, error) {
	waitUntil := playwright.WaitUntilStateDomcontentloaded
	if App.engine.IsDynamic {
		waitUntil = playwright.WaitUntilStateNetworkidle
	}
	_, err := page.Goto(url, playwright.PageGotoOptions{
		WaitUntil: waitUntil,
		Timeout:   playwright.Float(60000),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to navigate to Url: %w", err)
	}
	return GetPageDom(page)
}
