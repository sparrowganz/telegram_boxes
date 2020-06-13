package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"telegram_boxes/services/box/app/config"
)

func (b *botData) textsValidation(mess *tgbotapi.Message) {

	if len(b.Config().Keyboards().Buttons) > 0 {
		keyboards := b.Config().Keyboards()
		tp, err := keyboards.GetTypeForName(mess.Text)

		if err != nil {
			return
		}

		switch tp {
		case config.BalanceType:
			b.balanceCommandHandler(mess.Chat.ID)
		case config.ReferralsType:
			b.referralsCommandHandler(mess.Chat.ID)
		case config.HelpType:
			b.helpCommandHandler(mess.Chat.ID)
		}
	}
}
