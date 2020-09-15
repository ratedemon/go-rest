package auth

import (
	"context"
	"net/http"
	"testing"

	"github.com/go-kit/kit/log"
	"github.com/ratedemon/go-rest/apitest"
	protoauth "github.com/ratedemon/go-rest/proto/auth"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

type mockAuthClient struct {
	protoauth.AuthServiceClient

	request, response proto.Message
	err               error
}

func (c *mockAuthClient) reset() {
	c.request, c.response, c.err = nil, nil, nil
}

func (c *mockAuthClient) Login(
	ctx context.Context,
	in *protoauth.LoginRequest,
	opts ...grpc.CallOption,
) (*protoauth.LoginResponse, error) {
	c.request = in
	if c.err != nil {
		return nil, c.err
	}
	return c.response.(*protoauth.LoginResponse), nil
}

func (c *mockAuthClient) Signup(
	ctx context.Context,
	in *protoauth.SignupRequest,
	opts ...grpc.CallOption,
) (*protoauth.SignupResponse, error) {
	c.request = in
	if c.err != nil {
		return nil, c.err
	}
	return c.response.(*protoauth.SignupResponse), nil
}

func Test_AuthHandler(t *testing.T) {
	var c mockAuthClient

	mock := apitest.NewMockedServerRoutes(t, (&AuthHandler{context.Background(), log.NewNopLogger(), &c}).RegisterRoutes())
	defer mock.Close()

	t.Run("login", func(t *testing.T) {
		c.reset()
		c.response = &protoauth.LoginResponse{
			Id:       1,
			Username: "admin",
			Token:    "a1cs31.222.asd",
		}

		status, body, err := mock.Do(http.MethodPost, "/login", "{\"username\": \"test\", \"password\": \"password123\"}")
		mock.R.NoError(err)
		mock.R.Equal(http.StatusOK, status)

		mock.R.EqualValues(
			&protoauth.LoginRequest{
				Username: "test",
				Password: "password123",
			},
			c.request,
		)
		mock.R.JSONEq(
			`
			{
			"id": 1,
			"username": "admin",
			"token": "a1cs31.222.asd"
			}
		 `,
			body,
		)
	})

	t.Run("signup", func(t *testing.T) {
		c.reset()
		c.response = &protoauth.SignupResponse{
			Message: "Success",
		}

		status, body, err := mock.Do(http.MethodPost, "/signup", "{\"username\": \"test\", \"password\": \"password123\", \"confirm_password\": \"password123\"}")
		mock.R.NoError(err)
		mock.R.Equal(http.StatusOK, status)

		mock.R.EqualValues(
			&protoauth.SignupRequest{
				Username:        "test",
				Password:        "password123",
				ConfirmPassword: "password123",
			},
			c.request,
		)
		mock.R.JSONEq(
			`
			{
			"message": "Success"
			}
		 `,
			body,
		)
	})
}
