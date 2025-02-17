package main

import (
	"errors"
	"fmt"
	"github.com/sparrowganz/teleFly/telegram"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"runtime/debug"
	"strconv"
	"strings"
	"sync"
	"telegram_boxes/services/admin/app/admins"
	sLog "telegram_boxes/services/admin/app/log"
	"telegram_boxes/services/admin/app/servers"
	"telegram_boxes/services/admin/app/task"
	"telegram_boxes/services/admin/app/types"
	"telegram_boxes/services/admin/bot"
	"telegram_boxes/services/admin/protobuf"
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
	telegramSender, errSender := telegram.Create(isDebug, wg, os.Getenv("ADMIN_TOKEN"), 10.0)
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

	tasks , errCreateTasks := task.CreateTasks(os.Getenv("CORE_HOST"), os.Getenv("CORE_PORT"))
	if errCreateTasks != nil {
		_ = logger.System(errCreateTasks.Error())
		return
	}

	srv , errCreateServers := servers.CreateServers(os.Getenv("CORE_HOST"), os.Getenv("CORE_PORT"))
	if errCreateServers != nil {
		_ = logger.System(errCreateServers.Error())
		return
	}

	sender := bot.CreateBot(a, telegramSender, logger)
	sender.Methods().SetTasks(tasks)
	sender.Methods().SetTypes(types.CreateType())
	sender.Methods().SetServers(srv)

	wg.Add(1)
	go func() {
		defer recovery(logger)
		defer wg.Done()

		sender.StartReadErrors()
	}()

	wg.Add(5)
	for i := 0 ; i < 5 ; i ++ {
		go func() {
			defer recovery(logger)
			defer wg.Done()

			sender.StartHandle()
		}()
	}

	lis, errCreateConn := net.Listen("tcp", fmt.Sprintf(":%s", os.Getenv("ADMIN_PORT")))
	if errCreateConn != nil {
		_ = logger.System(fmt.Sprintf("failed to listen: %v", err))
		return
	}

	GRPCServer := grpc.NewServer(
		grpc.UnaryInterceptor(logger.Interceptor),
	)
	protobuf.RegisterAdminServer(GRPCServer, protobuf.CreateAdminService(sender))

	_ = logger.System(fmt.Sprintf("ADMIN started on  :%s", os.Getenv("ADMIN_PORT")))
	err = GRPCServer.Serve(lis)
	if err != nil {
		_ = logger.System(fmt.Sprintf("failed to serve: %s" + err.Error()))
	}

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
