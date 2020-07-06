package protobuf

import (
	"context"
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sparrowganz/teleFly/telegram"
	"github.com/sparrowganz/teleFly/telegram/keyboard"
	"github.com/sparrowganz/teleFly/telegram/limits"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"telegram_boxes/services/box/bot"
	"time"
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
	StartBroadcast(in *StartBroadcastRequest, stream Box_StartBroadcastServer) error
}

func (b *Box) StartBroadcast(in *StartBroadcastRequest, stream Box_StartBroadcastServer) error {

	if in.GetFileLink() != "" {

		res, err := http.Get(in.GetFileLink())
		if err != nil {
			_ = b.Bot.Methods().Log().Error(b.Bot.Methods().Username(), "StartBroadcast", err.Error())
			return err
		}

		parts := strings.Split(in.GetFileLink(), "/")
		filename := parts[len(parts)-1]

		f, errFile := os.Create(filepath.Join(os.TempDir(), filename))
		if errFile != nil {
			_ = b.Bot.Methods().Log().Error(b.Bot.Methods().Username(), "StartBroadcast", errFile.Error())
			return errFile
		}

		_, _ = io.Copy(f, res.Body)
		res.Body.Close()
		f.Close()

		if in.GetType() == bot.ImageType {
			resUpload, errSend := b.Bot.Methods().Telegram().API().Send(tgbotapi.NewPhotoUpload(in.GetChatID(),
				filepath.Join(os.TempDir(), filename)))
			if errSend == nil {
				photo := *resUpload.Photo
				in.FileLink = photo[0].FileID
			} else {
				_ = b.Bot.Methods().Log().Error(b.Bot.Methods().Username(), "StartBroadcast", errSend.Error())
				in.FileLink = ""
				in.Type = ""
			}
		} else if in.GetType() == bot.VideoType {
			resUpload, errSend := b.Bot.Methods().Telegram().API().Send(tgbotapi.NewVideoUpload(in.GetChatID(),
				filepath.Join(os.TempDir(), filename)))
			if errSend == nil {
				in.FileLink = resUpload.Video.FileID
			} else {
				_ = b.Bot.Methods().Log().Error(b.Bot.Methods().Username(), "StartBroadcast", errSend.Error())
				in.FileLink = ""
				in.Type = ""
			}
		} else {
			in.FileLink = ""
			in.Type = ""
		}

		_ = os.Remove(filepath.Join(os.TempDir(), filename))
	}

	var replyMarkUp interface{}
	if len(in.GetButtons()) > 0 {
		var rows [][]tgbotapi.InlineKeyboardButton

		for _, but := range in.GetButtons() {
			if but.GetName() != "" && but.GetUrl() != "" {
				inlineButton, err := keyboard.NewButton().SetText(
					but.GetName()).SetUrl(but.GetUrl()).ToUrl()
				if err != nil {
					_ = b.Bot.Methods().Log().Error(b.Bot.Methods().Username(), "StartBroadcast", err.Error())
					continue
				}
				rows = append(rows, tgbotapi.NewInlineKeyboardRow(inlineButton))
			}
		}

		if len(rows) > 0 {
			k := tgbotapi.NewInlineKeyboardMarkup(rows...)
			replyMarkUp = &k
		}
	}

	session := b.Bot.Methods().Database().GetMainSession().Clone()
	defer session.Close()


	users, err := b.Bot.Methods().Database().Models().Users().GetAllUsers(session)
	if err != nil {
		_ = b.Bot.Methods().Log().Error(b.Bot.Methods().Username(), "StartBroadcast", err.Error())
		return err
	}

	ch := make(chan struct{})
	broadcaster := &telegram.BroadcastData{}
	b.Bot.Methods().Telegram().Broadcast().StartListener(broadcaster, "", ch)
	broadcaster.SetAllCount(len(users))

	go func() {

		if in.GetType() == bot.ImageType {
			ticker := time.NewTicker(time.Second / time.Duration(limits.Broadcast()))
		Loop1:
			for {
				select {
				case <-ticker.C:
					if len(users) == 0 {
						ticker.Stop()
						break Loop1
					}

					userID := users[0].Account.Id

					mess := tgbotapi.NewPhotoShare(userID, in.GetFileLink())
					mess.Caption = in.GetText()
					mess.ParseMode = tgbotapi.ModeMarkdown
					mess.ReplyMarkup = replyMarkUp

					b.Bot.Methods().Telegram().ToQueue(&telegram.Message{
						Message: mess,
						Type:    broadcaster,
						UserId:  userID,
					})

					users = users[1:]
					broadcaster.IncAttemptCount()
				case <-stream.Context().Done():
					ticker.Stop()
					break Loop1
				}
			}
		} else if in.GetType() == bot.VideoType {
			ticker := time.NewTicker(time.Second / time.Duration(limits.Broadcast()))
		Loop2:
			for {
				select {
				case <-ticker.C:
					if len(users) == 0 {
						ticker.Stop()
						break Loop2
					}

					userID := users[0].Account.Id

					mess := tgbotapi.NewVideoShare(userID, in.GetFileLink())
					mess.Caption = in.GetText()
					mess.ParseMode = tgbotapi.ModeMarkdown
					mess.ReplyMarkup = replyMarkUp

					b.Bot.Methods().Telegram().ToQueue(&telegram.Message{
						Message: mess,
						Type:    broadcaster,
						UserId:  userID,
					})

					users = users[1:]
					broadcaster.IncAttemptCount()
				case <-stream.Context().Done():
					ticker.Stop()
					break Loop2
				}
			}
		} else {
			ticker := time.NewTicker(time.Second / time.Duration(limits.Broadcast()))
		Loop3:
			for {
				select {
				case <-ticker.C:
					if len(users) == 0 {
						ticker.Stop()
						break Loop3
					}

					userID := users[0].Account.Id

					mess := tgbotapi.NewMessage(userID, in.GetText())
					mess.ParseMode = tgbotapi.ModeMarkdown
					mess.ReplyMarkup = replyMarkUp

					b.Bot.Methods().Telegram().ToQueue(&telegram.Message{
						Message: mess,
						Type:    broadcaster,
						UserId:  userID,
					})

					users = users[1:]
					broadcaster.IncAttemptCount()
				case <-stream.Context().Done():
					ticker.Stop()
					break Loop3
				}
			}
		}
	}()

	t := time.NewTicker(time.Second * 5)
Loop4:
	for {
		select {
		case <-t.C:

			_ = stream.Send(&Stats{
				Success: broadcaster.GetTrueCounter(),
				Fail:    broadcaster.GetFalseCounter(),
			})
		case <-stream.Context().Done():
			t.Stop()
			break Loop4
		case <-ch:
			t.Stop()
			break Loop4
		}
	}
	_ = stream.Send(&Stats{
		Success: broadcaster.GetTrueCounter(),
		Fail:    broadcaster.GetFalseCounter(),
	})

	_ = b.Bot.Methods().Database().Models().Users().BlockListUsers(broadcaster.GetBlockUsers(), session)
	return nil
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

	return &CheckResponse{}, err

}
