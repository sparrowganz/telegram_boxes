package models

import (
	"gopkg.in/mgo.v2/bson"
	"sync"
)

type User interface {
	UserGetter
	UserSetter
	UserCheckData
}

func CreateUser(telegramID int64, username, firstName, lastName, email string) User {
	return &userData{
		Id:             bson.NewObjectId(),
		Account:        CreateAccount(telegramID, username, firstName, lastName, email),
		Balances:       CreateBalance(),
		ChecksData:     []string{},
		Time:           CreateTimestamp(),
		checkDataMutex: &sync.Mutex{},
	}
}

type userData struct {
	Id        bson.ObjectId `bson:"ID"`
	IsBlocked bool          `bson:"isBlocked"`

	InviterId string  `bson:"inviterID"`
	Balances  Balance `bson:"balance"`
	Account   Account `bson:"account"`

	ChecksData []string  `bson:"checks"`
	Time       Timestamp `bson:"timestamp"`

	checkDataMutex *sync.Mutex
}

type UserGetter interface {
	ID() string
	Blocked() bool
	InviterID() string
	Balance() Balance
	Telegram() Account
	Timestamp() Timestamp
}

func (u *userData) ID() string {
	return u.Id.Hex()
}

func (u *userData) Blocked() bool {
	return u.IsBlocked
}
func (u *userData) InviterID() string {
	return u.InviterId
}
func (u *userData) Balance() Balance {
	return u.Balances
}
func (u *userData) Telegram() Account {
	return u.Account
}

func (u *userData) Timestamp() Timestamp {
	return u.Time
}

type UserSetter interface {
	Blocker
	SetInviterID(id string)
}

func (u *userData) SetInviterID(id string) {
	u.InviterId = id
}

type Blocker interface {
	Block()
	Unblock()
}

func (u *userData) Block() {
	u.IsBlocked = true
}

func (u *userData) Unblock() {
	u.IsBlocked = false
}
