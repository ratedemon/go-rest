package db

import "gorm.io/gorm"

// DB represents db layer
type DB struct {
	db *gorm.DB
}

// NewDB creates new DB entity
func NewDB(db *gorm.DB) *DB {
	return &DB{db}
}
