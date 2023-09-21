package models

import "github.com/jinzhu/gorm"

const (
	Pending   string = "Pending"
	Completed string = "Completed"
	Canceled  string = "Canceled"
	Active    string = "Active"
)

type Task struct {
	gorm.Model
	Title  string `json:"title"`
	Status string `json:"status"`
	UserID uint   `json:"user_id"`
	User   User   `json:"user"` // Define the foreign key relationship
}
