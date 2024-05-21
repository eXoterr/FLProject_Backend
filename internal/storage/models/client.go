package models

import "gorm.io/gorm"

type Client struct {
	gorm.Model
	UserID          uint
	User            User
	Rating          int `gorm:"default:5"`
	OrdersCompleted int
}

type ClientStore interface {
}
