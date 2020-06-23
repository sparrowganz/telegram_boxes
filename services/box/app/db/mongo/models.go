package mongo

type Models interface {
	Users() Users
	Outputs() Outputs
}

type modelsData struct {
	users   Users
	outputs Outputs
}

func createModels(database string) Models {
	return &modelsData{
		users:   createUsersModel(database),
		outputs: createOutputsModel(database),
	}
}

func (m *modelsData) Users() Users {
	return m.users
}

func (m *modelsData) Outputs() Outputs {
	return m.outputs
}
