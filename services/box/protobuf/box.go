package protobuf

import (
	"context"
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sparrowganz/teleFly/telegram"
	"telegram_boxes/services/box/bot"
)

type Box struct {
	Bot bot.Bot
}

func CreateBoxService(b bot.Bot) BoxService {
	return &Box{Bot: b}
}

type BoxService interface {
	Check(ctx context.Context, r *CheckRequest) (*CheckResponse, error)
}

func (b *Box) Check(ctx context.Context, r *CheckRequest) (*CheckResponse, error) {
	ch := make(chan bool)
	var err error

	b.Bot.Methods().Telegram().ToQueue(&telegram.Message{
		Message: tgbotapi.NewMessage(r.GetChatID(), "CHECK is running bot"),
		Type:    ch,
		UserId:  r.GetChatID(),
	})

	if !<-ch {
		err = errors.New(" Message not send ")
	}
	close(ch)

	return &CheckResponse{}, err

}
