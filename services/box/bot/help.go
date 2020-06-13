package bot

import "telegram_boxes/services/box/app/config"

func (b *botData) help(telegramID int64, tp config.KeyboardType) (text string, keyb interface{}, err error) {

	text = b.GetHelpText()

	if tp == config.Inline {
		keyb = b.GetCancelKeyboard()
	}

	return
}
