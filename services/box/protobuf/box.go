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
	GetStatistics(ctx context.Context, r *GetStatisticsRequest) (*Statistic, error)
	RemoveCheckTask(ctx context.Context, r *RemoveCheckTaskRequest) (*RemoveCheckTaskResponse, error)
}

func (b *Box) RemoveCheckTask(_ context.Context, r *RemoveCheckTaskRequest) (*RemoveCheckTaskResponse, error) {
	out := &RemoveCheckTaskResponse{}

	session := b.Bot.Methods().Database().GetMainSession().Clone()
	defer session.Close()

	err := b.Bot.Methods().Database().Models().Users().RemoveCheck(r.GetTaskID(), session)
	if err != nil {
		return out, err
	}

	return out, nil
}

func (b *Box) GetStatistics(_ context.Context, _ *GetStatisticsRequest) (*Statistic, error) {
	out := &Statistic{}

	session := b.Bot.Methods().Database().GetMainSession().Clone()
	defer session.Close()

	out.All = int64(b.Bot.Methods().Database().Models().Users().GetAllCount(session))
	out.Blocked = int64(b.Bot.Methods().Database().Models().Users().GetBlockedCount(session))
	out.Current = int64(b.Bot.Methods().Telegram().GetCurrentUsers())
	return out, nil
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
