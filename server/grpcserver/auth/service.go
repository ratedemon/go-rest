package auth

import (
	"context"

	"github.com/go-kit/kit/log"
	pauth "github.com/ratedemon/go-rest/proto/auth"
	"google.golang.org/grpc"
)

type AuthService struct {
	log log.Logger
}

func NewAuthService(log log.Logger) *AuthService {
	return &AuthService{log}
}

func (as *AuthService) Login(ctx context.Context, req *pauth.LoginRequest) (*pauth.LoginResponse, error) {
	return &pauth.LoginResponse{
		Id:    1,
		Name:  "my_first_name",
		Email: "test_name@email.com",
	}, nil
}

func (as *AuthService) RegisterService(s *grpc.Server) {
	pauth.RegisterAuthServiceServer(s, as)
}
