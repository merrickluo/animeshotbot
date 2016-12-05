package shot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

const baseURL = "https://as.bitinn.net"
const shotPageSize = 20

//Shot shot
type Shot struct {
	Created      string `json:"created"`
	LargeURL     string `json:"image_large"`
	PreviewURL   string `json:"image_preview"`
	ThumbnailURL string `json:"image_thumbnail"`
	Sid          string `json:"sid"`
	Text         string `json:"text"`
	Updated      string `json:"updated"`
	URL          string `json:"url"`
}

//Search search animeshot
func Search(keyword string, offset int, count int) []Shot {
	page := (offset + shotPageSize) / shotPageSize
	keyword = url.QueryEscape(keyword)
	url := fmt.Sprintf("%s/api/shots?page=%d&q=%s", baseURL, page, keyword)

	response, err := http.Get(url)
	if err != nil {
		return nil
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return nil
	}

	var result []Shot

	err = json.Unmarshal(body, &result)

	if err != nil {
		fmt.Println(err)
	}

	off := offset % shotPageSize
	if len(result) < off {
		return []Shot{}
	}
	newer := result[off:]
	if len(newer) < count {
		return newer
	}
	return newer[:count]
}
