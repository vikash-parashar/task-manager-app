package models

import (
	"time"

	"github.com/google/uuid"
)

const (
	Pending   string = "Pending"
	Completed string = "Completed"
	Canceled  string = "Canceled"
	Active    string = "Active"
)

type Task struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;" json:"id"`
	Title     string    `json:"title"`
	Status    string    `json:"status"`
	UserID    uuid.UUID `json:"user_id"`
	User      User      `gorm:"foreignkey:UserID" json:"user"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}
