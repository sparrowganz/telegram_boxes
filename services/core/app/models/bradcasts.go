package models

import "gopkg.in/mgo.v2/bson"

type Broadcast interface {
	BroadcastGetter
	BroadcastSetter
}

type broadcastData struct {
	Id      bson.ObjectId `bson:"_id"`
	BotId   string        `bson:"botID"`
	AdminId string        `bson:"adminID"`
	Stats   Stats         `bson:"stats"`
	Info    Info          `bson:"Info"`
	Time    Timestamp     `bson:"timestamp"`
}

func CreateBroadcast() Broadcast {
	return &broadcastData{
		Id:    bson.NewObjectId(),
		Stats: CreateStats(),
		Info:  CreateInfo(),
		Time:  CreateTimestamp(),
	}
}

type BroadcastGetter interface {
	ID() bson.ObjectId
	BotID() string
	AdminID() string
	Statistics() Stats
	Information() Info
	Timestamp() Timestamp
}

func (br *broadcastData) ID() bson.ObjectId {
	return br.Id
}

func (br *broadcastData) BotID() string {
	return br.BotId
}

func (br *broadcastData) AdminID() string {
	return br.AdminId
}

func (br *broadcastData) Statistics() Stats {
	return br.Stats
}

func (br *broadcastData) Information() Info {
	return br.Info
}

func (br *broadcastData) Timestamp() Timestamp {
	return br.Time
}

type BroadcastSetter interface {
	SetBotID(id string)
	SetAdminID(id string)
}

func (br *broadcastData) SetBotID(id string) {
	br.BotId = id
}

func (br *broadcastData) SetAdminID(id string) {
	br.AdminId = id
}
