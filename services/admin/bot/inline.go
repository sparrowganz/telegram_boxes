package bot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sparrowganz/teleFly/telegram"
	"github.com/sparrowganz/teleFly/telegram/actions"
)

func (b *botData) inlineValidation(update *tgbotapi.CallbackQuery) {
	if update.Data == "" {
		return
	}

	callback := telegram.ParseCallBack(update.Data)
	switch callback.Type() {
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
		}
	case LastChoiceType:
		job, ok := b.Actions().Get(update.Message.Chat.ID)
		if !ok {
			b.Telegram().DeleteMessages(update.Message.Chat.ID, []int{update.Message.MessageID})
			_ = b.Log().Error("", "", "inlineValidation: job not found")
			b.Telegram().SendError(update.Message.Chat.ID, "Что-то пошло не так попробуйте снова", nil)
			return
		}

		if job.GetAction()+job.GetType() != callback.Action().String() {
			b.Telegram().DeleteMessages(update.Message.Chat.ID, []int{update.Message.MessageID})
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
					b.forceRemoveInlineTask(update.Message.Chat.ID, update.Message.MessageID, job.GetData().(string))
				case NoID:
					b.getTaskInlineHandler(
						update.Message.Chat.ID, update.Message.MessageID, job.GetData().(string))
				}
			}
		}

	}
}

//
//				TASKS
//

const (
	TaskNotFound = "Задание не найдено попробуйте снова"
	TaskTemplate = "Задание %s: \n????(Выводить ли само задание)\n\n%v"
)

func (b *botData) forceRemoveInlineTask(chatID int64, messageID int, actionID string) {
	b.Telegram().DeleteMessages(chatID, []int{messageID})

	err := b.Task().Delete(actionID)
	if err != nil {
		_ = b.Log().Error("", "", "forceRemoveInlineTask: "+err.Error())
		b.Telegram().SendError(chatID, TaskNotFound, nil)
		return
	}

	b.Actions().Delete(chatID)

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

	b.Actions().New(chatID, actions.NewJob(DeleteAction.String(), TaskType.String(), tsk.ID, messageID))

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
	//b.Telegram().DeleteMessages(chatID, []int{messageID})

	tsk, err := b.Task().CleanupRun(actionID)
	if err != nil {
		_ = b.Log().Error("", "", "cleanupRunTaskInlineHandler: "+err.Error())
		b.Telegram().SendError(chatID, TaskNotFound, nil)
		return
	}

	b.Telegram().ToQueue(&telegram.Message{
		Message: tgbotapi.NewCallbackWithAlert(queryID, "Выполнение "+tsk.Name+" очищено"),
		UserId:  chatID,
	})

	/*b.Telegram().ToQueue(&telegram.Message{
		Message: tgbotapi.MessageConfig{
			BaseChat: tgbotapi.BaseChat{
				ChatID:      chatID,
				ReplyMarkup: getTaskKeyboard(tsk),
			},
			Text:                  fmt.Sprintf("Выполнение очищено!!!\n\n"+TaskTemplate, tsk.Name),
			ParseMode:             tgbotapi.ModeMarkdown,
			DisableWebPagePreview: true,
		},
		UserId: chatID,
	})*/

}

func (b *botData) changePriorityTaskInlineHandler(chatID int64, messageID int, queryID, actionID string) {

	//b.Telegram().DeleteMessages(chatID, []int{messageID})

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

	b.getTaskInlineHandler(chatID, messageID, tsk.ID)
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
