package models

// UserImage represents `user_images` model in db
type UserImage struct {
	ID     uint   `gorm:"primaryKey"`
	Path   string `gorm:"column:image_path" json:"image_path"`
	UserID int64  `gorm:"column:user_id" json:"user_id"`
	User   User   `gorm:"foreignKey:UserID"`
}
