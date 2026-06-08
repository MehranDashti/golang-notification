package dto

import (
	"time"

	"github.com/gin-gonic/gin"
)

type APIResponse struct {
	Success bool        `json:"success"`
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

type NotificationResponse struct {
	ID        string            `json:"id"`
	UserID    string            `json:"user_id"`
	Channel   string            `json:"channel"`
	Status    string            `json:"status"`
	Title     string            `json:"title"`
	Body      string            `json:"body"`
	Metadata  map[string]string `json:"metadata,omitempty"`
	SentAt    *time.Time        `json:"sent_at,omitempty"`
	CreatedAt time.Time         `json:"created_at"`
}

func RespondSuccess(c *gin.Context, code int, message string, data interface{}) {
	respondData(c, code, message, data)
}

func respondData[T any](c *gin.Context, code int, message string, data T) {
	c.JSON(code, APIResponse{
		Success: true,
		Code:    code,
		Message: message,
		Data:    data,
	})
}
