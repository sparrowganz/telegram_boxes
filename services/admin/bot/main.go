package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sparrowganz/teleFly/telegram"
	"sync"
	"telegram_boxes/services/admin/app/admins"
	"telegram_boxes/services/admin/app/log"
)

type Bot interface {
	StartHandle(wg *sync.WaitGroup)
	StartReadErrors(wg *sync.WaitGroup)
}

type botData struct {
	admins  admins.Admin
	tSender telegram.Sender
	logger  log.Log
}

func CreateBot(a admins.Admin,t telegram.Sender, logs log.Log) Bot {
	return &botData{
		admins:  a,
		tSender: t,
		logger:  logs,
	}
}

func (b *botData) StartReadErrors(wg *sync.WaitGroup) {
	defer wg.Done()
	for err := range b.tSender.Errors().Ch() {
		if err.UserID != 0 {
			b.tSender.ToQueue(&telegram.Message{
				Message: tgbotapi.NewMessage(err.UserID, err.Err.Error()),
				UserId:  err.UserID,
			})
		} else {
			for _, id := range b.admins.GetAll() {
				b.tSender.ToQueue(&telegram.Message{
					Message: tgbotapi.NewMessage(id, err.Err.Error()),
					UserId:  id,
				})
			}
		}
	}
}

func (b *botData) StartHandle(wg *sync.WaitGroup) {
	defer wg.Done()
	for update := range b.tSender.Reader().Chan() {
		b.tSender.ToQueue(&telegram.Message{
			Message: tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text),
			UserId:  update.Message.Chat.ID,
		})
	}
}
