package protobuf

import (
	"telegram_boxes/services/core/app/db"
	"telegram_boxes/services/core/app/id"
	"telegram_boxes/services/core/app/log"
)

type Server interface {
	ServerGetter
}

type serverData struct {
	database db.Database
	logger   log.Log
	counter  id.Counter
}

func CreateServer(
	database db.Database,
	log log.Log,
	counter id.Counter) Server {

	return &serverData{
		database: database,
		logger:   log,
		counter:  counter,
	}
}

type ServerGetter interface {
	DB() db.Database
	Log() log.Log
	Counter() id.Counter
}

func (sd *serverData) DB() db.Database {
	return sd.database
}

func (sd *serverData) Log() log.Log {
	return sd.logger
}

func (sd *serverData) Counter() id.Counter {
	return sd.counter
}
