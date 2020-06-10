package models

type Balance interface {
	BalanceGetter
	BalanceSetter
	BalanceAdder
}

func CreateBalance() Balance {
	return &BalanceData{
		BotCost:     0,
		PaymentCost: 0,
	}
}

type BalanceData struct {
	BotCost     float64 `bson:"bot"`
	PaymentCost float64 `bson:"payment"`
}

type BalanceGetter interface {
	Bot() float64
	Payment() float64
}

func (b *BalanceData) Bot() float64 {
	return b.BotCost
}

func (b *BalanceData) Payment() float64 {
	return b.PaymentCost
}

type BalanceSetter interface {
	SetBot(float64)
	SetPayment(float64)
}

func (b *BalanceData) SetBot(cost float64) {
	b.BotCost = cost
}

func (b *BalanceData) SetPayment(cost float64) {
	b.PaymentCost = cost
}

type BalanceAdder interface {
	AddBot(float64)
	AddPayment(float64)
}

func (b *BalanceData) AddBot(cost float64) {
	b.BotCost += cost
}

func (b *BalanceData) AddPayment(cost float64) {
	b.PaymentCost += cost
}
