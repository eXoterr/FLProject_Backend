package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email        string
	Password     string
	Name         string
	Surname      string
	RegisteredAt time.Time
}
