package middleware

import (
	"errors"
	"log/slog"
	"net/http"
	"notification/internal/apperror"
	"notification/internal/dto"

	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		err := c.Errors.Last().Err

		var appErr *apperror.AppError
		if errors.As(err, &appErr) {
			if appErr.Code >= http.StatusInternalServerError && appErr.Err != nil {
				slog.Error("internal error",
					// "trace_id", trace.FromContext(c.Request.Context()), //TODO Add trace id
					"path", c.Request.URL.Path,
					"error", appErr.Err,
				)
			}
			dto.RespondError(c, appErr.Code, appErr.Message, appErr.Details)
			return
		}
		slog.Error("unhandled error — not an AppError",
			// "trace_id", trace.FromContext(c.Request.Context()),
			"path", c.Request.URL.Path,
			"error", err,
		)
		dto.RespondError(
			c,
			http.StatusInternalServerError,
			"internal server error",
			nil,
		)
	}
}
