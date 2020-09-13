package auth

import (
	"context"
	"time"

	"github.com/ratedemon/go-rest/datastore/models"
	"golang.org/x/crypto/bcrypt"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/go-kit/kit/log"
	"github.com/ratedemon/go-rest/config"
	"github.com/ratedemon/go-rest/datastore/db"
	pbauth "github.com/ratedemon/go-rest/proto/auth"
	"google.golang.org/grpc"
)

type AuthService struct {
	cfg *config.Config
	log log.Logger
	db  *db.DB
}

func NewAuthService(cfg *config.Config, log log.Logger, db *db.DB) *AuthService {
	return &AuthService{cfg, log, db}
}

func (as *AuthService) Signup(ctx context.Context, req *pbauth.SignupRequest) (*pbauth.SignupResponse, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.MinCost)
	if err != nil {
		return nil, err
	}

	user := models.User{
		Username: req.Username,
		Password: string(hash),
	}
	if err := as.db.CreateUser(&user); err != nil {
		return nil, err
	}

	return &pbauth.SignupResponse{
		Message: "Successfully created",
	}, nil
}

func (as *AuthService) Login(ctx context.Context, req *pbauth.LoginRequest) (*pbauth.LoginResponse, error) {
	var user models.User
	if err := as.db.FindUserByUsername(req.Username, &user); err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, err
	}

	token, err := as.createToken(int64(user.ID))
	if err != nil {
		return nil, err
	}

	return &pbauth.LoginResponse{
		Id:        int64(user.ID),
		FirstName: "my_first_name",
		Email:     "test_name@email.com",
		Token:     token,
	}, nil
}

func (as *AuthService) createToken(userID int64) (string, error) {
	atClaims := jwt.MapClaims{
		"authorized": true,
		"user_id":    userID,
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
	pbauth.RegisterAuthServiceServer(s, as)
}
