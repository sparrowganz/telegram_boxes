package models

import "gopkg.in/mgo.v2/bson"

//todo refactor TASKS model
type Task struct {
	Id      bson.ObjectId `bson:"_id"`
	Creator string        `bson:"creator"`
	Cost    float64       `bson:"cost"`

	Title string      `bson:"title"`
	Type  Type        `bson:"type"`
	Data  interface{} `bson:"data"`

	Condition string `bson:"condition"`
	UrlLabel  string `bson:"url_label"`

	IsActive  bool      `bson:"is_active"`
	ShowAdmin bool      `bson:"show_admin"`
	Timestamp Timestamp `bson:"timestamp"`
}

type Type struct {
	Smile    string `bson:"smile"`
	Platform string `json:"platform"`
}

type ChannelData struct {
	Id        int64  `bson:"id"`
	Username  string `bson:"username"`
	MessageId int    `bson:"message_id"`
}

type TaskData struct {
	Url string `bson:"url"`
}
