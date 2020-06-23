package models

import "strings"

type Bonus interface {
	BonusGetter
	BonusSetter
}

type BonusData struct {
	Active bool    `bson:"isActive"`
	Money  float64 `bson:"cost"`
	Time   string  `bson:"time"` //time in format 21:15
}

func CreateBonus() *BonusData {
	return &BonusData{}
}

type BonusGetter interface {
	IsActive() bool
	Cost() float64
	InTime() (hour, minutes string) // return (hour,minutes)
}

func (b *BonusData) IsActive() bool {
	return b.Active
}

func (b *BonusData) Cost() float64 {
	return b.Money
}

func (b *BonusData) InTime() (hour, minutes string) {
	parts := strings.Split(b.Time, ":")
	if len(parts) == 2 {
		return parts[0], parts[1]
	} else if len(parts) == 1 {
		return parts[0], "00"
	}
	return "00", "00"
}

type BonusSetter interface {
	SetActive()
	Inactive()
	SetCost(val float64)
	SetTime(time string)
}

func (b *BonusData) SetActive() {
	b.Active = true
}

func (b *BonusData) Inactive() {
	b.Active = false
}

func (b *BonusData) SetCost(val float64) {
	b.Money = val
}

func (b *BonusData) SetTime(time string) {
	if !strings.Contains(time, ":") {
		time += ":00"
	}
	b.Time = time
}
