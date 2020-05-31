package main

import (
	"github.com/sparrowganz/teleFly/telegram"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"telegram_boxes/services/admin/app/admins"
	sLog "telegram_boxes/services/admin/app/log"
	"telegram_boxes/services/admin/bot"
)

func main() {

	var isDebug bool
	if os.Getenv("APP_MODE") == "true" {
		isDebug = true
	}

	logger, err := sLog.CreateLogger(isDebug,os.Getenv("LOGS_HOST"), os.Getenv("LOGS_PORT"))
	if err != nil {
		log.Fatal(err)
	}

	wg := &sync.WaitGroup{}
	telegramSender, errSender := telegram.Create(isDebug, wg, os.Getenv("TOKEN"), 1.0)
	if errSender != nil {
		_ = logger.System(errSender.Error())
		return
	}

	a:= admins.CreateAdmin()
	adminsID := strings.Split(os.Getenv("ADMINS"),",")
	for _ , ad := range adminsID {
		id ,errID := strconv.Atoi(ad)
		if errID != nil {
			continue
		}
		a.Add(int64(id))
	}


	if len(adminsID) == 0 {
		_ = logger.System("len admins == 0")
		return
	}

	sender := bot.CreateBot(a,telegramSender, logger)

	wg.Add(1)
	go sender.StartReadErrors(wg)

	wg.Add(1)
	go sender.StartHandle(wg)

	_ = logger.System("Start admin bot")
	wg.Wait()
}
