package log

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"telegram_boxes/services/box/app"
	"telegram_boxes/services/box/protobuf/services/logs/protobuf"
	"time"
)

type Log interface {
	сonnector
	Getter
	logger
}

type сonnector interface {
	connect(host, port,username string) error
}

type Getter interface {
	client() protobuf.LoggerClient
	Interceptor(ctx context.Context,
		req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error)
}

type logger interface {
	Access(rID, method, user, duration string, time int64) (err error)
	System(msg string) error
	Error(rID, user, msg string) error
}

type logData struct {
	isDebug bool
	username string
	cl      protobuf.LoggerClient
}

func CreateLogger(isDebug bool, host, port, username string) (l Log, err error) {
	l = &logData{isDebug: isDebug, username: username}
	err = l.connect(host, port,username)
	return
}

func (l *logData) connect(host, port,username string) error {

	cnnLog, err := grpc.Dial(
		fmt.Sprintf("%s:%s", host, port),
		grpc.WithInsecure(),
	)
	if err != nil {
		return fmt.Errorf("%s.LogsConnect: %s", username, err.Error())
	}

	l.cl = protobuf.NewLoggerClient(cnnLog)
	_, cancel := context.WithTimeout(context.Background(), 10000*time.Millisecond)
	defer cancel()
	return nil
}

func (l *logData) client() protobuf.LoggerClient {
	return l.cl
}

func (l *logData) Interceptor(
	ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (interface{}, error) {

	start := time.Now()
	rID, userID := app.GetDataContext(ctx)

	// Calls the handler
	reply, err := handler(ctx, req)

	_ = l.Access(rID, info.FullMethod, userID, time.Since(start).String(), start.UnixNano())

	return reply, err
}

func (l *logData) Access(rID, method, user, duration string, time int64) (err error) {
	_, err = l.client().AccessLog(
		app.SetCallContext(rID, user),
		&protobuf.AccessLogRequest{
			RequestId:  rID,
			Method:     method,
			ServerName: l.username,
			Time:       time,
			Duration:   duration,
			User:       user,
		},
	)

	if l.isDebug && err != nil {
		fmt.Println(err)
	}

	return
}

func (l *logData) Error(rID, user, msg string) error {
	_, err := l.client().ErrorLog(
		app.SetCallContext(rID, user),
		&protobuf.ErrorLogRequest{
			RequestId:  rID,
			ServerName: l.username,
			Time:       time.Now().UnixNano(),
			Error:      msg,
		})

	if l.isDebug {
		if err != nil {
			fmt.Println("ERROR: ", err)
		}
		fmt.Println("ERROR: ", msg)
	}
	return err
}

func (l *logData) System(msg string) error {
	_, err := l.cl.SystemLog(
		context.Background(),
		&protobuf.SystemLogRequest{
			ServerName: l.username,
			Time:       time.Now().UnixNano(),
			Data:       msg,
		})
	if l.isDebug {
		if err != nil {
			_ = l.Error("", "system", err.Error())
		}
		fmt.Println("SYSTEM:", msg)
	}
	return err
}
