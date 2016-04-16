package main

import (
	"./search"
	"./telegram"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
)

func processInlineQuery(query telegram.InlineQuery) {
	if len(query.Id) == 0 {
		return
	}
	var images []telegram.InlineQueryResultPhoto
	photos := search.SearchImageForKeyword(query.Query)
	for _, photo := range photos {
		result := telegram.InlineQueryResultPhoto{"photo", photo.Image_large, photo.Image_large, photo.Image_thumbnail, photo.Text, photo.Text, photo.Text}
		images = append(images, result)
	}
	telegram.AnswerQuery(query.Id, images)
}

func processMessage(message telegram.Message) {
	if message.Message_id == 0 {
		return
	}
	query := message.Text
	full := false

	if strings.HasPrefix(query, "/full ") {
		query = query[6:]
		full = true
	}

	photos := search.SearchImageForKeyword(query)

	var text string = ""
	for _, photo := range photos {
		if full {
			telegram.SendMessage(message.Chat.Id, photo.Text+"%0A"+photo.Image_large)
		} else {
			text = text + photo.Text + "%0A"
		}
	}
	if len(photos) == 0 {
		text = "No Result. %0AWant to upload some shot? Go to https://as.bitinn.net"
	}
	if len(text) >= 0 {
		telegram.SendMessage(message.Chat.Id, text)
	}
}

func processUpdate(update telegram.Update) {
	processInlineQuery(update.Inline_query)
	processMessage(update.Message)
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

		for updates := range updatesch {
			for _, update := range updates {
				go processUpdate(update)
			}
		}
	}
}
