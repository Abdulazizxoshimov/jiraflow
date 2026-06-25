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
		// WebSocket clients cannot send custom headers during HTTP upgrade,
		// so we also accept the token as a ?token= query parameter.
		var tokenStr string
		if raw := c.GetHeader("Authorization"); raw != "" {
			s, ok := strings.CutPrefix(raw, "Bearer ")
			if !ok || s == "" {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error": "invalid authorization header format",
					"code":  "UNAUTHORIZED",
				})
				return
			}
			tokenStr = s
		} else if q := c.Query("token"); q != "" {
			tokenStr = q
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "missing authorization header",
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
