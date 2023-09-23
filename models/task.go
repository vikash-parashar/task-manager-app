package models

import (
	"time"

	"github.com/google/uuid"
)

const (
	Pending   string = "pending"
	Completed string = "completed"
	Canceled  string = "canceled"
	Active    string = "active"
)
const (
	Low    string = "low"
	Medium string = "medium"
	High   string = "high"
)

type Task struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;" json:"id"`
	Title     string    `json:"title"`
	Status    string    `json:"status"`
	Priority  string    `json:"priority"`
	UserID    uuid.UUID `json:"user_id" gorm:"foreignkey:UserID;references:ID"` // Foreign key to User table
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}
