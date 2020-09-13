package auth

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-kit/kit/log"
	"github.com/ratedemon/go-rest/api/helper"
	protoauth "github.com/ratedemon/go-rest/proto/auth"
	"google.golang.org/grpc"
)

// AuthHandler is gateway for auth endpoints
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

type loginBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (ah *AuthHandler) login(ctx context.Context, req *http.Request) (interface{}, error) {
	var body loginBody
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		return nil, err
	}
	if body.Username == "" || body.Password == "" {
		return nil, errors.New("Required field is missing")
	}
	res, err := ah.grpcClient.Login(ctx, &protoauth.LoginRequest{
		Username: body.Username,
		Password: body.Password,
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

type signUpBody struct {
	loginBody
	ConfirmPassword string `json:"confirm_password"`
}

func (ah *AuthHandler) signup(ctx context.Context, req *http.Request) (interface{}, error) {
	var body signUpBody
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		return nil, err
	}
	if body.Username == "" || body.Password == "" || body.ConfirmPassword == "" {
		return nil, errors.New("Required field is missing")
	}

	res, err := ah.grpcClient.Signup(ctx, &protoauth.SignupRequest{
		Username:        body.Username,
		Password:        body.Password,
		ConfirmPassword: body.ConfirmPassword,
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}
