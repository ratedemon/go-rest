package helper

import (
	"context"
	"strconv"

	"github.com/ratedemon/go-rest/api/middleware"
	"google.golang.org/grpc/metadata"
)

// contextWithUserID adds user_id to context for grpc calls
func contextWithUserID(ctx context.Context, userID int64) context.Context {
	id := strconv.Itoa(int(userID))
	md := metadata.Pairs(
		middleware.UserIDKey, id,
	)

	return metadata.NewOutgoingContext(ctx, md)
}
