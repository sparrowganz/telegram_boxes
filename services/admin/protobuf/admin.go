package protobuf

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sparrowganz/teleFly/telegram"
	"telegram_boxes/services/admin/bot"
)

type Admin struct {
	Bot bot.Bot
}

func CreateAdminService(b bot.Bot) AdminService {
	return &Admin{Bot: b}
}

type AdminService interface {
	SendError(ctx context.Context, r *SendErrorRequest) (*SendErrorResponse, error)
}

func (a *Admin) SendError(_ context.Context, r *SendErrorRequest) (*SendErrorResponse, error) {
	out := &SendErrorResponse{}

	for _, adminID := range a.Bot.Methods().Admins().GetAll() {
		a.Bot.Methods().Telegram().ToQueue(&telegram.Message{
			Message: tgbotapi.NewMessage(adminID,
				fmt.Sprintf("%v(%v) : %v", r.GetUsername(), r.GetStatus(), r.GetError())),
			UserId: adminID,
		})
	}
	return out, nil
}
