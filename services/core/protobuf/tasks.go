package protobuf

import (
	"context"
	"errors"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"telegram_boxes/services/core/app"
)

type Tasks interface {
	GetTask(ctx context.Context, r *GetTaskRequest) (*GetTaskResponse, error)
	FindTask(ctx context.Context, r *FindTaskRequest) (*FindTaskResponse, error)
	CheckTask(ctx context.Context, r *CheckTaskRequest) (*CheckTaskResponse, error)
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
