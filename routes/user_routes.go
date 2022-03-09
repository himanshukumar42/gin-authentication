package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/himanshuk42/gin-authentication/controllers"
	"github.com/himanshuk42/gin-authentication/middlewares"
)

func UserRoute(router *gin.Engine) {
	public := router.Group("/api")

	public.POST("/register", controllers.Register())
	public.POST("/login", controllers.Login())

	protected := router.Group("/api/admin")
	protected.Use(middlewares.JwtAuthMiddleware())
	{
		protected.GET("/users", controllers.GetAllUsers())
		protected.GET("/user/:userId", controllers.GetAUser())
		protected.POST("/user", controllers.CreateUser())
		protected.PUT("/user/:userId", controllers.EditUser())
		protected.DELETE("/user/:userId", controllers.DeleteUser())
	}
}
