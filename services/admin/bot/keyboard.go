package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sparrowganz/teleFly/telegram"
	"github.com/sparrowganz/teleFly/telegram/keyboard"
	"telegram_boxes/services/admin/app/task"
)

const (
	TaskType       telegram.Type = "t"
	LastChoiceType telegram.Type = "lc"

	CleanAction    telegram.Action = "c"
	DeleteAction   telegram.Action = "d"
	GetAction      telegram.Action = "g"
	PriorityAction telegram.Action = "pr"

	YesID = "y"
	NoID  = "n"
)

func lastChoiceKeyboard(action string) *tgbotapi.InlineKeyboardMarkup {
	yes, err := keyboard.NewButton().SetText("Да").SetData(
		LastChoiceType.String(), action, YesID).ToInline()
	if err != nil {
		return nil
	}

	no, errCreate := keyboard.NewButton().SetText("Нет").SetData(
		LastChoiceType.String(), action, NoID).ToInline()
	if errCreate != nil {
		return nil
	}

	keyb := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(yes, no))
	return &keyb
}

//Get Tasks inline keyboard
func getTasksKeyboard(tsks []*task.Task) *tgbotapi.InlineKeyboardMarkup {

	var rows [][]tgbotapi.InlineKeyboardButton

	for _, tsk := range tsks {

		nameB, err := keyboard.NewButton().SetText(tsk.Name).SetData(
			TaskType.String(), GetAction.String(), tsk.ID).ToInline()
		if err != nil {
			continue
		}

		rows = append(rows, tgbotapi.NewInlineKeyboardRow(nameB))
	}

	keyb := tgbotapi.NewInlineKeyboardMarkup(rows...)

	return &keyb
}

//Get Task action inline keyboard
func getTaskKeyboard(tsk *task.Task) *tgbotapi.InlineKeyboardMarkup {

	changePriorityButton, _ := keyboard.NewButton().SetText("Изменить приоритет").SetData(
		TaskType.String(), PriorityAction.String(), tsk.ID).ToInline()
	row1 := tgbotapi.NewInlineKeyboardRow(changePriorityButton)

	resetDoButton, _ := keyboard.NewButton().SetText("Очистить выполнение").SetData(
		TaskType.String(), CleanAction.String(), tsk.ID).ToInline()
	row2 := tgbotapi.NewInlineKeyboardRow(resetDoButton)

	deleteButton, _ := keyboard.NewButton().SetText("Удалить").SetData(
		TaskType.String(), DeleteAction.String(), tsk.ID).ToInline()
	row3 := tgbotapi.NewInlineKeyboardRow(deleteButton)

	keyb := tgbotapi.NewInlineKeyboardMarkup(row1, row2, row3)

	return &keyb
}
