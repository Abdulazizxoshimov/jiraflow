package v1

import (
	"github.com/gin-gonic/gin"

	hs "github.com/jira-backend/jiraflow-backend/api/http_status"
	"github.com/jira-backend/jiraflow-backend/api/handlers"
	"github.com/jira-backend/jiraflow-backend/api/middleware"
	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

func BulkUpdateIssues(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		actorID := c.GetString(middleware.CtxUserID)
		var req entity.BulkUpdateIssueReq
		if err := c.ShouldBindJSON(&req); err != nil {
			hs.BadRequest(c, err.Error())
			return
		}
		result, err := h.Issue.BulkUpdate(c.Request.Context(), &req, actorID)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Success(c, result)
	}
}

func BulkDeleteIssues(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		actorID := c.GetString(middleware.CtxUserID)
		var req entity.BulkDeleteIssueReq
		if err := c.ShouldBindJSON(&req); err != nil {
			hs.BadRequest(c, err.Error())
			return
		}
		if err := h.Issue.BulkDelete(c.Request.Context(), &req, actorID); err != nil {
			hs.Error(c, err)
			return
		}
		hs.NoContent(c)
	}
}
