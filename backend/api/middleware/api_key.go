package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	api_key "github.com/jira-backend/jiraflow-backend/internal/usecase/api_key"
)

// APIKeyAuth is an alternative auth middleware that accepts X-API-Key header.
// It sets the same context keys as JWT auth so downstream handlers work identically.
func APIKeyAuth(uc api_key.UseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		raw := c.GetHeader("X-API-Key")
		if raw == "" {
			// Try Bearer token with jfk_ prefix
			auth := c.GetHeader("Authorization")
			if after, ok := strings.CutPrefix(auth, "Bearer jfk_"); ok {
				raw = "jfk_" + after
			}
		}
		if raw == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    "UNAUTHORIZED",
				"message": "missing X-API-Key header",
			})
			return
		}

		key, err := uc.ValidateKey(c.Request.Context(), raw)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    "UNAUTHORIZED",
				"message": "invalid or expired api key",
			})
			return
		}

		c.Set(CtxUserID, key.UserID)
		c.Set(CtxRole, "member") // API keys act as member role
		c.Set(CtxSessionID, key.ID)
		c.Next()
	}
}
