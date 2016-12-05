package bot

import (
	"github.com/merrickluo/animeshotbot/config"
	"github.com/merrickluo/animeshotbot/shot"

	telegram "gopkg.in/telegram-bot-api.v4"
)

func searchPhotosByText(text string, offset int, count int) []shot.Shot {
	return shot.Search(text, offset, count)
}

func answerInlineQuery(queryID string, photos []shot.Shot) {
	var images []interface{}
	for _, photo := range photos {
		result := telegram.InlineQueryResultPhoto{
			Type:        "photo",
			ID:          photo.Sid,
			URL:         photo.LargeURL,
			ThumbURL:    photo.ThumbnailURL,
			Title:       photo.Text,
			Caption:     photo.Text,
			Description: photo.Text}
		images = append(images, result)
	}

	var config = telegram.InlineConfig{
		InlineQueryID: queryID,
		Results:       images,
		CacheTime:     0,
	}

	_, err := bot.AnswerInlineQuery(config)
	if err != nil {
		logger.Error(err)
	}
}

func sendText(chatID int64, text string) {
	msg := telegram.NewMessage(chatID, text)
	_, err := bot.Send(msg)
	if err != nil {
		logger.Error(err)
	}
}

func answerMessage(chatID int64, photos []shot.Shot) {
	if len(photos) == 0 {
		sendText(chatID, "No Result. \nWant to upload some shot? Go to https://as.bitinn.net")
		return
	}

	var text string
	for _, photo := range photos {
		text = text + photo.Text + "\n"
	}
	msg := telegram.NewMessage(chatID, text)
	buttons := []telegram.InlineKeyboardButton{
		telegram.NewInlineKeyboardButtonData("More", "next"),
	}
	msg.ReplyMarkup = telegram.NewInlineKeyboardMarkup(buttons)

	_, err := bot.Send(msg)
	if err != nil {
		logger.Error(err)
	}
	//	sendText(bot, chatID, text)
}

func answerMessageFullMode(chatID int64, photos []shot.Shot) {
	if len(photos) == 0 {
		sendText(chatID, "No Result. \nWant to upload some shot? Go to https://as.bitinn.net")
		return
	}

	for i, photo := range photos {
		msg := telegram.NewMessage(chatID, photo.Text+"\n"+photo.LargeURL)
		if len(photos) == i {
			buttons := []telegram.InlineKeyboardButton{
				telegram.NewInlineKeyboardButtonData("Next 3 photo", "next"),
			}
			msg.ReplyMarkup = telegram.NewInlineKeyboardMarkup(buttons)
		}
		_, err := bot.Send(msg)
		if err != nil {
			logger.Error(err)
		}
	}
}

func onCommand(message *telegram.Message) {
	userID := message.From.ID
	chatID := message.Chat.ID
	command := message.Command()

	switch command {
	case config.PagingCommand:
		paging(chatID, userID)
	case config.ModeCommand:
		mode(chatID, userID)
	case config.SettingsCommand:
		settings(chatID, userID)
	default:
		sendText(chatID, "Not a valid command")
	}
}

func onCommandValue(message *telegram.Message) {
	userID := message.From.ID
	chatID := message.Chat.ID

	command := config.GetConfig(userID).Command

	switch command {
	case config.PagingCommand:
		pagingValue(chatID, userID, message.Text)

	case config.ModeCommand:
		modeValue(chatID, userID, message.Text)
	}
}

func onInlineQuery(inlineQuery *telegram.InlineQuery) {
	queryText := inlineQuery.Query
	photos := searchPhotosByText(queryText, 0, 5)
	answerInlineQuery(inlineQuery.ID, photos)
}

func onMessage(message *telegram.Message) {
	user := message.From.ID

	c := config.GetConfig(user)

	if c.IsOnCommand() {
		onCommandValue(message)
		return
	}

	photos := searchPhotosByText(message.Text, 0, 5)

	if c.Mode == config.FullMode {
		answerMessageFullMode(message.Chat.ID, photos)
	} else {
		answerMessage(message.Chat.ID, photos)
	}
}

func onUpdate(update telegram.Update) {
	if update.InlineQuery != nil {
		onInlineQuery(update.InlineQuery)
	} else if update.Message != nil {
		if update.Message.IsCommand() {
			onCommand(update.Message)
		} else {
			onMessage(update.Message)
		}
	}
}
