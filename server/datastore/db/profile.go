package db

import (
	"errors"

	"github.com/ratedemon/go-rest/datastore/models"
)

func (db *DB) CreateProfile(userID int64, profile *models.UserProfile) error {
	var existProfile models.UserProfile
	db.db.First(&existProfile, "user_id = ?", userID)
	if existProfile.ID != 0 {
		return errors.New("Profile is already exist for current user")
	}

	result := db.db.Create(profile)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (db *DB) UpdateProfile(userID int64, profile *models.UserProfile) error {
	var existProfile models.UserProfile
	result := db.db.First(&existProfile, "user_id = ?", userID)
	if result.Error != nil {
		return result.Error
	}
	
	result = db.db.Model(existProfile).Updates(profile)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
