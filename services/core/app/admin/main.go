package admin

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"telegram_boxes/services/core/app"
	"telegram_boxes/services/core/protobuf/services/admin/protobuf"
	"time"
)

type Client interface {
	connector
	SendError(status, username, err string) error
	CheckExecution(url string, chatID int64) (bool, error)
}

type clientData struct {
	client protobuf.AdminClient
}

type connector interface {
	connect(host, port, username string) error
}

func (c *clientData) connect(host, port, username string) error {

	cnnServers, err := grpc.Dial(
		fmt.Sprintf("%s:%s", host, port),
		grpc.WithInsecure(),
	)
	if err != nil {
		return fmt.Errorf("%s.AdminConnect: %s", username, err.Error())
	}

	c.client = protobuf.NewAdminClient(cnnServers)
	_, cancel := context.WithTimeout(context.Background(), 10000*time.Millisecond)
	defer cancel()
	return nil
}

func CreateClient(host, port string) (Client, error) {
	d := &clientData{}

	err := d.connect(host, port, "core")
	if err != nil {
		return nil, err
	}
	return d, nil
}

func (c *clientData) SendError(status, username, err string) error {

	_, sendErr := c.client.SendError(
		app.SetCallContext("SendError", "core"),
		&protobuf.SendErrorRequest{
			Error:    err,
			Status:   status,
			Username: username,
		},
	)
	return sendErr
}

func (c *clientData) CheckExecution(url string, chatID int64) (bool, error) {
	res, sendErr := c.client.CheckExecution(
		app.SetCallContext("CheckExecution", "core"),
		&protobuf.CheckExecutionRequest{
			Url:    url,
			ChatID: chatID,
		},
	)
	if sendErr != nil {
		return false, sendErr
	}
	return res.GetIsCheck() , nil
}
