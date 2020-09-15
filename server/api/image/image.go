package image

import (
	"context"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/go-kit/kit/log"
	"github.com/gorilla/mux"
	"google.golang.org/grpc"

	"github.com/ratedemon/go-rest/api/helper"
	protoimage "github.com/ratedemon/go-rest/proto/image"
)

// ImageHandler is gateway for image endpoints
type ImageHandler struct {
	ctx        context.Context
	log        log.Logger
	grpcClient protoimage.ImageServiceClient
}

func NewImageHandler(ctx context.Context, log log.Logger, grpcConn *grpc.ClientConn) *ImageHandler {
	client := protoimage.NewImageServiceClient(grpcConn)
	return &ImageHandler{ctx, log, client}
}

func (ih *ImageHandler) upload(ctx context.Context, req *http.Request) (interface{}, error) {
	file, handler, err := req.FormFile("image")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	res, err := ih.grpcClient.Upload(ctx, &protoimage.UploadRequest{
		Image:    fileBytes,
		Filename: handler.Filename,
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (ih *ImageHandler) delete(ctx context.Context, req *http.Request) (interface{}, error) {
	vars := mux.Vars(req)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return nil, err
	}

	res, err := ih.grpcClient.Delete(ctx, &protoimage.DeleteRequest{
		Id: id,
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (ih *ImageHandler) RegisterRoutes() []helper.Route {
	return []helper.Route{
		{"/image", "POST", ih.upload},
		{"/image/{id:[0-9]+}", "DELETE", ih.delete},
	}
}
