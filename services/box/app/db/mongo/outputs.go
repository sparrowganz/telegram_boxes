package mongo

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"telegram_boxes/services/box/app/models"
)

type outputData struct {
	database   string
	collection string
}

func createOutputsModel(database string) Outputs {
	return &outputData{
		database:   database,
		collection: "Outputs",
	}
}

type Outputs interface {
	queryOutputs(session *mgo.Session) *mgo.Collection
	CreateOutput(out models.Output, session *mgo.Session) error
	UpdateOutput(out models.Output, session *mgo.Session) error
	RemoveOutput(out models.Output, session *mgo.Session) error
	FindOutputByUserID(userID string, session *mgo.Session) (models.Output, error)
}

func (output *outputData) queryOutputs(session *mgo.Session) *mgo.Collection {
	return session.DB(output.database).C(output.collection)
}

func (output *outputData) FindOutputByUserID(userID string, session *mgo.Session) (models.Output, error) {
	out := models.OutputData{}
	err := output.queryOutputs(session).Find(bson.M{
		"userID": userID,
	}).One(&out)
	return &out, err
}

func (output *outputData) CreateOutput(out models.Output, session *mgo.Session) error {
	out.Timestamp().SetCreateTime()
	return output.queryOutputs(session).Insert(out)
}

func (output *outputData) UpdateOutput(out models.Output, session *mgo.Session) error {
	out.Timestamp().SetUpdateTime()
	return output.queryOutputs(session).UpdateId(out.ID(), out)
}

func (output *outputData) RemoveOutput(out models.Output, session *mgo.Session) error {
	return output.queryOutputs(session).RemoveId(out.ID())
}
