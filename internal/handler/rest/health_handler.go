package rest

import (
	"context"
	"net/http"
	"notification/internal/dto"
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
		dto.RespondError(c, http.StatusServiceUnavailable, "db ping failed", err.Error())
		return
	}

	dto.RespondSuccess(c, http.StatusOK, "ok", gin.H{"db": "ok", "version": "1.0.0"})
}
