package task

import "errors"

type Tasks interface {
	Getter
	Changer
	Remover
}

type tasksData struct {
	//taskData
	//todo remove debug structure
	storage []*Task
}

//todo remove debug structure
type Task struct {
	ID         string
	Name       string
	IsPriority bool
}

func CreateTasks() Tasks {
	return &tasksData{
		//todo remove debug structure
		storage: []*Task{
			{"1", "1", false},
			{"2", "2", true},
			{"3", "3", false},
			{"4", "4", true},
		},
	}
}

type Getter interface {
	GetAllTasks() []*Task
	GetTask(id string) (*Task, error)
}

func (t *tasksData) GetAllTasks() []*Task {
	return t.storage
}

func (t *tasksData) GetTask(id string) (*Task, error) {
	for _, tsk := range t.storage {
		if tsk.ID == id {
			return tsk, nil
		}
	}
	return nil, errors.New(" Not found ")
}

type Changer interface {
	ChangePriority(id string) (*Task, error)
}

func (t *tasksData) ChangePriority(id string) (*Task, error) {
	for _, tsk := range t.storage {
		if tsk.ID == id {
			tsk.IsPriority = !tsk.IsPriority
			return tsk, nil
		}
	}
	return nil, errors.New(" Not found ")
}

type Remover interface {
	Delete(id string) error
	CleanupRun(id string) (*Task, error)
}

func (t *tasksData) Delete(id string) error {
	var newStorage []*Task
	var found bool

	for _, tsk := range t.storage {
		if tsk.ID != id {

			newStorage = append(newStorage, tsk)

		} else {
			found = true
		}
	}

	if !found {
		return errors.New(" Not found ")
	}

	t.storage = newStorage

	return nil
}

func (t *tasksData) CleanupRun(id string) (*Task, error) {

	return t.GetTask(id)
}
