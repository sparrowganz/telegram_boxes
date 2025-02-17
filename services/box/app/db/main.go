package db

import (
	"gopkg.in/mgo.v2"
	"telegram_boxes/services/box/app/db/mongo"
)

type Database interface {
	Connector
	Getter
}

func InitDatabaseConnect(host, port, username, password, database, mechanism string) (Database, error) {
	return mongo.CreateMongoDB(host, port, username, password, database, mechanism)
}

type Getter interface {
	GetMainSession() *mgo.Session
	GetDatabaseName() string
	Models() mongo.Models
}

type Connector interface {
	Connect(host, port, username, password, database, mechanism string) error
	Close()
}
