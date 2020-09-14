package helper

import (
	"context"
	"strconv"

	"google.golang.org/grpc/metadata"
)

// contextWithUserID adds user_id to context for grpc calls
func contextWithUserID(ctx context.Context, userID int64) context.Context {
	id := strconv.Itoa(int(userID))
	md := metadata.Pairs(
		"user_id", id,
	)

	return metadata.NewIncomingContext(ctx, md)
}
