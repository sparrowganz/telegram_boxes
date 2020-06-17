package task

import (
	"errors"
)

type Tasks interface {
	Getter
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
	TypeID     string
	Link       string
}

func CreateTasks() Tasks {
	return &tasksData{
		//todo remove debug structure
		storage: []*Task{
			{"1", "1", false, "1", "http://vk.com"},
			{"2", "2", true, "2", "http://vk.com"},
			{"3", "3", false, "3", "http://vk.com"},
			{"4", "4", true, "1", "http://vk.com"},
		},
	}
}

type Getter interface {
	GetTask(completedTask []string) (*Task, error)
}

func (t *tasksData) GetTask(completedTask []string) (*Task, error) {

	var notCompletedTask []*Task

	for _, tsk := range t.storage {
		for _, cId := range completedTask {
			if cId == tsk.ID {
				goto END
			}
		}
		notCompletedTask = append(notCompletedTask, tsk)
	END:
	}

	if len(notCompletedTask) == 0 {
		return nil, errors.New(" Not found ")
	}

	for _, task := range notCompletedTask {
		if task.IsPriority {
			return task, nil
		}
	}

	return notCompletedTask[0], nil
}

type Remover interface {
	CleanupRun(id string)
}

func (t *tasksData) CleanupRun(id string) {

}
