package mongo

import (
	"gopkg.in/mgo.v2"
	"telegram_boxes/services/box/app/models"
)

type usersData struct {
	database   string
	collection string
}

func createUsersModel(database string) Users {
	return &usersData{
		database:   database,
		collection: "Users",
	}
}

type Users interface {
	queryUsers(session *mgo.Session) *mgo.Collection
	CreateUser(user models.User, session *mgo.Session) error
	UpdateUser(user models.User, session *mgo.Session) error
	RemoveUser(user models.User, session *mgo.Session) error
}

func (users *usersData) queryUsers(session *mgo.Session) *mgo.Collection {
	return session.DB(users.database).C(users.collection)
}

func (users *usersData) CreateUser(user models.User, session *mgo.Session) error {
	user.Timestamp().SetCreateTime()
	return users.queryUsers(session).Insert(user)
}

func (users *usersData) UpdateUser(user models.User, session *mgo.Session) error {
	user.Timestamp().SetUpdateTime()
	return users.queryUsers(session).UpdateId(user.ID(), user)
}

func (users *usersData) RemoveUser(user models.User, session *mgo.Session) error {
	user.Timestamp().SetRemoveTime()
	return users.queryUsers(session).UpdateId(user.ID(), user)
}