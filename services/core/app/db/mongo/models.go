package mongo

type Models interface {
	Bots() Bots
	Tasks() Tasks
	Payments() Payments
	Broadcasts() Broadcasts
}

type modelsData struct {
	bots       Bots
	tasks      Tasks
	payments   Payments
	broadcasts Broadcasts
}

func createModels(database string) Models {
	return &modelsData{
		bots:       createBotModel(database),
		tasks:      createTaskModel(database),
		payments:   createPaymentsModel(database),
		broadcasts: createBroadcastModel(database),
	}
}

func (m *modelsData) Bots() Bots {
	return m.bots
}

func (m *modelsData) Tasks() Tasks {
	return m.tasks
}

func (m *modelsData) Broadcasts() Broadcasts {
	return m.broadcasts
}

func (m *modelsData) Payments() Payments {
	return m.payments
}
