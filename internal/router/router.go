package router

import (
	rest_handler "notification/internal/handler/rest"

	"github.com/gin-gonic/gin"
)

func Setup(
	healthHandler *rest_handler.HealthHandler,
) *gin.Engine {
	r := gin.New()

	r.GET("/health", healthHandler.Check)

	return r
}
