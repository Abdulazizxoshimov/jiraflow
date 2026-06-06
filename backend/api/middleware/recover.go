package middleware

import (
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"

	"github.com/jira-backend/jiraflow-backend/internal/pkg/logger"
)

// Recover catches panics, logs the stack trace, and returns 500.
func Recover(log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				ctx := c.Request.Context()
				log.Error(ctx, "panic recovered",
					logger.Any("error", r),
					logger.String("stack", string(debug.Stack())),
					logger.String("path", c.Request.URL.Path),
					logger.String("method", c.Request.Method),
				)
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"error": "internal server error",
					"code":  "INTERNAL_ERROR",
				})
			}
		}()
		c.Next()
	}
}
