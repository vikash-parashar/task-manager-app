package main

import (
	"github.com/jinzhu/gorm"
)

type Task struct {
	gorm.Model
	Title       string
	Description string
	Deadline    string
	Priority    int
	UserID      uint // For associating tasks with users
}

type User struct {
	gorm.Model
	FirstName string
	LastName  string
	Phone     string
	Email     string
	Password  string
}
