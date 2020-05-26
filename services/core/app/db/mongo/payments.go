package mongo

import (
	"gopkg.in/mgo.v2"
	"telegram_boxes/services/core/app/models"
)

type Payments interface {
	queryPayments(session *mgo.Session) *mgo.Collection
	CreatePayment(bot models.Payment, session *mgo.Session) error
	UpdatePayment(bot models.Payment, session *mgo.Session) error
	RemovePayment(bot models.Payment, session *mgo.Session) error
}

func (db *DB) queryPayments(session *mgo.Session) *mgo.Collection {
	return session.DB(db.DatabaseName).C("Payments")
}

func (db *DB) CreatePayment(pay models.Payment, session *mgo.Session) error {
	pay.Timestamp().SetCreateTime()
	return db.queryPayments(session).Insert(pay)
}

func (db *DB) UpdatePayment(pay models.Payment, session *mgo.Session) error {
	pay.Timestamp().SetUpdateTime()
	return db.queryPayments(session).UpdateId(pay.ID(), pay)
}

func (db *DB) RemovePayment(pay models.Payment, session *mgo.Session) error {
	pay.Timestamp().SetRemoveTime()
	return db.queryBot(session).UpdateId(pay.ID(), pay)
}
