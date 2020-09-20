package image

import (
	"bytes"
	"context"
	"io"
	"mime/multipart"
	"net/http"
	"testing"

	"github.com/go-kit/kit/log"
	"github.com/ratedemon/go-rest/apitest"
	protoimage "github.com/ratedemon/go-rest/proto/image"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

type mockImageClient struct {
	protoimage.ImageServiceClient

	request, response proto.Message
	err               error
}

func (c *mockImageClient) reset() {
	c.request, c.response, c.err = nil, nil, nil
}

func (c *mockImageClient) Upload(
	ctx context.Context,
	in *protoimage.UploadRequest,
	opts ...grpc.CallOption,
) (*protoimage.UploadResponse, error) {
	c.request = in
	if c.err != nil {
		return nil, c.err
	}
	return c.response.(*protoimage.UploadResponse), nil
}

func (c *mockImageClient) Delete(
	ctx context.Context,
	in *protoimage.DeleteRequest,
	opts ...grpc.CallOption,
) (*protoimage.DeleteResponse, error) {
	c.request = in
	if c.err != nil {
		return nil, c.err
	}
	return c.response.(*protoimage.DeleteResponse), nil
}

func Test_ImageHandler(t *testing.T) {
	var c mockImageClient

	mock := apitest.NewServer(t, (&ImageHandler{context.Background(), log.NewNopLogger(), &c}).RegisterRoutes())

	t.Run("upload", func(t *testing.T) {
		c.reset()
		c.response = &protoimage.UploadResponse{
			Id:        1,
			ImagePath: "files/a031asd.jpg",
		}
		image := []byte("data:image/gif;base64,R0lGODlhEAAOALMAAOazToeHh0tLS/7LZv/0jvb29t/f3//Ub//ge8WSLf/rhf/3kdbW1mxsbP//mf///yH5BAAAAAAALAAAAAAQAA4AAA				Re8L1Ekyky67QZ1hLnjM5UUde0ECwLJoExKcppV0aCcGCmTIHEIUEqjgaORCMxIC6e0CcguWw6aFjsVMkkIr7g77ZKPJjPZqIyd7sJAgVGoEGv2xsBxqNgYPj/gAwXEQA7")

		reqBody := new(bytes.Buffer)
		writer := multipart.NewWriter(reqBody)
		part, err := writer.CreateFormFile("image", "test.jpg")
		mock.R.NoError(err)
		r := bytes.NewReader(image)
		_, err = io.Copy(part, r)
		mock.R.NoError(err)

		mock.R.NoError(writer.Close())

		status, body, err := mock.DoFile(http.MethodPost, "/image", reqBody, writer)
		mock.R.NoError(err)
		mock.R.Equal(http.StatusOK, status)

		mock.R.EqualValues(
			&protoimage.UploadRequest{
				Image:    image,
				Filename: "test.jpg",
			},
			c.request,
		)
		mock.R.JSONEq(
			`{"id": 1, "image_path": "files/a031asd.jpg"}`,
			body,
		)
	})

	t.Run("delete", func(t *testing.T) {
		c.reset()
		c.response = &protoimage.DeleteResponse{}

		status, body, err := mock.Do(http.MethodDelete, "/image/1", "")
		mock.R.NoError(err)
		mock.R.Equal(http.StatusOK, status)

		mock.R.EqualValues(
			&protoimage.DeleteRequest{
				Id: 1,
			},
			c.request,
		)
		mock.R.JSONEq(
			`{}`,
			body,
		)
	})
}
