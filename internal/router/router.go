package router

import (
	rest_handler "notification/internal/handler/rest"
	"notification/internal/middleware"

	"github.com/gin-gonic/gin"
)

func Setup(
	healthHandler *rest_handler.HealthHandler,
	notificationHandler *rest_handler.NotificationHandler,
	providerHandler *rest_handler.ProviderHandler,
) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.ErrorHandler())

	r.GET("/health", healthHandler.Check)
	v1 := r.Group("/api/v1")
	{
		// Notifications
		v1.POST("/notifications", notificationHandler.Create)
		v1.GET("/notifications/:id", notificationHandler.GetByID)
		v1.GET("/users/:user_id/notifications", notificationHandler.ListByUser)

		// Provider Config
		provider := v1.Group("/providers")
		{
			provider.POST("", providerHandler.Create)
			provider.GET("", providerHandler.List)
			provider.GET("/:id", providerHandler.GetByID)
			provider.PUT("/:id", providerHandler.Update)
			provider.PATCH("/:id/activate", providerHandler.Activate)
			provider.PATCH("/:id/deactivate", providerHandler.Deactivate)
		}
	}

	return r
}
