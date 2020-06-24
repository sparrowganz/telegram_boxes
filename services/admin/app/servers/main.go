package servers

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"io"
	"telegram_boxes/services/admin/app"
	"telegram_boxes/services/admin/protobuf/services/core/protobuf"
	"time"
)

var (
	OK         = protobuf.Status_OK
	RECOVERING = protobuf.Status_Recovering
	FATAL      = protobuf.Status_Fatal
)

type Servers interface {
	connector
	Getter
	Checker
}

type connector interface {
	connect(host, port, username string) error
}

func (s *serversData) connect(host, port, username string) error {

	cnnServers, err := grpc.Dial(
		fmt.Sprintf("%s:%s", host, port),
		grpc.WithInsecure(),
	)
	if err != nil {
		return fmt.Errorf("%s.ServersConnect: %s", username, err.Error())
	}

	s.client = protobuf.NewServersClient(cnnServers)
	_, cancel := context.WithTimeout(context.Background(), 10000*time.Millisecond)
	defer cancel()
	return nil
}

type serversData struct {
	host        string
	port        string
	client      protobuf.ServersClient
	coefficient int64
}

func CreateServers(host, port string) Servers {
	return &serversData{
		host:        host,
		port:        port,
		coefficient: 11,
	}
}


type Getter interface {
	GetAllServers() ([]*protobuf.Server, error)
	GetServer(id string) (*protobuf.Server, error)
	GetAllUsersCount(isFake bool) ([]*protobuf.Counts, error)
	ChangeActiveAllBonuses() error
	ChangeActiveBonus(id string) error
}

func (s *serversData) ChangeActiveAllBonuses() error {

	_, err := s.client.ChangeBonusActive(
		app.SetCallContext("ChangeActiveAllBonuses", "admin"),
		&protobuf.ChangeBonusActiveRequest{Id: "all"})

	if err != nil {
		return err
	}
	return nil
}

func (s *serversData) ChangeActiveBonus(id string) error {
	_, err := s.client.ChangeBonusActive(
		app.SetCallContext("ChangeActiveAllBonuses", "admin"),
		&protobuf.ChangeBonusActiveRequest{Id: id})
	return err
}

func (s *serversData) GetAllServers() ([]*protobuf.Server, error) {
	res, err := s.client.GetListServers(
		app.SetCallContext("getAllServers", "admin"),
		&protobuf.GetListServersRequest{})
	if err != nil {
		return nil, err
	}

	return res.GetServers(), nil
}

func (s *serversData) GetServer(id string) (*protobuf.Server, error) {
	res, err := s.client.GetServer(
		app.SetCallContext("GetServer", "admin"),
		&protobuf.GetServerRequest{
			Id: id,
		})
	if err != nil {
		return nil, err
	}

	return res.GetServer(), nil
}

func (s *serversData) GetAllUsersCount(isFake bool) ([]*protobuf.Counts, error) {
	res, err := s.client.GetAllUsersCount(
		app.SetCallContext("GetAllUsersCount", "admin"),
		&protobuf.GetAllUsersCountRequest{})
	if err != nil {
		return []*protobuf.Counts{}, err
	}

	if isFake {
		for _, c := range res.GetCounts() {
			c.GetOld().All *= s.coefficient
			c.GetOld().Blocked *= s.coefficient
			c.GetNew().All *= s.coefficient
			c.GetNew().Blocked *= s.coefficient
		}
	}

	return res.GetCounts(), err
}

type Checker interface {
	HardCheckAll(ch chan *Check, userID int64)
}

type Check struct {
	ID       string
	Username string
	Status   protobuf.Status
	Error    string
}

func (s *serversData) HardCheckAll(ch chan *Check, userID int64) {

	defer close(ch)

	stream, err := s.client.HardCheck(
		app.SetCallContext("HardCheck", "admin"),
		&protobuf.HardCheckRequest{
			UserID:   userID,
		})

	if err != nil {
		ch <- &Check{
			Error: err.Error(),
		}
		return
	}

	for {
		status, errRecv := stream.Recv()
		if errRecv != nil {
			if errRecv == io.EOF || app.ParseGRPCError(errRecv) == context.Canceled.Error() {
				return
			}
			ch <- &Check{
				Error: errRecv.Error(),
			}
			return
		}

		ch <- &Check{
			ID:       status.Id,
			Username: status.Username,
			Status:   status.Status,
		}

	}
}

