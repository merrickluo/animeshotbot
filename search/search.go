package search

import (
	"fmt"
	"golang.org/x/net/html"
	"net/http"
	"net/url"
	"strings"
)

const BASE_URL = "https://as.bitinn.net"

type Photo struct {
	Photo_url string
	Thumb_url string
	Title     string
}

func SearchImageForKeyword(keyword string) []Photo {
	keyword = url.QueryEscape(keyword)
	realUrl := BASE_URL + "/search?q=" + keyword
	var result []Photo

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

	var hasclass func(attrs []html.Attribute, class string) bool
	hasclass = func(attrs []html.Attribute, class string) bool {
		for _, attr := range attrs {
			if attr.Key == "class" && attr.Val == class {
				return true
			}
		}
		return false
	}

	var travel func(*html.Node)
	travel = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "div" && hasclass(n.Attr, "main-preview") {
			photo := Photo{}
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				if c.Type == html.ElementNode && c.Data == "blockquote" {
					for d := c.FirstChild; d != nil; d = d.NextSibling {
						if d.Data == "p" {
							photo.Title = d.FirstChild.Data
						} else if d.Data == "cite" {
							in := d.FirstChild
							for _, attr := range in.Attr {
								if attr.Key == "href" {
									strs := strings.Split(attr.Val, "/")
									photo.Thumb_url = BASE_URL + "/upload/" + strs[len(strs)-1] + ".300.jpg"
									photo.Photo_url = BASE_URL + "/upload/" + strs[len(strs)-1] + ".1200.jpg"
									result = append(result, photo)
								}
							}
						}
					}
				}
			}
		} else {
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				travel(c)
			}
		}
	}

	travel(doc)

	return result
}
