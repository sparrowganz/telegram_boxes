package bot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sparrowganz/teleFly/telegram"
	"github.com/sparrowganz/teleFly/telegram/actions"
	"telegram_boxes/services/admin/app/servers"
	"telegram_boxes/services/admin/protobuf/services/core/protobuf"
)

func (b *botData) inlineValidation(update *tgbotapi.CallbackQuery) {
	if update.Data == "" {
		return
	}

	callback := telegram.ParseCallBack(update.Data)
	switch callback.Type() {
	case BonusType:
		switch callback.Action() {
		case ChooseAction:
			switch callback.ID() {
			case AllID:
				b.chooseAllBonusesHandler(update.Message.Chat.ID, update.Message.MessageID)
			default:
				b.chooseBonusHandler(update.Message.Chat.ID, update.Message.MessageID, callback.ID())
			}
		case ChangeActiveAction:
			switch callback.ID() {
			case AllID:
				b.changeActiveAllBonusesHandler(update.Message.Chat.ID, update.Message.MessageID)
			default:
				b.changeActiveBonusHandler(update.Message.Chat.ID, update.Message.MessageID, callback.ID())
			}
		}
	case ServerType:
		switch callback.Action() {
		case CheckAction:
			b.checkAllServers(update.Message.Chat.ID, update.Message.MessageID)
		case FakeAction:
			b.Telegram().DeleteMessages(update.Message.Chat.ID, []int{update.Message.MessageID})
			b.statisticsCommandHandler(update.Message.Chat.ID, true)
		}
	case TaskType:
		switch callback.Action() {
		case GetAction:
			b.getTaskInlineHandler(update.Message.Chat.ID, update.Message.MessageID, callback.ID())
		case PriorityAction:
			b.changePriorityTaskInlineHandler(update.Message.Chat.ID, update.Message.MessageID, update.ID, callback.ID())
		case CleanAction:
			b.cleanupRunTaskInlineHandler(update.Message.Chat.ID, update.Message.MessageID, update.ID, callback.ID())
		case DeleteAction:
			b.removeTaskInlineHandler(update.Message.Chat.ID, update.Message.MessageID, callback.ID())
		case AddAction:

			job, ok := b.Telegram().Actions().Get(update.Message.Chat.ID)
			if !ok {
				b.Telegram().DeleteMessages(update.Message.Chat.ID, []int{update.Message.MessageID})
				_ = b.Log().Error("", "", "inlineValidation: job not found")
				b.Telegram().SendError(update.Message.Chat.ID, "Что-то пошло не так попробуйте снова", nil)
				return
			}

			if job.GetAction() != callback.Action().String() || job.GetType() != callback.Type().String() {
				job.AddMessageId(update.Message.MessageID)
				b.Telegram().DeleteMessages(update.Message.Chat.ID, job.GetMessageIDs())
				b.Telegram().Actions().Delete(update.Message.Chat.ID)

				_ = b.Log().Error("", "", "inlineValidation: action or type is invalid")
				b.Telegram().SendError(update.Message.Chat.ID, "Что-то пошло не так попробуйте снова", nil)
				return
			}

			data := job.GetData().(*protobuf.Task)

			switch "" {
			case data.GetType():
				b.chooseTypeInTaskHandler(update.Message.Chat.ID, update.Message.MessageID, callback.ID(), job, data)
			default:
				b.Telegram().SendError(update.Message.Chat.ID, "Что-то пошло не так попробуйте снова", nil)
				return
			}
		}
	case CancelType:
		b.cancelHandler(update.Message.Chat.ID, update.ID, update.Message.MessageID)
	case LastChoiceType:
		job, ok := b.Telegram().Actions().Get(update.Message.Chat.ID)
		if !ok {
			b.Telegram().DeleteMessages(update.Message.Chat.ID, []int{update.Message.MessageID})
			_ = b.Log().Error("", "", "inlineValidation: job not found")
			b.Telegram().SendError(update.Message.Chat.ID, "Что-то пошло не так попробуйте снова", nil)
			return
		}

		if job.GetAction()+job.GetType() != callback.Action().String() {
			job.AddMessageId(update.Message.MessageID)
			b.Telegram().DeleteMessages(update.Message.Chat.ID, job.GetMessageIDs())
			b.Telegram().Actions().Delete(update.Message.Chat.ID)

			_ = b.Log().Error("", "", "inlineValidation: action or type is invalid")
			b.Telegram().SendError(update.Message.Chat.ID, "Что-то пошло не так попробуйте снова", nil)
			return
		}

		switch telegram.Type(job.GetType()) {
		case TaskType:
			switch telegram.Action(job.GetAction()) {
			case DeleteAction:
				switch callback.ID() {
				case YesID:
					b.forceRemoveInlineTask(update.Message.Chat.ID, update.Message.MessageID, job)
				case NoID:
					b.getTaskInlineHandler(
						update.Message.Chat.ID, update.Message.MessageID, job.GetData().(string))
				}
			case AddAction:
				switch callback.ID() {
				case YesID:
					b.createTaskInlineHandler(update.Message.Chat.ID, update.ID, update.Message.MessageID, job)
				case NoID:
					b.cancelHandler(update.Message.Chat.ID, update.ID, update.Message.MessageID)
				}
			}
		}

	}
}

//
//				BONUS
//

func (b *botData) changeActiveBonusHandler(chatID int64, messageID int, callbackID string) {
	err := b.Servers().ChangeActiveBonus(callbackID)
	if err != nil {
		_ = b.Log().Error("", "", "changeActiveBonusHandler: "+err.Error())
		b.Telegram().SendError(chatID, err.Error(), nil)
	}
	b.chooseBonusHandler(chatID, messageID, callbackID)
}

func (b *botData) changeActiveAllBonusesHandler(chatID int64, messageID int) {
	err := b.Servers().ChangeActiveAllBonuses()
	if err != nil {
		_ = b.Log().Error("", "", "changeActiveAllBonusesHandler: "+err.Error())
		b.Telegram().SendError(chatID, err.Error(), nil)
	}
	b.chooseAllBonusesHandler(chatID, messageID)
}

func (b *botData) chooseAllBonusesHandler(chatID int64, messageID int) {

	var (
		txt           string
		isSetInactive bool
	)

	srv , err := b.Servers().GetAllServers()
	if err != nil {
		_ = b.Log().Error("", "", "chooseBonusHandler: "+err.Error())
		b.Telegram().SendError(chatID, TaskNotFound, nil)
		return
	}

	for _, bon := range srv {

		var smile string
		if bon.GetIsActive() {
			smile += "✅ "
		} else {
			smile += "❌ "
			isSetInactive = true
		}

		txt += fmt.Sprintf("%s %s\n", smile, bon.GetUsername())
	}

	b.Telegram().ToQueue(
		&telegram.Message{
			Message: tgbotapi.EditMessageTextConfig{
				BaseEdit: tgbotapi.BaseEdit{
					ChatID:      chatID,
					MessageID:   messageID,
					ReplyMarkup: changeBonusKeyboard(AllID, isSetInactive),
				},
				Text:      txt,
				ParseMode: tgbotapi.ModeMarkdown,
			},
			UserId: chatID,
		})
}

func (b *botData) chooseBonusHandler(chatID int64, messageID int, callbackID string) {

	server, err := b.Servers().GetServer(callbackID)
	if err != nil {
		_ = b.Log().Error("", "", "chooseBonusHandler: "+err.Error())
		b.Telegram().SendError(chatID, TaskNotFound, nil)
		return
	}

	var smile string
	if server.GetIsActive() {
		smile += "✅ "
	} else {
		smile += "❌ "
	}

	txt := fmt.Sprintf("%s %s \n", smile, server.GetUsername())

	b.Telegram().ToQueue(
		&telegram.Message{
			Message: tgbotapi.EditMessageTextConfig{
				BaseEdit: tgbotapi.BaseEdit{
					ChatID:      chatID,
					MessageID:   messageID,
					ReplyMarkup: changeBonusKeyboard(server.GetId(), !server.GetIsActive()),
				},
				Text:      txt,
				ParseMode: tgbotapi.ModeMarkdown,
			},
			UserId: chatID,
		})
}

//
//				SERVERS
//
func (b *botData) checkAllServers(chatID int64, messageID int) {

	b.Telegram().ToQueue(&telegram.Message{
		Message: tgbotapi.EditMessageTextConfig{
			BaseEdit: tgbotapi.BaseEdit{
				ChatID:    chatID,
				MessageID: messageID,
			},
			Text:                  "Проверка началась",
			ParseMode:             tgbotapi.ModeMarkdown,
			DisableWebPagePreview: true,
		},
		UserId: chatID,
	})

	chResult := make(chan *servers.Check, 100)
	go b.Servers().HardCheckAll(chResult, chatID)

	for res := range chResult {
		b.Telegram().ToQueue(
			&telegram.Message{
				Message: tgbotapi.MessageConfig{
					BaseChat: tgbotapi.BaseChat{
						ChatID: chatID,
					},
					Text:      fmt.Sprintf("%v - %v", res.Username, res.Status),
					ParseMode: tgbotapi.ModeMarkdown,
				},
				UserId: chatID,
			})
	}

	b.Telegram().ToQueue(
		&telegram.Message{
			Message: tgbotapi.MessageConfig{
				BaseChat: tgbotapi.BaseChat{
					ChatID: chatID,
				},
				Text:      "Проверка закончена",
				ParseMode: tgbotapi.ModeMarkdown,
			},
			UserId: chatID,
		})
}

//
//				CANCEL
//
func (b *botData) cancelHandler(chatID int64, queryID string, messageID int) {
	b.Telegram().ToQueue(&telegram.Message{
		Message: tgbotapi.NewCallbackWithAlert(queryID, "Действие отменено"),
		UserId:  chatID,
	})

	job, ok := b.Telegram().Actions().Get(chatID)
	if ok {
		job.AddMessageId(messageID)
		b.Telegram().DeleteMessages(chatID, job.GetMessageIDs())
		job.FlushMessageId()
		b.Telegram().Actions().Delete(chatID)
	} else {
		b.Telegram().DeleteMessages(chatID, []int{messageID})
	}
}

//
//				/TASK
//
const (
	TaskNotFound = "Задание не найдено попробуйте снова"
	TaskTemplate = "Задание %s: \n????(Выводить ли само задание)\n\n%v"
)

func (b *botData) createTaskInlineHandler(chatID int64, queryID string, messageID int, act actions.Job) {
	_ = b.Task().Create(act.GetData().(*protobuf.Task))

	b.Telegram().ToQueue(&telegram.Message{
		Message: tgbotapi.NewCallbackWithAlert(queryID, "Задание создано"),
		UserId:  chatID,
	})

	act.AddMessageId(messageID)

	b.Telegram().DeleteMessages(chatID, act.GetMessageIDs())
	b.Telegram().Actions().Delete(chatID)
}

func (b *botData) chooseTypeInTaskHandler(chatID int64, messageID int, actionID string, act actions.Job, data *protobuf.Task) {

	data.Type = actionID
	data.WithCheck = b.Types().WithCheck(actionID)

	act.ChangeAutoAddMessages(true)

	keyboard, _ := cancelButton().ToKeyboard()

	b.Telegram().ToQueue(&telegram.Message{
		Message: tgbotapi.EditMessageTextConfig{
			BaseEdit: tgbotapi.BaseEdit{
				ChatID:      chatID,
				ReplyMarkup: keyboard,
				MessageID:   messageID,
			},
			Text:                  "Введите имя:",
			ParseMode:             tgbotapi.ModeMarkdown,
			DisableWebPagePreview: true,
		},
		UserId: chatID,
	})

}

func (b *botData) forceRemoveInlineTask(chatID int64, messageID int, job actions.Job) {

	job.AddMessageId(messageID)
	b.Telegram().DeleteMessages(chatID, job.GetMessageIDs())

	err := b.Task().Delete(job.GetData().(string))
	if err != nil {
		_ = b.Log().Error("", "", "forceRemoveInlineTask: "+err.Error())
		b.Telegram().SendError(chatID, TaskNotFound, nil)
		return
	}

	b.Telegram().Actions().Delete(chatID)

	b.tasksCommandHandler(chatID)
}

func (b *botData) removeTaskInlineHandler(chatID int64, messageID int, actionID string) {
	b.Telegram().DeleteMessages(chatID, []int{messageID})

	tsk, err := b.Task().GetTask(actionID)
	if err != nil {
		_ = b.Log().Error("", "", "removeTaskInlineHandler: "+err.Error())
		b.Telegram().SendError(chatID, TaskNotFound, nil)
		return
	}

	if j, ok := b.Telegram().Actions().Get(chatID); ok {
		b.Telegram().DeleteMessages(chatID, j.GetMessageIDs())
		b.Telegram().Actions().Delete(chatID)
	}
	b.Telegram().Actions().New(chatID,
		actions.NewJob(DeleteAction.String(), TaskType.String(), tsk.GetId(), messageID, false))

	b.Telegram().ToQueue(&telegram.Message{
		Message: tgbotapi.MessageConfig{
			BaseChat: tgbotapi.BaseChat{
				ChatID:      chatID,
				ReplyMarkup: lastChoiceKeyboard(DeleteAction.String() + TaskType.String()),
			},
			Text:                  `Вы уверены что хотите удалить задание: "` + tsk.Name + `" `,
			ParseMode:             tgbotapi.ModeMarkdown,
			DisableWebPagePreview: true,
		},
		UserId: chatID,
	})

}

func (b *botData) cleanupRunTaskInlineHandler(chatID int64, messageID int, queryID, actionID string) {

	b.Telegram().DeleteMessages(chatID, []int{messageID})
	err := b.Task().CleanupRun(actionID)
	if err != nil {
		_ = b.Log().Error("", "", "cleanupRunTaskInlineHandler: "+err.Error())
		b.Telegram().SendError(chatID, err.Error(), nil)
		return
	}

	tsk ,errGet := b.Task().GetTask(actionID)
	if errGet != nil {
		_ = b.Log().Error("", "", "cleanupRunTaskInlineHandler: "+errGet.Error())
		b.Telegram().SendError(chatID, errGet.Error(), nil)
		return
	}

	b.Telegram().ToQueue(&telegram.Message{
		Message: tgbotapi.NewCallbackWithAlert(queryID, "Выполнение "+tsk.GetName()+" очищено"),
		UserId:  chatID,
	})

}

func (b *botData) changePriorityTaskInlineHandler(chatID int64, messageID int, queryID, actionID string) {

	tsk, err := b.Task().ChangePriority(actionID)
	if err != nil {
		_ = b.Log().Error("", "", "changePriorityInlineHandler: "+err.Error())
		b.Telegram().SendError(chatID, TaskNotFound, nil)
		return
	}

	var txt string
	if tsk.IsPriority {
		txt = "Приоритет для " + tsk.Name + " активирован"
	} else {
		txt = "Приоритет для " + tsk.Name + " деактивирован"
	}

	b.Telegram().ToQueue(&telegram.Message{
		Message: tgbotapi.NewCallbackWithAlert(queryID, txt),
		UserId:  chatID,
	})

	b.getTaskInlineHandler(chatID, messageID, tsk.GetId())
}

func (b *botData) getTaskInlineHandler(chatID int64, messageID int, actionID string) {

	tsk, err := b.Task().GetTask(actionID)
	if err != nil {
		_ = b.Log().Error("", "", "getTaskInlineHandler: "+err.Error())
		b.Telegram().DeleteMessages(chatID, []int{messageID})
		b.Telegram().SendError(chatID, TaskNotFound, nil)
		return
	}

	var priorityText string
	if tsk.IsPriority {
		priorityText = "✅ Задание приоритетное"
	} else {
		priorityText = "❌ Задание не приоритетно"
	}

	b.Telegram().ToQueue(&telegram.Message{
		Message: tgbotapi.EditMessageTextConfig{
			BaseEdit: tgbotapi.BaseEdit{
				ChatID:      chatID,
				ReplyMarkup: getTaskKeyboard(tsk),
				MessageID:   messageID,
			},
			Text:                  fmt.Sprintf(TaskTemplate, tsk.Name, priorityText),
			ParseMode:             tgbotapi.ModeMarkdown,
			DisableWebPagePreview: true,
		},
		UserId: chatID,
	})
}
