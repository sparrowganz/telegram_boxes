package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sparrowganz/teleFly/telegram/actions"
	"telegram_boxes/services/admin/protobuf/services/core/protobuf"
)

func (b *botData) imagesValidation(update *tgbotapi.Message) {
	job, ok := b.Telegram().Actions().Get(update.Chat.ID)
	if !ok {
		return
	}

	switch job.GetType() {
	case BroadcastType.String():
		switch job.GetAction() {
		case AddAction.String():
			photos := *update.Photo
			b.addImageValidation(update.Chat.ID, photos[0].FileID, update.Caption, job)
		}
	}
}

func (b *botData) addImageValidation(chatID int64, photoID, caption string, job actions.Job) {
	b.Telegram().DeleteMessages(chatID, job.GetMessageIDs())
	job.FlushMessageId()

	data, ok := job.GetData().(*protobuf.StartBroadcastRequest)
	if !ok {

		b.Telegram().Actions().Delete(chatID)

		_ = b.Log().Error("", "", "addImageValidation: job not found")
		b.Telegram().SendError(chatID, "Что-то пошло не так попробуйте снова", nil)
		return
	}

	data.Type = imageType
	data.FileLink = photoID

	if caption != "" {
		data.Text = caption
	}

	b.broadcastSetData(chatID, data)
}
