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


type Servers interface {
	connector
	Getter
	Checker
	Broadcaster
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
	client      protobuf.ServersClient
	coefficient int64
}

func CreateServers(host, port string) (Servers, error) {
	d := &serversData{
		coefficient: 11,
	}

	err := d.connect(host, port, "admin")
	if err != nil {
		return nil, err
	}
	return d, nil
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
	Status   string
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

type Broadcaster interface {
	GetAllBroadcasts() ([]*protobuf.Stat, error)
	StartBroadcast(r *protobuf.StartBroadcastRequest) error
	StopBroadcast(id string) error
	GetStatisticsBroadcast(id string) ([]*protobuf.Stat, error)
}

func (s *serversData) GetAllBroadcasts() ([]*protobuf.Stat, error) {
	br, err := s.client.GetAllBroadcasts(
		app.SetCallContext("GetAllBroadcasts", "admin"), &protobuf.GetAllBroadcastsRequest{})
	if err != nil {
		return []*protobuf.Stat{}, err
	}

	return br.GetStats(), nil
}

func (s *serversData) StartBroadcast(r *protobuf.StartBroadcastRequest) error {
	_, err := s.client.StartBroadcast(
		app.SetCallContext("StartBroadcast", "admin"), r)
	if err != nil {
		return err
	}

	return nil
}

func (s *serversData) StopBroadcast(id string) error {
	_, err := s.client.StopBroadcast(
		app.SetCallContext("StopBroadcast", "admin"), &protobuf.StopBroadcastRequest{
			BroadcastID: id,
		})
	if err != nil {
		return err
	}

	return nil
}

func (s *serversData) GetStatisticsBroadcast(id string) ([]*protobuf.Stat, error) {
	res, err := s.client.GetStatisticsBroadcast(
		app.SetCallContext("GetStatisticsBroadcast", "admin"), &protobuf.GetStatisticsBroadcastRequest{
			BroadcastID: id,
		})
	if err != nil {
		return []*protobuf.Stat{}, err
	}

	return res.GetStats(), nil
}
