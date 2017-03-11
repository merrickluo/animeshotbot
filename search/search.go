package search

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

const BASE_URL = "https://as.bitinn.net"

type Photo struct {
	Created         string
	Image_large     string
	Image_preview   string
	Image_thumbnail string
	Sid             string
	Text            string
	updated         string
	url             string
}

//ImageForKeyword search image for keyword
func ImageForKeyword(keyword string, offset int) []Photo {
	keyword = url.QueryEscape(keyword)
	searchURL := fmt.Sprintf("%s/api/shots?q=%s&page=%d", BASE_URL, keyword, offset)

	response, err := http.Get(searchURL)
	if err != nil {
		return nil
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return nil
	}

	var result []Photo

	err = json.Unmarshal(body, &result)

	if err != nil {
		fmt.Println(err)
	}

	return result
}
