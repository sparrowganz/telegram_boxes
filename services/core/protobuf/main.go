package protobuf

import (
	"telegram_boxes/services/core/app/admin"
	"telegram_boxes/services/core/app/box"
	"telegram_boxes/services/core/app/db"
	"telegram_boxes/services/core/app/log"
)

type MainServer interface {
	Servers
	Tasks

	Admin() admin.Client
	DB() db.Client
	Log() log.Client
	Box() box.Clients
}

type serverData struct {
	database db.Client
	logger   log.Client
	admin    admin.Client
	boxes    box.Clients
}

func CreateServer(
	database db.Client,
	log log.Client, a admin.Client, c box.Clients) MainServer {

	return &serverData{
		database: database,
		logger:   log,
		admin:    a,
		boxes:    c,
	}
}


func (sd *serverData) Admin() admin.Client {
	return sd.admin
}

func (sd *serverData) Box() box.Clients {
	return sd.boxes
}

func (sd *serverData) DB() db.Client {
	return sd.database
}

func (sd *serverData) Log() log.Client {
	return sd.logger
}

