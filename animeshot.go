package main

import (
	"./search"
	"./telegram"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func processUpdate(update telegram.Update) {
	urls := search.SearchImageForKeyword(update.Inline_query.Query)
	if len(urls) > 0 {
		var images []telegram.InlineQueryResultPhoto
		for _, url := range urls {
			images = append(images, telegram.InlineQueryResultPhoto{"photo", url, url, url})
		}
		telegram.AnswerQuery(update.Inline_query.Id, images)
	}
}

func webhookHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		return
	}

	var result telegram.Update

	err = json.Unmarshal(body, &result)
	go processUpdate(result)
}

func main() {
	http.HandleFunc("/"+telegram.TOKEN, webhookHandler)
	http.ListenAndServe(":8185", nil)
}
