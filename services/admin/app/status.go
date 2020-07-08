package app

type Status string

func (s Status) String() string {
	return string(s)
}

var (
	StatusOK       Status = "ok"
	StatusRecovery Status = "recovery"
	StatusFatal    Status = "fatal"
)
