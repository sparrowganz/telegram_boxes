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
	WithCheck  bool
	Type       string
	Link       string
}

func CreateTasks() Tasks {
	return &tasksData{
		//todo remove debug structure
		storage: []*Task{
			{"1", "1", false, false, "channel", "http://vk.com"},
			{"2", "2", true, true, "checkChannel", "http://vk.com"},
			{"3", "3", false, false, "subscribeInstagram", "http://vk.com"},
			{"4", "4", true, false, "likeInstagram", "http://vk.com"},
			{"5", "5", false, false, "openWeb", "http://vk.com"},
			{"6", "6", false, false, "activateBot", "http://vk.com"},
		},
	}
}

type Getter interface {
	GetTask(completedTask []string) (*Task, error)
	FindTask(id string) (*Task, error)
	CheckTask(chatID int64, taskID string) (bool, error)
}

func (t *tasksData) CheckTask(chatID int64, taskID string) (bool, error) {
	for _, task := range t.storage {
		if task.ID == taskID {
			if task.WithCheck {
				return false, nil
			} else {
				return true, nil
			}
		}
	}
	return false, errors.New(" Task not found ")
}

func (t *tasksData) FindTask(id string) (*Task, error) {
	for _, task := range t.storage {
		if task.ID == id {
			return task, nil
		}
	}
	return nil, errors.New(" Task not found ")
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
