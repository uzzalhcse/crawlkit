package crawler

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/playwright-community/playwright-go"
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

	isLocalEnv := App.Config.Site.SiteEnv == Local
	var browserTypeLaunchOptions playwright.BrowserTypeLaunchOptions
	browserTypeLaunchOptions.Headless = playwright.Bool(!isLocalEnv)
	browserTypeLaunchOptions.Devtools = playwright.Bool(isLocalEnv)

	if App.Config.Site.Proxy != "" {
		// Set proxy options
		browserTypeLaunchOptions.Proxy = &playwright.Proxy{
			Server: App.Config.Site.Proxy,
		}
	}
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

	page, err := browser.NewPage(playwright.BrowserNewPageOptions{
		UserAgent: playwright.String(App.Config.Site.UserAgent),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create page: %w", err)
	}

	// Conditionally intercept and block resources based on configuration
	if App.engine.BlockResources {
		err := page.Route("**/*", func(route playwright.Route) {
			req := route.Request()
			resourceType := req.ResourceType()
			url := req.URL()

			// Check if the resource should be blocked based on resource type or URL
			if shouldBlockResource(resourceType, url) {
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

func GetBrowserPage(pw *playwright.Playwright, browserType string, proxy string) (playwright.Browser, playwright.Page, error) {
	var browser playwright.Browser
	var err error

	isLocalEnv := App.Config.Site.SiteEnv == Local
	var browserTypeLaunchOptions playwright.BrowserTypeLaunchOptions
	browserTypeLaunchOptions.Headless = playwright.Bool(!isLocalEnv)
	browserTypeLaunchOptions.Devtools = playwright.Bool(isLocalEnv)

	if App.Config.Site.Proxy != "" {
		// Set proxy options
		browserTypeLaunchOptions.Proxy = &playwright.Proxy{
			Server: proxy,
		}
	}
	switch browserType {
	case "chromium":
		browser, err = pw.Chromium.Launch(browserTypeLaunchOptions)
	case "firefox":
		browser, err = pw.Firefox.Launch(browserTypeLaunchOptions)
	case "webkit":
		browser, err = pw.WebKit.Launch(browserTypeLaunchOptions)
	default:
		return nil, nil, fmt.Errorf("unsupported browser type: %s", browserType)
	}

	if err != nil {
		return nil, nil, fmt.Errorf("failed to launch browser: %w", err)
	}

	page, err := browser.NewPage(playwright.BrowserNewPageOptions{
		UserAgent: playwright.String(App.Config.Site.UserAgent),
	})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create page: %w", err)
	}

	// Conditionally intercept and block resources based on configuration
	if App.engine.BlockResources {
		err := page.Route("**/*", func(route playwright.Route) {
			req := route.Request()
			resourceType := req.ResourceType()
			url := req.URL()

			// Check if the resource should be blocked based on resource type or URL
			if shouldBlockResource(resourceType, url) {
				route.Abort()
			} else {
				route.Continue()
			}
		})
		if err != nil {
			return nil, nil, fmt.Errorf("failed to set up request interception: %w", err)
		}
	}

	return browser, page, nil
}

func NavigateToURL(page playwright.Page, url string) (*goquery.Document, error) {
	waitUntil := playwright.WaitUntilStateDomcontentloaded
	if App.engine.IsDynamic {
		waitUntil = playwright.WaitUntilStateNetworkidle
	}

	_, err := page.Goto(url, playwright.PageGotoOptions{
		WaitUntil: waitUntil,
		//Timeout:   playwright.Float(0), // Increase timeout to 60 seconds
	})
	if err != nil {
		logErr := WritePageContentToFile(page)
		if logErr != nil {
			return nil, logErr
		}
		//return nil, fmt.Errorf("failed to navigate to Url: %w", err)
	}
	return GetPageDom(page)
}
