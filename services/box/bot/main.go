package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sparrowganz/teleFly/telegram"
	"telegram_boxes/services/box/app/config"
	"telegram_boxes/services/box/app/db"
	"telegram_boxes/services/box/app/log"
	"telegram_boxes/services/box/app/output"
	"telegram_boxes/services/box/app/servers"
	"telegram_boxes/services/box/app/task"
)

type Bot interface {
	StartHandle()
	StartReadErrors()
	Methods() Botter
}

func (b *botData) Methods() Botter {
	return b
}

type BotSetter interface {
	SetTasks(t task.Tasks)
	SetServers(s servers.Servers)
	SetOutput(o output.Outputs)
	SetConfig(c config.Config)
}

func (b *botData) SetTasks(t task.Tasks) {
	b.tasks = t
}


func (b *botData) SetServers(s servers.Servers) {
	b.servers = s
}

func (b *botData) SetConfig(c config.Config) {
	b.config = c
}

func (b *botData) SetOutput(o output.Outputs) {
	b.outputs = o
}

type Botter interface {
	BotSetter
	Telegram() telegram.Sender
	Log() log.Log
	Task() task.Tasks
	Output() output.Outputs
	Servers() servers.Servers
	Database() db.Database
	Config() config.Config
	Username() string
}

type botData struct {
	mongo    db.Database
	tSender  telegram.Sender
	logger   log.Log
	tasks    task.Tasks
	servers  servers.Servers
	config   config.Config
	outputs  output.Outputs
	username string
}

func (b *botData) Username() string {
	return b.username
}

func (b *botData) Config() config.Config {
	return b.config
}

func (b *botData) Output() output.Outputs {
	return b.outputs
}

func (b *botData) Telegram() telegram.Sender {
	return b.tSender
}

func (b *botData) Log() log.Log {
	return b.logger
}

func (b *botData) Task() task.Tasks {
	return b.tasks
}

func (b *botData) Servers() servers.Servers {
	return b.servers
}

func (b *botData) Database() db.Database {
	return b.mongo
}

func CreateBot(mongo db.Database, t telegram.Sender, logs log.Log, username string) Bot {
	return &botData{
		mongo:    mongo,
		tSender:  t,
		logger:   logs,
		username: username,
	}
}

func (b *botData) StartReadErrors() {
	for err := range b.Telegram().Errors().Ch() {
		if err.UserID != 0 {
			b.Telegram().ToQueue(&telegram.Message{
				Message: tgbotapi.NewMessage(err.UserID, "Что-то пошло не так. Попробуйте снова" /*err.Err.Error()*/),
				UserId:  err.UserID,
			})
		} else {
			_ = b.Log().Error("", "system", err.Err.Error())
			_ = b.Servers().SendError(err.Err.Error(), servers.OK)
		}
	}
}

func (b *botData) StartHandle() {
	for update := range b.Telegram().Reader().Chan() {

		if update.EditedMessage != nil &&
			update.ChannelPost != nil {
			continue
		}

		if update.Message == nil &&
			update.CallbackQuery == nil {
			continue
		}

		b.Validation(update)
	}
}
