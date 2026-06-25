package middleware

import (
	"context"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
	"github.com/jira-backend/jiraflow-backend/internal/infrastructure/repository"
)

// AuditLog automatically logs all mutating requests (POST/PUT/PATCH/DELETE) to
// audit_logs. Write happens in a goroutine so it never blocks the response.
// Only 2xx and 4xx responses are recorded; 5xx server errors are skipped.
func AuditLog(repo repository.AuditRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		if method != "POST" && method != "PUT" && method != "PATCH" && method != "DELETE" {
			c.Next()
			return
		}

		c.Next()

		status := c.Writer.Status()
		if status >= 500 {
			return
		}

		userID := c.GetString(CtxUserID)
		ip := c.ClientIP()
		ua := c.Request.UserAgent()
		action := method + " " + c.FullPath()
		entityType, entityID := extractEntity(c.FullPath(), c.Param("id"))

		log := &entity.AuditLog{
			ID:        uuid.NewString(),
			Action:    action,
			Details:   map[string]any{"status": status},
			CreatedAt: time.Now().UTC(),
		}
		if userID != "" {
			log.UserID = &userID
		}
		if ip != "" {
			log.IPAddress = &ip
		}
		if ua != "" {
			log.UserAgent = &ua
		}
		if entityType != "" {
			log.EntityType = &entityType
		}
		if entityID != "" {
			log.EntityID = &entityID
		}

		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()
			_ = repo.Create(ctx, log)
		}()
	}
}

// extractEntity derives entity type and entity ID from a Gin route pattern.
// e.g. "/api/v1/issues/:id" → ("issue", value-of-:id-param)
func extractEntity(fullPath, idParam string) (entityType, entityID string) {
	parts := strings.Split(strings.TrimPrefix(fullPath, "/api/v1/"), "/")
	if len(parts) == 0 {
		return
	}
	// first segment is the plural resource name, singularise naively
	resource := parts[0]
	entityType = strings.TrimSuffix(resource, "s")
	entityID = strings.TrimPrefix(idParam, "/")
	return
}
