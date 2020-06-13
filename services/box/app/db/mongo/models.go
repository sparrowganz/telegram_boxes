package mongo

type Models interface {
	Users() Users

}

type modelsData struct {
	users Users
}

func createModels(database string) Models {
	return &modelsData{
		users: createUsersModel(database),
	}
}

func (m *modelsData) Users() Users {
	return m.users
}
