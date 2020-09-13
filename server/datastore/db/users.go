package db

import (
	"github.com/ratedemon/go-rest/datastore/models"
)

func (db *DB) CreateUser(user *models.User) error {
	result := db.db.Create(user)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (db *DB) FindUserByUsername(username string, user *models.User) error {
	result := db.db.First(&user, "username = ?", username)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
