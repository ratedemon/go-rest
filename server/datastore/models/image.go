package models

type Image struct {
	ID   uint   `gorm:"primaryKey"`
	Path string `json:"image_path" gorm:"image_path"`
	User User   `gorm:"foreignKey:user_id,references:users"`
}
