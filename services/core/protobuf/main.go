package protobuf

import (
	"telegram_boxes/services/core/app/db"
	"telegram_boxes/services/core/app/log"
)

type MainServer interface {
	Servers
	Tasks

	DB() db.Database
	Log() log.Log
}

type serverData struct {
	database db.Database
	logger   log.Log
}

func CreateServer(
	database db.Database,
	log log.Log) MainServer {

	return &serverData{
		database: database,
		logger:   log,
	}
}


func (sd *serverData) DB() db.Database {
	return sd.database
}

func (sd *serverData) Log() log.Log {
	return sd.logger
}

