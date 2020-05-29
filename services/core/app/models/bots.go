package models

import (
	"gopkg.in/mgo.v2/bson"
)

type Bot interface {
	BotGetter
	BotSetter
}

type botData struct {
	Id       bson.ObjectId `bson:"_id"`
	Num      int           `bson:"number"`
	UserName string        `bson:"username"`
	Active   bool          `bson:"isActive"`
	Addr     Address       `bson:"address"`
	Bonus    Bonus         `bson:"bonus"`
	Times    Timestamp     `bson:"timestamp"`
}

func CreateBot(ip, port string) Bot {
	return &botData{
		Id:    bson.NewObjectId(),
		Addr:  CreateAddress(ip, port),
		Bonus: CreateBonus(),
		Times: CreateTimestamp(),
	}
}

type BotGetter interface {
	ID() bson.ObjectId
	Number() int
	Username() string
	IsActive() bool
	Address() Address
	Timestamp() Timestamp
}

func (b *botData) ID() bson.ObjectId {
	return b.Id
}

func (b *botData) Number() int {
	return b.Num
}

func (b *botData) Username() string {
	return "@" + b.UserName
}

func (b *botData) IsActive() bool {
	return b.Active
}

func (b *botData) Address() Address {
	return b.Addr
}

func (b *botData) Timestamp() Timestamp {
	return b.Times
}

type BotSetter interface {
	SetNumber(number int)
	SetUsername(username string)
	SetActive()
	InActive()
}

func (b *botData) SetNumber(number int) {
	b.Num = number
}

func (b *botData) SetUsername(username string) {
	b.UserName = username
}

func (b *botData) SetActive() {
	b.Active = true
}

func (b *botData) InActive() {
	b.Active = false
}
