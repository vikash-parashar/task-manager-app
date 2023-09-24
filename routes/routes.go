package routes

import (
	"task-manager-app/config"
	"task-manager-app/controllers"
	"task-manager-app/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, appConfig *config.AppConfig) {
	// Pass the database instance and config to your route handlers if needed
	r.GET("/", func(c *gin.Context) {
		controllers.HomePage(c)
	})
	r.GET("/register", func(c *gin.Context) {
		controllers.RegisterPage(c)
	})
	r.GET("/login", func(c *gin.Context) {
		controllers.LoginPage(c)
	})
	r.POST("/register", func(c *gin.Context) {
		controllers.Register(c, appConfig)
	})
	r.POST("/login", func(c *gin.Context) {
		controllers.Login(c, appConfig)
	})

	// Create a router group for protected routes that require authentication
	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware()) // Apply AuthMiddleware to this group
	{
		protected.GET("/todo", func(c *gin.Context) {
			controllers.TodoPage(c)
		})
		protected.GET("/profile", func(c *gin.Context) {
			controllers.GetCurrentUser(c, appConfig)
		})
		protected.GET("/tasks", func(c *gin.Context) {
			controllers.GetAllTasks(c, appConfig)
		})
		protected.POST("/task/create", func(c *gin.Context) {
			controllers.CreateTask(c, appConfig)
		})
		protected.PUT("/task/update/:id", func(c *gin.Context) {
			controllers.UpdateTask(c, appConfig)
		})
		protected.DELETE("/task/delete/:id", func(c *gin.Context) {
			controllers.DeleteTask(c, appConfig)
		})
		protected.GET("/logout", func(c *gin.Context) {
			controllers.Logout(c)
		})
	}
}
