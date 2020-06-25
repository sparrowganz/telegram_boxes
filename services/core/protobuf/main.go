package protobuf

import (
	"telegram_boxes/services/core/app/admin"
	"telegram_boxes/services/core/app/db"
	"telegram_boxes/services/core/app/log"
)

type MainServer interface {
	Servers
	Tasks

	Admin() admin.Client
	DB() db.Client
	Log() log.Client
	LeadUpConnects()
}

type serverData struct {
	database db.Client
	logger   log.Client
	admin    admin.Client
}

func CreateServer(
	database db.Client,
	log log.Client, a admin.Client) MainServer {

	return &serverData{
		database: database,
		logger:   log,
		admin:    a,
	}
}

func (sd *serverData) LeadUpConnects() {
	session := sd.DB().GetMainSession().Clone()
	defer session.Close()

	bots, err := sd.DB().Models().Bots().GetAll(session)
	if err != nil {
		return
	}
	for _, b := range bots {
		//todo check is up bot
		_ = b
	}
}

func (sd *serverData) Admin() admin.Client {
	return sd.admin
}

func (sd *serverData) DB() db.Client {
	return sd.database
}

func (sd *serverData) Log() log.Client {
	return sd.logger
}

