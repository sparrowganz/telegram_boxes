package mongo

import (
	"gopkg.in/mgo.v2"
	"telegram_boxes/services/core/app/models"
)

type Broadcasts interface {
	queryBroadcast(session *mgo.Session) *mgo.Collection
	CreateBroadcast(br models.Broadcast, session *mgo.Session) error
	UpdateBroadcast(br models.Broadcast, session *mgo.Session) error
	RemoveBroadcast(br models.Broadcast, session *mgo.Session) error
}

func (db *DB) queryBroadcast(session *mgo.Session) *mgo.Collection {
	return session.DB(db.DatabaseName).C("Broadcasts")
}

func (db *DB) CreateBroadcast(br models.Broadcast, session *mgo.Session) error {
	br.Timestamp().SetCreateTime()
	return db.queryBroadcast(session).Insert(br)
}

func (db *DB) UpdateBroadcast(br models.Broadcast, session *mgo.Session) error {
	br.Timestamp().SetUpdateTime()
	return db.queryBot(session).UpdateId(br.ID(), br)
}

func (db *DB) RemoveBroadcast(br models.Broadcast, session *mgo.Session) error {
	br.Timestamp().SetRemoveTime()
	return db.queryBot(session).UpdateId(br.ID(), br)
}
