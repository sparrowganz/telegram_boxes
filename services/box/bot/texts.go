package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sparrowganz/teleFly/telegram"
	"github.com/sparrowganz/teleFly/telegram/actions"
	"telegram_boxes/services/box/app/config"
	"telegram_boxes/services/box/app/output"
)

func (b *botData) textsValidation(mess *tgbotapi.Message) {

	if len(b.Config().Keyboards().Buttons) > 0 {
		keyboards := b.Config().Keyboards()
		tp, _ := keyboards.GetTypeForName(mess.Text)

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

			j, ok := b.Telegram().Actions().Get(mess.Chat.ID)
			if ok {

				if j.GetType() == OutputType.String() && j.GetAction() == AddAction.String() {
					data := j.GetData().(*output.Output)
					switch "" {
					case data.Data:
						b.setDataOutput(mess.Chat.ID, mess.Text, j)
					}
				}
			}
		}
	}
}

func (b *botData) setDataOutput(chatID int64, input string, j actions.Job) {
	b.Telegram().DeleteMessages(chatID, j.GetMessageIDs())

	txt, keyb, err := b.setOutputData(chatID, input, j.GetData().(*output.Output))
	if err != nil {
		_ = b.Log().Error(b.Username(), "setDataOutput", err.Error())
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
