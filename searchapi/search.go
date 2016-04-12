package searchapi

import (
	"fmt"
	"golang.org/x/net/html"
	"net/http"
	"net/url"
	"strings"
)

const BASE_URL = "https://as.bitinn.net"

type PhotUrls struct {
	Url []string
}

func SearchImageForKeyword(keyword string) []string {
	keyword = url.QueryEscape(keyword)
	realUrl := BASE_URL + "/search?q=" + keyword
	var result []string

	fmt.Println(realUrl)
	response, err := http.Get(realUrl)
	if err != nil {
		return nil
	}

	defer response.Body.Close()

	doc, err := html.Parse(response.Body)

	if err != nil {
		fmt.Println(err)
		return nil
	}

	var hasclass func(attrs []html.Attribute) bool
	hasclass = func(attrs []html.Attribute) bool {
		for _, attr := range attrs {
			if attr.Key == "class" {
				return true
			}
		}
		return false
	}

	var travel func(*html.Node)
	travel = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			if hasclass(n.Attr) {
				for _, attr := range n.Attr {
					if attr.Key == "href" {
						strs := strings.Split(attr.Val, "/")
						url := "/upload/" + strs[len(strs)-1] + ".1200.jpg"
						result = append(result, BASE_URL+url)
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			travel(c)
		}
	}

	travel(doc)

	return result
}
