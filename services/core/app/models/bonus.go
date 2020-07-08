package models

import (
	"strconv"
	"strings"
	"time"
)

type Bonus interface {
	BonusGetter
	BonusSetter
}

type BonusData struct {
	Active bool      `bson:"isActive"`
	Money  int64     `bson:"cost"`
	Time   time.Time `bson:"time"` //time in format 21:15
}

func CreateBonus() *BonusData {
	return &BonusData{}
}

type BonusGetter interface {
	IsActive() bool
	Cost() int64
	InTime() time.Time // return (hour,minutes)
}

func (b *BonusData) IsActive() bool {
	return b.Active
}

func (b *BonusData) Cost() int64 {
	return b.Money
}

func (b *BonusData) InTime() time.Time {
	return b.Time
}

type BonusSetter interface {
	SetActive()
	Inactive()
	SetCost(val int64)
	SetTime(time string)
}

func (b *BonusData) SetActive() {
	b.Active = true
	b.SetTime("14:26")
}

func (b *BonusData) Inactive() {
	b.Active = false
}

func (b *BonusData) SetCost(val int64) {
	b.Money = val
}

func (b *BonusData) SetTime(tm string) {

	var (
		hour    int
		minutes int
		err     error
	)

	parts := strings.Split(tm, ":")
	if len(parts) < 2 {
		minutes = 0
	}else{
		minutes, err = strconv.Atoi(parts[0])
		if err != nil || minutes > 60 {
			minutes = 0
			err = nil
		}
	}

	hour, err = strconv.Atoi(parts[0])
	if err != nil || hour > 24 {
		hour = 0
	}

	b.Time = time.Date(0, 0, 0, hour, minutes, 0, 0, time.Local)
}
