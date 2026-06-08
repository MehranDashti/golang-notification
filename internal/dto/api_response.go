package dto

import (
	"github.com/gin-gonic/gin"
)

type APIResponse struct {
	Success bool   `json:"success"`
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
	Error   any    `json:"error,omitempty"`
}

func RespondSuccess(c *gin.Context, code int, message string, data any) {
	c.JSON(code, APIResponse{
		Success: true,
		Code:    code,
		Message: message,
		Data:    data,
	})
}

func RespondError(c *gin.Context, code int, message string, details any) {
	c.JSON(code, APIResponse{
		Success: false,
		Code:    code,
		Message: message,
		Error:   details,
	})
}
