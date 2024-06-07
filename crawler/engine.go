package crawler

type Engine struct {
	BrowserType     string
	ConcurrentLimit int
	IsDynamic       bool
	DevCrawlLimit   int
	BlockResources  bool
	BlockedURLs     []string
}

func (e *Engine) SetBrowserType(browserType string) *Engine {
	e.BrowserType = browserType
	return e
}

func (e *Engine) SetConcurrentLimit(concurrentLimit int) *Engine {
	e.ConcurrentLimit = concurrentLimit
	return e
}

func (e *Engine) IsDynamicPage(isDynamic bool) *Engine {
	e.IsDynamic = isDynamic
	return e
}

func (e *Engine) SetCrawlLimit(crawlLimit int) *Engine {
	e.DevCrawlLimit = crawlLimit
	return e
}
func (e *Engine) SetBlockResources(block bool) *Engine {
	e.BlockResources = block
	return e
}
