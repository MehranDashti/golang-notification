package router

import (
	"notification/internal/handler"

	"github.com/gin-gonic/gin"
)

func Setup(
	healthHandler *handler.HealthHandler,
) *gin.Engine {
	r := gin.New()

	r.GET("/health", healthHandler.Check)

	return r
}
