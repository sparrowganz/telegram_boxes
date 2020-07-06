package bot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sparrowganz/teleFly/telegram"
	"github.com/sparrowganz/teleFly/telegram/actions"
	"net/url"
	"regexp"
	"telegram_boxes/services/admin/protobuf/services/core/protobuf"
)

func (b *botData) textValidation(update *tgbotapi.Message) {
	switch update.Text {
	default:

		job, ok := b.Telegram().Actions().Get(update.Chat.ID)
		if !ok {
			fmt.Println("!OK")
			return
		}

		fmt.Println(job.GetType(), job.GetAction(), job.GetObject())

		switch job.GetType() {
		case TaskType.String():
			switch job.GetAction() {
			case AddAction.String():

				data := job.GetData().(*protobuf.Task)

				switch "" {
				case data.Name:
					b.setNameToAddTaskHandler(update.Chat.ID, update.Text, job)
				case data.Link:
					b.setLinkToAddTaskHandler(update.Chat.ID, update.Text, job)
				}
			}
		case BroadcastType.String():
			switch job.GetAction() {
			case AddAction.String():
				b.setTextToBroadcastHandler(update.Chat.ID, update.Text, job)
			case AddButton.String():
				b.setButtonBroadcastHandler(update.Chat.ID, update.Text, job)
			}
		}
	}
}

func (b *botData) setButtonBroadcastHandler(chatID int64, message string, job actions.Job) {
	b.Telegram().DeleteMessages(chatID, job.GetMessageIDs())
	job.FlushMessageId()

	data, ok := job.GetData().(*protobuf.StartBroadcastRequest)
	if !ok {

		b.Telegram().Actions().Delete(chatID)

		_ = b.Log().Error("", "", "setTextToBroadcastHandler: job not found")
		b.Telegram().SendError(chatID, "Что-то пошло не так попробуйте снова", nil)
		return
	}

	r := regexp.MustCompile(`(.*) - (.*)`)
	parts := r.FindStringSubmatch(message)
	if len(parts) != 3 {
		b.Telegram().ToQueue(&telegram.Message{
			Message: tgbotapi.MessageConfig{
				BaseChat: tgbotapi.BaseChat{
					ChatID:      chatID,
					ReplyMarkup: backMenu(),
				},
				Text:                  "Некорректный формат!!\nВведите \"VK - https://vk.com\"",
				ParseMode:             tgbotapi.ModeMarkdown,
				DisableWebPagePreview: false,
			},
			UserId: chatID,
		})
		return
	}

	_, err := url.Parse(parts[2])
	if err != nil {
		b.Telegram().ToQueue(&telegram.Message{
			Message: tgbotapi.MessageConfig{
				BaseChat: tgbotapi.BaseChat{
					ChatID:      chatID,
					ReplyMarkup: backMenu(),
				},
				Text:                  `Некорректная ссылка`,
				ParseMode:             tgbotapi.ModeMarkdown,
				DisableWebPagePreview: false,
			},
			UserId: chatID,
		})
		return
	}

	data.Buttons = append(data.Buttons, &protobuf.Button{
		Name: parts[1],
		Url:  parts[2],
	})

	job.SetAction(AddAction.String())

	b.broadcastSetData(chatID, data)
}

func (b *botData) setTextToBroadcastHandler(chatID int64, message string, job actions.Job) {
	b.Telegram().DeleteMessages(chatID, job.GetMessageIDs())
	job.FlushMessageId()

	data, ok := job.GetData().(*protobuf.StartBroadcastRequest)
	if !ok {

		b.Telegram().Actions().Delete(chatID)

		_ = b.Log().Error("", "", "setTextToBroadcastHandler: job not found")
		b.Telegram().SendError(chatID, "Что-то пошло не так попробуйте снова", nil)
		return
	}

	data.Text = message

	b.broadcastSetData(chatID, data)
}

func (b *botData) setNameToAddTaskHandler(chatID int64, message string, job actions.Job) {

	b.Telegram().DeleteMessages(chatID, job.GetMessageIDs())
	job.FlushMessageId()

	data := job.GetData().(*protobuf.Task)
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

	data := job.GetData().(*protobuf.Task)
	data.Link = message

	tp, errGetType := b.Types().GetType(data.GetType())
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
		job.FlushMessageId()
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
