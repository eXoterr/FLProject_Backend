package models

import (
	"time"

	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	ClientID    uint
	Client      Client
	WorkerID    uint `gorm:"default:null"`
	Worker      Worker
	State       uint `gorm:"default:0"`
	CategoryID  int
	Category    Category
	Title       string
	Description string
	BudgetMin   uint
	BudgetMax   uint
	Deadline    time.Time
	Tags        []Tag `gorm:"many2many:order_tags"`
}
