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

func RespondSuccess(c *gin.Context, code int, message string, data interface{}) {
	c.JSON(code, APIResponse{
		Success: true,
		Code:    code,
		Message: message,
		Data:    data,
	})
}

func ResponseError(c *gin.Context, code int, message string, details interface{}) {
	c.JSON(code, APIResponse{
		Success: false,
		Code:    code,
		Message: message,
		Error:   details,
	})
}
