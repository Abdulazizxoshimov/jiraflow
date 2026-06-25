package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"github.com/jira-backend/jiraflow-backend/internal/pkg/logger"
	"github.com/jira-backend/jiraflow-backend/internal/pkg/token"
)

// Context keys for downstream handlers.
const (
	CtxUserID    = "user_id"
	CtxSessionID = "session_id"
	CtxRole      = "role"
)

// validRoles is the set of roles the system recognises.
var validRoles = map[string]struct{}{
	"admin":  {},
	"member": {},
	"viewer": {},
}

// EnforceCasbin checks Casbin policy using the role already injected into
// context by the Auth middleware. Must be placed AFTER Auth in the chain.
// Returns 403 when the role lacks permission for the route+method.
func EnforceCasbin(enforcer *casbin.Enforcer, log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.GetString(CtxRole)
		if role == "" {
			role = "viewer" // safe fallback
		}

		allowed, err := enforcer.Enforce(role, c.FullPath(), c.Request.Method)
		if err != nil {
			log.Error(c.Request.Context(), "casbin enforce error", logger.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error", "code": "INTERNAL_ERROR"})
			c.Abort()
			return
		}
		if !allowed {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "access denied",
				"code":  "FORBIDDEN",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RBAC validates JWT and enforces Casbin in one step (kept for compatibility).
func RBAC(enforcer *casbin.Enforcer, maker token.Maker, log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		role, claims, err := extractRole(c, maker, log)
		if err != nil {
			if errors.Is(err, jwt.ErrTokenExpired) {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "token expired", "code": "TOKEN_EXPIRED"})
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token", "code": "TOKEN_INVALID"})
			}
			c.Abort()
			return
		}

		allowed, err := enforcer.Enforce(role, c.FullPath(), c.Request.Method)
		if err != nil {
			log.Error(ctx, "casbin enforce error", logger.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
			c.Abort()
			return
		}
		if !allowed {
			c.JSON(http.StatusForbidden, gin.H{"error": "access denied", "code": "FORBIDDEN"})
			c.Abort()
			return
		}

		if claims != nil {
			c.Set(CtxUserID, claims.Sub)
			c.Set(CtxSessionID, claims.SessionID)
			c.Set(CtxRole, claims.Role)
		}

		c.Next()
	}
}

// extractRole parses the Authorization header and returns the role to enforce.
// Returns ("unauthorized", nil, nil) when no token is present.
// Returns ("", nil, err) when the token is present but invalid.
func extractRole(c *gin.Context, maker token.Maker, log logger.Logger) (string, *token.Claims, error) {
	raw := c.GetHeader("Authorization")
	if raw == "" {
		return "unauthorized", nil, nil
	}

	tokenStr, ok := strings.CutPrefix(raw, "Bearer ")
	if !ok || tokenStr == "" {
		return "unauthorized", nil, nil
	}

	claims, err := maker.ValidateAccess(c.Request.Context(), tokenStr)
	if err != nil {
		return "", nil, err
	}

	role := claims.Role
	if _, known := validRoles[role]; !known {
		log.Warn(c.Request.Context(), "rbac: unrecognised role in token",
			logger.String("role", role),
			logger.String("sub", claims.Sub),
		)
		role = "viewer" // safe fallback — least privilege
	}

	return role, claims, nil
}
