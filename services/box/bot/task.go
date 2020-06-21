package bot

import (
	"errors"
	"fmt"
	"telegram_boxes/services/box/app/models"
	"telegram_boxes/services/box/app/servers"
)

func (b *botData) getTask(chatID int64) (text string, keyb interface{}, err error) {
	session := b.Database().GetMainSession().Clone()
	defer session.Close()

	var currentUser models.User
	currentUser, err = b.Database().Models().Users().FindUserByTelegramID(chatID, session)
	if err != nil {
		return
	}

	task, errTask := b.Task().GetTask(currentUser.GetAll())
	if errTask != nil {
		text = b.TasksNotFoundText()
		return
	}

	kind, ok := b.Config().Kinds()[task.Type]
	if !ok {
		b.Servers().SendError("unknown task type "+task.Type, servers.OK)
		err = errors.New("unknown task type " + task.Type)
		return
	}

	text = b.TaskText(kind.Cost, kind.Condition)
	keyb = b.GetTaskKeyboard(task.Link, task.Type, task.ID)

	return

}

func (b *botData) checkTask(chatID int64, taskID string) (text string, keyb interface{}, err error) {
	session := b.Database().GetMainSession().Clone()
	defer session.Close()

	var currentUser models.User
	currentUser, err = b.Database().Models().Users().FindUserByTelegramID(chatID, session)
	if err != nil {
		return
	}

	if currentUser.Check(taskID) {
		text = b.TaskIsAlreadyCheck()
		keyb = b.NextTaskKeyboard()
		return
	}

	task, errFindTask := b.Task().FindTask(taskID)
	if errFindTask != nil {
		return "", nil, errFindTask
	}

	kind, ok := b.Config().Kinds()[task.Type]
	if !ok {
		b.Servers().SendError("unknown task type "+task.Type, servers.OK)
		err = errors.New("unknown task type " + task.Type)
		return
	}

	isCheckTask, errCheck := b.Task().CheckTask(chatID, taskID)
	if errCheck != nil {
		return "", nil, errCheck
	}
	if !isCheckTask {

		//ERROR CHECK
		text = b.ErrorCheckTask(kind.Cost, kind.Condition)
		keyb = b.GetTaskKeyboard(task.Link, task.Type, task.ID)
		return
	}

	//Success
	text = b.SuccessCheckTask(kind.Cost)
	keyb = b.NextTaskKeyboard()

	currentUser.Balance().AddBot(kind.Cost)
	currentUser.Add(taskID, models.StatusCheck)
	_ = b.Database().Models().Users().UpdateUser(currentUser, session)
	return

}

func (b *botData) skipTask(chatID int64, taskID string) (text string, keyb interface{}, err error) {
	session := b.Database().GetMainSession().Clone()
	defer session.Close()

	var currentUser models.User
	currentUser, err = b.Database().Models().Users().FindUserByTelegramID(chatID, session)
	if err != nil {
		return
	}

	//Success
	text = b.SkipTaskText()
	keyb = b.NextTaskKeyboard()

	currentUser.Add(taskID, models.StatusSkip)
	fmt.Println(b.Database().Models().Users().UpdateUser(currentUser, session))
	return
}
