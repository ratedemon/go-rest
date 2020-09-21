package user_test

import (
	"context"
	"strconv"
	"testing"

	"github.com/go-kit/kit/log"
	"github.com/ratedemon/go-rest/config"
	maindb "github.com/ratedemon/go-rest/datastore/db"
	"github.com/ratedemon/go-rest/datastore/models"
	"github.com/ratedemon/go-rest/dbtesting"
	svcuser "github.com/ratedemon/go-rest/grpcserver/user"
	pbuser "github.com/ratedemon/go-rest/proto/user"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/metadata"
)

type testUser struct {
	user    *models.User
	profile *models.UserProfile
	image   *models.UserImage
}

func TestUserService(t *testing.T) {
	users := []testUser{
		testUser{
			user: &models.User{
				Username: "admin",
				Password: "qwerty",
			},
			profile: &models.UserProfile{
				FirstName: "Tom",
				LastName:  "Cruz",
				Age:       50,
				Sex:       "male",
			},
		},
		testUser{
			user: &models.User{
				Username: "newadmin",
				Password: "newqwerty",
			},
			profile: &models.UserProfile{
				FirstName: "Johny",
				LastName:  "Depp",
				Age:       53,
				Sex:       "male",
			},
			image: &models.UserImage{
				Path: "test/photo.jpg",
			},
		},
	}

	fillDB := func(r *require.Assertions, db *maindb.DB) {
		for _, u := range users {
			hash, err := bcrypt.GenerateFromPassword([]byte(u.user.Password), bcrypt.MinCost)
			r.NoError(err)
			user := models.User{
				Username: u.user.Username,
				Password: string(hash),
			}
			r.NoError(db.CreateUser(&user))
			u.user.ID = user.ID

			if u.profile != nil {
				u.profile.UserID = int64(user.ID)
				r.NoError(db.CreateProfile(int64(user.ID), u.profile))
			}

			if u.image != nil {
				u.image.UserID = int64(user.ID)
				r.NoError(db.InsertImage(u.image))
			}
		}
	}

	t.Run("list", userTest(func(t *tester) {
		fillDB(t.r, t.db)

		ctx := contextWithUserID(t.ctx, int(t.user.ID))

		res, err := t.svc.List(ctx, &pbuser.ListRequest{})
		t.r.NoError(err)

		t.r.Len(res.Users, 2)
	}))

	t.Run("get", userTest(func(t *tester) {
		fillDB(t.r, t.db)
		ctx := contextWithUserID(t.ctx, int(t.user.ID))

		res, err := t.svc.Get(ctx, &pbuser.GetRequest{
			Id: int64(users[1].user.ID),
		})
		t.r.NoError(err)

		t.r.Equal("Johny", res.User.FirstName)
		t.r.Equal("Depp", res.User.LastName)
	}))
}

type tester struct {
	r    *require.Assertions
	ctx  context.Context
	db   *maindb.DB
	svc  *svcuser.UserService
	user *models.User
}

func userTest(f func(t *tester)) func(*testing.T) {
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

		svc := svcuser.NewUserService(cfg, logger, db)

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
