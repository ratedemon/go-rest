package user

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-kit/kit/log"
	"github.com/gorilla/mux"
	"google.golang.org/grpc"

	"github.com/ratedemon/go-rest/api/helper"
	protouser "github.com/ratedemon/go-rest/proto/user"
)

// UserHandler is gateway for user endpoints
type UserHandler struct {
	ctx        context.Context
	log        log.Logger
	grpcClient protouser.UserServiceClient
}

func (uh *UserHandler) RegisterRoutes() []helper.Route {
	return []helper.Route{
		{"/users", "GET", uh.list},
		{"/users/{id:[0-9]+}", "GET", uh.get},
	}
}

func (uh *UserHandler) list(ctx context.Context, _ *http.Request) (interface{}, error) {
	res, err := uh.grpcClient.List(ctx, &protouser.ListRequest{})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (uh *UserHandler) get(ctx context.Context, req *http.Request) (interface{}, error) {
	vars := mux.Vars(req)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return nil, err
	}

	res, err := uh.grpcClient.Get(ctx, &protouser.GetRequest{
		Id: id,
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func NewUserHandler(ctx context.Context, log log.Logger, grpcConn *grpc.ClientConn) *UserHandler {
	client := protouser.NewUserServiceClient(grpcConn)
	return &UserHandler{ctx, log, client}
}
