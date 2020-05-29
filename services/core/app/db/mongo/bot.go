package mongo

import (
	"gopkg.in/mgo.v2"
	"telegram_boxes/services/core/app/models"
)

type botsData struct {
	database   string
	collection string
}

func createBotModel(database string) Bots {
	return &botsData{
		database:   database,
		collection: "Bots",
	}
}

type Bots interface {
	queryBot(session *mgo.Session) *mgo.Collection
	CreateBot(bot models.Bot, session *mgo.Session) error
	UpdateBot(bot models.Bot, session *mgo.Session) error
	RemoveBot(bot models.Bot, session *mgo.Session) error
}

func (bd *botsData) queryBot(session *mgo.Session) *mgo.Collection {
	return session.DB(bd.database).C(bd.collection)
}

func (bd *botsData) CreateBot(bot models.Bot, session *mgo.Session) error {
	bot.Timestamp().SetCreateTime()
	return bd.queryBot(session).Insert(bot)
}

func (bd *botsData) UpdateBot(bot models.Bot, session *mgo.Session) error {
	bot.Timestamp().SetUpdateTime()
	return bd.queryBot(session).UpdateId(bot.ID(), bot)
}

func (bd *botsData) RemoveBot(bot models.Bot, session *mgo.Session) error {
	bot.Timestamp().SetRemoveTime()
	return bd.queryBot(session).UpdateId(bot.ID(), bot)
}
