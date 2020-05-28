package protobuf

import (
	"context"
	"telegram_boxes/services/logs/app"
)

type Server struct {
	Logger app.Logger
}

func (s *Server) AccessLog(ctx context.Context, r *AccessLogRequest) (*AccessLogResponse, error) {
	out := &AccessLogResponse{}

	s.Logger.Access(
		r.GetTime(),
		r.GetServerName(),
		r.GetMethod(), r.GetRequestId(), r.GetUser(), r.GetDuration())
	return out, nil
}

func (s *Server) ErrorLog(ctx context.Context, r *ErrorLogRequest) (*ErrorLogResponse, error) {
	out := &ErrorLogResponse{}
	s.Logger.Error(r.GetTime(), r.GetServerName(), r.GetRequestId(), r.GetError())
	return out, nil
}

func (s *Server) SystemLog(ctx context.Context, r *SystemLogRequest) (*SystemLogResponse, error) {
	out := &SystemLogResponse{}
	s.Logger.System(r.GetTime(), r.GetServerName(), r.GetData())
	return out, nil
}
