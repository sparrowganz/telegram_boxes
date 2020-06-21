package models

type Balance interface {
	BalanceGetter
	BalanceSetter
	BalanceAdder
}

func CreateBalance() *BalanceData {
	return &BalanceData{
		BotCost:     0.0,
		PaymentCost: 0.0,
	}
}

type BalanceData struct {
	BotCost     int `bson:"bot"`
	PaymentCost int `bson:"payment"`
}

type BalanceGetter interface {
	Bot() int
	Payment() int
}

func (b *BalanceData) Bot() int {
	return b.BotCost
}

func (b *BalanceData) Payment() int {
	return b.PaymentCost
}

type BalanceSetter interface {
	SetBot(int)
	SetPayment(int)
}

func (b *BalanceData) SetBot(cost int) {
	b.BotCost = cost
}

func (b *BalanceData) SetPayment(cost int) {
	b.PaymentCost = cost
}

type BalanceAdder interface {
	AddBot(int)
	AddPayment(int)
}

func (b *BalanceData) AddBot(cost int) {
	b.BotCost += cost
}

func (b *BalanceData) AddPayment(cost int) {
	b.PaymentCost += cost
}
