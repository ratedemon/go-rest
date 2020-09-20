package profile

import (
	"context"
	"net/http"
	"testing"

	"github.com/go-kit/kit/log"
	"github.com/ratedemon/go-rest/apitest"
	protoprofile "github.com/ratedemon/go-rest/proto/profile"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

type mockProfileClient struct {
	protoprofile.ProfileServiceClient

	request, response proto.Message
	err               error
}

func (c *mockProfileClient) reset() {
	c.request, c.response, c.err = nil, nil, nil
}

func (c *mockProfileClient) Create(
	ctx context.Context,
	in *protoprofile.CreateRequest,
	opts ...grpc.CallOption,
) (*protoprofile.CreateResponse, error) {
	c.request = in
	if c.err != nil {
		return nil, c.err
	}
	return c.response.(*protoprofile.CreateResponse), nil
}

func (c *mockProfileClient) Update(
	ctx context.Context,
	in *protoprofile.UpdateRequest,
	opts ...grpc.CallOption,
) (*protoprofile.UpdateResponse, error) {
	c.request = in
	if c.err != nil {
		return nil, c.err
	}
	return c.response.(*protoprofile.UpdateResponse), nil
}

func Test_ProfileHandler(t *testing.T) {
	var c mockProfileClient

	mock := apitest.NewServer(t, (&ProfileHandler{context.Background(), log.NewNopLogger(), &c}).RegisterRoutes())

	t.Run("create", func(t *testing.T) {
		c.reset()
		profile := &protoprofile.Profile{
			FirstName: "John",
			LastName:  "Doe",
			Age:       45,
			Email:     "john_doe@gmail.com",
			Sex:       protoprofile.Sex_MALE,
			Id:        1,
		}
		c.response = &protoprofile.CreateResponse{
			Profile:   profile,
			CreatedAt: "now",
		}

		status, body, err := mock.Do(http.MethodPost, "/profile", "{\"first_name\": \"John\", \"last_name\": \"Doe\", \"age\": 45, \"email\": \"john_doe@gmail.com\", \"sex\": \"male\"}")
		mock.R.NoError(err)
		mock.R.Equal(http.StatusOK, status)

		profile.Id = 0
		mock.R.EqualValues(
			&protoprofile.CreateRequest{
				Profile: profile,
			},
			c.request,
		)
		mock.R.JSONEq(
			`{"id":1,"first_name":"John","last_name":"Doe","age":45,"email":"john_doe@gmail.com","sex":1}`,
			body,
		)
	})

	t.Run("update", func(t *testing.T) {
		c.reset()
		profile := &protoprofile.Profile{
			FirstName: "John",
			LastName:  "Doe",
			Age:       45,
			Email:     "john_doe@gmail.com",
			Sex:       protoprofile.Sex_MALE,
			Id:        1,
		}
		c.response = &protoprofile.UpdateResponse{
			Profile: profile,
		}

		status, body, err := mock.Do(http.MethodPut, "/profile/1", "{\"first_name\": \"John\", \"last_name\": \"Doe\", \"age\": 45, \"email\": \"john_doe@gmail.com\", \"sex\": \"male\"}")
		mock.R.NoError(err)
		mock.R.Equal(http.StatusOK, status)

		mock.R.EqualValues(
			&protoprofile.UpdateRequest{
				Profile: profile,
			},
			c.request,
		)
		mock.R.JSONEq(
			`{
				"profile": {"id":1,"first_name":"John","last_name":"Doe","age":45,"email":"john_doe@gmail.com","sex":1}	
			}`,
			body,
		)
	})
}
