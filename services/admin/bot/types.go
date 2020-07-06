package bot

import "github.com/sparrowganz/teleFly/telegram"

const (
	TaskType       telegram.Type = "tsk"
	TypeType       telegram.Type = "tp"
	LastChoiceType telegram.Type = "lc"
	CancelType     telegram.Type = "cnl"
	ServerType     telegram.Type = "srv"
	BonusType      telegram.Type = "bns"
	BroadcastType  telegram.Type = "broadcast"

	ChangeActiveAction telegram.Action = "saa"
	ChooseAction       telegram.Action = "ch"
	SetAction          telegram.Action = "set"
	CleanAction        telegram.Action = "c"
	DeleteAction       telegram.Action = "d"
	GetAction          telegram.Action = "g"
	AddAction          telegram.Action = "a"
	PriorityAction     telegram.Action = "pr"
	CheckAction        telegram.Action = "check"
	FakeAction         telegram.Action = "fake"
	AddButton          telegram.Action = "but"
	SendAction         telegram.Action = "send"
	StopAction         telegram.Action = "stop"

	YesID    = "y"
	NoID     = "n"
	AllID    = "all"
	ButtonID = "but"

	imageType = "img"
	videoType = "vid"
)
