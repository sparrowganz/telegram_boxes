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
	storage map[string]*client //[bsonID.Hex()]protobuf.BoxClient
}

type client struct {
	client protobuf.BoxClient
	addr   string
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

func CreateClients(cl db.Client) Clients {
	c := &ClientsData{
		m:       &sync.Mutex{},
		storage: make(map[string]*client),
	}

	bots, _ := cl.Models().Bots().GetAll(cl.GetMainSession())
	for _, bot := range bots {
		conn, err := c.connect(bot.Address().IP(), bot.Address().Port(), bot.Username())
		if err != nil {
			bot.SetStatus("inactive")
			bot.InActive()
			_ = cl.Models().Bots().UpdateBot(bot, cl.GetMainSession())
		}

		bot.SetStatus(app.StatusOK.String())
		bot.SetActive()
		_ = cl.Models().Bots().UpdateBot(bot, cl.GetMainSession())

		c.m.Lock()
		c.storage[bot.ID().Hex()] = &client{
			client: conn,
			addr:   bot.Addr.Ip + ":" + bot.Addr.PortNum,
		}
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
	c.storage[bot.ID().Hex()] = &client{
		client: conn,
		addr:   bot.Address().IP() + ":" + bot.Address().Port(),
	}
	c.m.Unlock()
	return conn, nil
}

type Requester interface {
	CheckBox(bot models.Bot, chatID int64) (string, error)
	GetStats(bot models.Bot) (stats *protobuf.Statistic, err error)
	RemoveTask(bot models.Bot, taskID string) error
	StartBroadcast(bot models.Bot, ch chan *protobuf.Stats, ctx context.Context, r *protobuf.StartBroadcastRequest)
	check(client protobuf.BoxClient, chatID int64) (string, error)
}

func (c *ClientsData) RemoveTask(bot models.Bot, taskID string) (err error) {
	conn, ok := c.get(bot.ID().Hex())
	if !ok {
		conn, err = c.add(bot)
		if err != nil {
			return err
		}
	}
	_, err = conn.RemoveCheckTask(
		app.SetCallContext("RemoveTask", "core"),
		&protobuf.RemoveCheckTaskRequest{
			TaskID: taskID,
		},
	)
	return
}

func (c *ClientsData) GetStats(bot models.Bot) (stats *protobuf.Statistic, err error) {
	cl, ok := c.get(bot.ID().Hex())
	if !ok {
		cl, err = c.add(bot)
		if err != nil {
			return nil, err
		}
	}
	stats, err = cl.GetStatistics(
		app.SetCallContext("getStats", "core"),
		&protobuf.GetStatisticsRequest{},
	)
	return
}

func (c *ClientsData) CheckBox(bot models.Bot, chatID int64) (status string, err error) {
	cl, ok := c.get(bot.ID().Hex())
	if !ok {
		cl, err = c.add(bot)
		if err != nil {
			return "Not Connected", err
		}
	} else {
		c.updateAddr(bot)
	}

	status, err = c.check(cl, chatID)
	return
}

func (c *ClientsData) StartBroadcast(bot models.Bot, ch chan *protobuf.Stats, ctx context.Context, r *protobuf.StartBroadcastRequest) {
	defer close(ch)

	cl, ok := c.get(bot.ID().Hex())
	if !ok {
		var err error
		cl, err = c.add(bot)
		if err != nil {
			return
		}
	}

	stream, errStart := cl.StartBroadcast(app.SetCallContextWithContext(ctx, "StartBroadcast", "core"), r)
	if errStart != nil {
		return
	}

	for {
		stats, errRecv := stream.Recv()
		if errRecv != nil {
			break
		}
		ch <- stats
	}

}

func (c *ClientsData) check(client protobuf.BoxClient, chatID int64) (string, error) {
	_, err := client.Check(
		app.SetCallContext("check", "core"), &protobuf.CheckRequest{ChatID: chatID})
	if err != nil {
		return app.StatusFatal.String(), err
	}
	return app.StatusOK.String(), nil
}

type getter interface {
	get(botID string) (protobuf.BoxClient, bool)
	updateAddr(bot models.Bot)
}

func (c *ClientsData) updateAddr(bot models.Bot) {
	c.m.Lock()
	defer c.m.Unlock()

	cl, ok := c.storage[bot.ID().Hex()]
	if !ok {
		conn, _ := c.connect(bot.Address().IP(), bot.Address().Port(), bot.Username())
		c.storage[bot.ID().Hex()] = &client{
			client: conn,
			addr:   bot.Address().IP() + ":" + bot.Address().Port(),
		}
	}

	if cl.addr != bot.Address().IP()+":"+bot.Address().Port() {
		cl.addr = bot.Address().IP() + ":" + bot.Address().Port()
		cl.client, _ = c.connect(bot.Address().IP(), bot.Address().Port(), bot.Username())
	}

	return
}

func (c *ClientsData) get(botID string) (protobuf.BoxClient, bool) {
	c.m.Lock()
	defer c.m.Unlock()
	cl, ok := c.storage[botID]
	return cl.client, ok
}
