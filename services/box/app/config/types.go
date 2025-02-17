package config

type Type string

const (
	MainType      Type = "main"
	TaskType      Type = "tasks"
	BalanceType   Type = "balance"
	CancelType    Type = "cancel"
	OutputType    Type = "output"
	OutputGWType  Type = "outputGW"
	ReferralsType Type = "referrals"
	HelpType      Type = "help"
	BuyType       Type = "buy"

	SkipTaskType  Type = "skipTask"
	NextTaskType  Type = "nextTask"
	CheckTaskType Type = "checkTask"
)

func (t Type) ToString() string {
	return string(t)
}


