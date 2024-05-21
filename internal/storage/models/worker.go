package models

import "gorm.io/gorm"

type Worker struct {
	gorm.Model
	UserID           uint
	User             User
	SpecializationID int
	Specialization   Specialization
	Rating           int  `gorm:"default:5"`
	OrdersCompleted  int  `gorm:"default:0"`
	IsVerified       bool `gorm:"default:false"`
}
