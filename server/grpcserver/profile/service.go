package profile

import (
	"context"

	"github.com/ratedemon/go-rest/datastore/models"

	"github.com/go-kit/kit/log"

	"github.com/ratedemon/go-rest/config"
	"github.com/ratedemon/go-rest/datastore/db"
	"github.com/ratedemon/go-rest/grpcserver/helper"
	pbprofile "github.com/ratedemon/go-rest/proto/profile"
	"google.golang.org/grpc"
)

type ProfileService struct {
	cfg *config.Config
	log log.Logger
	db  *db.DB
}

func NewProfileService(cfg *config.Config, log log.Logger, db *db.DB) *ProfileService {
	return &ProfileService{cfg, log, db}
}

func (ps *ProfileService) Create(ctx context.Context, req *pbprofile.CreateRequest) (*pbprofile.CreateResponse, error) {
	userID, err := helper.GetUserID(ctx)
	if err != nil {
		return nil, err
	}
	profile := models.UserProfile{
		FirstName: req.Profile.FirstName,
		LastName:  req.Profile.LastName,
		Email:     req.Profile.Email,
		Age:       int16(req.Profile.Age),
		UserID:    userID,
		User: models.User{
			ID: uint(userID),
		},
	}
	if err := ps.db.CreateProfile(userID, &profile); err != nil {
		return nil, err
	}

	return &pbprofile.CreateResponse{
		Profile: &pbprofile.Profile{
			Id:        int64(profile.ID),
			FirstName: profile.FirstName,
			LastName:  profile.LastName,
			Email:     profile.Email,
			Age:       int64(profile.Age),
		},
	}, nil
}

func (ps *ProfileService) Update(ctx context.Context, req *pbprofile.UpdateRequest) (*pbprofile.UpdateResponse, error) {
	return &pbprofile.UpdateResponse{}, nil
}

func (ps *ProfileService) RegisterService(s *grpc.Server) {
	pbprofile.RegisterProfileServiceServer(s, ps)
}
