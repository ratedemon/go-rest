package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/ratedemon/go-rest/datastore/models"
	"golang.org/x/crypto/bcrypt"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/go-kit/kit/log"
	"github.com/ratedemon/go-rest/config"
	"github.com/ratedemon/go-rest/datastore/db"
	pbauth "github.com/ratedemon/go-rest/proto/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	if req.Password != req.ConfirmPassword {
		return nil, status.Errorf(codes.InvalidArgument, "`confirm_password` and`password` do not match")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.MinCost)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("Failed to hash password: %v", err))
	}

	user := models.User{
		Username: req.Username,
		Password: string(hash),
	}
	if err := as.db.CreateUser(&user); err != nil {
		return nil, status.Errorf(codes.Unknown, fmt.Sprintf("Failed to create new user: %v", err))
	}

	return &pbauth.SignupResponse{
		Message: "Successfully created",
	}, nil
}

func (as *AuthService) Login(ctx context.Context, req *pbauth.LoginRequest) (*pbauth.LoginResponse, error) {
	var user models.User
	if err := as.db.FindUserByUsername(req.Username, &user); err != nil {
		return nil, status.Errorf(codes.Unknown, fmt.Sprintf("Failed to find the user: %v", err))
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, status.Errorf(codes.Unknown, fmt.Sprint("Password do not match to exist"))
	}

	token, err := as.createToken(int64(user.ID))
	if err != nil {
		return nil, status.Errorf(codes.Unknown, fmt.Sprintf("Failed to create a token: %v", err))
	}

	return &pbauth.LoginResponse{
		Id:       int64(user.ID),
		Username: user.Username,
		Token:    token,
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
