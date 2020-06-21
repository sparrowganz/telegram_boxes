package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sparrowganz/teleFly/telegram"
	"telegram_boxes/services/box/app/config"
)

func (b *botData) inlineValidation(update *tgbotapi.CallbackQuery) {
	if update.Data == "" {
		return
	}

	callback := telegram.ParseCallBack(update.Data)

	switch callback.Type().String() {
	case config.BalanceType.ToString():
		b.balanceInlineHandler(update.Message.Chat.ID, update.Message.MessageID)
	case config.ReferralsType.ToString():
		b.referralsInlineHandler(update.Message.Chat.ID, update.Message.MessageID)
	case config.HelpType.ToString():
		b.helpInlineHandler(update.Message.Chat.ID, update.Message.MessageID)
	case config.CancelType.ToString():
		b.cancelInlineHandler(update.Message.Chat.ID, update.Message.MessageID)
	case config.OutputType.ToString():
		b.outputInlineHandler(update.Message.Chat.ID, update.Message.MessageID)
	case config.CheckTaskType.ToString():
		b.checkTaskInlineHandler(update.Message.Chat.ID, update.Message.MessageID, callback.Action().String())
	case config.SkipTaskType.ToString():
		b.skipTaskInlineHandler(update.Message.Chat.ID, update.Message.MessageID, callback.Action().String())
	case config.NextTaskType.ToString():
		b.nextTaskInlineHandler(update.Message.Chat.ID, update.Message.MessageID)
	default:
		if b.IsOutputGWButton(callback.Type().String()) {
			b.outputGWInlineHandler(update.Message.Chat.ID, update.Message.MessageID, callback.Action().String())
		}
	}
}

func (b *botData) nextTaskInlineHandler(chatID int64, messageID int) {
	text, keyb, err := b.getTask(chatID)
	if err != nil {
		_ = b.Log().Error(b.Username(), "nextTaskInlineHandler", err.Error())
		b.Telegram().SendError(chatID, b.GetErrorText(), nil)
		return
	}

	message := tgbotapi.EditMessageTextConfig{
		BaseEdit: tgbotapi.BaseEdit{
			ChatID:    chatID,
			MessageID: messageID,
		},
		Text:      text,
		ParseMode: tgbotapi.ModeHTML,
	}

	switch val := keyb.(type) {
	case *tgbotapi.InlineKeyboardMarkup:
		message.ReplyMarkup = val
	}

	b.Telegram().ToQueue(&telegram.Message{
		Message: message,
		UserId:  chatID,
	})
}

func (b *botData) skipTaskInlineHandler(chatID int64, messageID int, id string) {
	text, keyb, err := b.skipTask(chatID, id)
	if err != nil {
		_ = b.Log().Error(b.Username(), "skipTaskInlineHandler", err.Error())
		b.Telegram().SendError(chatID, b.GetErrorText(), nil)
		return
	}

	message := tgbotapi.EditMessageTextConfig{
		BaseEdit: tgbotapi.BaseEdit{
			ChatID:    chatID,
			MessageID: messageID,
		},
		Text:      text,
		ParseMode: tgbotapi.ModeHTML,
	}

	switch val := keyb.(type) {
	case *tgbotapi.InlineKeyboardMarkup:
		message.ReplyMarkup = val
	}

	b.Telegram().ToQueue(&telegram.Message{
		Message: message,
		UserId:  chatID,
	})
}

func (b *botData) checkTaskInlineHandler(chatID int64, messageID int, id string) {
	text, keyb, err := b.checkTask(chatID, id)
	if err != nil {
		_ = b.Log().Error(b.Username(), "checkTaskInlineHandler", err.Error())
		b.Telegram().SendError(chatID, b.GetErrorText(), nil)
		return
	}

	message := tgbotapi.EditMessageTextConfig{
		BaseEdit: tgbotapi.BaseEdit{
			ChatID:    chatID,
			MessageID: messageID,
		},
		Text:      text,
		ParseMode: tgbotapi.ModeHTML,
	}

	switch val := keyb.(type) {
	case *tgbotapi.InlineKeyboardMarkup:
		message.ReplyMarkup = val
	}

	b.Telegram().ToQueue(&telegram.Message{
		Message: message,
		UserId:  chatID,
	})
}

func (b *botData) outputGWInlineHandler(chatID int64, messageID int, outputGW string) {
	//b.Telegram().DeleteMessages(chatID, []int{messageID})
	text, keyb, err := b.chooseOutputGW(chatID, config.Inline, outputGW)
	if err != nil {
		_ = b.Log().Error(b.Username(), "outputGWInlineHandler", err.Error())
		b.Telegram().SendError(chatID, b.GetErrorText(), nil)
		return
	}

	message := tgbotapi.EditMessageTextConfig{
		BaseEdit: tgbotapi.BaseEdit{
			ChatID:    chatID,
			MessageID: messageID,
		},
		Text:      text,
		ParseMode: tgbotapi.ModeHTML,
	}

	switch val := keyb.(type) {
	case *tgbotapi.InlineKeyboardMarkup:
		message.ReplyMarkup = val
	}

	b.Telegram().ToQueue(&telegram.Message{
		Message: message,
		UserId:  chatID,
	})
}

func (b *botData) outputInlineHandler(chatID int64, messageID int) {
	text, keyb, err := b.output(chatID, config.Inline)
	if err != nil {
		b.Telegram().DeleteMessages(chatID, []int{messageID})
		_ = b.Log().Error(b.Username(), "outputInlineHandler", err.Error())
		b.Telegram().SendError(chatID, b.GetErrorText(), nil)
		return
	}

	message := tgbotapi.EditMessageTextConfig{
		BaseEdit: tgbotapi.BaseEdit{
			ChatID:    chatID,
			MessageID: messageID,
		},
		Text:      text,
		ParseMode: tgbotapi.ModeHTML,
	}

	switch val := keyb.(type) {
	case *tgbotapi.InlineKeyboardMarkup:
		message.ReplyMarkup = val
	}

	b.Telegram().ToQueue(&telegram.Message{
		Message: message,
		UserId:  chatID,
	})

}

func (b *botData) helpInlineHandler(chatID int64, messageID int) {
	text, keyb, err := b.help(chatID, config.Inline)
	if err != nil {
		b.Telegram().DeleteMessages(chatID, []int{messageID})
		_ = b.Log().Error(b.Username(), "helpInlineHandler", err.Error())
		b.Telegram().SendError(chatID, b.GetErrorText(), nil)
		return
	}

	message := tgbotapi.EditMessageTextConfig{
		BaseEdit: tgbotapi.BaseEdit{
			ChatID:    chatID,
			MessageID: messageID,
		},
		Text:      text,
		ParseMode: tgbotapi.ModeHTML,
	}

	switch val := keyb.(type) {
	case *tgbotapi.InlineKeyboardMarkup:
		message.ReplyMarkup = val
	}

	b.Telegram().ToQueue(&telegram.Message{
		Message: message,
		UserId:  chatID,
	})

}

func (b *botData) referralsInlineHandler(chatID int64, messageID int) {
	text, keyb, err := b.referrals(chatID, config.Inline)
	if err != nil {
		b.Telegram().DeleteMessages(chatID, []int{messageID})
		_ = b.Log().Error(b.Username(), "referralsInlineHandler", err.Error())
		b.Telegram().SendError(chatID, b.GetErrorText(), nil)
		return
	}

	message := tgbotapi.EditMessageTextConfig{
		BaseEdit: tgbotapi.BaseEdit{
			ChatID:    chatID,
			MessageID: messageID,
		},
		Text:      text,
		ParseMode: tgbotapi.ModeHTML,
	}

	switch val := keyb.(type) {
	case *tgbotapi.InlineKeyboardMarkup:
		message.ReplyMarkup = val
	}

	b.Telegram().ToQueue(&telegram.Message{
		Message: message,
		UserId:  chatID,
	})

}

func (b *botData) cancelInlineHandler(chatID int64, messageID int) {

	j, ok := b.Telegram().Actions().Get(chatID)
	if ok {
		j.AddMessageId(messageID)
		b.Telegram().DeleteMessages(chatID, j.GetMessageIDs())
		b.Telegram().Actions().Delete(chatID)
	} else {
		b.Telegram().DeleteMessages(chatID, []int{messageID})
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

func (b *botData) balanceInlineHandler(chatID int64, messageID int) {
	text, keyb, err := b.balance(chatID, config.Inline)
	if err != nil {
		b.Telegram().DeleteMessages(chatID, []int{messageID})
		b.Telegram().DeleteMessages(chatID, []int{messageID})
		b.Telegram().SendError(chatID, "Что-то пошло не так попробуйте снова", nil)
		return
	}

	message := tgbotapi.EditMessageTextConfig{
		BaseEdit: tgbotapi.BaseEdit{
			ChatID:    chatID,
			MessageID: messageID,
		},
		Text:      text,
		ParseMode: tgbotapi.ModeHTML,
	}

	switch val := keyb.(type) {
	case *tgbotapi.InlineKeyboardMarkup:
		message.ReplyMarkup = val
	}

	b.Telegram().ToQueue(&telegram.Message{
		Message: message,
		UserId:  chatID,
	})
	return

}
