package servers

import (
	"errors"
	"time"
)

type Status int

var (
	OK         Status = 1
	RECOVERING Status = 2
	FATAL      Status = 3
)

func (s Status) String() string {
	switch s {
	case OK:
		return "OK"
	case RECOVERING:
		return "RECOVERING"
	case FATAL:
		return "FATAL"
	}
	return ""
}

type Servers interface {
	Getter
	Checker
}

type serversData struct {
	//taskData
	//todo remove debug structure
	storage     []*Server
	coefficient int
}

//todo remove debug structure
type Server struct {
	ID       string
	Username string

	All     int
	Blocked int
	UseNow  int

	IsActiveBonus bool

	Status Status
}

func CreateServers() Servers {
	return &serversData{
		//todo remove debug structure
		coefficient: 11,
		storage: []*Server{
			{"1", "@username1", 100, 2, 5, true, 1},
			{"2", "@username2", 200, 4, 10, false, 2},
			{"3", "@username3", 300, 6, 15, true, 3},
			{"4", "@username4", 400, 8, 20, false, 1},
		},
	}
}

type StatusData struct {
	ID       string
	Username string
	Status   Status
}

type Count struct {
	ID       string
	Username string
	All      int
	Blocked  int
	UseNow   int
}

type Bonus struct {
	ID       string
	Username string
	IsActive bool
}

type Getter interface {
	GetAllServers() []*Server
	GetAllServersStatus() []*StatusData
	GetAllUsersCount() []*Count
	GetAllUsersFakeCount() []*Count
	GetAllServersBonuses() []*Bonus
	GetServerBonus(id string) (*Bonus, error)
	ChangeActiveAllBonuses(isSetInactive bool)
	ChangeActiveBonus(id string, isSetInactive bool)
}

func (s *serversData) ChangeActiveAllBonuses(isSetInactive bool) {
	for _, b := range s.storage {
		b.IsActiveBonus = isSetInactive
	}
}

func (s *serversData) ChangeActiveBonus(id string, isSetInactive bool) {
	for _, b := range s.storage {
		if id == b.ID {
			b.IsActiveBonus = isSetInactive
			return
		}
	}
}

func (s *serversData) GetAllServers() []*Server {
	return s.storage
}

func (s *serversData) GetServerBonus(id string) (*Bonus, error) {
	for _, server := range s.storage {
		if id == server.ID {
			return &Bonus{
				ID:       server.ID,
				Username: server.Username,
				IsActive: server.IsActiveBonus,
			}, nil
		}
	}
	return nil, errors.New(" Server not found")
}

func (s *serversData) GetAllServersBonuses() (all []*Bonus) {
	for _, server := range s.storage {
		all = append(all, &Bonus{
			ID:       server.ID,
			Username: server.Username,
			IsActive: server.IsActiveBonus,
		})
	}
	return all
}

func (s *serversData) GetAllUsersCount() (all []*Count) {
	for _, server := range s.storage {
		all = append(all, &Count{
			ID:       server.ID,
			Username: server.Username,
			All:      server.All,
			Blocked:  server.Blocked,
			UseNow:   server.UseNow,
		})
	}
	return all
}

func (s *serversData) GetAllUsersFakeCount() (all []*Count) {
	for _, server := range s.storage {
		all = append(all, &Count{
			ID:       server.ID,
			Username: server.Username,
			All:      server.All * s.coefficient,
			Blocked:  server.Blocked,
			UseNow:   server.UseNow * s.coefficient,
		})
	}
	return all
}

func (s *serversData) GetAllServersStatus() (all []*StatusData) {
	for _, server := range s.storage {
		all = append(all, &StatusData{
			ID:       server.ID,
			Username: server.Username,
			Status:   server.Status,
		})
	}
	return
}

type Checker interface {
	HardCheckAll(ch chan *StatusData, userID int64)
}

func (s *serversData) HardCheckAll(ch chan *StatusData, userID int64) {
	for _, server := range s.storage {
		time.Sleep(time.Second * 2)

		ch <- &StatusData{
			ID:       server.ID,
			Username: server.Username,
			Status:   OK,
		}

		server.Status = OK
	}
	close(ch)
}

func (s *serversData) HardCheck(id string) (bool, error) {
	for _, server := range s.storage {
		if server.ID == id {
			return true, nil
		}
	}
	return false, errors.New(" Server not found")
}
