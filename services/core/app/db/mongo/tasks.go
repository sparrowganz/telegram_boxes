package mongo

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"telegram_boxes/services/core/app/models"
)

type tasksData struct {
	database   string
	collection string
}

func createTaskModel(database string) Tasks {
	return &tasksData{
		database:   database,
		collection: "Tasks",
	}
}

type Tasks interface {
	queryTasks(session *mgo.Session) *mgo.Collection
	CreateTask(tsk models.Task, session *mgo.Session) error
	UpdateTask(tsk models.Task, session *mgo.Session) error
	RemoveTask(tsk models.Task, session *mgo.Session) error
	GetNextTask(currentTaskID []bson.ObjectId, session *mgo.Session) (models.Task, error)
	FindTask(taskID bson.ObjectId, session *mgo.Session) (models.Task, error)
}

func (td *tasksData) queryTasks(session *mgo.Session) *mgo.Collection {
	return session.DB(td.database).C(td.collection)
}

func (td *tasksData) FindTask(taskID bson.ObjectId, session *mgo.Session) (models.Task, error) {
	tsk := &models.TaskData{}
	err := td.queryTasks(session).FindId(taskID).One(&tsk)
	return tsk, err
}

func (td *tasksData) GetNextTask(currentTaskID []bson.ObjectId, session *mgo.Session) (models.Task, error) {
	tsk := &models.TaskData{}
	err := td.queryTasks(session).Find(bson.M{
		"_id": bson.M{
			"$ne": currentTaskID,
		},
	}).Sort("isPriority").One(&tsk)
	return tsk, err
}

func (td *tasksData) CreateTask(tsk models.Task, session *mgo.Session) error {
	tsk.Timestamp().SetCreateTime()
	return td.queryTasks(session).Insert(tsk)
}

func (td *tasksData) UpdateTask(tsk models.Task, session *mgo.Session) error {
	tsk.Timestamp().SetUpdateTime()
	return td.queryTasks(session).UpdateId(tsk.ID(), tsk)
}

func (td *tasksData) RemoveTask(tsk models.Task, session *mgo.Session) error {
	tsk.Timestamp().SetRemoveTime()
	return td.queryTasks(session).UpdateId(tsk.ID(), tsk)
}
