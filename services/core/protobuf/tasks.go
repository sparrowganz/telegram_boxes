package protobuf

import (
	"context"
	"errors"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"telegram_boxes/services/core/app"
	"telegram_boxes/services/core/app/models"
)

type Tasks interface {
	GetTask(ctx context.Context, r *GetTaskRequest) (*GetTaskResponse, error)
	FindTask(ctx context.Context, r *FindTaskRequest) (*FindTaskResponse, error)
	CheckTask(ctx context.Context, r *CheckTaskRequest) (*CheckTaskResponse, error)
	GetAllTask(ctx context.Context, r *GetAllTaskRequest) (*GetAllTaskResponse, error)
	ChangePriorityTask(ctx context.Context, r *ChangePriorityTaskRequest) (*ChangePriorityTaskResponse, error)
	DeleteTask(ctx context.Context, r *DeleteTaskRequest) (*DeleteTaskResponse, error)
	CleanupRunTask(ctx context.Context, r *CleanupRunTaskRequest) (*CleanupRunTaskResponse, error)
	CreateTask(ctx context.Context, r *Task) (*CreateTaskResponse, error)
}

func (sd *serverData) CreateTask(ctx context.Context, r *Task) (*CreateTaskResponse, error) {
	out := &CreateTaskResponse{}

	action, username := app.GetDataContext(ctx)

	session := sd.DB().GetMainSession().Clone()
	defer session.Close()

	task := models.CreateTask(r.GetName(), r.GetType(), r.GetLink(), r.GetIsPriority(), r.GetWithCheck())
	err := sd.DB().Models().Tasks().CreateTask(task, session)
	if err != nil {
		_ = sd.Log().Error(action, username, err.Error())
		return out, err
	}

	return out, nil
}

func (sd *serverData) CleanupRunTask(ctx context.Context,
	r *CleanupRunTaskRequest) (*CleanupRunTaskResponse, error) {
	out := &CleanupRunTaskResponse{}

	action, username := app.GetDataContext(ctx)

	session := sd.DB().GetMainSession().Clone()
	defer session.Close()

	if !bson.IsObjectIdHex(r.GetTaskID()) {
		return out, errors.New(" taskID is not bson objectID")
	}

	boxes, err := sd.DB().Models().Bots().GetAll(session)
	if err != nil {
		if err != mgo.ErrNotFound {
			_ = sd.Log().Error(action, username, err.Error())
		}

		return out, err
	}

	for _, box := range boxes {
		err = sd.Box().RemoveTask(box, r.GetTaskID())
		if err != nil {
			_ = sd.Admin().SendError("OK", "core", err.Error())
		}
	}


	return out, nil
}

func (sd *serverData) DeleteTask(ctx context.Context,
	r *DeleteTaskRequest) (*DeleteTaskResponse, error) {
	out := &DeleteTaskResponse{}

	action, username := app.GetDataContext(ctx)

	session := sd.DB().GetMainSession().Clone()
	defer session.Close()

	if !bson.IsObjectIdHex(r.GetTaskID()) {
		return out, errors.New(" taskID is not bson objectID")
	}

	boxes, err := sd.DB().Models().Bots().GetAll(session)
	if err != nil {
		if err != mgo.ErrNotFound {
			_ = sd.Log().Error(action, username, err.Error())
		}

		return out, err
	}

	for _, box := range boxes {
		err = sd.Box().RemoveTask(box, r.GetTaskID())
		if err != nil {
			_ = sd.Admin().SendError("OK", "core", err.Error())
		}
	}

	err = sd.DB().Models().Tasks().RemoveTask(bson.ObjectIdHex(r.GetTaskID()), session)
	if err != nil {
		_ = sd.Log().Error(action, username, err.Error())
		return out, err
	}

	return out, nil
}

func (sd *serverData) ChangePriorityTask(ctx context.Context,
	r *ChangePriorityTaskRequest) (*ChangePriorityTaskResponse, error) {
	out := &ChangePriorityTaskResponse{}

	action, username := app.GetDataContext(ctx)

	session := sd.DB().GetMainSession().Clone()
	defer session.Close()

	if !bson.IsObjectIdHex(r.GetTaskID()) {
		return out, errors.New(" taskID is not bson objectID")
	}

	task, err := sd.DB().Models().Tasks().FindTask(bson.ObjectIdHex(r.GetTaskID()), session)
	if err != nil {
		if err != mgo.ErrNotFound {
			_ = sd.Log().Error(action, username, err.Error())
		}

		return out, err
	}

	task.ChangePriority()
	err = sd.DB().Models().Tasks().UpdateTask(task, session)
	if err != nil {
		return out, err
	}

	out.Task = &Task{
		Id:         task.ID().Hex(),
		Name:       task.Title(),
		Type:       task.Type(),
		Link:       task.URL(),
		IsPriority: task.IsPriority(),
	}

	return out, nil
}

func (sd *serverData) GetAllTask(ctx context.Context, _ *GetAllTaskRequest) (*GetAllTaskResponse, error) {
	out := &GetAllTaskResponse{}

	action, username := app.GetDataContext(ctx)

	session := sd.DB().GetMainSession().Clone()
	defer session.Close()

	tasks, err := sd.DB().Models().Tasks().GetAllTasks(session)
	if err != nil {
		if err != mgo.ErrNotFound {
			_ = sd.Log().Error(action, username, err.Error())
		}

		return out, err
	}

	for _, task := range tasks {
		out.Tasks = append(out.Tasks, &Task{
			Id:   task.ID().Hex(),
			Name: task.Title(),
			Type: task.Type(),
			Link: task.URL(),
		})
	}

	return out, nil
}

func (sd *serverData) GetTask(ctx context.Context, r *GetTaskRequest) (*GetTaskResponse, error) {
	out := &GetTaskResponse{}

	action, username := app.GetDataContext(ctx)

	session := sd.DB().GetMainSession().Clone()
	defer session.Close()

	var bsonTasksId []bson.ObjectId

	for _, id := range r.GetTasksID() {
		if !bson.IsObjectIdHex(id) {
			continue
		}

		bsonTasksId = append(bsonTasksId, bson.ObjectIdHex(id))
	}

	task, err := sd.DB().Models().Tasks().GetNextTask(bsonTasksId, session)
	if err != nil {
		if err != mgo.ErrNotFound {
			_ = sd.Log().Error(action, username, err.Error())
		}

		return out, err
	}

	out.Task = &Task{
		Id:   task.ID().Hex(),
		Name: task.Title(),
		Type: task.Type(),
		Link: task.URL(),
	}

	return out, nil
}

func (sd *serverData) FindTask(ctx context.Context, r *FindTaskRequest) (*FindTaskResponse, error) {
	out := &FindTaskResponse{}

	action, username := app.GetDataContext(ctx)

	session := sd.DB().GetMainSession().Clone()
	defer session.Close()

	if !bson.IsObjectIdHex(r.GetId()) {
		return out, errors.New(" id is not bson objectID")
	}

	task, err := sd.DB().Models().Tasks().FindTask(bson.ObjectIdHex(r.GetId()), session)
	if err != nil {
		if err != mgo.ErrNotFound {
			_ = sd.Log().Error(action, username, err.Error())
		}

		return out, err
	}

	out.Task = &Task{
		Id:   task.ID().Hex(),
		Name: task.Title(),
		Type: task.Type(),
		Link: task.URL(),
	}

	return out, nil
}

func (sd *serverData) CheckTask(ctx context.Context, r *CheckTaskRequest) (*CheckTaskResponse, error) {
	out := &CheckTaskResponse{}

	action, username := app.GetDataContext(ctx)

	session := sd.DB().GetMainSession().Clone()
	defer session.Close()

	if !bson.IsObjectIdHex(r.GetTaskID()) {
		return out, errors.New(" id is not bson objectID")
	}

	task, err := sd.DB().Models().Tasks().FindTask(bson.ObjectIdHex(r.GetTaskID()), session)
	if err != nil {
		if err != mgo.ErrNotFound {
			_ = sd.Log().Error(action, username, err.Error())
		}

		return out, err
	}

	if task.WithCheck() {
		//todo check task
	}else{
		out.IsCheck = true
	}

	return out, nil
}
