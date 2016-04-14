package main

import (
	"./search"
	"./telegram"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func processUpdate(update telegram.Update) {
	if len(update.Inline_query.Query) > 0 {
		photos := search.SearchImageForKeyword(update.Inline_query.Query)
		if len(photos) > 0 {
			var images []telegram.InlineQueryResultPhoto
			for _, photo := range photos {
				result := telegram.InlineQueryResultPhoto{"photo", photo.Photo_url, photo.Photo_url, photo.Thumb_url, photo.Title, photo.Title, photo.Title}
				images = append(images, result)
			}
			telegram.AnswerQuery(update.Inline_query.Id, images)
		}
	} else if len(update.Message.Text) > 0 {
		photos := search.SearchImageForKeyword(update.Message.Text)
		if len(photos) > 0 {
			for _, photo := range photos {
				telegram.SendMessage(update.Message.Chat.Id, photo.Title+" "+photo.Photo_url)
			}
		} else {
			telegram.SendMessage(update.Message.Chat.Id, "No Search result > <")
		}
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
