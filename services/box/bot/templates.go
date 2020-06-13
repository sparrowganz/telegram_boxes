package bot

import (
	"strconv"
	"strings"
)

func (b *botData) GetErrorCommandText() string {
	return b.Config().Texts().Errors.IncorrectCommand
}

func (b *botData) GetErrorText() string {
	return b.Config().Texts().Errors.Error
}

func (b *botData) GetStartText() string {
	return b.Config().Texts().StartText
}

func (b *botData) BonusForReferralText(username string, referralLink string) string {
	txt := b.Config().Texts().NotifyForReferral

	if strings.Contains(txt, "@username") {
		txt = strings.Replace(txt, "@username", username, -1)
	}

	if strings.Contains(txt, "@costForReferral") {
		txt = strings.Replace(txt, "@costForReferral", strconv.Itoa(b.Config().Counts().CostForReferral), -1)
	}

	if strings.Contains(txt, "@referralLink") {
		txt = strings.Replace(txt, "@referralLink", referralLink, -1)
	}

	return txt
}

func (b *botData) BalanceText(balance, paymentBalance, countInvitedUsers int) string {
	txt := b.Config().Texts().Balance

	if strings.Contains(txt, "@balance") {
		txt = strings.Replace(txt, "@balance", strconv.Itoa(balance), -1)
	}

	if strings.Contains(txt, "@countInvitedUsers") {
		txt = strings.Replace(txt, "@countInvitedUsers", strconv.Itoa(countInvitedUsers), -1)
	}

	if strings.Contains(txt, "@paymentsBalance") {
		txt = strings.Replace(txt, "@paymentsBalance", strconv.Itoa(paymentBalance), -1)
	}

	if strings.Contains(txt, "@costForReferral") {
		txt = strings.Replace(txt, "@costForReferral", strconv.Itoa(b.Config().Counts().CostForReferral), -1)
	}

	return txt
}

func (b *botData) ReferralsText(countInvitedUsers int, referralLink string, status bool) string {
	txt := b.Config().Texts().Referrals

	if strings.Contains(txt, "@countInvitedUsers") {
		txt = strings.Replace(txt, "@countInvitedUsers", strconv.Itoa(countInvitedUsers), -1)
	}

	if strings.Contains(txt, "@costForReferral") {
		txt = strings.Replace(txt, "@costForReferral", strconv.Itoa(b.Config().Counts().CostForReferral), -1)
	}

	if strings.Contains(txt, "@referralLink") {
		txt = strings.Replace(txt, "@referralLink", referralLink, -1)
	}

	if strings.Contains(txt, "@referralStatus") {
		txt = strings.Replace(txt, "@referralStatus", b.GetStatusReferral(status), -1)
	}

	return txt
}

func (b *botData) GetStatusReferral(status bool) string {
	if status {
		return b.Config().Texts().StatusReferralOk
	} else {
		txt := b.Config().Texts().StatusReferralFalse

		if strings.Contains(txt, "@verifiedCount") {
			txt = strings.Replace(txt, "@verifiedCount", strconv.Itoa(b.Config().Counts().VerifiedCount), -1)
		}

		return txt

	}
}

func (b *botData) GetHelpText() string {
	return b.Config().Texts().Help
}
