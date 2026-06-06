package middleware

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jira-backend/jiraflow-backend/internal/infrastructure/repository"
)

// RequireProjectMember checks that the authenticated user is a member of the
// project whose ID is at the given URL param name (e.g. "id").
func RequireProjectMember(repo repository.ProjectMemberRepository, paramName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, ok := c.Get(CtxUserID)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		projectID := c.Param(paramName)
		if projectID == "" {
			c.Next()
			return
		}
		isMember, err := repo.IsMember(context.Background(), projectID, userID.(string))
		if err != nil || !isMember {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "you are not a member of this project"})
			return
		}
		c.Next()
	}
}

// RequireSpaceMember checks that the authenticated user is a member of the
// space whose ID is at the given URL param name (e.g. "id").
func RequireSpaceMember(repo repository.SpaceRepository, paramName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		spaceID := c.Param(paramName)
		if spaceID == "" {
			c.Next()
			return
		}
		userID, ok := c.Get(CtxUserID)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		isMember, err := repo.IsMember(context.Background(), spaceID, userID.(string))
		if err != nil || !isMember {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "you are not a member of this space"})
			return
		}
		c.Next()
	}
}
