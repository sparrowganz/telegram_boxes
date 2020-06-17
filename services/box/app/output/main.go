package output

import (
	"errors"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Outputs interface {
	Getter
	Setter
}

type outputData struct {
	//taskData
	//todo remove debug structure
	botID   string
	storage []*Output
}

//todo remove debug structure
type Output struct {
	ID             string
	BotID          string
	UserID         string
	Cost           int
	PaymentGateway string
	Data           string
	Tasks          []string
	Timestamp      time.Time
}

func CreateOutput(botID string) Outputs {
	return &outputData{
		botID:   botID,
		storage: []*Output{},
	}
}

type Getter interface {
	GetOutput(userID string) (*Output, error)
}

func (t *outputData) GetOutput(userID string) (res *Output, err error) {

	for _, out := range t.storage {
		if out.UserID == userID && out.BotID == t.botID {
			return out, nil
		}
	}

	return nil, errors.New(" Not Found ")
}

type Setter interface {
	Set(userID, pGW, data string, cost int, tasks []string)
}

func (t *outputData) Set(userID, pGW, data string, cost int, tasks []string) {
	t.storage = append(t.storage, &Output{
		ID:             bson.NewObjectId().Hex(),
		BotID:          t.botID,
		UserID:         userID,
		Cost:           cost,
		PaymentGateway: pGW,
		Data:           data,
		Tasks:          tasks,
		Timestamp:      time.Now(),
	})
}
