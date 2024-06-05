package helper

import (
	"crawlkit/crawler"
	"github.com/PuerkitoBio/goquery"
	"regexp"
	"strings"
)

func GetText(selection *goquery.Selection) string {
	trimIndex := []int{}

	textSlice := selection.Contents().Map(func(i int, s *goquery.Selection) string {
		tagName := strings.ToLower(goquery.NodeName(s))
		if s.ChildrenFiltered("*").Length() == 0 {
			if tagName == "br" {
				return "\n"
			}
			if tagName == "img" {
				src, ok := s.Attr("src")
				if ok {
					if !strings.HasPrefix(src, "http") {
						src = crawler.App.Config.Site.BaseUrl + src
					}
				}
				return src + "\n"
			}
			if tagName == "style" {
				return "\n"
			}

			text := s.Text()
			text = strings.ReplaceAll(text, "\u00A0", " ")
			text = strings.Trim(text, " \t\n")
			if len(text) > 0 {
				text += "\n"
			}
			if tagName == "sup" || tagName == "sub" {
				if i-1 > 0 {
					trimIndex = append(trimIndex, i-1)
				}
				text = strings.Trim(text, "\n")
			}

			return text
		} else if tagName == "table" {
			return GetTableDataAsString(s) + "\n"
		} else {
			return GetText(s) + "\n"
		}
	})

	for _, i := range trimIndex {
		textSlice[i] = strings.Trim(textSlice[i], " \n")
	}

	return strings.Trim(strings.Join(textSlice, ""), " \n")
}

func GetTableDataAsString(selection *goquery.Selection) string {
	data := ""

	selection.Find("tr").Each(func(i int, s *goquery.Selection) {
		items := []string{}
		s.Children().Each(func(j int, h *goquery.Selection) {
			items = append(items, GetText(h))
		})

		if len(items) > 0 {
			data += strings.Join(items, "/") + "\n"
			return
		}
	})
	data = strings.Trim(data, " \n")

	return data
}
func GetNumericsAsString(str string) string {
	reg := regexp.MustCompile(`[^0-9]`)
	str = reg.ReplaceAllString(str, "")

	return str
}
