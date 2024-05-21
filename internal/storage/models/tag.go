package models

import "gorm.io/gorm"

type Tag struct {
	gorm.Model
	ID   uint   `gorm:"primarykey"`
	Name string `gorm:"unique"`
}
