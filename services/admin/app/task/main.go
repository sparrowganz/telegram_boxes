package task

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"telegram_boxes/services/admin/app"
	"telegram_boxes/services/admin/protobuf/services/core/protobuf"
	"time"
)

type Tasks interface {
	Getter
	Changer
	Remover
	Creator
	connector
}

type connector interface {
	connect(host, port, username string) error
}

func (t *tasksData) connect(host, port, username string) error {

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

type tasksData struct {
	host   string
	port   string
	client protobuf.TasksClient
}

func CreateTasks(host, port string) (Tasks, error) {
	d := &tasksData{
		host: host,
		port: port,
	}

	err := d.connect(host, port, "admin")
	if err != nil {
		return nil, err
	}
	return d, nil
}

type Getter interface {
	GetAllTasks() ([]*protobuf.Task, error)
	GetTask(id string) (*protobuf.Task, error)
}

func (t *tasksData) GetAllTasks() ([]*protobuf.Task, error) {
	res, err := t.client.GetAllTask(
		app.SetCallContext("GetAllTasks", "admin"),
		&protobuf.GetAllTaskRequest{})
	if err != nil {
		return []*protobuf.Task{}, err
	}

	return res.GetTasks(), nil
}

func (t *tasksData) GetTask(id string) (*protobuf.Task, error) {
	res, err := t.client.FindTask(
		app.SetCallContext("GetTask", "admin"),
		&protobuf.FindTaskRequest{Id: id},
	)
	if err != nil {
		return &protobuf.Task{}, err
	}

	return res.GetTask(), nil
}

type Changer interface {
	ChangePriority(id string) (*protobuf.Task, error)
}

func (t *tasksData) ChangePriority(id string) (*protobuf.Task, error) {
	res, err := t.client.ChangePriorityTask(
		app.SetCallContext("ChangePriority", "admin"),
		&protobuf.ChangePriorityTaskRequest{
			TaskID: id,
		})
	if err != nil {
		return &protobuf.Task{}, err
	}

	return res.GetTask(), nil
}

type Remover interface {
	Delete(id string) error
	CleanupRun(id string) (*protobuf.Task, error)
}

func (t *tasksData) Delete(id string) error {

	_, err := t.client.DeleteTask(
		app.SetCallContext("Delete", "admin"),
		&protobuf.DeleteTaskRequest{
			TaskID: id,
		})
	return err
}

func (t *tasksData) CleanupRun(id string) (*protobuf.Task, error) {
	res, err := t.client.CleanupRunTask(
		app.SetCallContext("ChangePriority", "admin"),
		&protobuf.CleanupRunTaskRequest{
			TaskID: id,
		})
	if err != nil {
		return &protobuf.Task{}, err
	}

	return res.GetTask(), nil
}

type Creator interface {
	Create(t *protobuf.Task) error
}

func (t *tasksData) Create(tsk *protobuf.Task) error {
	_, err := t.client.CreateTask(
		app.SetCallContext("CreateTask", "admin"), tsk)
	if err != nil {
		return err
	}

	return nil
}
