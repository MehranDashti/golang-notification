package router

import (
	rest_handler "notification/internal/handler/rest"
	"notification/internal/middleware"

	"github.com/gin-gonic/gin"
)

func Setup(
	healthHandler *rest_handler.HealthHandler,
	notificationHandler *rest_handler.NotificationHandler,
) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.ErrorHandler())

	r.GET("/health", healthHandler.Check)
	v1 := r.Group("/api/v1")
	{
		v1.POST("/notifications", notificationHandler.Create)
		v1.GET("/notifications/:id", notificationHandler.GetByID)
		v1.GET("/users/:user_id/notifications", notificationHandler.ListByUser)
	}

	return r
}
