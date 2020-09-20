package helper

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// GetUserID returns userID from context
func GetUserID(ctx context.Context) (int64, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return 0, status.Errorf(codes.NotFound, "Failed to get metadata")
	}

	var userIDStr string
	if userID, ok := md["user_id"]; ok {
		userIDStr = strings.Join(userID, "")
	} else {
		return 0, status.Errorf(codes.NotFound, "'user_id' is not found in context")
	}
	id, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		return 0, status.Errorf(codes.InvalidArgument, fmt.Sprintf("Failed to create new user: %v", err))
	}

	return id, nil
}
