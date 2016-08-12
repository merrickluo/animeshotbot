package main

import (
	"github.com/op/go-logging"
	"gopkg.in/alecthomas/kingpin.v2"
	telegram "gopkg.in/telegram-bot-api.v4"
	"net/http"
	"os"
	"strconv"
)

var (
	log    = logging.MustGetLogger("animeshotbot")
	format = logging.MustStringFormatter("%{color}%{time:15:04:05.000}" +
		"%{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}")
	mode = kingpin.
		Flag("mode", "start bot in \"fetch\" mode or \"webhook\" mode").
		Short('m').Default("fetch").String()
	port = kingpin.
		Flag("port", "listen port, use only with webhook mode").
		Short('p').Default("8185").Int()
	debug = kingpin.
		Flag("debug", "enable debug mode").
		Short('d').Bool()
)

const fetchMode = "fetch"
const webhookMode = "webhook"

func main() {
	var token = os.Getenv("ANIMESHOTBOT_TG_TOKEN")
	if len(token) == 0 {
		log.Fatal("Telegram bot token not set, " +
			"please set environment variable ANIMESHOTBOT_TG_TOKEN")
	}

	kingpin.Version("0.0.1")
	kingpin.Parse()

	logBackend := logging.NewLogBackend(os.Stdout, "", 0)
	if *debug {
		logging.SetBackend(logBackend)
	} else {
		logBackendLeveled := logging.AddModuleLevel(logBackend)
		logBackendLeveled.SetLevel(logging.ERROR, "")
		logging.SetBackend(logBackendLeveled)
	}

	bot, err := telegram.NewBotAPI(token)
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = *debug
	log.Info("Authorized on account " + bot.Self.UserName)

	if *mode == fetchMode {
		config := telegram.NewWebhook("")
		bot.SetWebhook(config)

		u := telegram.NewUpdate(0)
		u.Timeout = 60
		updates, err := bot.GetUpdatesChan(u)
		if err != nil {
			log.Fatal(err)
		}

		for update := range updates {
			go processUpdate(*bot, update)
		}
	} else if *mode == webhookMode {
		url := "https://merrick.luois.ninja/" + bot.Token
		config := telegram.NewWebhook(url)
		_, err := bot.SetWebhook(config)
		if err != nil {
			log.Fatal(err)
		}
		log.Info("webhook set to " + url)

		updates := bot.ListenForWebhook("/" + bot.Token)
		go http.ListenAndServe(":"+strconv.Itoa(*port), nil)

		for update := range updates {
			log.Info(update)
			go processUpdate(*bot, update)
		}
	} else {
		log.Fatal("mode must be \"fetch\" or \"webhook\"")
	}
}
