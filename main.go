package main

import (
	"log"
	"os"

	"github.com/merrickluo/animeshotbot/bot"
	"github.com/merrickluo/animeshotbot/config"
	"github.com/merrickluo/animeshotbot/logging"
	"gopkg.in/alecthomas/kingpin.v2"
	mgo "gopkg.in/mgo.v2"
)

var (
	logger = logging.Logger()
	mode   = kingpin.
		Flag("mode", "start bot in \"fetch\" mode or \"webhook\" mode").
		Short('m').Default("fetch").String()
	port = kingpin.
		Flag("port", "listen port, use only with webhook mode").
		Short('p').Default("8185").String()
	debug = kingpin.
		Flag("debug", "enable debug mode").
		Short('d').Bool()
)

func main() {
	token := os.Getenv("ANIMESHOTBOT_TG_TOKEN")
	if len(token) == 0 {
		log.Fatal("Telegram bot token not set, " +
			"please set environment variable ANIMESHOTBOT_TG_TOKEN")
	}

	kingpin.Version("0.0.1")
	kingpin.Parse()

	//logging
	logging.Setup(os.Stdout, *debug)

	//config
	db, err := mgo.Dial("localhost")
	if err != nil {
		log.Fatal(err)
	}
	config.Setup(db.DB("animeshotbot").C("config"))

	if *mode != "fetch" && *mode != "webhook" {
		log.Fatal("mode must be \"fetch\" or \"webhook\"")
	}

	bot.Serve(*port, token, *mode, *debug)
}
