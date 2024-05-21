package payloads

import "time"

type Order struct {
	CategoryID  int       `json:"category_id" validate:"required"`
	Category    string    `json:"category"`
	Title       string    `json:"title" validate:"required,min=20,max=120"`
	Description string    `json:"description" validate:"required"`
	BudgetMin   uint      `json:"budget_min"`
	BudgetMax   uint      `json:"budget_max"`
	Deadline    time.Time `json:"deadline" validate:"required"`
	Tags        []uint    `json:"tags"`
}

type Tags []uint
