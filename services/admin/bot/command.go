package bot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sparrowganz/teleFly/telegram"
	"github.com/sparrowganz/teleFly/telegram/actions"
	"telegram_boxes/services/admin/protobuf/services/core/protobuf"
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

	allTasks ,err := b.Task().GetAllTasks()
	if err != nil {
		b.Telegram().SendError(chatID, err.Error(), nil)
		return
	}

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
				DisableWebPagePreview: false,
			},
			UserId: chatID,
		})

	return
}

func (b *botData) addTaskCommandHandler(chatID int64) {
	if j, ok := b.Telegram().Actions().Get(chatID); ok {
		b.Telegram().DeleteMessages(chatID, j.GetMessageIDs())
		j.FlushMessageId()
		b.Telegram().Actions().Delete(chatID)
	}
	b.Telegram().Actions().New(chatID,
		actions.NewJob(AddAction.String(), TaskType.String(), &protobuf.Task{}, 0, false),
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
				DisableWebPagePreview: false,
			},
			UserId: chatID,
		})
}

func (b *botData) statisticsCommandHandler(chatID int64, isFake bool) {

	var (
		data []*protobuf.Counts
		err  error
		keyb interface{}
	)

	if !isFake {
		data, err = b.Servers().GetAllUsersCount(false)
		keyb, _ = fakeDataButton().ToKeyboard()
	} else {
		data, err = b.Servers().GetAllUsersCount(true)
	}

	if err != nil {
		b.Telegram().SendError(chatID, err.Error(), nil)
		return
	}

	if len(data) == 0 {
		b.Telegram().ToQueue(
			&telegram.Message{
				Message: tgbotapi.MessageConfig{
					BaseChat: tgbotapi.BaseChat{
						ChatID:      chatID,
					},
					Text:                  "Телеграмм боты не найдены",
					ParseMode:             tgbotapi.ModeMarkdown,
					DisableWebPagePreview: false,
				},
				UserId: chatID,
			})
		return
	}

	var txt string

	for _, status := range data {
		txt += fmt.Sprintf("%s %v (%v)\n", status.GetUsername(),
			status.GetNew().GetAll(), status.GetNew().GetAll()-status.GetOld().GetAll())
		txt += fmt.Sprintf("Активных: %v\n", status.GetNew().GetAll()-status.GetNew().GetBlocked())
		txt += fmt.Sprintf("Заблокировали: %v (%v)\n", status.GetNew().GetBlocked(),
			status.GetNew().GetBlocked()-status.GetOld().GetBlocked())
		txt += fmt.Sprintf("Пользуются сейчас: %v\n", status.GetCurrent())
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

	var isSetBroadcast bool
	stats, err := b.Servers().GetAllBroadcasts()
	if err == nil && len(stats) > 0 {
		isSetBroadcast = true
	}

	b.Telegram().ToQueue(
		&telegram.Message{
			Message: tgbotapi.MessageConfig{
				BaseChat: tgbotapi.BaseChat{
					ChatID:      chatID,
					ReplyMarkup: getMainBroadcastKeyboard(isSetBroadcast),
				},
				Text:                  "Меню управления рассылок: ",
				ParseMode:             tgbotapi.ModeMarkdown,
				DisableWebPagePreview: false,
			},
			UserId: chatID,
		})

	b.Telegram().Actions().New(chatID, actions.NewJob(
		AddAction.String(),BroadcastType.String() , &protobuf.StartBroadcastRequest{}, 0, true))
	return
}

func (b *botData) diagnosticsCommandHandler(chatID int64) {

	servers, err := b.Servers().GetAllServers()
	if err != nil {
		b.Telegram().SendError(chatID, "Что-то пошло не так попробуйте снова", nil)
		return
	}

	if len(servers) == 0 {
		b.Telegram().ToQueue(
			&telegram.Message{
				Message: tgbotapi.MessageConfig{
					BaseChat: tgbotapi.BaseChat{
						ChatID: chatID,
					},
					Text:                  "Телеграмм боты не найдены",
					ParseMode:             tgbotapi.ModeMarkdown,
					DisableWebPagePreview: false,
				},
				UserId: chatID,
			})
		return
	}

	var txt string

	for _, status := range servers {
		txt += status.GetUsername() + " (" + status.GetStatus() + ")\n"
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

	serv, err := b.Servers().GetAllServers()
	if err != nil {
		b.Telegram().SendError(chatID, "Что-то пошло не так попробуйте снова", nil)
		return
	}

	if len(serv) == 0 {
		b.Telegram().ToQueue(
			&telegram.Message{
				Message: tgbotapi.MessageConfig{
					BaseChat: tgbotapi.BaseChat{
						ChatID:      chatID,
					},
					Text:                  "Телеграмм боты не найдены",
					ParseMode:             tgbotapi.ModeMarkdown,
					DisableWebPagePreview: false,
				},
				UserId: chatID,
			})
		return
	}

	b.Telegram().ToQueue(
		&telegram.Message{
			Message: tgbotapi.MessageConfig{
				BaseChat: tgbotapi.BaseChat{
					ChatID:      chatID,
					ReplyMarkup: getBonusServersKeyboard(serv),
				},
				Text:                  "Выберите бот для управлением бонуса:",
				ParseMode:             tgbotapi.ModeMarkdown,
				DisableWebPagePreview: false,
			},
			UserId: chatID,
		})
	return
}

func (b *botData) sendUnknownCommand(chatID int64) {
	b.Telegram().ToQueue(
		&telegram.Message{
			Message: tgbotapi.NewMessage(chatID, "Некорректная команда, попробуйте снова"),
			UserId:  chatID,
		})
}
