package rest

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type HealthHandler struct {
	db *mongo.Client
}

func NewHealthHandler(db *mongo.Client) *HealthHandler {
	return &HealthHandler{db: db}
}

func (h *HealthHandler) Check(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancel()

	if err := h.db.Ping(ctx, nil); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status": "error",
			"db":     "ping failed: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"db":      "ok",
		"version": "1.0.0",
	})
}
