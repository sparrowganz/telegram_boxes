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

func addButtonBroadcastKeyboard(isHasContent bool) *tgbotapi.InlineKeyboardMarkup {

	var rows [][]tgbotapi.InlineKeyboardButton

	if isHasContent {
		nextButton, _ := keyboard.NewButton().SetText("Добавить клавишу").SetData(
			BroadcastType.String(), AddAction.String(), ButtonID).ToInline()

		rows = append(rows, tgbotapi.NewInlineKeyboardRow(nextButton))

		sendButton, _ := keyboard.NewButton().SetText("Начать рассылку").SetData(
			BroadcastType.String(), SendAction.String()).ToInline()
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(sendButton))
	}

	cancelB, _ := cancelButton().ToInline()
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(cancelB))

	k := tgbotapi.NewInlineKeyboardMarkup(rows...)
	return &k
}

func chooseServersKeyboard(chooseServersID []string, servs []*protobuf.Server) *tgbotapi.InlineKeyboardMarkup {
	var row []tgbotapi.InlineKeyboardButton
	var rows [][]tgbotapi.InlineKeyboardButton

	var hasChoose bool
	for _, srv := range servs {

		var smile = "❌ "
		for _, id := range chooseServersID {
			if id == srv.Id {
				smile = "✅ "
				hasChoose = true
				break
			}
		}

		but, _ := keyboard.NewButton().SetText(smile+srv.Username).SetData(
			BroadcastType.String(), ChooseAction.String(), srv.Id).ToInline()

		row = append(row, but)
		if len(row) == 2 {
			rows = append(rows, row)
			row = []tgbotapi.InlineKeyboardButton{}
		}
	}
	if len(row) != 0 {
		rows = append(rows, row)
	}

	if hasChoose {
		nextButton, _ := keyboard.NewButton().SetText("Сохранить").SetData(
			BroadcastType.String(), AddAction.String()).ToInline()
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(nextButton))
	}
	cancelB, _ := cancelButton().ToInline()

	rows = append(rows, tgbotapi.NewInlineKeyboardRow(cancelB))

	k := tgbotapi.NewInlineKeyboardMarkup(rows...)
	return &k
}

func getMainBroadcastKeyboard(isSetBroadcast bool) *tgbotapi.InlineKeyboardMarkup {

	var rows [][]tgbotapi.InlineKeyboardButton

	if isSetBroadcast {
		getBroadcasts, _ := keyboard.NewButton().SetText("Текущие рассылки").SetData(
			BroadcastType.String(), GetAction.String(), "all").ToInline()
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(getBroadcasts))
	}

	addBroadcast, _ := keyboard.NewButton().SetText("Новая рассылка").SetData(
		BroadcastType.String(), AddAction.String()).ToInline()
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(addBroadcast))

	keyb := tgbotapi.NewInlineKeyboardMarkup(rows...)

	return &keyb
}

func backMenu() *tgbotapi.InlineKeyboardMarkup {
	cancel, _ := keyboard.NewButton().SetText("Отмена").SetData(
		BroadcastType.String(), DeleteAction.String(), ButtonID).ToInline()
	keyb := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(cancel))
	return &keyb
}

func chooseBroadcastBot(data map[string]string) *tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton

	for username, id := range data {
		bot, _ := keyboard.NewButton().SetText(username).SetData(
			BroadcastType.String(), GetAction.String(), id).ToInline()
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(bot))
	}

	cancelB, _ := cancelButton().ToInline()
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(cancelB))

	keyb := tgbotapi.NewInlineKeyboardMarkup(rows...)
	return &keyb
}

func actionsBroadcastBot(data map[string]string) *tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton

	for id, title := range data {
		bot, _ := keyboard.NewButton().SetText("Остановить от "+title).SetData(
			BroadcastType.String(), StopAction.String(), id).ToInline()
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(bot))
	}

	cancelB, _ := cancelButton().ToInline()
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(cancelB))

	keyb := tgbotapi.NewInlineKeyboardMarkup(rows...)
	return &keyb
}
