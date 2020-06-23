package mongo

import "gopkg.in/mgo.v2"

type payData struct {
	database   string
	collection string
}

func createPaymentsModel(database string) Payments {
	return &payData{
		database:   database,
		collection: "Payments",
	}
}

type Payments interface {
	queryPayments(session *mgo.Session) *mgo.Collection
	//CreatePayment(bot models.Payment, session *mgo.Session) error
	//UpdatePayment(bot models.Payment, session *mgo.Session) error
	//RemovePayment(bot models.Payment, session *mgo.Session) error
}

func (pd *payData) queryPayments(session *mgo.Session) *mgo.Collection {
	return session.DB(pd.database).C(pd.collection)
}

/*func (db *DB) CreatePayment(pay models.Payment, session *mgo.Session) error {
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
*/
