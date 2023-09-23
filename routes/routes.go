package routes

import (
	"task-manager-app/controllers"
	"task-manager-app/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	// Pass the database instance to your route handlers if needed
	r.GET("/", controllers.HomePage)
	r.GET("/register", controllers.RegisterPage)
	r.GET("/login", controllers.LoginPage)

	// Create a router group for protected routes that require authentication
	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware()) // Apply AuthMiddleware to this group
	{
		protected.POST("/register", controllers.Register)
		protected.POST("/login", controllers.Login)
		protected.GET("/todo", controllers.TodoPage)
		protected.GET("/profile", controllers.GetCurrentUser)
		protected.GET("/tasks", controllers.GetAllTasks)
		protected.POST("/task/create", controllers.CreateTask)
		protected.PUT("/task/update/:id", controllers.UpdateTask)
		protected.DELETE("/task/delete/:id", controllers.DeleteTask)
	}
}
