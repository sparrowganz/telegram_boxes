package models

import "gopkg.in/mgo.v2/bson"

type Output interface {
	OutputGetter
}

func CreateOutput(userID, pg, data string, cost int, tasks []string) Output {
	return &OutputData{
		Id:             bson.NewObjectId(),
		UserId:         userID,
		CostData:       cost,
		PaymentGateway: pg,
		PaymentData:    data,
		ChecksTasks:    tasks,
		TimestampData:  CreateTimestamp(),
	}
}

type OutputData struct {
	Id             bson.ObjectId  `bson:"_id"`
	UserId         string         `bson:"userID"`
	CostData       int            `bson:"cost"`
	PaymentGateway string         `bson:"gateway"`
	PaymentData    string         `bson:"data"`
	ChecksTasks    []string       `bson:"tasks"`
	TimestampData  *TimestampData `bson:"timestamp"`
}

type OutputGetter interface {
	ID() bson.ObjectId
	UserID() string
	Cost() int
	PaymentGW() string
	Data() string
	Tasks() []string
	Timestamp() Timestamp
}

func (o *OutputData) ID() bson.ObjectId {
	return o.Id
}

func (o *OutputData) UserID() string {
	return o.UserId
}

func (o *OutputData) Cost() int {
	return o.CostData
}

func (o *OutputData) PaymentGW() string {
	return o.PaymentGateway
}

func (o *OutputData) Data() string {
	return o.PaymentData
}

func (o *OutputData) Tasks() []string {
	return o.ChecksTasks
}

func (o *OutputData) Timestamp() Timestamp {
	return o.TimestampData
}
