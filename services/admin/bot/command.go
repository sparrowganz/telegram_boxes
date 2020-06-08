package bot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sparrowganz/teleFly/telegram"
	"github.com/sparrowganz/teleFly/telegram/actions"
	"telegram_boxes/services/admin/app/servers"
	"telegram_boxes/services/admin/app/task"
)

type Command string

var (
	start       Command = "start"
	tasks       Command = "task"
	addtask     Command = "addtask"
	statistics  Command = "statistics"
	diagnostics Command = "diagnostics"
	bonus       Command = "bonus"
	broadcast   Command = "broadcast"
)

func (b *botData) commandValidation(message *tgbotapi.Message) {
	switch Command(message.Command()) {
	case start:
		b.startCommandHandler(message.Chat.ID)
	case tasks:
		b.tasksCommandHandler(message.Chat.ID)
	case addtask:
		b.addTaskCommandHandler(message.Chat.ID)
	case statistics:
		b.statisticsCommandHandler(message.Chat.ID, false)
	case diagnostics:
		b.diagnosticsCommandHandler(message.Chat.ID)
	case bonus:
		b.bonusCommandHandler(message.Chat.ID)
	case broadcast:
		b.broadcastCommandHandler(message.Chat.ID)
	default:
		b.sendUnknownCommand(message.Chat.ID)
	}
}

func (b *botData) startCommandHandler(chatID int64) {
	b.Telegram().ToQueue(
		&telegram.Message{
			Message: tgbotapi.NewMessage(chatID, "Добро пожаловать (?? username), вы имеете доступ к боту (список команд ??)"),
			UserId:  chatID,
		})
	return
}

func (b *botData) tasksCommandHandler(chatID int64) {

	allTasks := b.Task().GetAllTasks()

	var (
		txt  string
		keyb interface{}
	)
	if len(allTasks) == 0 {
		txt = "Задания не найдены\nИспользуйте команду /addtask для добавления"
	} else {
		txt = "Информация по заданиям:"
		keyb = getTasksKeyboard(allTasks)
	}

	b.Telegram().ToQueue(
		&telegram.Message{
			Message: tgbotapi.MessageConfig{
				BaseChat: tgbotapi.BaseChat{
					ChatID:      chatID,
					ReplyMarkup: keyb,
				},
				Text:                  txt,
				ParseMode:             tgbotapi.ModeMarkdown,
				DisableWebPagePreview: true,
			},
			UserId: chatID,
		})

	return
}

func (b *botData) addTaskCommandHandler(chatID int64) {
	if j, ok := b.Telegram().Actions().Get(chatID); ok {
		b.Telegram().DeleteMessages(chatID, j.GetMessageIDs())
		b.Telegram().Actions().Delete(chatID)
	}
	b.Telegram().Actions().New(chatID,
		actions.NewJob(AddAction.String(), TaskType.String(), &task.Task{}, 0, false),
	)

	b.Telegram().ToQueue(
		&telegram.Message{
			Message: tgbotapi.MessageConfig{
				BaseChat: tgbotapi.BaseChat{
					ChatID:      chatID,
					ReplyMarkup: getTypesKeyboard(b.Types().GetAllTypes()),
				},
				Text:                  "Выберите тип",
				ParseMode:             tgbotapi.ModeMarkdown,
				DisableWebPagePreview: true,
			},
			UserId: chatID,
		})
}

func (b *botData) statisticsCommandHandler(chatID int64, isFake bool) {

	var (
		data []*servers.Count
		keyb interface{}
	)

	if !isFake {
		data = b.Servers().GetAllUsersCount()
		keyb, _ = fakeDataButton().ToKeyboard()
	} else {
		data = b.Servers().GetAllUsersFakeCount()
	}

	var txt string

	for _, status := range data {
		txt += fmt.Sprintf("%s (%v)\n", status.Username, status.All)
		txt += fmt.Sprintf("Активных: %v\n", status.All-status.Blocked)
		txt += fmt.Sprintf("Заблокировали: %v\n", status.Blocked)
		txt += fmt.Sprintf("Пользуются сейчас: %v\n", status.UseNow)
		txt += "\n\n"
	}

	b.Telegram().ToQueue(
		&telegram.Message{
			Message: tgbotapi.MessageConfig{
				BaseChat: tgbotapi.BaseChat{
					ChatID:      chatID,
					ReplyMarkup: keyb,
				},
				Text:                  txt,
				ParseMode:             tgbotapi.ModeMarkdown,
				DisableWebPagePreview: false,
			},
			UserId: chatID,
		})
	return
}

func (b *botData) broadcastCommandHandler(chatID int64) {
	b.Telegram().ToQueue(
		&telegram.Message{
			Message: tgbotapi.NewMessage(chatID, "Здесь будет отображены текущие рассылки и создание новой"),
			UserId:  chatID,
		})
	return
}

func (b *botData) diagnosticsCommandHandler(chatID int64) {

	sStatus := b.Servers().GetAllServersStatus()

	var txt string

	for _, status := range sStatus {
		txt += status.Username + " (" + status.Status.String() + ")\n"
	}

	keyb, _ := hardCheckButton().ToKeyboard()

	b.Telegram().ToQueue(
		&telegram.Message{
			Message: tgbotapi.MessageConfig{
				BaseChat: tgbotapi.BaseChat{
					ChatID:      chatID,
					ReplyMarkup: keyb,
				},
				Text:                  txt,
				ParseMode:             tgbotapi.ModeMarkdown,
				DisableWebPagePreview: false,
			},
			UserId: chatID,
		})
	return
}

func (b *botData) bonusCommandHandler(chatID int64) {
	b.tSender.ToQueue(
		&telegram.Message{
			Message: tgbotapi.NewMessage(chatID, "Здесь будет система управления бонусами"),
			UserId:  chatID,
		})
	return
}

func (b *botData) sendUnknownCommand(chatID int64) {
	b.tSender.ToQueue(
		&telegram.Message{
			Message: tgbotapi.NewMessage(chatID, "Некорректная команда, попробуйте снова"),
			UserId:  chatID,
		})
}
