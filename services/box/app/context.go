package app

import (
	"context"
	"google.golang.org/grpc/metadata"
)

func SetCallContext(rId, userId string) context.Context {
	callContext := context.Background()
	mdOut := metadata.Pairs(
		"rid", rId,
		"userid", userId,
	)
	callContext = metadata.NewOutgoingContext(callContext, mdOut)
	return callContext
}

func GetDataContext(ctx context.Context) (rId, userId string) {
	mdIn, _ := metadata.FromIncomingContext(ctx)
	if len(mdIn["rid"]) > 0 {
		rId = mdIn["rid"][0]
	}
	if len(mdIn["userid"]) > 0 {
		userId = mdIn["userid"][0]
	}
	return
}