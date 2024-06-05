package crawler

import "github.com/PuerkitoBio/goquery"

type UrlSelector struct {
	Selector       string
	SingleResult   bool
	FindSelector   string
	Attr           string
	ToCollection   string
	FromCollection string
	Handler        func(urlCollection UrlCollection, fullUrl string, a *goquery.Selection) (string, map[string]interface{})
}

type TestSelector struct {
	Selector  string
	Selector2 string
	Selector3 string
	Attr      string
	UseFirst  bool
	UseFirst2 bool
	UseFirst3 bool
	Handler   func(urlCollection UrlCollection, data string, a *goquery.Selection) any
}

type CategorySelector struct {
	*DataSelector
}
type DataSelector struct {
	Query    string
	Attr     string
	UseFirst bool
	Handler  func(urlCollection UrlCollection, data string, a *goquery.Selection) any
}

type JanSelector struct {
	Selector string
}

type SingleSelector struct {
	Selector string
}

type Selector struct {
	Query        string // CSS selector
	Attr         string // Attribute to extract (e.g., "src" or "href")
	SingleResult bool
}

type MultiSelectors struct {
	Selectors []Selector // Array of selectors
}

type ProductDetailSelector struct {
	Jan              interface{}
	PageTitle        interface{}
	Url              interface{}
	Images           interface{}
	ProductCodes     interface{}
	Maker            interface{}
	Brand            interface{}
	ProductName      interface{}
	Category         interface{}
	Description      interface{}
	Reviews          interface{}
	ItemTypes        interface{}
	ItemSizes        interface{}
	ItemWeights      interface{}
	SingleItemSize   interface{}
	SingleItemWeight interface{}
	NumOfItems       interface{}
	ListPrice        interface{}
	SellingPrice     interface{}
	Attributes       interface{}
}
