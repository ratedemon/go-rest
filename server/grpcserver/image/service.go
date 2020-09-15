package image

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/jpeg"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/nfnt/resize"

	"github.com/go-kit/kit/log"
	"google.golang.org/grpc"

	"github.com/ratedemon/go-rest/config"
	"github.com/ratedemon/go-rest/datastore/db"
	"github.com/ratedemon/go-rest/datastore/models"
	"github.com/ratedemon/go-rest/grpcserver/helper"
	pbimage "github.com/ratedemon/go-rest/proto/image"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

type ImageService struct {
	cfg *config.Config
	log log.Logger
	db  *db.DB
}

func NewImageService(cfg *config.Config, log log.Logger, db *db.DB) *ImageService {
	return &ImageService{cfg, log, db}
}

func (is *ImageService) Upload(ctx context.Context, req *pbimage.UploadRequest) (*pbimage.UploadResponse, error) {
	userID, err := helper.GetUserID(ctx)
	if err != nil {
		return nil, err
	}

	r := bytes.NewReader(req.Image)
	img, _, err := image.Decode(r)
	if err != nil {
		return nil, err
	}

	ext := filepath.Ext(req.Filename)
	newFilePath := path.Join(is.cfg.Image.ImagePrefixPath, fmt.Sprintf("%s%s", randStringRunes(16), ext))
	out, err := os.Create(newFilePath)
	if err != nil {
		return nil, err
	}
	defer out.Close()

	m := resize.Thumbnail(uint(is.cfg.Image.SideMeasure), uint(is.cfg.Image.SideMeasure), img, resize.Lanczos3)
	if err = jpeg.Encode(out, m, nil); err != nil {
		return nil, err
	}

	imageModel := models.UserImage{
		Path:   newFilePath,
		UserID: userID,
	}

	if err := is.db.InsertImage(&imageModel); err != nil {
		return nil, err
	}

	return &pbimage.UploadResponse{
		Id:        int64(imageModel.ID),
		ImagePath: imageModel.Path,
	}, nil
}

func (is *ImageService) Delete(ctx context.Context, req *pbimage.DeleteRequest) (*pbimage.DeleteResponse, error) {
	userID, err := helper.GetUserID(ctx)
	if err != nil {
		return nil, err
	}

	imageModel := models.UserImage{
		ID:     uint(req.Id),
		UserID: userID,
	}
	imageSrc, err := is.db.DeleteImage(&imageModel)
	if err != nil {
		return nil, err
	}
	if err := os.Remove(imageSrc); err != nil {
		return nil, err
	}

	return &pbimage.DeleteResponse{}, nil
}

func (is *ImageService) RegisterService(s *grpc.Server) {
	pbimage.RegisterImageServiceServer(s, is)
}

func randStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
