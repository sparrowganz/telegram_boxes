package bot

import (
	"fmt"
	"telegram_boxes/services/box/app/config"
	"telegram_boxes/services/box/app/models"
)

func (b *botData) getReferralLink(chatID int64) string {
	return fmt.Sprintf("https://t.me/%v?start=%v", b.Username(), chatID)
}

func (b *botData) referrals(telegramID int64, tp config.KeyboardType) (text string, keyb interface{}, err error) {

	session := b.Database().GetMainSession().Clone()
	defer session.Close()

	var currentUser models.User
	currentUser, err = b.Database().Models().Users().FindUserByTelegramID(telegramID, session)
	if err != nil {
		return
	}

	countInvitedUsers := b.Database().Models().Users().GetCountInvitedUsers(currentUser.ID(), session)

	text = b.ReferralsText(countInvitedUsers, b.getReferralLink(telegramID), currentUser.Verified())

	if tp == config.Inline {
		keyb = b.GetCancelKeyboard(config.Inline)
	}

	return
}
