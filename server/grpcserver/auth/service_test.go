package auth_test

import (
	"context"
	"testing"

	"github.com/go-kit/kit/log"
	"github.com/ratedemon/go-rest/config"
	maindb "github.com/ratedemon/go-rest/datastore/db"
	"github.com/ratedemon/go-rest/datastore/models"
	"github.com/ratedemon/go-rest/dbtesting"
	"github.com/ratedemon/go-rest/grpcserver/auth"
	pbauth "github.com/ratedemon/go-rest/proto/auth"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestAuthService(t *testing.T) {
	t.Run("login", authTest(func(t *tester) {
		hash, err := bcrypt.GenerateFromPassword([]byte("qwerty1234"), bcrypt.MinCost)
		t.r.NoError(err)

		user := models.User{
			Username: "admin",
			Password: string(hash),
		}
		t.r.NoError(t.db.CreateUser(&user))

		res, err := t.svc.Login(t.ctx, &pbauth.LoginRequest{
			Username: "admin",
			Password: "qwerty1234",
		})
		t.r.NoError(err)

		t.r.Equal("admin", res.Username)
		t.r.NotEmpty(res.Token)
	}))

	t.Run("login: not found", authTest(func(t *tester) {
		_, err := t.svc.Login(t.ctx, &pbauth.LoginRequest{
			Username: "admin",
			Password: "qwerty1234",
		})
		t.r.Error(err)
		t.r.EqualError(err, (status.Errorf(codes.Unknown, "Failed to find the user: record not found")).Error())
	}))

	t.Run("signup", authTest(func(t *tester) {
		res, err := t.svc.Signup(t.ctx, &pbauth.SignupRequest{
			Username:        "admin",
			Password:        "qwerty1234",
			ConfirmPassword: "qwerty1234",
		})
		t.r.NoError(err)

		t.r.EqualValues(&pbauth.SignupResponse{
			Message: "Successfully created",
		}, res)
	}))

	t.Run("signup: not equal passwords", authTest(func(t *tester) {
		_, err := t.svc.Signup(t.ctx, &pbauth.SignupRequest{
			Username:        "admin",
			Password:        "qwerty1234",
			ConfirmPassword: "1234qwerty",
		})
		t.r.Error(err)
		t.r.EqualError(err, (status.Errorf(codes.InvalidArgument, "`confirm_password` and`password` do not match")).Error())
	}))
}

type tester struct {
	r   *require.Assertions
	ctx context.Context
	db  *maindb.DB
	svc *auth.AuthService
}

func authTest(f func(t *tester)) func(*testing.T) {
	return dbtesting.Inject(func(t *testing.T, db *maindb.DB) {
		r := require.New(t)
		ctx := context.Background()

		logger := log.NewNopLogger()
		cfg := &config.Config{}

		svc := auth.NewAuthService(cfg, logger, db)

		f(&tester{
			r:   r,
			ctx: ctx,
			db:  db,
			svc: svc,
		})
	})
}
