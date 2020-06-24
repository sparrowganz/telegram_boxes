package mongo

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
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
	FindByUsername(username string, session *mgo.Session) (models.Bot, error)
	FindByID(id bson.ObjectId, session *mgo.Session) (models.Bot, error)
	GetAll(session *mgo.Session) ([]*models.BotData, error)
}

func (bd *botsData) queryBot(session *mgo.Session) *mgo.Collection {
	return session.DB(bd.database).C(bd.collection)
}

func (bd *botsData) GetAll(session *mgo.Session) (bots []*models.BotData, err error) {
	err = bd.queryBot(session).Find(nil).All(&bots)
	return
}

func (bd *botsData) FindByID(id bson.ObjectId, session *mgo.Session) (models.Bot, error) {
	bot := &models.BotData{}
	err := bd.queryBot(session).FindId(id).One(&bot)
	return bot, err
}

func (bd *botsData) FindByUsername(username string, session *mgo.Session) (models.Bot, error) {
	bot := &models.BotData{}
	err := bd.queryBot(session).Find(bson.M{
		"username": username,
	}).One(&bot)
	return bot, err
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
