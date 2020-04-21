package models

import "github.com/jinzhu/gorm"

// Gallery is the image container resource that visitors
// view
type Gallery struct {
	gorm.Model
	UserID uint   `gorm:"not_null;index"`
	Title  string `gorm:"not_null"`
}
