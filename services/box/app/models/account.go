package models

type Account interface {
	AccountGetter
	AccountSetter
}

type AccountData struct {
	Id        int64  `bson:"id"`
	Username  string `bson:"username"`
	Firstname string `bson:"firstname"`
	Lastname  string `bson:"lastname"`
	Mail      string `bson:"email"`
}

func CreateAccount(telegramID int64, username, firstName, lastName, email string) *AccountData {
	return &AccountData{
		Id:        telegramID,
		Username:  username,
		Firstname: firstName,
		Lastname:  lastName,
		Mail:      email,
	}
}

type AccountGetter interface {
	ID() int64
	UserName() string
	FirstName() string
	LastName() string
	Email() string
}

func (a *AccountData) ID() int64 {
	return a.Id
}

func (a *AccountData) UserName() string {
	return a.Username
}

func (a *AccountData) FirstName() string {
	return a.Firstname
}

func (a *AccountData) LastName() string {
	return a.Lastname
}

func (a *AccountData) Email() string {
	return a.Mail
}

type AccountSetter interface {
	SetID(int64)
	SetUserName(string)
	SetFirstName(string)
	SetLastName(string)
	SetEmail(string)
}

func (a *AccountData) SetID(id int64) {
	if id != 0 {
		a.Id = id
	}
}

func (a *AccountData) SetUserName(name string) {
	if name != "" {
		a.Username = name
	}
}

func (a *AccountData) SetFirstName(name string) {
	if name != "" {
		a.Firstname = name
	}
}

func (a *AccountData) SetLastName(name string) {
	if name != "" {
		a.Lastname = name
	}
}

func (a *AccountData) SetEmail(email string) {
	if email != "" {
		a.Mail = email
	}
}
