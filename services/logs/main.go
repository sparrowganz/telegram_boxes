package main

import (
	"errors"
	"telegram_boxes/services/logs/app"
	"telegram_boxes/services/logs/protobuf"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"time"
)

func main() {
	s := new(protobuf.Server)
	defer recovery(s)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", os.Getenv("LOGS_PORT")))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	_ = os.Mkdir(os.Getenv("LOGS_PATH"), os.ModePerm)
	s.Logger = app.InitLogger(os.Getenv("LOGS_PATH"))

	grpcServer := grpc.NewServer()
	protobuf.RegisterLoggerServer(grpcServer, s)

	s.Logger.System(time.Now().UnixNano(), app.ServiceName,
		fmt.Sprintf("Protobuf started on  :%s", os.Getenv("LOGS_PORT")))

	err = grpcServer.Serve(lis)
	if err != nil {
		s.Logger.System(time.Now().UnixNano(), app.ServiceName, fmt.Sprintf("failed to serve: %s", err.Error()))
	}
}

func recovery(s *protobuf.Server) {
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
		s.Logger.System(time.Now().UnixNano(), app.ServiceName, fmt.Sprintf("RECOVERY SERVER: %v", err.Error()))
	}
}
