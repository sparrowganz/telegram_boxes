package servers

import "fmt"

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
	Initer
	Getter
	Sender
}

type Getter interface {
	ID() string
}

func (s *ServersData) ID() string {
	return s.serverID
}

type ServersData struct {
	serverID string
}

func CreateServers() Servers {
	return &ServersData{}
}

type Initer interface {
	Init()
}

func (s *ServersData) Init() {
	fmt.Println("SEND INIT BOX")
	s.serverID = "123"
}

type Sender interface {
	SendError(err string, status Status)
}

func (s *ServersData) SendError(err string, status Status) {
	fmt.Println(s.serverID, err, status)
}
