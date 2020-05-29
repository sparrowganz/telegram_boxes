package id

import (
	"sync/atomic"
	"time"
)

//------------------------------------------------------------------------------------------------------------------
//  REQUEST
//------------------------------------------------------------------------------------------------------------------
/*
	Counter for request
	GetIdDirectory() uint64

	increment DirectoryID
*/

type Counter interface {
	ID() uint64
}

func NewCounter() Counter {
	return &id{
		oldDay: time.Now().Day(),
	}
}

type id struct {
	oldDay  int
	counter uint64
}

func (r *id) ID() uint64 {
	if time.Now().Day() != r.oldDay {
		r.oldDay = time.Now().Day()
		atomic.StoreUint64(&r.counter, 1)
		return 1
	} else {
		return atomic.AddUint64(&r.counter, 1)
	}
}
