package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func main() {
	// Initialize the database
	db, err := gorm.Open("sqlite3", "tasks.db")
	if err != nil {
		panic("Failed to connect to database")
	}
	defer db.Close()

	// Auto-migrate the database
	db.AutoMigrate(&Task{})

	r := mux.NewRouter()

	// Define routes
	r.HandleFunc("/", taskListHandler).Methods("GET")
	r.HandleFunc("/task/create", createTaskHandler).Methods("POST")
	r.HandleFunc("/task/update/{id}", updateTaskHandler).Methods("POST")
	r.HandleFunc("/task/delete/{id}", deleteTaskHandler).Methods("POST")

	// Add authentication middleware and user routes here

	http.Handle("/", r)

	// Start the server
	http.ListenAndServe(":8080", nil)
}
