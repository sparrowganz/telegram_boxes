package protobuf

import (
	"context"
	"google.golang.org/grpc/metadata"
)

func SetCallContext(rID string) context.Context {
	callContext := context.Background()
	mdOut := metadata.Pairs(
		"rid", rID,
	)
	callContext = metadata.NewOutgoingContext(callContext, mdOut)
	return callContext
}

func GetDataContext(ctx context.Context) (rID string) {
	mdIn, _ := metadata.FromIncomingContext(ctx)
	if len(mdIn["rid"]) > 0 {
		rID = mdIn["rid"][0]
	}
	return
}
