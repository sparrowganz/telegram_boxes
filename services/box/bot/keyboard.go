package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sparrowganz/teleFly/telegram/keyboard"
	"strings"
	"telegram_boxes/services/box/app/config"
)

func (b *botData) GetMainKeyboard() interface{} {
	if len(b.Config().Keyboards().Buttons) == 0 {
		return nil
	}

	keyboards := b.Config().Keyboards()
	tp, keyb := keyboards.GetMain()

	return b.getKeyboard(tp, keyb)
}

func (b *botData) GetOutputKeyboard() (config.KeyboardType, interface{}) {
	if len(b.Config().Keyboards().Buttons) == 0 {
		return "", nil
	}

	keyboards := b.Config().Keyboards()
	tp, keyb := keyboards.GetOutput()

	newKeyb := b.getKeyboard(tp, keyb)

	for _, row := range keyb {
		for _, but := range row {
			but.Value = but.Text
		}
	}

	return tp, newKeyb
}

func (b *botData) GetCancelKeyboard(tp config.KeyboardType) interface{} {
	if len(b.Config().Keyboards().Buttons) == 0 {
		return nil
	}

	keyboards := b.Config().Keyboards()
	_, keyb := keyboards.GetCancel()

	return b.getKeyboard(tp, keyb)
}

func (b *botData) GetTaskKeyboard(taskURL, tp, id string) interface{} {
	if len(b.Config().Keyboards().Buttons) == 0 {
		return nil
	}

	keyboards := b.Config().Keyboards()
	keybType, keyb := keyboards.GetFromType(tp)

	for _, row := range keyb {
		for _, but := range row {
			if strings.Contains(but.Type.ToString(), "URL") {
				but.Value = taskURL
			} else {
				but.Value = id
			}
		}
	}

	return b.getKeyboard(keybType, keyb)
}

func (b *botData) NextTaskKeyboard() interface{} {
	if len(b.Config().Keyboards().Buttons) == 0 {
		return nil
	}

	keyboards := b.Config().Keyboards()
	keybType, keyb := keyboards.GetNextTask()
	return b.getKeyboard(keybType, keyb)
}

func (b *botData) GetCancelButton(tp config.KeyboardType) interface{} {
	if len(b.Config().Keyboards().Buttons) == 0 {
		return nil
	}

	name, ok := b.Config().Keyboards().Buttons[config.CancelType]
	if !ok {
		return nil
	}

	return b.getButton(tp, name, config.CancelType.ToString())
}

func (b *botData) getButton(tp config.KeyboardType, value, data string) interface{} {
	switch tp {
	case config.Inline:
		inlineButton, err := keyboard.NewButton().SetText(
			value).SetData(data).ToInline()
		if err != nil {
			return nil
		}
		return inlineButton
	case config.Static:
		tgbotapi.NewKeyboardButton(value)
	default:
		return nil
	}
	return nil
}

func (b *botData) getKeyboard(tp config.KeyboardType, keyb [][]*config.Result) interface{} {
	switch tp {
	case config.Inline:

		var tRows [][]tgbotapi.InlineKeyboardButton
		for _, rows := range keyb {

			var tRow []tgbotapi.InlineKeyboardButton
			for _, configType := range rows {

				var inlineButton tgbotapi.InlineKeyboardButton

				if strings.Contains(configType.Type.ToString(), "URL") && configType.Value != "" {

					var err error
					inlineButton, err = keyboard.NewButton().SetText(
						configType.Text).SetUrl(configType.Value).ToUrl()
					if err != nil {
						continue
					}
				} else {
					var err error
					inlineButton, err = keyboard.NewButton().SetText(
						configType.Text).SetData(configType.Type.ToString(), configType.Value).ToInline()
					if err != nil {
						continue
					}
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
				tRow = append(tRow, tgbotapi.NewKeyboardButton(but.Text))
			}

			tRows = append(tRows, tgbotapi.NewKeyboardButtonRow(tRow...))
		}

		return tgbotapi.NewReplyKeyboard(tRows...)
	}
	return nil
}
