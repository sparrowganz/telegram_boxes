package admins

import "sync"

type Admin interface {
	Add(id int64)
	GetAll() []int64
	IsSet(id int64) bool
}

type adminData struct {
	m       *sync.Mutex
	storage []int64
}

func CreateAdmin() Admin {
	return &adminData{
		m:       &sync.Mutex{},
		storage: make([]int64, 10),
	}
}

func (a *adminData) Add(id int64) {
	if a.IsSet(id) {
		return
	}

	a.m.Lock()
	defer a.m.Unlock()
	a.storage = append(a.storage, id)
}

func (a *adminData) GetAll() []int64 {
	a.m.Lock()
	defer a.m.Unlock()
	return a.storage
}

func (a *adminData) IsSet(id int64) bool {
	a.m.Lock()
	defer a.m.Unlock()
	for _, sID := range a.storage {
		if sID == id {
			return true
		}
	}
	return false
}
