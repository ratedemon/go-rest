package image_test

import (
	"bytes"
	"context"
	"io/ioutil"
	"strconv"
	"testing"

	"github.com/go-kit/kit/log"
	"github.com/ratedemon/go-rest/config"
	maindb "github.com/ratedemon/go-rest/datastore/db"
	"github.com/ratedemon/go-rest/datastore/models"
	"github.com/ratedemon/go-rest/dbtesting"
	"github.com/ratedemon/go-rest/grpcserver/image"

	gimage "image"

	pbimage "github.com/ratedemon/go-rest/proto/image"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/metadata"
)

const sideMeasure = 160

func TestImageService(t *testing.T) {
	deleteImage := ""

	t.Run("upload", imageTest(func(t *tester) {
		ctx := contextWithUserID(t.ctx, int(t.user.ID))

		b, err := ioutil.ReadFile("test/original.jpg")
		t.r.NoError(err)

		res, err := t.svc.Upload(ctx, &pbimage.UploadRequest{
			Image:    b,
			Filename: "empty.jpg",
		})
		t.r.NoError(err)

		deleteImage = res.ImagePath

		newImage, err := ioutil.ReadFile(res.ImagePath)
		t.r.NoError(err)

		r := bytes.NewReader(newImage)
		img, _, err := gimage.DecodeConfig(r)

		t.r.LessOrEqual(img.Height, sideMeasure)
		t.r.LessOrEqual(img.Width, sideMeasure)
	}))

	t.Run("delete", imageTest(func(t *tester) {
		ctx := contextWithUserID(t.ctx, int(t.user.ID))

		imageObj := models.UserImage{
			Path:   deleteImage,
			UserID: int64(t.user.ID),
		}
		t.r.NoError(t.db.InsertImage(&imageObj))

		_, err := t.svc.Delete(ctx, &pbimage.DeleteRequest{
			Id: int64(imageObj.ID),
		})
		t.r.NoError(err)
	}))
}

type tester struct {
	r    *require.Assertions
	ctx  context.Context
	db   *maindb.DB
	svc  *image.ImageService
	user *models.User
}

func imageTest(f func(t *tester)) func(*testing.T) {
	return dbtesting.Inject(func(t *testing.T, db *maindb.DB) {
		r := require.New(t)
		ctx := context.Background()

		logger := log.NewNopLogger()
		cfg := &config.Config{
			Image: config.Image{
				ImagePrefixPath: "test",
				SideMeasure:     sideMeasure,
			},
		}

		hash, err := bcrypt.GenerateFromPassword([]byte("qwerty1234"), bcrypt.MinCost)
		r.NoError(err)
		user := models.User{
			Username: "testuser",
			Password: string(hash),
		}
		r.NoError(db.CreateUser(&user))

		svc := image.NewImageService(cfg, logger, db)

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
