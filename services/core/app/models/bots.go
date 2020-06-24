package models

import (
	"gopkg.in/mgo.v2/bson"
)

type Bot interface {
	BotGetter
	BotSetter
}

type BotData struct {
	Id bson.ObjectId `bson:"_id"`

	UserName string `bson:"username"`
	//Num      int           `bson:"number"`
	Status string `bson:"status"`
	Active bool   `bson:"isActive"`

	StatisticsData *StatisticsData `bson:"statistics"`
	Addr           *AddressData    `bson:"address"`
	BonusData      *BonusData      `bson:"bonus"`
	Times          *TimestampData  `bson:"timestamp"`
}

func CreateBot(ip, port string) Bot {
	return &BotData{
		Id:             bson.NewObjectId(),
		Addr:           CreateAddress(ip, port),
		StatisticsData: CreateStatistics(),
		BonusData:      CreateBonus(),
		Times:          CreateTimestamp(),
	}
}

type BotGetter interface {
	ID() bson.ObjectId
	Object() *BotData
	//Number() int
	BotStatus() string
	Username() string
	IsActive() bool
	Statistics() Statistics
	Bonus() Bonus
	Address() Address
	Timestamp() Timestamp
}

func (b *BotData) Object() *BotData {
	return b
}

func (b *BotData) Statistics() Statistics {
	return b.StatisticsData
}

func (b *BotData) Bonus() Bonus {
	return b.BonusData
}

func (b *BotData) BotStatus() string {
	return b.Status
}

func (b *BotData) ID() bson.ObjectId {
	return b.Id
}

/*func (b *BotData) Number() int {
	return b.Num
}*/

func (b *BotData) Username() string {
	return "@" + b.UserName
}

func (b *BotData) IsActive() bool {
	return b.Active
}

func (b *BotData) Address() Address {
	return b.Addr
}

func (b *BotData) Timestamp() Timestamp {
	return b.Times
}

type BotSetter interface {
	//SetNumber(number int)
	SetUsername(username string)
	SetStatus(status string)
	SetBonus(data *BonusData)
	SetActive()
	InActive()
}

func (b *BotData) SetBonus(data *BonusData) {
	b.BonusData = data
}

func (b *BotData) SetStatus(status string) {
	b.Status = status
}

/*func (b *BotData) SetNumber(number int) {
	b.Num = number
}*/

func (b *BotData) SetUsername(username string) {
	b.UserName = username
}

func (b *BotData) SetActive() {
	b.Active = true
}

func (b *BotData) InActive() {
	b.Active = false
}
