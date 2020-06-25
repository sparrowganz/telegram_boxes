package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sparrowganz/teleFly/telegram/keyboard"
	"telegram_boxes/services/admin/app/types"
	"telegram_boxes/services/admin/protobuf/services/core/protobuf"
)

func cancelButton() (b keyboard.Button) {
	return keyboard.NewButton().SetText("Отмена").SetData(CancelType.String())
}

func hardCheckButton() (b keyboard.Button) {
	return keyboard.NewButton().SetText("Проверить").SetData(ServerType.String(), CheckAction.String())
}

func fakeDataButton() (b keyboard.Button) {
	return keyboard.NewButton().SetText("Фейковая статистика").SetData(ServerType.String(), FakeAction.String())
}

func allServersBonusButton() (b keyboard.Button) {
	return keyboard.NewButton().SetText("Все боты").SetData(BonusType.String(), ChooseAction.String(), AllID)
}

func changeBonusKeyboard(id string, isActive bool) *tgbotapi.InlineKeyboardMarkup {

	activeB := tgbotapi.InlineKeyboardButton{}
	if !isActive {
		activeB, _ = keyboard.NewButton().SetText("Деактивировать").SetData(
			BonusType.String(), ChangeActiveAction.String(), id).ToInline()
	} else {
		activeB, _ = keyboard.NewButton().SetText("Активировать").SetData(
			BonusType.String(), ChangeActiveAction.String(), id).ToInline()
	}

	k := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(activeB))
	return &k
}

func getBonusServersKeyboard(servers []*protobuf.Server) *tgbotapi.InlineKeyboardMarkup {

	var rows [][]tgbotapi.InlineKeyboardButton
	for _, server := range servers {
		but, err := keyboard.NewButton().SetText(server.GetUsername()).SetData(
			BonusType.String(), ChooseAction.String(), server.GetId()).ToInline()
		if err != nil {
			return nil
		}
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(but))
	}

	all, _ := allServersBonusButton().ToInline()

	k := tgbotapi.NewInlineKeyboardMarkup(append(rows, tgbotapi.NewInlineKeyboardRow(all))...)
	return &k
}

func getTypesKeyboard(tps []*types.Type) *tgbotapi.InlineKeyboardMarkup {

	var row []tgbotapi.InlineKeyboardButton
	var rows [][]tgbotapi.InlineKeyboardButton

	for _, tp := range tps {

		but, err := keyboard.NewButton().SetText(tp.Name).SetData(TaskType.String(),
			AddAction.String(), tp.ID).ToInline()
		if err != nil {
			return nil
		}

		row = append(row, but)
		if len(row) == 2 {
			rows = append(rows, row)
			row = []tgbotapi.InlineKeyboardButton{}
		}
	}
	if len(row) != 0 {
		rows = append(rows, row)
	}

	cancelB, _ := cancelButton().ToInline()
	k := tgbotapi.NewInlineKeyboardMarkup(append(rows, tgbotapi.NewInlineKeyboardRow(cancelB))...)
	return &k
}

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
func getTasksKeyboard(tsks []*protobuf.Task) *tgbotapi.InlineKeyboardMarkup {

	var rows [][]tgbotapi.InlineKeyboardButton

	for _, tsk := range tsks {

		nameB, err := keyboard.NewButton().SetText(tsk.GetName()).SetData(
			TaskType.String(), GetAction.String(), tsk.GetId()).ToInline()
		if err != nil {
			continue
		}

		rows = append(rows, tgbotapi.NewInlineKeyboardRow(nameB))
	}

	keyb := tgbotapi.NewInlineKeyboardMarkup(rows...)

	return &keyb
}

//Get Task action inline keyboard
func getTaskKeyboard(tsk *protobuf.Task) *tgbotapi.InlineKeyboardMarkup {

	changePriorityButton, _ := keyboard.NewButton().SetText("Изменить приоритет").SetData(
		TaskType.String(), PriorityAction.String(), tsk.GetId()).ToInline()
	row1 := tgbotapi.NewInlineKeyboardRow(changePriorityButton)

	resetDoButton, _ := keyboard.NewButton().SetText("Очистить выполнение").SetData(
		TaskType.String(), CleanAction.String(), tsk.GetId()).ToInline()
	row2 := tgbotapi.NewInlineKeyboardRow(resetDoButton)

	deleteButton, _ := keyboard.NewButton().SetText("Удалить").SetData(
		TaskType.String(), DeleteAction.String(), tsk.GetId()).ToInline()
	row3 := tgbotapi.NewInlineKeyboardRow(deleteButton)

	keyb := tgbotapi.NewInlineKeyboardMarkup(row1, row2, row3)

	return &keyb
}
