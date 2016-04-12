package main

import (
	"./searchapi"
	"./telegramapi"
)

func main() {
	updatesch := make(chan []telegramapi.Update)

	go telegramapi.StartFetchUpdates(&updatesch)

	for updates := range updatesch {
		for _, update := range updates {
			urls := searchapi.SearchImageForKeyword(update.Inline_query.Query)
			if len(urls) > 0 {
				var images []telegramapi.InlineQueryResultPhoto
				for _, url := range urls {
					images = append(images, telegramapi.InlineQueryResultPhoto{"photo", url, url, url})
				}
				telegramapi.AnswerQuery(update.Inline_query.Id, images)
			}
		}
	}
}
