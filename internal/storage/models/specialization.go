package models

import "gorm.io/gorm"

type Specialization struct {
	gorm.Model
	Name string `gorm:"unique"`
}
