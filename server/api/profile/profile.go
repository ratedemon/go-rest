package profile

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-kit/kit/log"
	"github.com/gorilla/mux"

	"github.com/ratedemon/go-rest/api/helper"
	protoprofile "github.com/ratedemon/go-rest/proto/profile"
	"google.golang.org/grpc"
)

// ProfileHandler is gateway for profile endpoints
type ProfileHandler struct {
	ctx        context.Context
	log        log.Logger
	grpcClient protoprofile.ProfileServiceClient
}

type profileBody struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Age       int64  `json:"age"`
	Email     string `json:"email"`
	Sex       string `json:"sex,omitempty"`
}

func (ph *ProfileHandler) RegisterRoutes() []helper.Route {
	return []helper.Route{
		{"/profile", "POST", ph.create},
		{"/profile/{id:[0-9]+}", "PUT", ph.update},
	}
}

func (ph *ProfileHandler) create(ctx context.Context, req *http.Request) (interface{}, error) {
	var body profileBody
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		return nil, err
	}
	profileSex := protoprofile.Sex_UNKNOWN
	if body.Sex != "" {
		ps, ok := protoprofile.Sex_value[strings.ToLower(body.Sex)]
		if !ok {
			return nil, errors.New("Sex is not defined")
		}
		profileSex = protoprofile.Sex(ps)
	}

	res, err := ph.grpcClient.Create(ctx, &protoprofile.CreateRequest{
		Profile: &protoprofile.Profile{
			FirstName: body.FirstName,
			LastName:  body.LastName,
			Age:       body.Age,
			Email:     body.Email,
			Sex:       profileSex,
		},
	})
	if err != nil {
		return nil, err
	}

	return res.Profile, nil
}

func (ph *ProfileHandler) update(ctx context.Context, req *http.Request) (interface{}, error) {
	vars := mux.Vars(req)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return nil, err
	}

	var body profileBody
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		return nil, err
	}
	res, err := ph.grpcClient.Update(ctx, &protoprofile.UpdateRequest{
		Profile: &protoprofile.Profile{
			Id:        id,
			FirstName: body.FirstName,
			LastName:  body.LastName,
			Age:       body.Age,
			Email:     body.Email,
		},
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func NewProfileHandler(ctx context.Context, log log.Logger, grpcConn *grpc.ClientConn) *ProfileHandler {
	client := protoprofile.NewProfileServiceClient(grpcConn)
	return &ProfileHandler{ctx, log, client}
}
