package models

type Image struct {
	ID     uint   `gorm:"primaryKey"`
	Path   string `json:"image_path" gorm:"image_path"`
	UserID int64  `gorm:"column:user_id" json:"user_id"`
	User   User   `gorm:"foreignKey:UserID"`
}
