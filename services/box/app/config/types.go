package config

type Type string

const (
	TaskType      Type = "tasks"
	BalanceType   Type = "balance"
	CancelType    Type = "cancel"
	OutputType    Type = "output"
	ReferralsType Type = "referrals"
	HelpType      Type = "help"
	BuyType       Type = "buy"
)

func (t Type) ToString() string {
	return string(t)
}
