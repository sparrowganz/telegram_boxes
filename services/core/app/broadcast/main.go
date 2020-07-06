package broadcast

import (
	"context"
	"gopkg.in/mgo.v2/bson"
	"sync"
	"time"
)

type Broadcaster interface {
	Adder
	Getter
	Remover
}

type Data struct {
	m       *sync.Mutex
	storage map[string]Sender //[id]Sender
}

type Getter interface {
	GetAll() map[string]Sender
	GetAllByBotID(bot string) map[string]Sender
	Get(id string) (Sender, bool)
}

func (d *Data) GetAllByBotID(bot string) map[string]Sender {
	d.m.Lock()
	defer d.m.Unlock()

	dataForBot := map[string]Sender{}

	for id, s := range d.storage {
		if s.Bot() == bot {
			dataForBot[id] = s
		}
	}

	return dataForBot
}

func (d *Data) GetAll() map[string]Sender {
	d.m.Lock()
	defer d.m.Unlock()
	newMap := d.storage
	return newMap
}

func (d *Data) Get(id string) (Sender, bool) {
	d.m.Lock()
	defer d.m.Unlock()

	s, ok := d.storage[id]
	return s, ok
}

type Adder interface {
	Add(botID string, ctx context.Context, cancel context.CancelFunc) (id string, s Sender)
}

func (d *Data) Add(botID string, ctx context.Context, cancel context.CancelFunc) (id string, s Sender) {

	id = bson.NewObjectId().Hex()

	s = &SenderData{
		BotID:   botID,
		Success: 0,
		Fail:    0,
		ctx:     ctx,
		cancel:  cancel,
		time:    time.Now(),
	}

	d.m.Lock()
	defer d.m.Unlock()

	d.storage[id] = s
	return
}

type Remover interface {
	Remove(id string)
	RemoveAll()
}

func (d *Data) Remove(id string) {
	d.m.Lock()
	d.m.Unlock()

	s, ok := d.storage[id]
	if !ok {
		return
	}

	s.Stop()
	delete(d.storage, id)
}

func (d *Data) RemoveAll() {
	d.m.Lock()
	d.m.Unlock()

	for _, s := range d.storage {
		s.Stop()
	}
	d.storage = map[string]Sender{}
}

func Create() Broadcaster {
	return &Data{
		m:       &sync.Mutex{},
		storage: map[string]Sender{},
	}
}

type Sender interface {
	Bot() string
	Stop()
	Stats() (success, fail int64)
	SetAccess(count int64)
	SetFail(count int64)
	StartTime() time.Time
}

type SenderData struct {
	BotID   string
	Success int64
	Fail    int64
	ctx     context.Context
	cancel  context.CancelFunc
	time    time.Time
}

func (s *SenderData) StartTime() time.Time {
	return s.time
}

func (s *SenderData) Bot() string {
	return s.BotID
}

func (s *SenderData) Stop() {
	s.cancel()
}

func (s *SenderData) Stats() (success, fail int64) {
	return s.Success, s.Fail
}

func (s *SenderData) SetAccess(count int64) {
	s.Success = count
}

func (s *SenderData) SetFail(count int64) {
	s.Fail = count
}
