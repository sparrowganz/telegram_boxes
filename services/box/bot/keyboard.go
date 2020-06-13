package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sparrowganz/teleFly/telegram/keyboard"
	"telegram_boxes/services/box/app/config"
)

func (b *botData) GetMainKeyboard() interface{} {
	if len(b.Config().Keyboards().Buttons) == 0 {
		return nil
	}

	keyboards := b.Config().Keyboards()
	tp, keyb := keyboards.GetMain()

	return getKeyboard(tp, keyb)
}

func (b *botData) GetCancelKeyboard() interface{} {
	if len(b.Config().Keyboards().Buttons) == 0 {
		return nil
	}

	keyboards := b.Config().Keyboards()
	tp, keyb := keyboards.GetCancel()

	return getKeyboard(tp, keyb)
}

func getKeyboard(tp config.KeyboardType, keyb [][]config.Result) interface{} {
	switch tp {
	case "":
		return nil
	case config.Inline:

		var tRows [][]tgbotapi.InlineKeyboardButton
		for _, rows := range keyb {

			var tRow []tgbotapi.InlineKeyboardButton
			for _, configType := range rows {

				inlineButton, err := keyboard.NewButton().SetText(
					configType.Value).SetData(configType.Type.ToString()).ToInline()
				if err != nil {
					continue
				}
				tRow = append(tRow, inlineButton)

			}

			if len(tRow) > 0 {
				tRows = append(tRows, tRow)
			}
		}

		if len(tRows) > 0 {
			k := tgbotapi.NewInlineKeyboardMarkup(tRows...)
			return &k
		} else {
			return nil
		}

	case config.Static:

		var tRows [][]tgbotapi.KeyboardButton
		for _, rows := range keyb {

			var tRow []tgbotapi.KeyboardButton
			for _, but := range rows {
				tRow = append(tRow, tgbotapi.NewKeyboardButton(but.Value))
			}

			tRows = append(tRows, tgbotapi.NewKeyboardButtonRow(tRow...))
		}

		return tgbotapi.NewReplyKeyboard(tRows...)

	}
	return nil
}
