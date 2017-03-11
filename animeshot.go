package main

import (
	"fmt"
	"github.com/merrickluo/animeshotbot/search"
	telegram "gopkg.in/telegram-bot-api.v4"
	"strconv"
	"strings"
)

func searchPhotosByText(text string, offset int) []search.Photo {
	return search.ImageForKeyword(text, offset)
}

func answerInlineQuery(bot telegram.BotAPI, queryId string, photos []search.Photo, offset int) {
	var images []interface{}
	for _, photo := range photos {
		result := telegram.InlineQueryResultPhoto{
			Type:        "photo",
			ID:          photo.Sid,
			URL:         photo.Image_large,
			ThumbURL:    photo.Image_thumbnail,
			Title:       photo.Text,
			Caption:     photo.Text,
			Description: photo.Text}
		images = append(images, result)
	}

	var config = telegram.InlineConfig{
		InlineQueryID: queryId,
		Results:       images,
		CacheTime:     0,
		NextOffset:    fmt.Sprintf("%d", offset+1),
	}

	bot.AnswerInlineQuery(config)
}

func answerMessage(bot telegram.BotAPI, chatId int64, photos []search.Photo) {
	if len(photos) == 0 {
		msg := telegram.NewMessage(chatId, "No Result. \nWant to upload some shot? Go to https://as.bitinn.net")
		bot.Send(msg)
		return

	}

	var text string = ""
	for _, photo := range photos {
		text = text + photo.Text + "\n"
	}

	if len(text) >= 0 {
		msg := telegram.NewMessage(chatId, text)
		bot.Send(msg)
	}
}

func answerMessageFullMode(bot telegram.BotAPI, chatId int64, photos []search.Photo) {
	if len(photos) == 0 {
		msg := telegram.NewMessage(chatId, "No Result. \nWant to upload some shot? Go to https://as.bitinn.net")
		bot.Send(msg)
		return
	}

	for _, photo := range photos {
		msg := telegram.NewMessage(chatId, photo.Text+"\n"+photo.Image_large)
		bot.Send(msg)
	}

}

func processUpdate(bot telegram.BotAPI, update telegram.Update) {
	if update.InlineQuery != nil {
		queryText := update.InlineQuery.Query
		offset, err := strconv.ParseInt(update.InlineQuery.Offset, 10, 0)
		if err != nil {
			offset = 1
		}
		photos := searchPhotosByText(queryText, int(offset))
		answerInlineQuery(bot, update.InlineQuery.ID, photos, int(offset))
	} else if update.Message != nil {
		if strings.HasPrefix(update.Message.Text, "/full ") {
			queryText := update.Message.Text[6:]
			photos := searchPhotosByText(queryText, 1)

			answerMessageFullMode(bot, update.Message.Chat.ID, photos)
		} else {
			queryText := update.Message.Text
			photos := searchPhotosByText(queryText, 1)
			answerMessage(bot, update.Message.Chat.ID, photos)
		}
	}
}
