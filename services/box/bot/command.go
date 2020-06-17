package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sparrowganz/teleFly/telegram"
	"gopkg.in/mgo.v2"
	"strconv"
	"telegram_boxes/services/box/app/config"
	"telegram_boxes/services/box/app/models"
)

type Command string

var (
	start Command = "start"
)

func (b *botData) commandValidation(mess *tgbotapi.Message) {

	switch Command(mess.Command()) {
	case start:
		if mess.CommandArguments() != "" {
			b.startReferralCommandHandler(
				mess.Chat.ID,
				mess.Chat.UserName, mess.Chat.FirstName, mess.Chat.LastName,
				"", mess.CommandArguments())
			return
		} else {
			b.startCommandHandler(mess.Chat.ID, mess.Chat.UserName, mess.Chat.FirstName, mess.Chat.LastName, "")
			return
		}
	default:
		if len(b.Config().Commands()) > 0 {

			tp, err := b.Config().GetTypeForNameCommands(mess.Command())
			if err != nil {
				b.sendUnknownCommand(mess.Chat.ID)
				return
			}

			switch tp {
			case config.BalanceType:
				b.balanceCommandHandler(mess.Chat.ID)
			case config.ReferralsType:
				b.referralsCommandHandler(mess.Chat.ID)
			case config.HelpType:
				b.helpCommandHandler(mess.Chat.ID)
			case config.OutputType:
				b.outputCommandHandler(mess.Chat.ID)
			default:
				b.sendUnknownCommand(mess.Chat.ID)
				return
			}
		} else {
			b.sendUnknownCommand(mess.Chat.ID)
			return
		}
	}

}

func (b *botData) sendUnknownCommand(chatID int64) {
	b.Telegram().ToQueue(
		&telegram.Message{
			Message: tgbotapi.NewMessage(chatID, b.GetErrorCommandText()),
			UserId:  chatID,
		})
}

func (b *botData) outputCommandHandler(chatID int64) {
	txt, keyb, err := b.output(chatID, config.Static)
	if err != nil {
		_ = b.Log().Error(b.Username(), "outputCommandHandler", err.Error())
		b.Telegram().SendError(chatID, b.GetErrorText(), nil)
		return
	}

	b.Telegram().ToQueue(&telegram.Message{
		Message: tgbotapi.MessageConfig{
			BaseChat: tgbotapi.BaseChat{
				ChatID:      chatID,
				ReplyMarkup: keyb,
			},
			Text:                  txt,
			ParseMode:             tgbotapi.ModeHTML,
			DisableWebPagePreview: true,
		},
		UserId: chatID,
	})

}

func (b *botData) helpCommandHandler(chatID int64) {
	txt, keyb, err := b.help(chatID, config.Static)
	if err != nil {
		_ = b.Log().Error(b.Username(), "referralsCommandHandler", err.Error())
		b.Telegram().SendError(chatID, b.GetErrorText(), nil)
		return
	}

	b.Telegram().ToQueue(&telegram.Message{
		Message: tgbotapi.MessageConfig{
			BaseChat: tgbotapi.BaseChat{
				ChatID:      chatID,
				ReplyMarkup: keyb,
			},
			Text:                  txt,
			ParseMode:             tgbotapi.ModeHTML,
			DisableWebPagePreview: true,
		},
		UserId: chatID,
	})

}

func (b *botData) referralsCommandHandler(chatID int64) {
	txt, keyb, err := b.referrals(chatID, config.Static)
	if err != nil {
		_ = b.Log().Error(b.Username(), "referralsCommandHandler", err.Error())
		b.Telegram().SendError(chatID, b.GetErrorText(), nil)
		return
	}

	b.Telegram().ToQueue(&telegram.Message{
		Message: tgbotapi.MessageConfig{
			BaseChat: tgbotapi.BaseChat{
				ChatID:      chatID,
				ReplyMarkup: keyb,
			},
			Text:                  txt,
			ParseMode:             tgbotapi.ModeHTML,
			DisableWebPagePreview: true,
		},
		UserId: chatID,
	})

}

func (b *botData) balanceCommandHandler(chatID int64) {
	txt, keyb, err := b.balance(chatID, config.Static)
	if err != nil {
		_ = b.Log().Error(b.Username(), "balanceCommandHandler", err.Error())
		b.Telegram().SendError(chatID, b.GetErrorText(), nil)
		return
	}

	b.Telegram().ToQueue(&telegram.Message{
		Message: tgbotapi.MessageConfig{
			BaseChat: tgbotapi.BaseChat{
				ChatID:      chatID,
				ReplyMarkup: keyb,
			},
			Text:                  txt,
			ParseMode:             tgbotapi.ModeHTML,
			DisableWebPagePreview: true,
		},
		UserId: chatID,
	})

}

func (b *botData) startCommandHandler(chatID int64, username, firstname, lastname, email string) {

	session := b.Database().GetMainSession().Clone()
	defer session.Close()

	_, err := b.Database().Models().Users().FindUserByTelegramID(chatID, session)
	if err != nil {
		if err != mgo.ErrNotFound {
			_ = b.Log().Error(b.Username(), "startCommandHandler", err.Error())
			b.Telegram().SendError(chatID, b.GetErrorText(), nil)
			return
		}
		_ = b.Database().Models().Users().CreateUser(
			models.CreateUser(chatID, username, firstname, lastname, email), session)
	}

	b.Telegram().ToQueue(&telegram.Message{
		Message: tgbotapi.MessageConfig{
			BaseChat: tgbotapi.BaseChat{
				ChatID:      chatID,
				ReplyMarkup: b.GetMainKeyboard(),
			},
			Text:                  b.GetStartText(),
			ParseMode:             tgbotapi.ModeHTML,
			DisableWebPagePreview: true,
		},
		UserId: chatID,
	})
}

func (b *botData) startReferralCommandHandler(chatID int64, username, firstname, lastname, email string, args string) {

	session := b.Database().GetMainSession().Clone()
	defer session.Close()

	inviterID, errConv := strconv.Atoi(args)
	if errConv != nil {
		b.startCommandHandler(chatID, username, firstname, lastname, email)
		return
	}

	_, err := b.Database().Models().Users().FindUserByTelegramID(chatID, session)
	if err != nil {
		if err != mgo.ErrNotFound {
			_ = b.Log().Error(b.Username(), "startCommandHandler", err.Error())
			b.Telegram().SendError(chatID, b.GetErrorText(), nil)
			return
		}

		u := models.CreateUser(chatID, username, firstname, lastname, email)

		inviter, errFindInviter := b.Database().Models().Users().FindUserByTelegramID(int64(inviterID), session)
		if errFindInviter == nil {
			u.SetInviterID(inviter.ID())

			var referralUsername string
			if u.Telegram().UserName() != "" {
				referralUsername = "@" + u.Telegram().UserName()
			} else {
				referralUsername = strconv.Itoa(int(u.Telegram().ID()))
			}

			inviter.Balance().AddBot(float64(b.Config().Counts().CostForReferral))

			countReferral := b.Database().Models().Users().GetCountInvitedUsers(inviter.ID(), session)
			if countReferral >= b.Config().Counts().VerifiedCount {
				inviter.SetVerified()
			}

			_ = b.Database().Models().Users().UpdateUser(inviter, session)

			b.Telegram().ToQueue(&telegram.Message{
				Message: tgbotapi.MessageConfig{
					BaseChat: tgbotapi.BaseChat{
						ChatID: inviter.Telegram().ID(),
					},
					Text: b.BonusForReferralText(
						referralUsername, b.getReferralLink(inviter.Telegram().ID())),
					ParseMode:             tgbotapi.ModeMarkdown,
					DisableWebPagePreview: true,
				},
				UserId: inviter.Telegram().ID(),
			})

		}

		_ = b.Database().Models().Users().CreateUser(u, session)

	}

	b.Telegram().ToQueue(&telegram.Message{
		Message: tgbotapi.MessageConfig{
			BaseChat: tgbotapi.BaseChat{
				ChatID:      chatID,
				ReplyMarkup: b.GetMainKeyboard(),
			},
			Text:                  b.GetStartText(),
			ParseMode:             tgbotapi.ModeHTML,
			DisableWebPagePreview: true,
		},
		UserId: chatID,
	})

}
