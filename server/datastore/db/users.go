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

func (db *DB) FindUserById(userID int64) (*models.User, error) {
	user := &models.User{}

	result := db.db.Joins("Image").Joins("Profile").First(&user, userID)
	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

func (db *DB) FindAllUsers(userID int64) ([]models.User, error) {
	users := []models.User{}

	result := db.db.Joins("Image").Joins("Profile").Where("users.id != ?", userID).Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}

	return users, nil
}
