package v1

import (
	"github.com/gin-gonic/gin"

	hs "github.com/jira-backend/jiraflow-backend/api/http_status"
	"github.com/jira-backend/jiraflow-backend/api/handlers"
	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

func SetIssueAssignees(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req entity.SetIssueAssigneesReq
		if err := c.ShouldBindJSON(&req); err != nil {
			hs.BadRequest(c, err.Error())
			return
		}
		if err := h.IssueAssignee.Set(c.Request.Context(), c.Param("id"), &req); err != nil {
			hs.Error(c, err)
			return
		}
		hs.NoContent(c)
	}
}

func ListIssueAssignees(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		assignees, err := h.IssueAssignee.List(c.Request.Context(), c.Param("id"))
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Success(c, assignees)
	}
}

func RemoveIssueAssignee(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := h.IssueAssignee.Remove(c.Request.Context(), c.Param("id"), c.Param("user_id")); err != nil {
			hs.Error(c, err)
			return
		}
		hs.NoContent(c)
	}
}
