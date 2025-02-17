package log

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"telegram_boxes/services/core/app"
	"telegram_boxes/services/core/protobuf/services/logs/protobuf"
	"time"
)

type Client interface {
	Connector
	Getter
	logger
}

type Connector interface {
	connect(host, port string) error
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
	cl protobuf.LoggerClient
}

func CreateLogger(host, port string) (l Client, err error) {
	l = &logData{}
	err = l.connect(host, port)
	return
}

func (l *logData) connect(host, port string) error {

	cnnLog, err := grpc.Dial(
		fmt.Sprintf("%s:%s", host, port),
		grpc.WithInsecure(),
	)
	if err != nil {
		return fmt.Errorf("%s.LogsConnect: %s", app.ServiceName, err.Error())
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
			ServerName: app.ServiceName,
			Time:       time,
			Duration:   duration,
			User:       user,
		},
	)
	return
}

func (l *logData) Error(rID, user, msg string) error {
	_, err := l.client().ErrorLog(
		app.SetCallContext(rID, user),
		&protobuf.ErrorLogRequest{
			RequestId:  rID,
			ServerName: app.ServiceName,
			Time:       time.Now().UnixNano(),
			Error:      msg,
		})
	return err
}

func (l *logData) System(msg string) error {
	_, err := l.cl.SystemLog(
		context.Background(),
		&protobuf.SystemLogRequest{
			ServerName: app.ServiceName,
			Time:       time.Now().UnixNano(),
			Data:       msg,
		})
	return err
}
