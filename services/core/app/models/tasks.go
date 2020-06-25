package models

import "gopkg.in/mgo.v2/bson"

type Task interface {
	TaskGetter
	TaskSetter
}

type TaskData struct {
	Id            bson.ObjectId  `bson:"_id"`
	Name          string         `bson:"name"`
	Priority      bool           `bson:"isPriority"`
	Check         bool           `bson:"withCheck"`
	TypeName      string         `bson:"type"`
	Link          string         `bson:"link"`
	TimestampData *TimestampData `bson:"timestamp"`
}

func CreateTask(name, typeName, link string, isPriority, withCheck bool) *TaskData {
	return &TaskData{
		Id:            bson.NewObjectId(),
		Name:          name,
		TypeName:      typeName,
		Link:          link,
		Priority:      isPriority,
		Check:         withCheck,
		TimestampData: CreateTimestamp(),
	}
}

type TaskGetter interface {
	ID() bson.ObjectId
	Title() string
	Type() string
	URL() string
	IsPriority() bool
	WithCheck() bool
	Timestamp() Timestamp
}

func (t *TaskData) ID() bson.ObjectId {
	return t.Id
}

func (t *TaskData) Title() string {
	return t.Name
}

func (t *TaskData) Type() string {
	return t.TypeName
}

func (t *TaskData) URL() string {
	return t.Link
}

func (t *TaskData) IsPriority() bool {
	return t.Priority
}

func (t *TaskData) WithCheck() bool {
	return t.Check
}

func (t *TaskData) Timestamp() Timestamp {
	return t.TimestampData
}
 type TaskSetter interface {
 	ChangePriority()
 }

 func (t *TaskData) ChangePriority() {
 	t.Priority = !t.Priority
 }