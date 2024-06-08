package crawler

type Proxy struct {
	Server   string
	Username string
	Password string
}
type Engine struct {
	BrowserType     string
	ConcurrentLimit int
	IsDynamic       bool
	DevCrawlLimit   int
	BlockResources  bool
	BlockedURLs     []string
	BoostCrawling   bool
	ProxyServers    []Proxy
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

func (e *Engine) EnableBoostCrawling() *Engine {
	e.BoostCrawling = true
	e.ProxyServers = e.getProxyList()
	return e
}

// Todo: getProxyList should be generate dynamically in future
func (e *Engine) getProxyList() []Proxy {
	var proxies []Proxy
	proxies = append(proxies, Proxy{
		Server: "http://34.146.11.125:3000", // proxy-server-1
	}, Proxy{
		Server: "http://34.146.155.165:3000", // proxy-server-2
	}, Proxy{
		Server: "http://34.143.176.68:3000", // proxy-server-3
	})
	return proxies
}
