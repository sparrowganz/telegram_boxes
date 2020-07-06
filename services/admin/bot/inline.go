package bot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sparrowganz/teleFly/telegram"
	"github.com/sparrowganz/teleFly/telegram/actions"
	"github.com/sparrowganz/teleFly/telegram/limits"
	"telegram_boxes/services/admin/app/servers"
	"telegram_boxes/services/admin/protobuf/services/core/protobuf"
	"time"
)

func (b *botData) inlineValidation(update *tgbotapi.CallbackQuery) {
	if update.Data == "" {
		return
	}

	callback := telegram.ParseCallBack(update.Data)
	switch callback.Type() {
	case BroadcastType:
		switch callback.Action() {
		case StopAction:
			b.stopBroadcast(update.Message.Chat.ID, update.ID, update.Message.MessageID, callback.ID())
		case SendAction:
			b.sendBroadcastHandler(update.Message.Chat.ID, update.ID, update.Message.MessageID)
		case DeleteAction:
			switch callback.ID() {
			case ButtonID:
				b.cancelButtonID(update.Message.Chat.ID, update.Message.MessageID)
			}
		case GetAction:
			switch callback.ID() {
			case "all":
				b.getListBroadcastsHandler(update.Message.Chat.ID, update.Message.MessageID)
			default:
				b.getBroadcastHandler(update.Message.Chat.ID, update.Message.MessageID, callback.ID())
			}
		case ChooseAction:
			b.chooseBoxBroadcastsHandler(update.Message.Chat.ID, update.Message.MessageID, callback.ID())
		case AddAction:
			switch callback.ID() {
			case ButtonID:
				b.addButtonBroadcast(update.Message.Chat.ID, update.Message.MessageID)
			default:
				job, ok := b.Telegram().Actions().Get(update.Message.Chat.ID)
				if !ok {
					b.Telegram().DeleteMessages(update.Message.Chat.ID, []int{update.Message.MessageID})
					_ = b.Log().Error("", "", "inlineValidation: job not found")
					b.Telegram().SendError(update.Message.Chat.ID, "Что-то пошло не так попробуйте снова", nil)
					return
				}

				data, isNormalData := job.GetData().(*protobuf.StartBroadcastRequest)
				if !isNormalData {
					b.Telegram().DeleteMessages(update.Message.Chat.ID, append(job.GetMessageIDs(), update.Message.MessageID))
					job.FlushMessageId()

					b.Telegram().Actions().Delete(update.Message.Chat.ID)

					_ = b.Log().Error("", "", "inlineValidation: job not found")
					b.Telegram().SendError(update.Message.Chat.ID, "Что-то пошло не так попробуйте снова", nil)
					return
				}

				if len(data.GetBotIDs()) == 0 {
					b.broadcastBotsHandler(update.Message.Chat.ID, update.Message.MessageID, data)
				} else {
					b.Telegram().DeleteMessages(update.Message.Chat.ID, []int{update.Message.MessageID})
					b.broadcastSetData(update.Message.Chat.ID, data)
				}
			}
		}
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
				job.FlushMessageId()

				b.Telegram().Actions().Delete(update.Message.Chat.ID)

				_ = b.Log().Error("", "", "inlineValidation: action or type is invalid")
				b.Telegram().SendError(update.Message.Chat.ID, "Что-то пошло не так попробуйте снова", nil)
				return
			}

			data, isNormalData := job.GetData().(*protobuf.Task)
			if !isNormalData {
				b.Telegram().DeleteMessages(update.Message.Chat.ID, []int{update.Message.MessageID})
				_ = b.Log().Error("", "", "inlineValidation: job not found")
				b.Telegram().SendError(update.Message.Chat.ID, "Что-то пошло не так попробуйте снова", nil)
				return
			}

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
			job.FlushMessageId()

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

					data, isNormalData := job.GetData().(string)
					if !isNormalData {
						b.Telegram().DeleteMessages(update.Message.Chat.ID, []int{update.Message.MessageID})
						_ = b.Log().Error("", "", "inlineValidation: job not found")
						b.Telegram().SendError(update.Message.Chat.ID, "Что-то пошло не так попробуйте снова", nil)
						return
					}

					b.getTaskInlineHandler(
						update.Message.Chat.ID, update.Message.MessageID, data)
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
//				BROADCAST
//

func (b *botData) stopBroadcast(chatID int64, queryID string, messageID int, id string) {
	b.Telegram().DeleteMessages(chatID, []int{messageID})
	err := b.Servers().StopBroadcast(id)
	if err != nil {
		_ = b.Log().Error("", "", "stopBroadcast: "+err.Error())
		b.Telegram().SendError(chatID, err.Error(), nil)
		return
	}

	b.Telegram().ToQueue(
		&telegram.Message{
			Message: tgbotapi.NewCallbackWithAlert(queryID, "Рассылка остановлена"),
			UserId:  chatID,
		})
	return
}

func (b *botData) getBroadcastHandler(chatID int64, messageID int, id string) {
	stats, err := b.Servers().GetAllBroadcasts()
	if err != nil {
		_ = b.Log().Error("", "", "getListBroadcastsHandler: "+err.Error())
		b.Telegram().SendError(chatID, err.Error(), nil)
		return
	}

	var botsData = map[string]string{}

	var txt = "Текущие рассылки:\n"

	for _, stat := range stats {
		if stat.BotID == id {

			tm := time.Unix(0, stat.GetTime()).Format("02.Jan.06 15:04")
			txt += fmt.Sprintf("%v - %v / %v\n", tm, stat.Success, stat.Fail)

			botsData[stat.Id] = tm
		}
	}

	b.Telegram().ToQueue(
		&telegram.Message{
			Message: tgbotapi.EditMessageTextConfig{
				BaseEdit: tgbotapi.BaseEdit{
					ChatID:      chatID,
					MessageID:   messageID,
					ReplyMarkup: actionsBroadcastBot(botsData),
				},
				Text:                  txt,
				ParseMode:             tgbotapi.ModeMarkdown,
				DisableWebPagePreview: false,
			},
			UserId: chatID,
		})
	return

}

func (b *botData) getListBroadcastsHandler(chatID int64, messageID int) {

	stats, err := b.Servers().GetAllBroadcasts()
	if err != nil {
		_ = b.Log().Error("", "", "getListBroadcastsHandler: "+err.Error())
		b.Telegram().SendError(chatID, err.Error(), nil)
		return
	}

	var botsHash = map[string]string{}
	for _, stat := range stats {
		botsHash[stat.BotUsername] = stat.BotID
	}

	b.Telegram().ToQueue(
		&telegram.Message{
			Message: tgbotapi.EditMessageTextConfig{
				BaseEdit: tgbotapi.BaseEdit{
					ChatID:      chatID,
					MessageID:   messageID,
					ReplyMarkup: chooseBroadcastBot(botsHash),
				},
				Text:                  `Выберите бот:`,
				ParseMode:             tgbotapi.ModeMarkdown,
				DisableWebPagePreview: false,
			},
			UserId: chatID,
		})
	return
}

func (b *botData) sendBroadcastHandler(chatID int64, queryID string, messageID int) {
	job, ok := b.Telegram().Actions().Get(chatID)
	if !ok {
		b.Telegram().DeleteMessages(chatID, []int{messageID})
		_ = b.Log().Error("", "", "addButtonBroadcast: job not found")
		b.Telegram().SendError(chatID, "Что-то пошло не так попробуйте снова", nil)
		return
	}

	b.Telegram().DeleteMessages(chatID, append(job.GetMessageIDs(), messageID))

	data, isNormalData := job.GetData().(*protobuf.StartBroadcastRequest)
	if !isNormalData {
		b.Telegram().DeleteMessages(chatID, append(job.GetMessageIDs(), messageID))
		job.FlushMessageId()

		b.Telegram().Actions().Delete(chatID)

		_ = b.Log().Error("", "", "addButtonBroadcast: job not found")
		b.Telegram().SendError(chatID, "Что-то пошло не так попробуйте снова", nil)
		return
	}
	if data.GetFileLink() != "" {
		var errGetFileLink error
		data.FileLink, errGetFileLink = b.Telegram().API().GetFileDirectURL(data.GetFileLink())
		if errGetFileLink != nil {
			_ = b.Log().Error("", "", errGetFileLink.Error())
			b.Telegram().SendError(chatID, "Что-то пошло не так попробуйте снова", nil)
			return
		}
	}

	data.ChatID = chatID
	err := b.Servers().StartBroadcast(data)
	if err != nil {
		b.Telegram().SendError(chatID, "Ошибка рассылки: "+err.Error(), nil)
		return
	}

	b.Telegram().ToQueue(
		&telegram.Message{
			Message: tgbotapi.NewCallbackWithAlert(queryID, "Рассылка началась"),
			UserId:  chatID,
		})

	b.Telegram().Actions().Delete(chatID)
}

func (b *botData) cancelButtonID(chatID int64, messageID int) {
	job, ok := b.Telegram().Actions().Get(chatID)
	if !ok {
		b.Telegram().DeleteMessages(chatID, []int{messageID})
		_ = b.Log().Error("", "", "addButtonBroadcast: job not found")
		b.Telegram().SendError(chatID, "Что-то пошло не так попробуйте снова", nil)
		return
	}

	data, isNormalData := job.GetData().(*protobuf.StartBroadcastRequest)
	if !isNormalData {
		b.Telegram().DeleteMessages(chatID, append(job.GetMessageIDs(), messageID))
		job.FlushMessageId()

		b.Telegram().Actions().Delete(chatID)

		_ = b.Log().Error("", "", "addButtonBroadcast: job not found")
		b.Telegram().SendError(chatID, "Что-то пошло не так попробуйте снова", nil)
		return
	}

	job.SetAction(AddAction.String())

	b.Telegram().DeleteMessages(chatID, []int{messageID})
	b.broadcastSetData(chatID, data)
}

func (b *botData) addButtonBroadcast(chatID int64, messageID int) {

	b.Telegram().DeleteMessages(chatID, []int{messageID})
	job, ok := b.Telegram().Actions().Get(chatID)
	if !ok {
		//b.Telegram().DeleteMessages(chatID, []int{messageID})
		_ = b.Log().Error("", "", "addButtonBroadcast: job not found")
		b.Telegram().SendError(chatID, "Что-то пошло не так попробуйте снова", nil)
		return
	}

	job.SetAction(AddButton.String())

	b.Telegram().ToQueue(
		&telegram.Message{
			Message: tgbotapi.MessageConfig{
				BaseChat: tgbotapi.BaseChat{
					ChatID:      chatID,
					ReplyMarkup: backMenu(),
				},
				Text:                  `Введите данные клавиши в формате "VK - https://vk.com"`,
				ParseMode:             tgbotapi.ModeMarkdown,
				DisableWebPagePreview: false,
			},
			UserId: chatID,
		})
	return

}

func (b *botData) broadcastSetData(chatID int64, data *protobuf.StartBroadcastRequest) {

	var txt string

	if data.Text == "" && data.FileLink == "" {
		txt = "Пришлите данные для рассылки:"
	} else if data.Text != "" {
		txt = data.Text
	}

	if len(data.Buttons) > 0 {
		txt += "\n\n"
		for _, but := range data.Buttons {
			txt += fmt.Sprintf("(%v)\n", but.Name)
		}
	}

	var mess interface{}
	if data.Type == imageType {

		if len([]rune(data.Text)) > limits.Caption() {
			txt = fmt.Sprintf("*ПРЕВЫШЕН ЛИМИТ СИМВОЛОВ НА ОПИСАНИЕ (%v/1024)*", len([]rune(data.Text)))
		}

		conf := tgbotapi.NewPhotoShare(chatID, data.FileLink)
		conf.Caption = txt
		conf.ParseMode = tgbotapi.ModeMarkdown
		conf.ReplyMarkup = addButtonBroadcastKeyboard(true)
		mess = conf
	} else if data.Type == videoType {

		if len([]rune(data.Text)) > limits.Caption() {
			txt = fmt.Sprintf("*ПРЕВЫШЕН ЛИМИТ СИМВОЛОВ НА ОПИСАНИЕ (%v/1024)*", len([]rune(data.Text)))
		}

		conf := tgbotapi.NewVideoShare(chatID, data.FileLink)
		conf.Caption = txt
		conf.ParseMode = tgbotapi.ModeMarkdown
		conf.ReplyMarkup = addButtonBroadcastKeyboard(true)
		mess = conf
	} else {

		var hasContent bool
		if data.Text != "" || data.FileLink != "" {
			hasContent = true
		}

		mess = tgbotapi.MessageConfig{
			BaseChat: tgbotapi.BaseChat{
				ChatID:      chatID,
				ReplyMarkup: addButtonBroadcastKeyboard(hasContent),
			},
			Text:                  txt,
			ParseMode:             tgbotapi.ModeMarkdown,
			DisableWebPagePreview: false,
		}
	}

	b.Telegram().ToQueue(
		&telegram.Message{
			Message: mess,
			UserId:  chatID,
		})
}

func (b *botData) chooseBoxBroadcastsHandler(chatID int64, messageID int, botID string) {
	job, ok := b.Telegram().Actions().Get(chatID)
	if !ok {
		b.Telegram().DeleteMessages(chatID, []int{messageID})
		_ = b.Log().Error("", "", "chooseBoxBroadcastsHandler: job not found")
		b.Telegram().SendError(chatID, "Что-то пошло не так попробуйте снова", nil)
		return
	}

	data, isNormalData := job.GetData().(*protobuf.StartBroadcastRequest)
	if !isNormalData {
		b.Telegram().DeleteMessages(chatID, append(job.GetMessageIDs(), messageID))
		job.FlushMessageId()
		b.Telegram().Actions().Delete(chatID)

		_ = b.Log().Error("", "", "chooseBoxBroadcastsHandler: job not found")
		b.Telegram().SendError(chatID, "Что-то пошло не так попробуйте снова", nil)
		return
	}

	var isSet bool
	for i, id := range data.BotIDs {
		if id == botID {
			data.BotIDs[i] = data.BotIDs[len(data.BotIDs)-1] // Copy last element to index i.
			data.BotIDs[len(data.BotIDs)-1] = ""             // Erase last element (write zero value).
			data.BotIDs = data.BotIDs[:len(data.BotIDs)-1]
			isSet = true
			break
		}
	}

	if !isSet {
		data.BotIDs = append(data.BotIDs, botID)
	}

	b.broadcastBotsHandler(chatID, messageID, data)
}

func (b *botData) broadcastBotsHandler(chatID int64, messageID int, data *protobuf.StartBroadcastRequest) {
	srvs, err := b.Servers().GetAllServers()
	if err != nil {
		_ = b.Log().Error("", "", "broadcastBotsHandler: "+err.Error())
		b.Telegram().SendError(chatID, err.Error(), nil)
		return
	}

	var onlineServers []*protobuf.Server
	for _, s := range srvs {
		if s.Status == servers.OK.String() {
			onlineServers = append(onlineServers, s)
		}
	}

	if len(onlineServers) == 0 {
		b.Telegram().ToQueue(
			&telegram.Message{
				Message: tgbotapi.EditMessageTextConfig{
					BaseEdit: tgbotapi.BaseEdit{
						ChatID:    chatID,
						MessageID: messageID,
					},
					Text:      "Работающие коробки для рассылки не найдены ",
					ParseMode: tgbotapi.ModeMarkdown,
				},
				UserId: chatID,
			})
		return
	}

	var txt = "Выберите коробки в которых будет рассылка"
	if len(data.BotIDs) > 0 {
		txt += "\n"
		txt += "Выбранные коробки:\n"
		for _, id := range data.BotIDs {

			for _, s := range srvs {
				if s.Id == id {
					txt += s.Username + "\n"
					break
				}
			}
		}
	}

	b.Telegram().ToQueue(
		&telegram.Message{
			Message: tgbotapi.EditMessageTextConfig{
				BaseEdit: tgbotapi.BaseEdit{
					ChatID:      chatID,
					MessageID:   messageID,
					ReplyMarkup: chooseServersKeyboard(data.BotIDs, onlineServers),
				},
				Text:      txt,
				ParseMode: tgbotapi.ModeMarkdown,
			},
			UserId: chatID,
		})
	return

}

//
//				BONUS
//

func (b *botData) changeActiveBonusHandler(chatID int64, messageID int, callbackID string) {
	err := b.Servers().ChangeActiveBonus(callbackID)
	if err != nil {
		_ = b.Log().Error("", "", "changeActiveBonusHandler: "+err.Error())
		b.Telegram().SendError(chatID, err.Error(), nil)
		return
	}
	b.chooseBonusHandler(chatID, messageID, callbackID)
}

func (b *botData) changeActiveAllBonusesHandler(chatID int64, messageID int) {
	err := b.Servers().ChangeActiveAllBonuses()
	if err != nil {
		_ = b.Log().Error("", "", "changeActiveAllBonusesHandler: "+err.Error())
		b.Telegram().SendError(chatID, err.Error(), nil)
		return
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
	act.FlushMessageId()
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
		j.FlushMessageId()
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
