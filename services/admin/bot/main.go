package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sparrowganz/teleFly/telegram"
	"telegram_boxes/services/admin/app/admins"
	"telegram_boxes/services/admin/app/log"
	"telegram_boxes/services/admin/app/servers"
	"telegram_boxes/services/admin/app/task"
	"telegram_boxes/services/admin/app/types"
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
	SetTypes(t types.Types)
	SetServers(s servers.Servers)
}

func (b *botData) SetTasks(t task.Tasks) {
	b.tasks = t
}

func (b *botData) SetTypes(t types.Types) {
	b.types = t
}

func (b *botData) SetServers(s servers.Servers) {
	b.servers = s
}

type Botter interface {
	BotSetter
	Admins() admins.Admin
	Telegram() telegram.Sender
	Log() log.Log
	Task() task.Tasks
	Types() types.Types
	Servers() servers.Servers
}

type botData struct {
	admins  admins.Admin
	tSender telegram.Sender
	logger  log.Log
	tasks   task.Tasks
	types   types.Types
	servers servers.Servers
}

func (b *botData) Admins() admins.Admin {
	return b.admins
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

func (b *botData) Types() types.Types {
	return b.types
}

func (b *botData) Servers() servers.Servers {
	return b.servers
}

func CreateBot(a admins.Admin, t telegram.Sender, logs log.Log) Bot {
	return &botData{
		admins:  a,
		tSender: t,
		logger:  logs,
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
			for _, id := range b.Admins().GetAll() {
				b.Telegram().ToQueue(&telegram.Message{
					Message: tgbotapi.NewMessage(id, err.Err.Error()),
					UserId:  id,
				})
			}
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

		var chatID int64
		if update.Message != nil {
			chatID = update.Message.Chat.ID
		} else if update.CallbackQuery != nil {
			chatID = update.CallbackQuery.Message.Chat.ID
		}

		if !b.Admins().IsSet(chatID) {
			continue
		}

		b.Validation(update)
	}
}
