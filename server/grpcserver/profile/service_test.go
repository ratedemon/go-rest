package profile_test

import (
	"context"
	"strconv"
	"testing"

	"github.com/go-kit/kit/log"
	"github.com/ratedemon/go-rest/config"
	maindb "github.com/ratedemon/go-rest/datastore/db"
	"github.com/ratedemon/go-rest/datastore/models"
	"github.com/ratedemon/go-rest/dbtesting"
	"github.com/ratedemon/go-rest/grpcserver/profile"
	pbprofile "github.com/ratedemon/go-rest/proto/profile"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/metadata"
)

func TestProfileService(t *testing.T) {
	t.Run("create", profileTest(func(t *tester) {
		ctx := contextWithUserID(t.ctx, int(t.user.ID))

		res, err := t.svc.Create(ctx, &pbprofile.CreateRequest{
			Profile: &pbprofile.Profile{
				FirstName: "John",
				LastName:  "Doe",
				Age:       48,
				Email:     "johndoe@gmail.com",
				Sex:       pbprofile.Sex_MALE,
			},
		})
		t.r.NoError(err)

		t.r.EqualValues(&pbprofile.Profile{
			Id:        1,
			FirstName: "John",
			LastName:  "Doe",
			Age:       48,
			Email:     "johndoe@gmail.com",
			Sex:       pbprofile.Sex_MALE,
		}, res.Profile)
	}))

	t.Run("update", profileTest(func(t *tester) {
		ctx := contextWithUserID(t.ctx, int(t.user.ID))

		profile := models.UserProfile{
			FirstName: "David",
			LastName:  "Beckham",
			Age:       21,
			Email:     "david.beckham@gmail.com",
			Sex:       "male",
			UserID:    int64(t.user.ID),
		}
		t.r.NoError(t.db.CreateProfile(int64(t.user.ID), &profile))

		res, err := t.svc.Update(ctx, &pbprofile.UpdateRequest{
			Profile: &pbprofile.Profile{
				FirstName: "John",
				LastName:  "Doe",
				Age:       48,
				Email:     "johndoe@gmail.com",
				Sex:       pbprofile.Sex_MALE,
			},
		})
		t.r.NoError(err)

		t.r.EqualValues(&pbprofile.Profile{
			Id:        1,
			FirstName: "John",
			LastName:  "Doe",
			Age:       48,
			Email:     "johndoe@gmail.com",
			Sex:       pbprofile.Sex_MALE,
		}, res.Profile)
	}))
}

type tester struct {
	r    *require.Assertions
	ctx  context.Context
	db   *maindb.DB
	svc  *profile.ProfileService
	user *models.User
}

func profileTest(f func(t *tester)) func(*testing.T) {
	return dbtesting.Inject(func(t *testing.T, db *maindb.DB) {
		r := require.New(t)
		ctx := context.Background()

		logger := log.NewNopLogger()
		cfg := &config.Config{}

		hash, err := bcrypt.GenerateFromPassword([]byte("qwerty1234"), bcrypt.MinCost)
		r.NoError(err)
		user := models.User{
			Username: "testuser",
			Password: string(hash),
		}
		r.NoError(db.CreateUser(&user))

		svc := profile.NewProfileService(cfg, logger, db)

		f(&tester{
			r:    r,
			ctx:  ctx,
			db:   db,
			svc:  svc,
			user: &user,
		})
	})
}

func contextWithUserID(ctx context.Context, userID int) context.Context {
	id := strconv.Itoa(userID)
	md := metadata.Pairs(
		"user_id", id,
	)

	return metadata.NewIncomingContext(ctx, md)
}
