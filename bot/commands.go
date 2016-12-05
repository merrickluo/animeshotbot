package bot

import (
	"fmt"
	"strconv"

	"github.com/merrickluo/animeshotbot/config"
	telegram "gopkg.in/telegram-bot-api.v4"
)

func settings(chatID int64, userID int) {
	c := config.GetConfig(userID)
	mode := "Caption"
	if c.Mode == config.FullMode {
		mode = "Full"
	}
	sendText(chatID, fmt.Sprintf("Mode: %s,\nPaging:%d.", mode, c.Paging))
}

func paging(chatID int64, userID int) {
	config.GetConfig(userID).OnCommand(config.PagingCommand)

	msg := telegram.NewMessage(chatID,
		"How many shots will return per request? Currently support 1, 2, 5, 10.")
	msg.ReplyMarkup = telegram.NewReplyKeyboard(
		telegram.NewKeyboardButtonRow(
			telegram.NewKeyboardButton("1"),
			telegram.NewKeyboardButton("2"),
		),
		telegram.NewKeyboardButtonRow(
			telegram.NewKeyboardButton("5"),
			telegram.NewKeyboardButton("10"),
		),
	)
	_, err := bot.Send(msg)
	if err != nil {
		logger.Error(err)
	}
}

func mode(chatID int64, userID int) {
	config.GetConfig(userID).OnCommand(config.ModeCommand)

	msg := telegram.NewMessage(chatID,
		"Should I return full image or just caption? 1 for full and 2 for caption.")
	msg.ReplyMarkup = telegram.NewReplyKeyboard(
		telegram.NewKeyboardButtonRow(
			telegram.NewKeyboardButton("Full"),
		),
		telegram.NewKeyboardButtonRow(
			telegram.NewKeyboardButton("Caption"),
		),
	)
}

func pagingValue(chatID int64, userID int, value string) {
	count, err := strconv.Atoi(value)
	if err != nil {
		sendText(chatID, "must be a number")
		return
	}
	if count != 1 && count != 2 && count != 5 && count != 10 {
		sendText(chatID, "must be 1, 2, 5, or 10")
		return
	}

	c := config.GetConfig(userID)
	c.Paging = count
	c.Command = ""
	c.Save()

	msg := telegram.NewMessage(chatID,
		fmt.Sprintf("We are now returning %d shots every time", count))
	msg.ReplyMarkup = struct{}{}
	bot.Send(msg)
}

func modeValue(chatID int64, userID int, value string) {
	if value != "Full" && value != "Caption" {
		sendText(chatID, "must be \"Full\" or \"Caption\"")
		return
	}
	c := config.GetConfig(userID)
	if value == "Full" {
		c.Mode = config.FullMode
	} else {
		c.Mode = config.CaptionMode
	}
	c.Command = ""
	c.Save()
	sendText(chatID, fmt.Sprintf("We are now in %s mode", value))
}
