package main

import (
	"errors"
	"fmt"
	"github.com/sparrowganz/teleFly/telegram"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"runtime/debug"
	"sync"
	"syscall"
	"telegram_boxes/services/box/app"
	"telegram_boxes/services/box/app/config"
	"telegram_boxes/services/box/app/db"
	sLog "telegram_boxes/services/box/app/log"
	"telegram_boxes/services/box/app/servers"
	"telegram_boxes/services/box/app/task"
	"telegram_boxes/services/box/bot"
	boxProto "telegram_boxes/services/box/protobuf"
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
	defer dbConnect.Close()

	conf, errRead := config.CreateConfig(
		filepath.Join(os.Getenv("BOX_PATH_DIR"), os.Getenv("NAME_CONFIG_FILE")))
	if errRead != nil {
		_ = logger.System(errRead.Error())
		return
	}

	wg := &sync.WaitGroup{}
	telegramSender, errSender := telegram.Create(isDebug, wg, os.Getenv("TOKEN"), 10)
	if errSender != nil {
		_ = logger.System(errSender.Error())
		return
	}

	sender := bot.CreateBot(dbConnect, telegramSender, logger, os.Getenv("BOT_USERNAME"))
	defer recovery(sender)
	servData, errCreateServers := servers.CreateServers(
		os.Getenv("APP_IP"),
		os.Getenv("APP_PORT"),
		os.Getenv("CORE_HOST"),
		os.Getenv("CORE_PORT"),
		sender.Methods().Username(),
	)
	if errCreateServers != nil {
		_ = logger.System(errCreateServers.Error())
		return
	}

	taskData, errCreateConnection := task.CreateTasks(
		os.Getenv("CORE_HOST"),
		os.Getenv("CORE_PORT"),
		sender.Methods().Username(),
	)
	if errCreateConnection != nil {
		_ = logger.System(errCreateConnection.Error())
		return
	}

	sender.Methods().SetServers(servData)
	sender.Methods().SetTasks(taskData)

	sender.Methods().SetConfig(conf)

	wg.Add(1)
	go func() {
		defer recovery(sender)
		//defer wg.Done()

		sender.StartReadErrors()
	}()

	wg.Add(5)
	for i := 0 ; i < 5 ; i ++ {
		go func() {
			defer recovery(sender)
			defer wg.Done()

			sender.StartHandle()
		}()
	}

	go waitForShutdown(sender)

	lis, errCreateConn := net.Listen("tcp", fmt.Sprintf(":%s", os.Getenv("APP_PORT")))
	if errCreateConn != nil {
		_ = logger.System(fmt.Sprintf("failed to listen: %v", err))
		return
	}

	GRPCServer := grpc.NewServer(
		grpc.UnaryInterceptor(logger.Interceptor),
	)
	boxProto.RegisterBoxServer(GRPCServer, boxProto.CreateBoxService(sender))

	_ = logger.System(fmt.Sprintf("Protobuf %v started on  :%s",
		os.Getenv("BOT_USERNAME"), os.Getenv("APP_PORT")))
	err = GRPCServer.Serve(lis)
	if err != nil {
		_ = logger.System(fmt.Sprintf("failed to serve: %s" + err.Error()))
	}

	wg.Wait()
}

func waitForShutdown(b bot.Bot) {
	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Block until we receive our signal.
	<-interruptChan

	b.Close()

	_ = b.Methods().Servers().SendError("Сервер выключился", app.StatusFatal.String())
	os.Exit(0)
}

func recovery(b bot.Bot) {
	var err error
	r := recover()
	if r != nil {
		_ = b.Methods().Servers().SendError("Критическая ошибка", app.StatusRecovery.String())
		switch t := r.(type) {
		case string:
			err = errors.New(t)
		case error:
			err = t
		default:
			err = errors.New("Unknown error ")
		}

		_ = b.Methods().Log().System(
			fmt.Sprintf("RECOVERY : %s \n %s", err.Error(), debug.Stack()))
		_ = b.Methods().Servers().SendError("Работа восстановлена", app.StatusOK.String())
	}
}
