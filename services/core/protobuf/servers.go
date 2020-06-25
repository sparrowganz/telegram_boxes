package protobuf

import (
	"context"
	"errors"
	"google.golang.org/grpc/peer"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"strings"
	"sync"
	"telegram_boxes/services/core/app"
	"telegram_boxes/services/core/app/models"
)

type Servers interface {
	InitBox(ctx context.Context, r *InitBoxRequest) (*InitBoxResponse, error)
	SendError(ctx context.Context, r *SendErrorRequest) (*SendErrorResponse, error)
	GetListServers(ctx context.Context, r *GetListServersRequest) (*GetListServersResponse, error)
	GetServer(ctx context.Context, r *GetServerRequest) (*GetServerResponse, error)
	GetAllUsersCount(ctx context.Context, r *GetAllUsersCountRequest) (*GetAllUsersCountResponse, error)
	ChangeBonusActive(ctx context.Context, r *ChangeBonusActiveRequest) (*ChangeBonusActiveResponse, error)
	HardCheck(r *HardCheckRequest, server Servers_HardCheckServer) error
}

func (sd *serverData) HardCheck(r *HardCheckRequest, stream Servers_HardCheckServer) error {

	action, username := app.GetDataContext(stream.Context())

	session := sd.DB().GetMainSession().Clone()
	defer session.Close()

	bots, err := sd.DB().Models().Bots().GetAll(session)
	if err != nil {
		_ = sd.Log().Error(action, username, err.Error())
		return err
	}

	chResult := make(chan *Check, 100)

	wg := &sync.WaitGroup{}

	for  range bots {

		wg.Add(1)
		go func() {
			defer wg.Done()
			//todo send request to box


		}()
	}

	readWg := &sync.WaitGroup{}
	readWg.Add(1)
	go func() {
		defer readWg.Done()
		for res := range chResult {
			_ = stream.Send(res)
		}
	}()

	wg.Wait()
	close(chResult)
	readWg.Wait()

	return nil
}

func (sd *serverData) ChangeBonusActive(ctx context.Context,
	r *ChangeBonusActiveRequest) (*ChangeBonusActiveResponse, error) {
	out := &ChangeBonusActiveResponse{}

	action, username := app.GetDataContext(ctx)

	session := sd.DB().GetMainSession().Clone()
	defer session.Close()

	var bots []*models.BotData

	if r.GetId() == "all" {
		var err error
		bots, err = sd.DB().Models().Bots().GetAll(session)
		if err != nil {
			_ = sd.Log().Error(action, username, err.Error())
			return out, err
		}
	} else {
		bot, err := sd.DB().Models().Bots().FindByID(bson.ObjectIdHex(r.GetId()), session)
		if err != nil {
			_ = sd.Log().Error(action, username, err.Error())
			return out, err
		}
		bots = append(bots, bot.Object())
	}

	var setStatus bool

	for _, bot := range bots {
		if !bot.Bonus().IsActive() {
			setStatus = true
		}
	}

	for _, bot := range bots {
		if setStatus {
			bot.Bonus().SetActive()
		} else {
			bot.Bonus().Inactive()
		}
		_ = sd.DB().Models().Bots().UpdateBot(bot, session)
	}

	return out, nil
}

func (sd *serverData) GetAllUsersCount(ctx context.Context,
	r *GetAllUsersCountRequest) (*GetAllUsersCountResponse, error) {
	out := &GetAllUsersCountResponse{}

	action, username := app.GetDataContext(ctx)

	session := sd.DB().GetMainSession().Clone()
	defer session.Close()

	bots, err := sd.DB().Models().Bots().GetAll(session)
	if err != nil {
		_ = sd.Log().Error(action, username, err.Error())
		return out, err
	}

	for _, bot := range bots {
		//todo get new info from bots

		out.Counts = append(out.Counts, &Counts{
			Id:       bot.ID().Hex(),
			Username: bot.Username(),
			Old: &Count{
				All:     bot.Statistics().GetAll(),
				Blocked: bot.Statistics().GetAll(),
			},
			New: &Count{
				All:     0,
				Blocked: 0,
			},
			Current: 0,
		})

		bot.Statistics().SetAll(0)
		bot.Statistics().SetBlocked(0)
		_ = sd.DB().Models().Bots().UpdateBot(bot, session)
	}

	return out, nil
}

func (sd *serverData) GetServer(ctx context.Context, r *GetServerRequest) (*GetServerResponse, error) {
	out := &GetServerResponse{}

	action, username := app.GetDataContext(ctx)

	session := sd.DB().GetMainSession().Clone()
	defer session.Close()

	bot, err := sd.DB().Models().Bots().FindByID(bson.ObjectIdHex(r.GetId()), session)
	if err != nil {
		_ = sd.Log().Error(action, username, err.Error())
		return out, err
	}

	out.Server = &Server{
		Id:       bot.ID().Hex(),
		Username: bot.Username(),
		Status:   bot.BotStatus(),
		IsActive: bot.IsActive(),
		Bonus: &Bonus{
			IsActive: bot.Bonus().IsActive(),
			Cost:     bot.Bonus().Cost(),
			Time:     bot.Bonus().InTime().UnixNano(),
		},
	}

	return out, nil
}

func (sd *serverData) GetListServers(ctx context.Context, r *GetListServersRequest) (*GetListServersResponse, error) {
	out := &GetListServersResponse{}

	action, username := app.GetDataContext(ctx)

	session := sd.DB().GetMainSession().Clone()
	defer session.Close()

	bots, err := sd.DB().Models().Bots().GetAll(session)
	if err != nil {
		_ = sd.Log().Error(action, username, err.Error())
		return out, err
	}

	for _, bot := range bots {
		out.Servers = append(out.Servers, &Server{
			Id:       bot.ID().Hex(),
			Username: bot.Username(),
			Status:   bot.BotStatus(),
			IsActive: bot.IsActive(),
			Bonus: &Bonus{
				IsActive: bot.Bonus().IsActive(),
				Cost:     bot.Bonus().Cost(),
				Time:     bot.Bonus().InTime().UnixNano(),
			},
		})
	}

	return out, nil
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

	_ = sd.Admin().SendError(r.GetStatus().String(), username, r.GetError())

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

		_ = sd.Admin().SendError("START", r.GetUsername(), "New box in system")
	}

	bot.SetActive()
	bot.SetStatus(Status_OK.String())
	if bot.Address().Addr() != parts[0]+":"+parts[1] {
		bot.Address().SetIP(parts[0])
		bot.Address().SetPort(parts[1])
	}
	err = sd.DB().Models().Bots().UpdateBot(bot, session)

	out.ID = bot.ID().Hex()
	_ = sd.Admin().SendError("UP", r.GetUsername(), "OLD Box start again")

	return out, nil
}
