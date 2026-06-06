package v1

import (
	"github.com/gin-gonic/gin"

	hs "github.com/jira-backend/jiraflow-backend/api/http_status"
	"github.com/jira-backend/jiraflow-backend/api/handlers"
	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

func GetRoadmap(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		projectID := c.Param("id")
		items, err := h.Issue.GetRoadmap(c.Request.Context(), projectID)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Success(c, items)
	}
}

func GetBacklog(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		projectID := c.Param("id")
		var filter entity.IssueFilter
		if err := c.ShouldBindQuery(&filter); err != nil {
			hs.BadRequest(c, err.Error())
			return
		}
		issues, total, err := h.Issue.GetBacklog(c.Request.Context(), projectID, &filter)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.List(c, issues, total, filter.Page, filter.GetLimit())
	}
}

func GetEpicProgress(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		epicID := c.Param("id")
		progress, err := h.Issue.GetEpicProgress(c.Request.Context(), epicID)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Success(c, progress)
	}
}
