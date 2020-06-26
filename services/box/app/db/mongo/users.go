package mongo

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
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
	FindUserByTelegramID(id int64, session *mgo.Session) (user models.User, err error)
	GetCountInvitedUsers(id string, session *mgo.Session) int
	GetAllCount(session *mgo.Session) int
	GetBlockedCount(session *mgo.Session) int
	RemoveCheck(taskID string, session *mgo.Session) error
}

func (users *usersData) queryUsers(session *mgo.Session) *mgo.Collection {
	return session.DB(users.database).C(users.collection)
}

func (users *usersData) RemoveCheck(taskID string, session *mgo.Session) error {
	_, err := users.queryUsers(session).UpdateAll(
		bson.M{
			"checks." + taskID: bson.M{"$exists": true},
		}, bson.M{
			"$unset": bson.M{"checks." + taskID: ""},
		})
	if err != nil {
		return err
	}
	return nil
}

func (users *usersData) GetAllCount(session *mgo.Session) int {
	count, err := users.queryUsers(session).Find(nil).Count()
	if err != nil {
		return 0
	}
	return count
}

func (users *usersData) GetBlockedCount(session *mgo.Session) int {
	count, err := users.queryUsers(session).Find(bson.M{
		"isBlocked": true,
	}).Count()
	if err != nil {
		return 0
	}
	return count
}

func (users *usersData) FindUserByTelegramID(id int64, session *mgo.Session) (models.User, error) {
	user := models.UserData{}
	err := users.queryUsers(session).Find(bson.M{
		"account.id": id,
	}).One(&user)
	return &user, err
}

func (users *usersData) GetCountInvitedUsers(id string, session *mgo.Session) int {
	count, _ := users.queryUsers(session).Find(bson.M{
		"inviterID": id,
	}).Count()
	return count
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
