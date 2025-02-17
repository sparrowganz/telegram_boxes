package app

import (
	"context"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	ServiceName      = "Core"
	ServiceShortName = "CORE-"
)

func SetCallContext(action, username string) context.Context {
	callContext := context.Background()
	mdOut := metadata.Pairs(
		"action", action,
		"username", username,
	)
	callContext = metadata.NewOutgoingContext(callContext, mdOut)
	return callContext
}

func SetCallContextWithContext(callContext context.Context,action, username string) context.Context {
	mdOut := metadata.Pairs(
		"action", action,
		"username", username,
	)
	callContext = metadata.NewOutgoingContext(callContext, mdOut)
	return callContext
}

func GetDataContext(ctx context.Context) (action, username string) {
	mdIn, _ := metadata.FromIncomingContext(ctx)
	if len(mdIn["action"]) > 0 {
		action = mdIn["action"][0]
	}
	if len(mdIn["username"]) > 0 {
		username = mdIn["username"][0]
	}
	return
}

func ParseGRPCError(err error) string {
	return status.Convert(err).Message()
}

