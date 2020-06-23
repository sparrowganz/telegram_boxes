package servers

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"telegram_boxes/services/box/app"
	"telegram_boxes/services/box/protobuf/services/core/protobuf"
	"time"
)

type Servers interface {
	connector
	Initer
	Getter
	Sender
}

type connector interface {
	connect(host, port, username string) error
}

func (data *Data) connect(host, port, username string) error {

	cnnServers, err := grpc.Dial(
		fmt.Sprintf("%s:%s", host, port),
		grpc.WithInsecure(),
	)
	if err != nil {
		return fmt.Errorf("%s.ServersConnect: %s", username, err.Error())
	}

	data.client = protobuf.NewServersClient(cnnServers)
	_, cancel := context.WithTimeout(context.Background(), 10000*time.Millisecond)
	defer cancel()
	return nil
}

type Getter interface {
	ID() string
}

func (data *Data) ID() string {
	return data.serverID
}

type Data struct {
	host     string
	port     string
	username string
	client   protobuf.ServersClient
	serverID string
}

func CreateServers(host, port, username string) (Servers, error) {
	d := &Data{
		host:     host,
		port:     port,
		username: username,
	}

	err := d.connect(host, port, username)
	if err != nil {
		return nil, err
	}

	err = d.Init(username)
	if err != nil {
		return nil, err
	}

	return d, nil
}

type Initer interface {
	Init(username string) error
}

func (data *Data) Init(username string) error {
	if data.client == nil {
		return errors.New("client not initialize")
	}

	res, err := data.client.InitBox(
		app.SetCallContext("init", username),
		&protobuf.InitBoxRequest{
			Username: username,
		})
	if err != nil {
		return err
	}

	data.serverID = res.ID
	return nil
}

type Sender interface {
	SendError(err string, status protobuf.Status) error
}

func (data *Data) SendError(err string, status protobuf.Status) error {

	if data.client == nil {
		return errors.New("client not initialize")
	}

	_, errSend := data.client.SendError(
		app.SetCallContext("error", data.username),
		&protobuf.SendErrorRequest{
			Error:  err,
			Status: status,
		})
	if errSend != nil {
		return errSend
	}

	return nil
}
