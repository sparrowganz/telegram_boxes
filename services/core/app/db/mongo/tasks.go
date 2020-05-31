package mongo

type tasksData struct {
	database   string
	collection string
}

func createTaskModel(database string) Tasks {
	return &tasksData{
		database:   database,
		collection: "Bots",
	}
}

type Tasks interface {
	queryTasks(session *mgo.Session) *mgo.Collection
}

func (td *tasksData) queryTasks(session *mgo.Session) *mgo.Collection {
	return session.DB(td.database).C(td.collection)
}
