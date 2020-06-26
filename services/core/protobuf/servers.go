package protobuf

import (
	"context"
	"errors"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"sync"
	"telegram_boxes/services/core/app"
	"telegram_boxes/services/core/app/models"
	boxProto "telegram_boxes/services/core/protobuf/services/box/protobuf"
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

	for _, bot := range bots {

		wg.Add(1)
		currentBot := bot
		go func() {
			defer wg.Done()
			status, errCheck := sd.Box().CheckBox(currentBot, r.GetUserID())
			if errCheck != nil {
				currentBot.SetStatus(status)
				_ = sd.DB().Models().Bots().UpdateBot(currentBot, session)
				_ = sd.Admin().SendError(status, currentBot.Username(), errCheck.Error())
				return
			}

			currentBot.SetStatus(status)
			_ = sd.DB().Models().Bots().UpdateBot(currentBot, session)
			_ = sd.Admin().SendError(status, currentBot.Username(), "Check Successful")

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

		var stats = &boxProto.Statistic{}
		if bot.IsActive() {
			stats, _ = sd.Box().GetStats(bot)
		}

		out.Counts = append(out.Counts, &Counts{
			Id:       bot.ID().Hex(),
			Username: bot.Username(),
			Old: &Count{
				All:     bot.Statistics().GetAll(),
				Blocked: bot.Statistics().GetAll(),
			},
			New: &Count{
				All:     stats.GetAll(),
				Blocked: stats.GetBlocked(),
			},
			Current: stats.GetCurrent(),
		})

		if stats.GetAll() > 0 {
			bot.Statistics().SetAll(stats.GetAll())
		}
		if stats.GetBlocked() > 0 {
			bot.Statistics().SetBlocked(stats.GetBlocked())
		}
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

	if r.GetStatus() != Status_OK {
		bot.InActive()
	}

	bot.SetStatus(r.GetStatus().String())
	err = sd.DB().Models().Bots().UpdateBot(bot, session)

	_ = sd.Admin().SendError(r.GetStatus().String(), username, r.GetError())

	return out, nil
}

func (sd *serverData) InitBox(ctx context.Context, r *InitBoxRequest) (*InitBoxResponse, error) {

	out := &InitBoxResponse{}

	action, username := app.GetDataContext(ctx)

	if r.GetUsername() == "" || r.GetHost() == "" || r.GetPort() == "" {
		return out, errors.New("init data does not exist")
	}

	session := sd.DB().GetMainSession().Clone()
	defer session.Close()

	bot, err := sd.DB().Models().Bots().FindByUsername(r.GetUsername(), session)

	if err != nil {
		if err != mgo.ErrNotFound {
			_ = sd.Log().Error(action, username, err.Error())
			return out, err
		}

		bot = models.CreateBot(r.GetHost(), r.GetPort())
		bot.SetUsername(r.GetUsername())

		err = sd.DB().Models().Bots().CreateBot(bot, session)
		if err != nil {
			_ = sd.Log().Error(action, username, err.Error())
			return out, err
		}

		_ = sd.Box().AddBox(bot)
		_ = sd.Admin().SendError("START", r.GetUsername(), "New box in system")
		return out, nil
	}

	bot.SetActive()
	bot.SetStatus(Status_OK.String())
	if bot.Address().Addr() != r.GetHost()+":"+r.GetPort() {
		bot.Address().SetIP(r.GetHost())
		bot.Address().SetPort(r.GetPort())
	}
	_ = sd.DB().Models().Bots().UpdateBot(bot, session)
	_ = sd.Box().AddBox(bot)

	out.ID = bot.ID().Hex()
	_ = sd.Admin().SendError("UP", r.GetUsername(), "OLD Box start again")

	return out, nil
}
