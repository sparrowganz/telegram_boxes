package protobuf

import "context"

func (s *Server) AccessLog(ctx context.Context, r *AccessLogRequest) (*AccessLogResponse, error) {
	out := &AccessLogResponse{}
	s.Logger.Access(r.Time, r.ServerName, r.Method, r.RequestId, r.User, r.Duration)
	return out, nil
}

func (s *Server) ErrorLog(ctx context.Context, r *ErrorLogRequest) (*ErrorLogResponse, error) {
	out := &ErrorLogResponse{}
	s.Logger.Error(r.Time, r.ServerName, r.RequestId, r.Error)
	return out, nil
}

func (s *Server) SystemLog(ctx context.Context, r *SystemLogRequest) (*SystemLogResponse, error) {
	out := &SystemLogResponse{}
	s.Logger.System(r.Time, r.ServerName, r.Data)
	return out, nil
}
