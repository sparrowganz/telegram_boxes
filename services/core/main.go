package main

import (
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
	"telegram_boxes/services/core/app/admin"
	"telegram_boxes/services/core/app/box"
	"telegram_boxes/services/core/app/db"
	slog "telegram_boxes/services/core/app/log"
	"telegram_boxes/services/core/protobuf"
)

func main() {

	logger, errLogger := slog.CreateLogger(os.Getenv("LOGS_HOST"), os.Getenv("LOGS_PORT"))
	if errLogger != nil {
		log.Fatal(errLogger)
		return
	}

	dbConnect, err := db.InitDatabaseConnect(
		os.Getenv("MONGO_HOST"), os.Getenv("MONGO_PORT"),
		os.Getenv("MONGO_USERNAME"), os.Getenv("MONGO_PASSWORD"),
		os.Getenv("MONGO_DATABASE"), os.Getenv("MONGO_MECHANISM"),
	)
	if err != nil {
		_ = logger.System(err.Error())
		return
	}
	defer dbConnect.Close()

	adminClient, errAdmin := admin.CreateClient(os.Getenv("ADMIN_HOST"), os.Getenv("ADMIN_PORT"))
	if errAdmin != nil {
		_ = logger.System(errAdmin.Error())
		return
	}

	//Create new server
	s := protobuf.CreateServer(dbConnect, logger, adminClient, box.CreateClients(dbConnect))

	//
	defer recovery(s.Log())
	go waitForShutdown(s)

	lis, errCreateConn := net.Listen("tcp", fmt.Sprintf(":%s", os.Getenv("CORE_PORT")))
	if errCreateConn != nil {
		_ = logger.System(fmt.Sprintf("failed to listen: %v", err))
		return
	}

	GRPCServer := grpc.NewServer(
		grpc.UnaryInterceptor(logger.Interceptor),
	)
	protobuf.RegisterServersServer(GRPCServer, s)
	protobuf.RegisterTasksServer(GRPCServer, s)

	_ = logger.System(fmt.Sprintf("Protobuf CORE started on  :%s", os.Getenv("CORE_PORT")))
	err = GRPCServer.Serve(lis)
	if err != nil {
		_ = logger.System(fmt.Sprintf("failed to serve: %s" + err.Error()))
	}

	return
}

func waitForShutdown(b protobuf.MainServer) {
	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Block until we receive our signal.
	<-interruptChan

	session := b.DB().GetMainSession().Clone()
	defer session.Close()

	servers , _ := b.DB().Models().Bots().GetAll(session)
	for _ , s := range servers {
		s.Status = protobuf.Status_Fatal.String()
		b.DB().Models().Bots().UpdateBot(s,session)
		_ = b.Admin().SendError(s.Status, s.UserName, "Shutdown core")
	}

	os.Exit(0)
}

//Recovery application out of panic
func recovery(l slog.Client) {
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

		_ = l.System("RECOVERY :" + err.Error() + "\n" + string(debug.Stack()))
	}
}
