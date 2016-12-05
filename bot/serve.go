package bot

import (
	"net/http"

	"github.com/merrickluo/animeshotbot/logging"
	telegram "gopkg.in/telegram-bot-api.v4"
)

var (
	bot    *telegram.BotAPI
	logger = logging.Logger()
)

//Serve start bot service
func Serve(port string, token string, mode string, debug bool) {
	b, err := telegram.NewBotAPI(token)
	if err != nil {
		logger.Fatal(err)
	}
	bot = b

	bot.Debug = debug
	logger.Info("Authorized on account " + bot.Self.UserName)

	var updates <-chan telegram.Update

	if mode == "fetch" {
		config := telegram.NewWebhook("")
		_, err = bot.SetWebhook(config)
		if err != nil {
			logger.Error(err)
		}

		u := telegram.NewUpdate(0)
		u.Timeout = 60
		updates, err = bot.GetUpdatesChan(u)

		if err != nil {
			logger.Fatal(err)
		}
	} else if mode == "webhook" {
		url := "https://merrick.luois.ninja/" + bot.Token
		config := telegram.NewWebhook(url)
		_, err := bot.SetWebhook(config)
		if err != nil {
			logger.Fatal(err)
		}
		logger.Info("webhook set to " + url)

		updates = bot.ListenForWebhook("/" + bot.Token)
		go http.ListenAndServe(":"+port, nil)
	}

	for update := range updates {
		go onUpdate(update)
	}
}
