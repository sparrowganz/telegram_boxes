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
	return &UserData{
		Id:             bson.NewObjectId(),
		Account:        CreateAccount(telegramID, username, firstName, lastName, email),
		Balances:       CreateBalance(),
		ChecksData:     []string{},
		Time:           CreateTimestamp(),
		checkDataMutex: &sync.Mutex{},
	}
}

type UserData struct {
	Id        bson.ObjectId `bson:"_id"`
	IsBlocked bool          `bson:"isBlocked"`

	InviterId string `bson:"inviterID"`

	IsVerified bool         `bson:"isVerified"`
	Account    *AccountData `bson:"account"`
	Balances   *BalanceData `bson:"balance"`

	ChecksData []string       `bson:"checks"`
	Time       *TimestampData `bson:"timestamp"`

	checkDataMutex *sync.Mutex
}

type UserGetter interface {
	ID() string
	Blocked() bool
	InviterID() string
	Balance() Balance
	Telegram() Account
	Timestamp() Timestamp
	Verified() bool
}

func (u *UserData) ID() string {
	return u.Id.Hex()
}

func (u *UserData) Blocked() bool {
	return u.IsBlocked
}
func (u *UserData) InviterID() string {
	return u.InviterId
}
func (u *UserData) Balance() Balance {
	return u.Balances
}
func (u *UserData) Telegram() Account {
	return u.Account
}

func (u *UserData) Verified() bool {
	return u.IsVerified
}

func (u *UserData) Timestamp() Timestamp {
	return u.Time
}

type UserSetter interface {
	Blocker
	SetInviterID(id string)
	SetVerified()
}

func (u *UserData) SetInviterID(id string) {
	u.InviterId = id
}

func (u *UserData) SetVerified() {
	u.IsVerified = true
}

type Blocker interface {
	Block()
	Unblock()
}

func (u *UserData) Block() {
	u.IsBlocked = true
}

func (u *UserData) Unblock() {
	u.IsBlocked = false
}
