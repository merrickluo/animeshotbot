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

func SearchImageForKeyword(keyword string) []Photo {
	keyword = url.QueryEscape(keyword)
	search_url := BASE_URL + "/api/shots?q=" + keyword

	response, err := http.Get(search_url)
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
