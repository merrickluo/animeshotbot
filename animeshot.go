package main

import (
	"./search"
	"./telegram"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func processUpdate(update telegram.Update) {
	var query string = ""
	var inline bool = false
	var full bool = false
	if len(update.Inline_query.Query) > 0 {
		query = update.Inline_query.Query
		inline = true
	} else if len(update.Message.Text) > 0 {
		text := update.Message.Text
		if strings.HasPrefix(text, "/full ") {
			query = text[6:]
			full = true
		} else {
			query = text
		}
	}

	if len(query) <= 0 {
		return
	}

	photos := search.SearchImageForKeyword(query)
	if len(photos) > 0 {
		if inline {
			var images []telegram.InlineQueryResultPhoto
			for _, photo := range photos {
				result := telegram.InlineQueryResultPhoto{"photo", photo.Photo_url, photo.Photo_url, photo.Thumb_url, photo.Title, photo.Title, photo.Title}
				images = append(images, result)
			}
			telegram.AnswerQuery(update.Inline_query.Id, images)
		} else {
			var text string = ""
			for _, photo := range photos {
				if full {
					telegram.SendMessage(update.Message.Chat.Id, photo.Title+"%0A"+photo.Photo_url)
				} else {
					text = text + photo.Title + "%0A"
				}
			}
			if len(text) >= 0 {
				telegram.SendMessage(update.Message.Chat.Id, text)
			}
		}
	} else if !inline {
		telegram.SendMessage(update.Message.Chat.Id, "No Search result > <")
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

var mode = "Release"

func main() {
	if mode == "Release" {
		http.HandleFunc("/"+telegram.Token, webhookHandler)
		http.ListenAndServe(":8185", nil)
	} else {
		updatesch := make(chan []telegram.Update)

		go telegram.StartFetchUpdates(&updatesch)
		fmt.Println("eee")

		for updates := range updatesch {
			for _, update := range updates {
				go processUpdate(update)
			}
		}
	}
}
