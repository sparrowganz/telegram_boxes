package models

type UserCheckData interface {
	Check(id string) bool
	Add(id string)
	Delete(id string)
	GetAll() []string
}

func (u *UserData) Check(id string) bool {
	u.checkDataMutex.Lock()
	defer u.checkDataMutex.Unlock()

	for _, cID := range u.ChecksData {
		if cID == id {
			return true
		}
	}
	return false
}

func (u *UserData) Add(id string) {
	if !u.Check(id) {

		u.checkDataMutex.Lock()
		defer u.checkDataMutex.Unlock()

		u.ChecksData = append(u.ChecksData, id)
	}
}

func (u *UserData) Delete(id string) {

	u.checkDataMutex.Lock()
	defer u.checkDataMutex.Unlock()

	var newSlice []string

	for _, cID := range u.ChecksData {
		if cID != id {
			newSlice = append(newSlice, cID)
		}
	}

	u.ChecksData = newSlice
	return
}

func (u *UserData) GetAll() []string {

	return u.ChecksData
}
