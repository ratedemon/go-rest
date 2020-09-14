package helper

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"google.golang.org/grpc/metadata"
)

// GetUserID returns userID from context
func GetUserID(ctx context.Context) (int64, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		return 0, errors.New("Failed to get metadata")
	}

	var userIdStr string
	if userID, ok := md["user_id"]; ok {
		userIdStr = strings.Join(userID, "")
	} else {
		return 0, errors.New("'user_id' is not found in context")
	}
	id, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		return 0, err
	}

	return id, nil
}
