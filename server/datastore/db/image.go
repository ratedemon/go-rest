package db

import (
	"errors"

	"github.com/ratedemon/go-rest/datastore/models"
)

func (db *DB) InsertImage(image *models.UserImage) error {
	result := db.db.Create(image)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (db *DB) DeleteImage(image *models.UserImage) (string, error) {
	result := db.db.First(image)
	if result.Error != nil {
		return "", result.Error
	}

	imageSrc := image.Path
	if imageSrc == "" {
		return "", errors.New("image path is empty")
	}

	result = db.db.Delete(image)
	if result.Error != nil {
		return "", result.Error
	}

	return imageSrc, nil
}
