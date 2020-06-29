package models

type Status string

const (
	StatusCheck Status = "check"
	StatusSkip  Status = "skip"
)

func (status Status) String() string {
	return string(status)
}

type UserCheckData interface {
	Check(id string) bool
	Add(id string, status Status)
	Delete(id string)
	GetAll() []string
	GetAllChecks() (out []string)
	CleanChecks()
}

func (u *UserData) Check(id string) bool {
	//u.checkDataMutex.Lock()
	//defer u.checkDataMutex.Unlock()

	_, ok := u.ChecksData[id]
	return ok
}

func (u *UserData) Add(id string, status Status) {
	if !u.Check(id) {

		if len(u.ChecksData) == 0 {
			u.ChecksData = make(map[string]string)
		}

		//u.checkDataMutex.Lock()
		//defer u.checkDataMutex.Unlock()

		u.ChecksData[id] = status.String()
	}
}

func (u *UserData) Delete(id string) {

	//u.checkDataMutex.Lock()
	//defer u.checkDataMutex.Unlock()

	delete(u.ChecksData, id)
	return
}

func (u *UserData) GetAll() (out []string) {
	//u.checkDataMutex.Lock()
	//defer u.checkDataMutex.Unlock()
	for id := range u.ChecksData {
		out = append(out, id)
	}

	return
}

func (u *UserData) GetAllChecks() (out []string) {
	//u.checkDataMutex.Lock()
	//defer u.checkDataMutex.Unlock()
	for id, status := range u.ChecksData {
		if status == StatusCheck.String() {
			out = append(out, id)
		}
	}

	return
}

func (u *UserData) CleanChecks() {
	u.ChecksData = map[string]string{}
}
