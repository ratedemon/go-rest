package models

import "time"

// Profile is a base model for user entity
type Profile struct {
	ID        uint   `gorm:"primaryKey"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Age       int16  `json:"age"`
	Email     string `json:"email"`
	CreatedAt time.Time
	User      User `gorm:"foreignKey:user_id,references:users"`
}
