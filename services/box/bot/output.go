package bot

import (
	"github.com/sparrowganz/teleFly/telegram/actions"
	"gopkg.in/mgo.v2"
	"strings"
	"telegram_boxes/services/box/app/config"
	"telegram_boxes/services/box/app/models"
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

	//find is set in DB output
	out, errOutput := b.Database().Models().Outputs().FindOutputByUserID(currentUser.ID().Hex(), session)
	if errOutput != nil {

		if errOutput != mgo.ErrNotFound {
			return "", nil, errOutput
		}

		text = b.GetOutputText()
		_, keyb = b.GetOutputKeyboard()

	} else if time.Now().Sub(out.Timestamp().Created()) >= time.Hour*5*24 {

		//todo delete task
		//Check all already checked tasks
		for _, task := range out.Tasks() {
			isCheck, _ := b.Task().CheckTask(currentUser.Telegram().ID(), task)
			if !isCheck {

				tsk ,_ := b.Task().FindTask(task)

				text = b.GetErrorTaskOutputText(tsk.GetLink())
				return
			}
		}

		text = b.GetErrorOutputText()

	} else {
		//View current  output text
		text = b.GetCurrentOutputText(out.Cost(), out.PaymentGW(), out.Data(), out.Timestamp().Created().Add(time.Hour*24*5))
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

	if len(currentUser.GetAllChecks()) == 0 {
		text = b.ChecksNotFoundText()
	}

	if currentUser.Balance().Bot() < b.Config().Counts().MinOutput {
		text = b.GetNotMinOutputText()
		return
	}

	b.Telegram().Actions().New(telegramID, actions.NewJob(AddAction.String(), OutputType.String(), &models.OutputData{
		PaymentGateway: nameGW,
	}, 0, true))

	text = b.GetSettingDataOutputText(nameGW)
	keyb = b.GetCancelKeyboard(tp)

	return
}

func (b *botData) setOutputData(telegramID int64, data string, out *models.OutputData) (text string, keyb interface{}, err error) {

	session := b.Database().GetMainSession().Clone()
	defer session.Close()

	var currentUser models.User
	currentUser, err = b.Database().Models().Users().FindUserByTelegramID(telegramID, session)
	if err != nil {
		return
	}

	output := models.CreateOutput(
		currentUser.ID().Hex(), out.PaymentGateway, data, currentUser.Balance().Bot(), currentUser.GetAllChecks())

	_ = b.Database().Models().Outputs().CreateOutput(output, session)

	b.Telegram().Actions().Delete(telegramID)

	text = b.GetFinalOutputText(currentUser.Balance().Bot(), out.PaymentGateway, data, time.Now().Add(time.Hour*24*5))
	keyb = b.GetMainKeyboard()

	currentUser.Balance().SetBot(0)
	currentUser.CleanChecks()
	_ = b.Database().Models().Users().UpdateUser(currentUser, session)

	return
}
