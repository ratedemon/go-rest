package models

import "time"

// UserProfile is a base model for user entity
type UserProfile struct {
	ID        uint   `gorm:"primaryKey"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Age       int16  `json:"age"`
	Email     string `json:"email"`
	Sex       string `json:"sex"`
	CreatedAt time.Time
	UserID    int64 `gorm:"column:user_id" json:"user_id"`
	// User      User  `gorm:"foreignKey:UserID"`
}
