package task

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"telegram_boxes/services/box/app"
	"telegram_boxes/services/box/protobuf/services/core/protobuf"
	"time"
)

type Tasks interface {
	connector
	Getter
	Remover
}

type connector interface {
	connect(host, port, username string) error
}

func (t *Data) connect(host, port, username string) error {

	cnnServers, err := grpc.Dial(
		fmt.Sprintf("%s:%s", host, port),
		grpc.WithInsecure(),
	)
	if err != nil {
		return fmt.Errorf("%s.ServersConnect: %s", username, err.Error())
	}

	t.client = protobuf.NewTasksClient(cnnServers)
	_, cancel := context.WithTimeout(context.Background(), 10000*time.Millisecond)
	defer cancel()
	return nil
}

type Data struct {
	username string
	client   protobuf.TasksClient
	serverID string
}

func CreateTasks(host, port, username string) (Tasks, error) {
	d := &Data{
		username: username,
	}

	err := d.connect(host, port, username)
	if err != nil {
		return nil, err
	}

	return d, nil
}

type Getter interface {
	GetTask(completedTask []string) (*protobuf.Task, error)
	FindTask(id string) (*protobuf.Task, error)
	CheckTask(chatID int64, taskID string) (bool, error)
}

func (t *Data) CheckTask(chatID int64, taskID string) (bool, error) {

	res, err := t.client.CheckTask(
		app.SetCallContext("checkTask", t.username),
		&protobuf.CheckTaskRequest{
			ChatID: chatID,
			TaskID: taskID,
		})
	if err != nil {
		return false, err
	}

	return res.GetIsCheck(), nil
}

func (t *Data) FindTask(id string) (*protobuf.Task, error) {
	res, err := t.client.FindTask(
		app.SetCallContext("findTask", t.username),
		&protobuf.FindTaskRequest{
			Id: id,
		})
	if err != nil {
		return nil, err
	}

	return res.GetTask(), nil
}

func (t *Data) GetTask(completedTask []string) (*protobuf.Task, error) {
	res, err := t.client.GetTask(
		app.SetCallContext("getTask", t.username),
		&protobuf.GetTaskRequest{
			TasksID: completedTask,
		})
	if err != nil {
		return nil, err
	}

	return res.GetTask(), nil

}

type Remover interface {
	CleanupRun(id string)
}

func (t *Data) CleanupRun(id string) {

}
