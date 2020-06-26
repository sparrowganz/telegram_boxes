package box

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"sync"
	"telegram_boxes/services/core/app"
	"telegram_boxes/services/core/app/db"
	"telegram_boxes/services/core/app/models"
	"telegram_boxes/services/core/protobuf/services/box/protobuf"
	"time"
)

type Clients interface {
	connector
	getter
	Requester
	Adder
}

type ClientsData struct {
	m       *sync.Mutex
	storage map[string]protobuf.BoxClient //[bsonID.Hex()]protobuf.BoxClient
}

type connector interface {
	connect(host, port, username string) (protobuf.BoxClient, error)
}

func (c *ClientsData) connect(host, port, username string) (protobuf.BoxClient, error) {

	cnnServers, err := grpc.Dial(
		fmt.Sprintf("%s:%s", host, port),
		grpc.WithInsecure(),
	)
	if err != nil {
		return nil, fmt.Errorf("%s.BoxConnect: %s", username, err.Error())
	}

	_, cancel := context.WithTimeout(context.Background(), 10000*time.Millisecond)
	defer cancel()

	return protobuf.NewBoxClient(cnnServers), nil
}

func CreateClients(client db.Client) Clients {
	c := &ClientsData{
		m:       &sync.Mutex{},
		storage: map[string]protobuf.BoxClient{},
	}

	bots , _ := client.Models().Bots().GetAll(client.GetMainSession())
	for _, bot := range bots {
		conn, err := c.connect(bot.Address().IP(), bot.Address().Port(), bot.Username())
		if err != nil {
			bot.SetStatus("inactive")
			_ = client.Models().Bots().UpdateBot(bot, client.GetMainSession())
		}
		c.m.Lock()
		c.storage[bot.ID().Hex()] = conn
		c.m.Unlock()
	}

	return c
}

type Adder interface {
	AddBox(bot models.Bot) error
	add(bot models.Bot) (protobuf.BoxClient, error)
}

func (c *ClientsData) AddBox(bot models.Bot) (err error) {
	_, err = c.add(bot)
	return
}

func (c *ClientsData) add(bot models.Bot) (protobuf.BoxClient, error) {
	conn, err := c.connect(bot.Address().IP(), bot.Address().Port(), bot.Username())
	if err != nil {
		return nil, err
	}
	c.m.Lock()
	c.storage[bot.ID().Hex()] = conn
	c.m.Unlock()
	return conn, nil
}

type Requester interface {
	CheckBox(bot models.Bot, chatID int64) (string, error)
	GetStats(bot models.Bot) (stats *protobuf.Statistic, err error)
	RemoveTask(bot models.Bot, taskID string) error
	check(client protobuf.BoxClient, chatID int64) (string, error)
}

func (c *ClientsData) RemoveTask(bot models.Bot, taskID string) (err error) {
	client, ok := c.get(bot.ID().Hex())
	if !ok {
		client, err = c.add(bot)
		if err != nil {
			return err
		}
	}
	_, err = client.RemoveCheckTask(
		app.SetCallContext("RemoveTask", "core"),
		&protobuf.RemoveCheckTaskRequest{
			TaskID: taskID,
		},
	)
	return
}

func (c *ClientsData) GetStats(bot models.Bot) (stats *protobuf.Statistic, err error) {
	client, ok := c.get(bot.ID().Hex())
	if !ok {
		client, err = c.add(bot)
		if err != nil {
			return nil, err
		}
	}
	stats, err = client.GetStatistics(
		app.SetCallContext("getStats", "core"),
		&protobuf.GetStatisticsRequest{},
	)
	return
}

func (c *ClientsData) CheckBox(bot models.Bot, chatID int64) (status string, err error) {
	client, ok := c.get(bot.ID().Hex())
	if !ok {
		client, err = c.add(bot)
		if err != nil {
			return "Not Connected", err
		}
	}
	status, err = c.check(client, chatID)
	return
}

func (c *ClientsData) check(client protobuf.BoxClient, chatID int64) (string, error) {
	_, err := client.Check(
		app.SetCallContext("check", "core"), &protobuf.CheckRequest{ChatID: chatID})
	if err != nil {
		return "FATAl", err
	}
	return "OK", nil
}

type getter interface {
	get(botID string) (protobuf.BoxClient, bool)
}

func (c *ClientsData) get(botID string) (protobuf.BoxClient, bool) {
	c.m.Lock()
	defer c.m.Unlock()
	client, ok := c.storage[botID]
	return client, ok
}
