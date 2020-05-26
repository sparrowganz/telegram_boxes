package models

import "gopkg.in/mgo.v2/bson"


//todo REFACTOR payments model
var (
	ValidateStatus Status = "validate"
	OKStatus       Status = "ok"
)

type Status string

func (s Status) String() string {
	return string(s)
}

type Payment interface {
	PaymentGetter
	PaymentSetter
}

type paymentData struct {
	Id bson.ObjectId `bson:"_id"`

	UserData    User `bson:"user"`
	PaymentData Data `bson:"data"`

	Value float64 `bson:"cost"`

	State Status    `bson:"status"`
	Time  Timestamp `bson:"timestamp"`
}

func CreatePayment(u User, d Data) Payment {
	return &paymentData{
		Id:          bson.NewObjectId(),
		UserData:    u,
		PaymentData: d,
		Time:        CreateTimestamp(),
	}
}

type PaymentGetter interface {
	ID() bson.ObjectId
	User() User
	Data() Data
	Cost() float64
	Status() Status
	Timestamp() Timestamp
}

func (p *paymentData) ID() bson.ObjectId {
	return p.Id
}

func (p *paymentData) User() User {
	return p.UserData
}
func (p *paymentData) Data() Data {
	return p.PaymentData
}

func (p *paymentData) Cost() float64 {
	return p.Value
}

func (p *paymentData) Status() Status {
	return p.State
}
func (p *paymentData) Timestamp() Timestamp {
	return p.Time
}

type PaymentSetter interface {
	SetCost(cost float64)
}

func (p *paymentData) SetCost(cost float64) {
	p.Value = cost
}
