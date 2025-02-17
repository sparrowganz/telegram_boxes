package protobuf

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sparrowganz/teleFly/telegram"
	"net/url"
	"strings"
	"telegram_boxes/services/admin/app"
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
	SendMessage(ctx context.Context, r *SendMessageRequest) (*SendMessageResponse, error)
	CheckExecution(ctx context.Context, r *CheckExecutionRequest) (*CheckExecutionResponse, error)
}

func (a *Admin) SendMessage(_ context.Context, r *SendMessageRequest) (*SendMessageResponse, error) {
	out := &SendMessageResponse{}

	for _, adminID := range a.Bot.Methods().Admins().GetAll() {
		a.Bot.Methods().Telegram().ToQueue(&telegram.Message{
			Message: tgbotapi.NewMessage(adminID,
				fmt.Sprintf("%v %v", r.GetUsername(), r.GetMessage())),
			UserId: adminID,
		})
	}
	return out, nil
}

func (a *Admin) SendError(_ context.Context, r *SendErrorRequest) (*SendErrorResponse, error) {
	out := &SendErrorResponse{}

	for _, adminID := range a.Bot.Methods().Admins().GetAll() {

		var smile string
		if r.GetStatus() == app.StatusOK.String() {
			smile = "✅"
		} else if r.GetStatus() == app.StatusFatal.String() {
			smile = "❌"
		} else if r.GetStatus() == app.StatusRecovery.String() {
			smile = "⚠"
		} else {
			smile = "❔"
		}

		var err string
		if r.GetError() != "" {
			err = fmt.Sprintf("%v", r.GetError())
		}

		a.Bot.Methods().Telegram().ToQueue(&telegram.Message{
			Message: tgbotapi.NewMessage(adminID,
				fmt.Sprintf("@%v %v %v", r.GetUsername(), smile, err)),
			UserId: adminID,
		})
	}
	return out, nil
}

func (a *Admin) CheckExecution(_ context.Context, r *CheckExecutionRequest) (*CheckExecutionResponse, error) {
	out := &CheckExecutionResponse{}

	ch := make(chan bool)

	urlParsed, _ := url.Parse(r.GetUrl())
	groupUsername := fmt.Sprintf("@%v", strings.Replace(urlParsed.Path, "/", "", -1))


	a.Bot.Methods().Telegram().ToQueue(
		&telegram.Message{
			Message: tgbotapi.ChatConfigWithUser{
				SuperGroupUsername: groupUsername,
				UserID:             int(r.GetChatID()),
			},
			Type:   ch,
			UserId: r.GetChatID(),
		})


	res := <-ch

	if res {
		out.IsCheck = true
	}

	return out, nil
}
