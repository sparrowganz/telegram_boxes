package main

import (
	"errors"
	"fmt"
	"github.com/sparrowganz/teleFly/telegram"
	"log"
	"os"
	"path/filepath"
	"runtime/debug"
	"sync"
	"telegram_boxes/services/box/app/config"
	"telegram_boxes/services/box/app/db"
	sLog "telegram_boxes/services/box/app/log"
	"telegram_boxes/services/box/app/output"
	"telegram_boxes/services/box/app/servers"
	"telegram_boxes/services/box/app/task"
	"telegram_boxes/services/box/bot"
)

func main() {

	var isDebug bool
	if os.Getenv("APP_MODE") == "debug" {
		isDebug = true
	}

	logger, err := sLog.CreateLogger(isDebug,
		os.Getenv("LOGS_HOST"), os.Getenv("LOGS_PORT"), os.Getenv("BOT_USERNAME"))
	if err != nil {
		log.Fatal(err)
	}

	dbConnect, errInitDB := db.InitDatabaseConnect(
		os.Getenv("MONGO_HOST"), os.Getenv("MONGO_PORT"),
		os.Getenv("MONGO_USERNAME"), os.Getenv("MONGO_PASSWORD"),
		os.Getenv("MONGO_DATABASE"), os.Getenv("MONGO_MECHANISM"),
	)
	if errInitDB != nil {
		_ = logger.System(errInitDB.Error())
		return
	}

	conf, errRead := config.CreateConfig(
		filepath.Join(os.Getenv("BOX_PATH_DIR"), os.Getenv("NAME_CONFIG_FILE")))
	if errRead != nil {
		_ = logger.System(errRead.Error())
		return
	}

	wg := &sync.WaitGroup{}
	telegramSender, errSender := telegram.Create(isDebug, wg, os.Getenv("TOKEN"), 10.0)
	if errSender != nil {
		_ = logger.System(errSender.Error())
		return
	}

	sender := bot.CreateBot(dbConnect, telegramSender, logger, os.Getenv("BOT_USERNAME"))

	servData, errCreateServers := servers.CreateServers(
		os.Getenv("CORE_HOST"),
		os.Getenv("CORE_PORT"),
		sender.Methods().Username(),
	)
	if errCreateServers != nil {
		_ = logger.System(errCreateServers.Error())
		return
	}

	sender.Methods().SetServers(servData)
	sender.Methods().SetTasks(task.CreateTasks())
	sender.Methods().SetOutput(output.CreateOutput(sender.Methods().Servers().ID()))

	sender.Methods().SetConfig(conf)

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

	_ = logger.System("Start @" + sender.Methods().Username() + " bot")

	//sender.Methods().Servers().Init()
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
