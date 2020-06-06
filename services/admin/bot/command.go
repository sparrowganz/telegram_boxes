package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sparrowganz/teleFly/telegram"
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
		b.statisticsCommandHandler(message.Chat.ID)
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
	b.tSender.ToQueue(
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


	b.tSender.ToQueue(
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
	b.tSender.ToQueue(
		&telegram.Message{
			Message: tgbotapi.NewMessage(chatID, "Здесь будут отображен функционал по добавлению заданий"),
			UserId:  chatID,
		})
	return
}

func (b *botData) statisticsCommandHandler(chatID int64) {
	b.tSender.ToQueue(
		&telegram.Message{
			Message: tgbotapi.NewMessage(chatID, "Здесь будут отображены статистика ботов"),
			UserId:  chatID,
		})
	return
}

func (b *botData) broadcastCommandHandler(chatID int64) {
	b.tSender.ToQueue(
		&telegram.Message{
			Message: tgbotapi.NewMessage(chatID, "Здесь будет отображены текущие рассылки и создание новой"),
			UserId:  chatID,
		})
	return
}

func (b *botData) diagnosticsCommandHandler(chatID int64) {
	b.tSender.ToQueue(
		&telegram.Message{
			Message: tgbotapi.NewMessage(chatID, "Здесь будет (??????)"),
			UserId:  chatID,
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
