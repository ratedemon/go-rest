package auth

import (
	"context"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/go-kit/kit/log"
	"github.com/ratedemon/go-rest/config"
	pauth "github.com/ratedemon/go-rest/proto/auth"
	"google.golang.org/grpc"
)

type AuthService struct {
	cfg *config.Config
	log log.Logger
}

func NewAuthService(cfg *config.Config, log log.Logger) *AuthService {
	return &AuthService{cfg, log}
}

func (as *AuthService) Login(ctx context.Context, req *pauth.LoginRequest) (*pauth.LoginResponse, error) {
	return &pauth.LoginResponse{
		Id:    1,
		Name:  "my_first_name",
		Email: "test_name@email.com",
	}, nil
}

func (as *AuthService) createToken(userId int64) (string, error) {
	atClaims := jwt.MapClaims{
		"authorized": true,
		"user_id":    userId,
		"exp":        time.Now().Add(time.Minute * time.Duration(as.cfg.ExpLoginTimeout)).Unix(),
	}
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(as.cfg.JWTSecret))
	if err != nil {
		return "", err
	}

	return token, nil
}

func (as *AuthService) RegisterService(s *grpc.Server) {
	pauth.RegisterAuthServiceServer(s, as)
}
