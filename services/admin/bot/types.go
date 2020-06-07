package bot

import "github.com/sparrowganz/teleFly/telegram"

const (
	TaskType       telegram.Type = "tsk"
	TypeType       telegram.Type = "tp"
	LastChoiceType telegram.Type = "lc"
	CancelType     telegram.Type = "cnl"

	ChooseAction   telegram.Action = "ch"
	CleanAction    telegram.Action = "c"
	DeleteAction   telegram.Action = "d"
	GetAction      telegram.Action = "g"
	AddAction      telegram.Action = "a"
	PriorityAction telegram.Action = "pr"

	YesID = "y"
	NoID  = "n"
)
