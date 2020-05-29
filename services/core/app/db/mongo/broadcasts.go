package mongo

import (
	"gopkg.in/mgo.v2"
	"telegram_boxes/services/core/app/models"
)

type broadcastData struct {
	database   string
	collection string
}

func createBroadcastModel(database string) Broadcasts {
	return &broadcastData{
		database:   database,
		collection: "Bots",
	}
}

type Broadcasts interface {
	queryBroadcast(session *mgo.Session) *mgo.Collection
	CreateBroadcast(br models.Broadcast, session *mgo.Session) error
	UpdateBroadcast(br models.Broadcast, session *mgo.Session) error
	RemoveBroadcast(br models.Broadcast, session *mgo.Session) error
}

func (br *broadcastData) queryBroadcast(session *mgo.Session) *mgo.Collection {
	return session.DB(br.database).C(br.collection)
}

func (br *broadcastData) CreateBroadcast(model models.Broadcast, session *mgo.Session) error {
	model.Timestamp().SetCreateTime()
	return br.queryBroadcast(session).Insert(model)
}

func (br *broadcastData) UpdateBroadcast(model models.Broadcast, session *mgo.Session) error {
	model.Timestamp().SetUpdateTime()
	return br.queryBroadcast(session).UpdateId(model.ID(), model)
}

func (br *broadcastData) RemoveBroadcast(model models.Broadcast, session *mgo.Session) error {
	model.Timestamp().SetRemoveTime()
	return br.queryBroadcast(session).UpdateId(model.ID(), model)
}
