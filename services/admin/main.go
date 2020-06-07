package main

import (
	"errors"
	"fmt"
	"github.com/sparrowganz/teleFly/telegram"
	"log"
	"os"
	"runtime/debug"
	"strconv"
	"strings"
	"sync"
	"telegram_boxes/services/admin/app/admins"
	sLog "telegram_boxes/services/admin/app/log"
	"telegram_boxes/services/admin/app/task"
	"telegram_boxes/services/admin/app/types"
	"telegram_boxes/services/admin/bot"
)

func main() {

	var isDebug bool
	if os.Getenv("APP_MODE") == "debug" {
		isDebug = true
	}

	logger, err := sLog.CreateLogger(isDebug, os.Getenv("LOGS_HOST"), os.Getenv("LOGS_PORT"))
	if err != nil {
		log.Fatal(err)
	}

	wg := &sync.WaitGroup{}
	telegramSender, errSender := telegram.Create(isDebug, wg, os.Getenv("TOKEN"), 10.0)
	if errSender != nil {
		_ = logger.System(errSender.Error())
		return
	}

	a := admins.CreateAdmin()
	adminsID := strings.Split(os.Getenv("ADMINS"), ",")
	for _, ad := range adminsID {
		id, errID := strconv.Atoi(ad)
		if errID != nil {
			continue
		}
		a.Add(int64(id))
	}

	if len(a.GetAll()) == 0 {
		_ = logger.System("len admins == 0")
		return
	}

	sender := bot.CreateBot(a, telegramSender, logger)
	sender.Methods().SetTasks(task.CreateTasks())
	sender.Methods().SetTypes(types.CreateType())

	wg.Add(1)
	go func() {
		defer recovery(logger)
		defer wg.Done()

		sender.StartReadErrors()
	}()

	wg.Add(1)
	go func() {
		defer recovery(logger)
		defer wg.Done()

		sender.StartHandle()
	}()

	_ = logger.System("Start admin bot")
	wg.Wait()
}

func recovery(l sLog.Log) {
	var err error
	r := recover()
	if r != nil {
		switch t := r.(type) {
		case string:
			err = errors.New(t)
		case error:
			err = t
		default:
			err = errors.New("Unknown error ")
		}

		_ = l.System(
			fmt.Sprintf("RECOVERY : %s \n %s", err.Error(), debug.Stack()))
	}
}
