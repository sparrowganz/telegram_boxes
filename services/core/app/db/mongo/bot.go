package mongo

import (
	"gopkg.in/mgo.v2"
	"telegram_boxes/services/core/app/models"
)

type Bots interface {
	queryBot(session *mgo.Session) *mgo.Collection
	CreateBot(bot models.Bot, session *mgo.Session) error
	UpdateBot(bot models.Bot, session *mgo.Session) error
	RemoveBot(bot models.Bot, session *mgo.Session) error
}

func (db *DB) queryBot(session *mgo.Session) *mgo.Collection {
	return session.DB(db.DatabaseName).C("Bots")
}

func (db *DB) CreateBot(bot models.Bot, session *mgo.Session) error {
	bot.Timestamp().SetCreateTime()
	return db.queryBot(session).Insert(bot)
}

func (db *DB) UpdateBot(bot models.Bot, session *mgo.Session) error {
	bot.Timestamp().SetUpdateTime()
	return db.queryBot(session).UpdateId(bot.ID(), bot)
}

func (db *DB) RemoveBot(bot models.Bot, session *mgo.Session) error {
	bot.Timestamp().SetRemoveTime()
	return db.queryBot(session).UpdateId(bot.ID(), bot)
}
