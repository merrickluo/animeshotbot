package config

import (
	"github.com/merrickluo/animeshotbot/logging"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	log "qiniupkg.com/x/log.v7"
)

//Config as is
type Config struct {
	ID      bson.ObjectId `bson:"_id"`
	UserID  int           `bson:"user_id"`
	Mode    int           `bson:"mode"`
	Paging  int           `bson:"paging"`
	Command string        `bson:"command"`
}

var (
	logger = logging.Logger()
	db     *mgo.Collection
)

const (
	//FullMode as is
	FullMode = 1
	//CaptionMode as is
	CaptionMode = 2

	//PagingCommand as is
	PagingCommand = "setpaging"
	//ModeCommand as is
	ModeCommand = "setmode"
	//SettingsCommand as is
	SettingsCommand = "settings"
)

//Setup setup config db
func Setup(collection *mgo.Collection) {
	db = collection
}

//OnCommand put user in command mode
func (c *Config) OnCommand(command string) {
	c.Command = command
	c.Save()
}

//IsOnCommand is user in command mode
func (c *Config) IsOnCommand() bool {
	return c.Command != ""
}

//ExitCommand put user to normal mode
func (c *Config) ExitCommand() {
	c.Command = ""
	c.Save()
}

//Save save current config
func (c *Config) Save() {
	err := db.UpdateId(c.ID, c)
	if err != nil {
		log.Error(err)
	}
}

//GetConfig get config for userID
func GetConfig(userID int) *Config {
	c := &Config{}
	n, err := db.Find(bson.M{"user_id": userID}).Count()
	if err != nil {
		logger.Error(err)
	}

	if n == 0 {
		err := db.Insert(bson.M{
			"user_id": userID,
			"mode":    CaptionMode,
			"paging":  5,
			"command": "",
		})
		if err != nil {
			logger.Error(err)
		}
		logger.Info(c)
	}

	err = db.Find(bson.M{"user_id": userID}).One(&c)
	if err != nil {
		logger.Error(err)
	}

	return c
}
