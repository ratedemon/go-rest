package auth

import (
	"context"
	"net/http"

	"github.com/go-kit/kit/log"
	"github.com/ratedemon/go-rest/api/helper"
	protoauth "github.com/ratedemon/go-rest/proto/auth"
	"google.golang.org/grpc"
)

type AuthHandler struct {
	ctx        context.Context
	log        log.Logger
	grpcClient protoauth.AuthServiceClient
}

func NewAuthHandler(ctx context.Context, log log.Logger, grpcConn *grpc.ClientConn) *AuthHandler {
	client := protoauth.NewAuthServiceClient(grpcConn)
	return &AuthHandler{ctx, log, client}
}

func (ah *AuthHandler) RegisterRoutes() []helper.Route {
	return []helper.Route{
		{"/login", "POST", ah.login},
		{"/signup", "POST", ah.signup},
	}
}

func (ah *AuthHandler) login(ctx context.Context, req *http.Request) (interface{}, error) {
	res, err := ah.grpcClient.Login(ctx, &protoauth.LoginRequest{})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (ah *AuthHandler) signup(ctx context.Context, req *http.Request) (interface{}, error) {
	return "method is not implemented", nil
}
