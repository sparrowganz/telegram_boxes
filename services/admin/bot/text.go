package bot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sparrowganz/teleFly/telegram"
	"github.com/sparrowganz/teleFly/telegram/actions"
	"net/url"
	"telegram_boxes/services/admin/app/task"
)

func (b *botData) textValidation(update *tgbotapi.Message) {
	switch update.Text {
	default:

		job, ok := b.Telegram().Actions().Get(update.Chat.ID)
		if !ok {
			return
		}

		switch job.GetType() {
		case TaskType.String():
			switch job.GetAction() {
			case AddAction.String():

				data := job.GetData().(*task.Task)

				switch "" {
				case data.Name:
					b.setNameToAddTaskHandler(update.Chat.ID, update.Text, job)
				case data.Link:
					b.setLinkToAddTaskHandler(update.Chat.ID, update.Text, job)
				}
			}
		}
	}
}

func (b *botData) setNameToAddTaskHandler(chatID int64, message string, job actions.Job) {

	b.Telegram().DeleteMessages(chatID, job.GetMessageIDs())
	job.FlushMessageId()

	data := job.GetData().(*task.Task)
	data.Name = message

	keyboard, _ := cancelButton().ToKeyboard()

	b.Telegram().ToQueue(&telegram.Message{
		Message: tgbotapi.MessageConfig{
			BaseChat: tgbotapi.BaseChat{
				ChatID:      chatID,
				ReplyMarkup: keyboard,
			},
			Text:                  "Введите ссылку:",
			ParseMode:             tgbotapi.ModeMarkdown,
			DisableWebPagePreview: true,
		},
		UserId: chatID,
	})
}

func (b *botData) setLinkToAddTaskHandler(chatID int64, message string, job actions.Job) {

	_, err := url.ParseRequestURI(message)
	if err != nil {
		b.Telegram().ToQueue(&telegram.Message{
			Message: tgbotapi.MessageConfig{
				BaseChat: tgbotapi.BaseChat{
					ChatID: chatID,
				},
				Text:                  "Некорректная ссылка",
				ParseMode:             tgbotapi.ModeMarkdown,
				DisableWebPagePreview: true,
			},
			UserId: chatID,
		})
		return
	}

	b.Telegram().DeleteMessages(chatID, job.GetMessageIDs())
	job.FlushMessageId()

	data := job.GetData().(*task.Task)
	data.Link = message

	tp, errGetType := b.Types().GetType(data.TypeID)
	if errGetType != nil {
		b.Telegram().ToQueue(&telegram.Message{
			Message: tgbotapi.MessageConfig{
				BaseChat: tgbotapi.BaseChat{
					ChatID: chatID,
				},
				Text:                  "Произошла ошибка: Выбранный тип не найден. Попробуйте снова",
				ParseMode:             tgbotapi.ModeMarkdown,
				DisableWebPagePreview: true,
			},
			UserId: chatID,
		})
		b.Telegram().DeleteMessages(chatID, job.GetMessageIDs())
		b.Telegram().Actions().Delete(chatID)
		return
	}
	job.ChangeAutoAddMessages(false)

	b.Telegram().ToQueue(&telegram.Message{
		Message: tgbotapi.MessageConfig{
			BaseChat: tgbotapi.BaseChat{
				ChatID:      chatID,
				ReplyMarkup: lastChoiceKeyboard(AddAction.String() + TaskType.String()),
			},
			Text:                  fmt.Sprintf("%s, %s\n%s\n\nСохранить?", data.Name, tp.Name, data.Link),
			ParseMode:             tgbotapi.ModeMarkdown,
			DisableWebPagePreview: true,
		},
		UserId: chatID,
	})
}
