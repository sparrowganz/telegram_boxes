package bot

import (
	"telegram_boxes/services/box/app/config"
	"telegram_boxes/services/box/app/models"
)

func (b *botData) balance(telegramID int64, tp config.KeyboardType) (text string, keyb interface{}, err error) {

	session := b.Database().GetMainSession().Clone()
	defer session.Close()

	var currentUser models.User
	currentUser, err = b.Database().Models().Users().FindUserByTelegramID(telegramID, session)
	if err != nil {
		return
	}

	countInvitedUsers := b.Database().Models().Users().GetCountInvitedUsers(currentUser.ID(), session)

	text = b.BalanceText(int(currentUser.Balance().Bot()), int(currentUser.Balance().Payment()),countInvitedUsers)

	if tp == config.Inline {
		keyb = b.GetCancelKeyboard(config.Inline)
	}

	return
}
