package routes

import (
	"task-manager-app/controllers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(r *gin.Engine, dbInstance *gorm.DB) {
	// Pass the database instance to your route handlers if needed
	r.GET("/", controllers.HomePage)
	r.GET("/register", controllers.RegisterPage)
	r.GET("/login", controllers.LoginPage)
	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)
	r.GET("/todo", controllers.TodoPage)
	r.GET("/profile", controllers.GetCurrentUser)
	r.GET("/tasks", controllers.GetAllTasks)
	r.POST("/task/create", controllers.CreateTask)
	r.PUT("/task/update/:id", controllers.UpdateTask)
	r.DELETE("/task/delete/:id", controllers.DeleteTask)
}
