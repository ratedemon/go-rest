package user

import (
	"context"

	"github.com/go-kit/kit/log"
	"google.golang.org/grpc"

	"github.com/ratedemon/go-rest/config"
	"github.com/ratedemon/go-rest/datastore/db"
	"github.com/ratedemon/go-rest/datastore/models"
	"github.com/ratedemon/go-rest/grpcserver/helper"
	pbuser "github.com/ratedemon/go-rest/proto/user"
)

type UserService struct {
	cfg *config.Config
	log log.Logger
	db  *db.DB
}

func (us *UserService) Get(ctx context.Context, req *pbuser.GetRequest) (*pbuser.GetResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	user, err := us.db.FindUserById(req.Id)
	if err != nil {
		return nil, err
	}

	return &pbuser.GetResponse{
		User: toUserResponse(user),
	}, nil
}

func (us *UserService) List(ctx context.Context, _ *pbuser.ListRequest) (*pbuser.ListResponse, error) {
	userID, err := helper.GetUserID(ctx)
	if err != nil {
		return nil, err
	}

	users, err := us.db.FindAllUsers(userID)
	if err != nil {
		return nil, err
	}

	result := make([]*pbuser.User, len(users))
	for i, v := range users {
		result[i] = toUserResponse(&v)
	}

	return &pbuser.ListResponse{
		Users: result,
	}, nil
}

func toUserResponse(u *models.User) *pbuser.User {
	result := &pbuser.User{
		Id:       int64(u.ID),
		Username: u.Username,
	}
	result.FirstName = u.Profile.FirstName
	result.LastName = u.Profile.LastName
	result.Age = int64(u.Profile.Age)
	result.Email = u.Profile.Email
	result.Sex = u.Profile.Sex
	result.ImagePath = u.Image.Path

	return result
}

func NewUserService(cfg *config.Config, log log.Logger, db *db.DB) *UserService {
	return &UserService{cfg, log, db}
}

func (us *UserService) RegisterService(s *grpc.Server) {
	pbuser.RegisterUserServiceServer(s, us)
}
