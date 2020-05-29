package models

type User interface {
	UserGetter
	UserSetter
}

func CreateUser(id, botID string, telegramID int64) User {
	return &userData{
		Id:         id,
		TelegramId: telegramID,
		BotId:      botID,
	}
}

type userData struct {
	Id         string `bson:"ID"`
	TelegramId int64  `bson:"telegramID"`
	BotId      string `bson:"botID"`
}

type UserGetter interface {
	ID() string
	TelegramID() int64
	BotID() string
}

func (u *userData) ID() string {
	return u.Id
}

func (u *userData) TelegramID() int64 {
	return u.TelegramId
}

func (u *userData) BotID() string {
	return u.BotId
}

type UserSetter interface {
	SetID(id string)
	SetTelegramID(id int64)
	SetBotID(id string)
}

func (u *userData) SetID(id string) {
	u.Id = id
}

func (u *userData) SetTelegramID(id int64) {
	u.TelegramId = id
}

func (u *userData) SetBotID(id string) {
	u.BotId = id
}
