package protobuf

import (
	"context"
	"errors"
	"google.golang.org/grpc/peer"
	"gopkg.in/mgo.v2"
	"strings"
	"telegram_boxes/services/core/app"
	"telegram_boxes/services/core/app/models"
)

type Servers interface {
	InitBox(ctx context.Context, r *InitBoxRequest) (*InitBoxResponse, error)
	SendError(ctx context.Context, r *SendErrorRequest) (*SendErrorResponse, error)
}

func (sd *serverData) SendError(ctx context.Context, r *SendErrorRequest) (*SendErrorResponse, error) {
	out := &SendErrorResponse{}

	action, username := app.GetDataContext(ctx)

	session := sd.DB().GetMainSession().Clone()
	defer session.Close()

	bot, err := sd.DB().Models().Bots().FindByUsername(username, session)
	if err != nil {
		_ = sd.Log().Error(action, username, err.Error())
		return out, err
	}

	bot.SetStatus(r.GetStatus().String())
	err = sd.DB().Models().Bots().UpdateBot(bot, session)

	//todo Send ADMIN error
	return out, nil
}

func (sd *serverData) InitBox(ctx context.Context, r *InitBoxRequest) (*InitBoxResponse, error) {

	out := &InitBoxResponse{}

	action, username := app.GetDataContext(ctx)

	if r.GetUsername() == "" {
		return out, errors.New("@username does not exist")
	}

	p, ok := peer.FromContext(ctx)
	if !ok {
		return out, errors.New(" peer not found ")
	}

	parts := strings.Split(p.Addr.String(), ":")
	if len(parts) != 2 {
		return out, errors.New(p.Addr.String() + " is not valid address")
	}

	session := sd.DB().GetMainSession().Clone()
	defer session.Close()

	bot, err := sd.DB().Models().Bots().FindByUsername(r.GetUsername(), session)

	if err != nil {
		if err != mgo.ErrNotFound {
			_ = sd.Log().Error(action, username, err.Error())
			return out, err
		}

		bot = models.CreateBot(parts[0], parts[1])
		bot.SetUsername(r.GetUsername())

		err = sd.DB().Models().Bots().CreateBot(bot, session)
		if err != nil {
			_ = sd.Log().Error(action, username, err.Error())
			return out, err
		}
	}

	bot.SetActive()
	bot.SetStatus(Status_OK.String())
	if bot.Address().Addr() != parts[0]+":"+parts[1] {
		bot.Address().SetIP(parts[0])
		bot.Address().SetPort(parts[1])
	}
	err = sd.DB().Models().Bots().UpdateBot(bot, session)

	out.ID = bot.ID().Hex()

	return out, nil
}
