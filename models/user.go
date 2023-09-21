package models

import "github.com/jinzhu/gorm"

type User struct {
	gorm.Model
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
	Password  []byte `json:"password"`
	Tasks     []Task `json:"tasks"` // Define the one-to-many relationship
}
