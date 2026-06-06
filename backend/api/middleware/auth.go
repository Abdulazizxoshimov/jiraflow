package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"github.com/jira-backend/jiraflow-backend/internal/pkg/token"
)

// Auth validates the JWT and injects user identity into the Gin context.
// Unlike RBAC, it does not consult Casbin — use it for routes that only
// need authentication, not role-based access control.
func Auth(maker token.Maker) gin.HandlerFunc {
	return func(c *gin.Context) {
		raw := c.GetHeader("Authorization")
		if raw == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "missing authorization header",
				"code":  "UNAUTHORIZED",
			})
			return
		}

		tokenStr, ok := strings.CutPrefix(raw, "Bearer ")
		if !ok || tokenStr == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid authorization header format",
				"code":  "UNAUTHORIZED",
			})
			return
		}

		claims, err := maker.ValidateAccess(c.Request.Context(), tokenStr)
		if err != nil {
			if errors.Is(err, jwt.ErrTokenExpired) {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error": "token expired",
					"code":  "TOKEN_EXPIRED",
				})
			} else {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error": "invalid token",
					"code":  "TOKEN_INVALID",
				})
			}
			return
		}

		c.Set(CtxUserID, claims.Sub)
		c.Set(CtxSessionID, claims.SessionID)
		c.Set(CtxRole, claims.Role)
		c.Next()
	}
}
