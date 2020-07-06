package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (b *botData) Validation(update tgbotapi.Update) {
	if update.CallbackQuery != nil {
		b.inlineValidation(update.CallbackQuery)
	} else if update.Message != nil {
		if update.Message.Photo != nil {
			b.imagesValidation(update.Message)
		} else if update.Message.Video != nil {
			b.videosValidation(update.Message)
		} else if update.Message.IsCommand() {
			b.commandValidation(update.Message)
		} else if update.Message.Text != "" {
			b.textValidation(update.Message)
		}
	}
}
