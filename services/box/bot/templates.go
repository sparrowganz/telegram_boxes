package bot

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func (b *botData) GetErrorCommandText() string {
	return b.Config().Texts().IncorrectCommand
}

func (b *botData) GetErrorText() string {
	return b.Config().Texts().Error
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

func (b *botData) GetNotVerifiedOutputText(referralLink string) string {
	txt := b.Config().Texts().NotVerifiedOutput

	if strings.Contains(txt, "@verifiedCount") {
		txt = strings.Replace(txt, "@verifiedCount", strconv.Itoa(b.Config().Counts().VerifiedCount), -1)
	}

	if strings.Contains(txt, "@referralLink") {
		txt = strings.Replace(txt, "@referralLink", referralLink, -1)
	}

	return txt
}

func (b *botData) GetSettingDataOutputText(paymentGW string) string {
	txt := b.Config().Texts().SetDataOutput

	if strings.Contains(txt, "@paymentGW") {
		txt = strings.Replace(txt, "@paymentGW", paymentGW, -1)
	}
	return txt
}

func (b *botData) GetNotMinOutputText() string {
	txt := b.Config().Texts().OutputErrorBalance

	if strings.Contains(txt, "@minOutput") {
		txt = strings.Replace(txt, "@minOutput", strconv.Itoa(b.Config().Counts().MinOutput), -1)
	}

	return txt
}

func (b *botData) GetOutputText() string {
	return b.Config().Texts().Output
}

func (b *botData) GetErrorOutputText() string {
	return b.Config().Texts().ErrorOutput
}

func (b *botData) GetErrorTaskOutputText(link string) string {
	txt := b.Config().Texts().ErrorTaskOutput

	if strings.Contains(txt, "@link") {
		txt = strings.Replace(txt, "@link", link, -1)
	}
	return txt
}

func (b *botData) ChecksNotFoundText() string {
	return b.Config().Texts().ChecksNotFound
}

func (b *botData) GetCurrentOutputText(cost int, gateway, data string, tm time.Time) string {

	txt := b.Config().Texts().CurrentOutput

	if strings.Contains(txt, "@cost") {
		txt = strings.Replace(txt, "@cost", strconv.Itoa(cost), -1)
	}

	if strings.Contains(txt, "@gateway") {
		txt = strings.Replace(txt, "@gateway", gateway, -1)
	}

	if strings.Contains(txt, "@data") {
		txt = strings.Replace(txt, "@data", data, -1)
	}

	if strings.Contains(txt, "@time") {

		var tmString string
		res := tm.Sub(time.Now())
		if res.Hours() < 24 {
			tmString = fmt.Sprintf("%v часов", res.Hours())
		} else {
			tmString = fmt.Sprintf("%v дней %v часов", int(res.Hours()/24), int(res.Hours())-int(res.Hours()/24))
		}

		txt = strings.Replace(txt, "@time", tmString, -1)
	}
	return txt
}

func (b *botData) GetFinalOutputText(cost int, gateway, data string, tm time.Time) string {
	txt := b.Config().Texts().FinalOutput

	if strings.Contains(txt, "@currentOutput") {
		txt = strings.Replace(txt, "@currentOutput", b.GetCurrentOutputText(cost, gateway, data, tm), -1)
	}
	return txt
}

func (b *botData) TasksNotFoundText() string {
	return b.Config().Texts().TasksNotFound
}

func (b *botData) SkipTaskText() string {
	return b.Config().Texts().SkipTask
}

func (b *botData) TaskText(cost int, condition string) string {
	txt := b.Config().Texts().TaskTemplate

	if strings.Contains(txt, "@cost") {
		txt = strings.Replace(txt, "@cost", strconv.Itoa(cost), -1)
	}

	if strings.Contains(txt, "@condition") {
		txt = strings.Replace(txt, "@condition", condition, -1)
	}
	return txt
}

func (b *botData) ErrorCheckTask(cost int, condition string) string {
	txt := b.Config().Texts().TaskWrongCheck

	if strings.Contains(txt, "@taskTemplate") {
		txt = strings.Replace(txt, "@taskTemplate", b.TaskText(cost, condition), -1)
	}

	return txt
}

func (b *botData) TaskIsAlreadyCheck() string {
	return b.Config().Texts().TaskIsAlreadyCheck
}

func (b *botData) SuccessCheckTask(cost int) string {

	txt := b.Config().Texts().SuccessCheckTask

	if strings.Contains(txt, "@cost") {
		txt = strings.Replace(txt, "@cost", strconv.Itoa(cost), -1)
	}

	return txt
}

func (b *botData) GetHelpText() string {
	return b.Config().Texts().Help
}
