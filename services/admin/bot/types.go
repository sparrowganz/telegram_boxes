package bot

import "github.com/sparrowganz/teleFly/telegram"

const (
	TaskType       telegram.Type = "tsk"
	TypeType       telegram.Type = "tp"
	LastChoiceType telegram.Type = "lc"
	CancelType     telegram.Type = "cnl"
	ServerType     telegram.Type = "srv"
	BonusType      telegram.Type = "bns"

	ChangeActiveAction telegram.Action = "saa"
	ChooseAction       telegram.Action = "ch"
	CleanAction        telegram.Action = "c"
	DeleteAction       telegram.Action = "d"
	GetAction          telegram.Action = "g"
	AddAction          telegram.Action = "a"
	PriorityAction     telegram.Action = "pr"
	CheckAction        telegram.Action = "check"
	FakeAction         telegram.Action = "fake"

	YesID = "y"
	NoID  = "n"
	AllID = "all"
)
