package models

import "time"

type Timestamp interface {
	TimestampSetter
	TimestampAborter
	TimestampGetter
}

func CreateTimestamp() Timestamp {
	return &timestampData{
		Created: time.Time{},
		Updated: time.Time{},
		Removed: time.Time{},
	}
}

type timestampData struct {
	Created time.Time `bson:"created"`
	Updated time.Time `bson:"updated"`
	Removed time.Time `bson:"removed" `
}

type TimestampSetter interface {
	SetCreateTime()
	SetUpdateTime()
	SetRemoveTime()
}

func (tm *timestampData) SetCreateTime() {
	tm.Created = time.Now()
}

func (tm *timestampData) SetUpdateTime() {
	tm.Updated = time.Now()
}

func (tm *timestampData) SetRemoveTime() {
	tm.Removed = time.Now()
}

type TimestampAborter interface {
	AbortUpdatedTime()
	AbortRemoveTime()
}

func (tm *timestampData) AbortUpdatedTime() {
	tm.Updated = time.Time{}
}

func (tm *timestampData) AbortRemoveTime() {
	tm.Removed = time.Time{}
}

type TimestampGetter interface {
	CreatedNotZeroUnixNano() int64
	UpdatedNotZeroUnixNano() int64
	RemovedNotZeroUnixNano() int64
}

func (tm *timestampData) CreatedNotZeroUnixNano() int64 {
	if !tm.Created.IsZero() {
		return tm.Created.UnixNano()
	} else {
		return 0
	}
}

func (tm *timestampData) UpdatedNotZeroUnixNano() int64 {
	if !tm.Updated.IsZero() {
		return tm.Updated.UnixNano()
	} else {
		return 0
	}
}

func (tm *timestampData) RemovedNotZeroUnixNano() int64 {
	if !tm.Removed.IsZero() {
		return tm.Removed.UnixNano()
	} else {
		return 0
	}
}
