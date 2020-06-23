package models

import "time"

type Timestamp interface {
	TimestampSetter
	TimestampAborter
	TimestampGetter
}

func CreateTimestamp() *TimestampData {
	return &TimestampData{
		CreatedData: time.Time{},
		Updated:     time.Time{},
		Removed:     time.Time{},
	}
}

type TimestampData struct {
	CreatedData time.Time `bson:"created"`
	Updated     time.Time `bson:"updated"`
	Removed     time.Time `bson:"removed" `
}

type TimestampSetter interface {
	SetCreateTime()
	SetUpdateTime()
	SetRemoveTime()
}

func (tm *TimestampData) SetCreateTime() {
	tm.CreatedData = time.Now()
}

func (tm *TimestampData) SetUpdateTime() {
	tm.Updated = time.Now()
}

func (tm *TimestampData) SetRemoveTime() {
	tm.Removed = time.Now()
}

type TimestampAborter interface {
	AbortUpdatedTime()
	AbortRemoveTime()
}

func (tm *TimestampData) AbortUpdatedTime() {
	tm.Updated = time.Time{}
}

func (tm *TimestampData) AbortRemoveTime() {
	tm.Removed = time.Time{}
}

type TimestampGetter interface {
	Created() time.Time
	CreatedNotZeroUnixNano() int64
	UpdatedNotZeroUnixNano() int64
	RemovedNotZeroUnixNano() int64
}

func (tm *TimestampData) Created() time.Time {
	return tm.CreatedData
}

func (tm *TimestampData) CreatedNotZeroUnixNano() int64 {
	if !tm.Created().IsZero() {
		return tm.Created().UnixNano()
	} else {
		return 0
	}
}

func (tm *TimestampData) UpdatedNotZeroUnixNano() int64 {
	if !tm.Updated.IsZero() {
		return tm.Updated.UnixNano()
	} else {
		return 0
	}
}

func (tm *TimestampData) RemovedNotZeroUnixNano() int64 {
	if !tm.Removed.IsZero() {
		return tm.Removed.UnixNano()
	} else {
		return 0
	}
}
