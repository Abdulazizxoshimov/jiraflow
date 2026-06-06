package middleware

import (
	"time"

	"github.com/gin-gonic/gin"

	"github.com/jira-backend/jiraflow-backend/internal/pkg/logger"
)

// Logger logs only 5xx server errors to the log file.
// All request traffic is handled by gin.Logger() which writes to stdout.
func Logger(log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path

		c.Next()

		if c.Writer.Status() < 500 {
			return
		}

		reqID, _ := c.Get(HeaderRequestID)
		fields := []logger.Field{
			logger.String("method", c.Request.Method),
			logger.String("path", path),
			logger.Int("status", c.Writer.Status()),
			logger.Int64("latency_ms", time.Since(start).Milliseconds()),
			logger.String("ip", c.ClientIP()),
			logger.Any("request_id", reqID),
		}
		if q := c.Request.URL.RawQuery; q != "" {
			fields = append(fields, logger.String("query", q))
		}
		if len(c.Errors) > 0 {
			fields = append(fields, logger.String("errors", c.Errors.String()))
		}

		log.Error(c.Request.Context(), "server error", fields...)
	}
}
