package models

import "time"

// User is a base model for user entity
type User struct {
	ID        uint   `gorm:"primaryKey"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	CreatedAt time.Time

	Image   UserImage   `gorm:"foreignKey:UserID"`
	Profile UserProfile `gorm:"foreignKey:UserID"`
}
