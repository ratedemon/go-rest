package profile

import (
	"context"
	"strings"

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
	if err := req.Validate(); err != nil {
		return nil, err
	}

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
		Sex:       strings.ToLower(req.Profile.Sex.String()),
	}
	if err := ps.db.CreateProfile(userID, &profile); err != nil {
		return nil, err
	}

	createdProfile := &pbprofile.Profile{
		Id:        int64(profile.ID),
		FirstName: profile.FirstName,
		LastName:  profile.LastName,
		Email:     profile.Email,
		Age:       int64(profile.Age),
	}

	if profile.Sex != "" {
		createdProfile.Sex = pbprofile.Sex(pbprofile.Sex_value[strings.ToUpper(profile.Sex)])
	}

	return &pbprofile.CreateResponse{
		Profile: createdProfile,
	}, nil
}

func (ps *ProfileService) Update(ctx context.Context, req *pbprofile.UpdateRequest) (*pbprofile.UpdateResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	userID, err := helper.GetUserID(ctx)
	if err != nil {
		return nil, err
	}

	profile, err := ps.db.FindProfile(userID)
	if err != nil {
		return nil, err
	}

	if req.Profile.FirstName != "" {
		profile.FirstName = req.Profile.FirstName
	}
	if req.Profile.LastName != "" {
		profile.LastName = req.Profile.LastName
	}
	if req.Profile.Email != "" {
		profile.Email = req.Profile.Email
	}
	if req.Profile.Age != 0 {
		profile.Age = int16(req.Profile.Age)
	}
	if req.Profile.Sex.String() != "" {
		profile.Sex = strings.ToLower(req.Profile.Sex.String())
	}

	if err := ps.db.UpdateProfile(profile); err != nil {
		return nil, err
	}

	updatedProfile := &pbprofile.Profile{
		Id:        int64(profile.ID),
		FirstName: profile.FirstName,
		LastName:  profile.LastName,
		Email:     profile.Email,
		Age:       int64(profile.Age),
	}

	if profile.Sex != "" {
		updatedProfile.Sex = pbprofile.Sex(pbprofile.Sex_value[strings.ToUpper(profile.Sex)])
	}

	return &pbprofile.UpdateResponse{
		Profile: updatedProfile,
	}, nil
}

func (ps *ProfileService) RegisterService(s *grpc.Server) {
	pbprofile.RegisterProfileServiceServer(s, ps)
}
