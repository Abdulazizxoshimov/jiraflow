package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// CORS sets permissive CORS headers. Adjust AllowOrigins for production.
func CORS(allowOrigins ...string) gin.HandlerFunc {
	origins := map[string]struct{}{}
	for _, o := range allowOrigins {
		origins[o] = struct{}{}
	}

	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")

		allowed := origin
		if len(origins) > 0 {
			if _, ok := origins[origin]; !ok {
				allowed = ""
			}
		}

		if allowed != "" {
			c.Header("Access-Control-Allow-Origin", allowed)
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Request-ID")
			c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
			c.Header("Access-Control-Max-Age", "86400")
		}

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
