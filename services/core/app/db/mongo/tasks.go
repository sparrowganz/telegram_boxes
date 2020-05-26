package mongo

import "gopkg.in/mgo.v2"

type Tasks interface {
	queryTasks(session *mgo.Session) *mgo.Collection
}
func (db *DB) queryTasks(session *mgo.Session) *mgo.Collection {
	return session.DB(db.DatabaseName).C("Tasks")
}