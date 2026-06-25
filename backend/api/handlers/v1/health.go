package v1

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func HealthCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"time":   time.Now().UTC().Format(time.RFC3339),
		})
	}
}

// ReadyCheck returns 200 when DB and Redis are reachable, 503 otherwise.
// Used as a Kubernetes readiness probe (/ready).
func ReadyCheck(check func() error) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := check(); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status": "unavailable",
				"error":  err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status": "ready",
			"time":   time.Now().UTC().Format(time.RFC3339),
		})
	}
}
