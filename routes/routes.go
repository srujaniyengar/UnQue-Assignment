package routes

import (
	"UnQue/control"
	"UnQue/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes() *gin.Engine {
	router := gin.Default()

	// Public route
	router.POST("/login", control.Login)

	// Protected routes group
	protected := router.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.POST("/availability", control.SetAvailability)
		protected.GET("/availability", control.GetAvailability)
		protected.POST("/appointments", control.BookAppointment)
		protected.DELETE("/appointments/:id", control.CancelAppointment)
		protected.GET("/appointments", control.GetAppointments)
	}

	return router
}
