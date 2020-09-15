package user

import (
	"context"
	"net/http"
	"testing"

	"github.com/go-kit/kit/log"
	"github.com/ratedemon/go-rest/apitest"
	protouser "github.com/ratedemon/go-rest/proto/user"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

type mockUserClient struct {
	protouser.UserServiceClient

	request, response proto.Message
	err               error
}

func (c *mockUserClient) reset() {
	c.request, c.response, c.err = nil, nil, nil
}

func (c *mockUserClient) Get(
	ctx context.Context,
	in *protouser.GetRequest,
	opts ...grpc.CallOption,
) (*protouser.GetResponse, error) {
	c.request = in
	if c.err != nil {
		return nil, c.err
	}
	return c.response.(*protouser.GetResponse), nil
}

func (c *mockUserClient) List(
	ctx context.Context,
	in *protouser.ListRequest,
	opts ...grpc.CallOption,
) (*protouser.ListResponse, error) {
	c.request = in
	if c.err != nil {
		return nil, c.err
	}
	return c.response.(*protouser.ListResponse), nil
}

func Test_UserHandler(t *testing.T) {
	var c mockUserClient

	mock := apitest.NewMockedServerRoutes(t, (&UserHandler{context.Background(), log.NewNopLogger(), &c}).RegisterRoutes())
	defer mock.Close()

	t.Run("list", func(t *testing.T) {
		c.reset()
		c.response = &protouser.ListResponse{
			Users: []*protouser.User{
				&protouser.User{
					Id:        1,
					Username:  "john",
					ImagePath: "files/admin.jpg",
				},
				&protouser.User{
					Id:        2,
					Username:  "sam",
					FirstName: "adam",
					LastName:  "friend",
				},
			},
		}

		status, body, err := mock.Do(http.MethodGet, "/users", "")
		mock.R.NoError(err)
		mock.R.Equal(http.StatusOK, status)

		mock.R.EqualValues(
			&protouser.ListRequest{},
			c.request,
		)
		mock.R.JSONEq(
			`{"users":[{"id":1,"username":"john","image_path":"files/admin.jpg"},{"id":2,"username":"sam","first_name":"adam","last_name":"friend"}]}`,
			body,
		)
	})
}
