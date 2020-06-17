package bot

import (
	"github.com/sparrowganz/teleFly/telegram/actions"
	"strings"
	"telegram_boxes/services/box/app/config"
	"telegram_boxes/services/box/app/models"
	"telegram_boxes/services/box/app/output"
	"time"
)

type OutputData struct {
	PaymentGW string
	Data      string
}

func (b *botData) IsOutputGWButton(tp string) bool {
	return strings.Contains(tp, config.OutputGWType.ToString())
}

func (b *botData) output(telegramID int64, tp config.KeyboardType) (text string, keyb interface{}, err error) {
	session := b.Database().GetMainSession().Clone()
	defer session.Close()

	var currentUser models.User
	currentUser, err = b.Database().Models().Users().FindUserByTelegramID(telegramID, session)
	if err != nil {
		return
	}

	out, errOutput := b.Output().GetOutput(currentUser.ID())
	if errOutput != nil {

		text = b.GetOutputText()
		_, keyb = b.GetOutputKeyboard()

	} else {
		text = b.GetCurrentOutputText(out.Cost, out.PaymentGateway, out.Data, out.Timestamp.Add(time.Hour*24*5))
		if tp == config.Inline {
			keyb = b.GetCancelKeyboard(config.Inline)
		}
	}

	return
}

func (b *botData) chooseOutputGW(telegramID int64, tp config.KeyboardType,
	nameGW string) (text string, keyb interface{}, err error) {

	session := b.Database().GetMainSession().Clone()
	defer session.Close()

	var currentUser models.User
	currentUser, err = b.Database().Models().Users().FindUserByTelegramID(telegramID, session)
	if err != nil {
		return
	}

	if !currentUser.Verified() {
		text = b.GetNotVerifiedOutputText(b.getReferralLink(telegramID))
		return
	}

	if currentUser.Balance().Bot() < float64(b.Config().Counts().MinOutput) {
		text = b.GetNotMinOutputText()
		return
	}

	b.Telegram().Actions().New(telegramID, actions.NewJob(AddAction.String(), OutputType.String(), &output.Output{
		PaymentGateway: nameGW,
	}, 0, false))

	text = b.GetSettingDataOutputText(nameGW)
	keyb = b.GetCancelKeyboard(tp)

	return
}

func (b *botData) setOutputData(telegramID int64, data string, out *output.Output) (text string, keyb interface{}, err error) {

	session := b.Database().GetMainSession().Clone()
	defer session.Close()

	var currentUser models.User
	currentUser, err = b.Database().Models().Users().FindUserByTelegramID(telegramID, session)
	if err != nil {
		return
	}

	b.Output().Set(currentUser.ID(), out.PaymentGateway, data, int(currentUser.Balance().Bot()), currentUser.GetAll())

	b.Telegram().Actions().Delete(telegramID)

	text = b.GetFinalOutputText(int(currentUser.Balance().Bot()), out.PaymentGateway, data, time.Now().Add(time.Hour*24*5))
	keyb = b.GetMainKeyboard()

	currentUser.Balance().SetBot(0)
	_ = b.Database().Models().Users().UpdateUser(currentUser, session)

	return
}
