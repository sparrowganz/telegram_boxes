package main

import (
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"runtime/debug"
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

	//Create new server
	s := protobuf.CreateServer(dbConnect, logger)

	//
	defer recovery(s.Log())

	lis, errCreateConn := net.Listen("tcp", fmt.Sprintf(":%s", os.Getenv("CORE_PORT")))
	if errCreateConn != nil {
		_ = logger.System(fmt.Sprintf("failed to listen: %v", err))
		return
	}

	GRPCServer := grpc.NewServer(
		grpc.UnaryInterceptor(logger.Interceptor),
	)
	protobuf.RegisterServersServer(GRPCServer, s)


	_ = logger.System(fmt.Sprintf("Protobuf CORE started on  :%s", os.Getenv("CORE_PORT")))
	err = GRPCServer.Serve(lis)
	if err != nil {
		_ = logger.System(fmt.Sprintf("failed to serve: %s" + err.Error()))
	}

	return
}

//Recovery application out of panic
func recovery(l slog.Log) {
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
