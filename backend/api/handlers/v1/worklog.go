package v1

import (
	"github.com/gin-gonic/gin"

	hs "github.com/jira-backend/jiraflow-backend/api/http_status"
	"github.com/jira-backend/jiraflow-backend/api/handlers"
	"github.com/jira-backend/jiraflow-backend/api/middleware"
	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

func CreateWorklog(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		issueID := c.Param("id")
		userID := c.GetString(middleware.CtxUserID)
		var req entity.CreateWorklogReq
		if err := c.ShouldBindJSON(&req); err != nil {
			hs.BadRequest(c, err.Error())
			return
		}
		w, err := h.Worklog.Create(c.Request.Context(), issueID, userID, &req)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Created(c, w)
	}
}

func ListWorklogs(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		issueID := c.Param("id")
		var filter entity.WorklogFilter
		if err := c.ShouldBindQuery(&filter); err != nil {
			hs.BadRequest(c, err.Error())
			return
		}
		filter.IssueID = issueID
		worklogs, total, err := h.Worklog.List(c.Request.Context(), &filter)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.List(c, worklogs, total, filter.Page, filter.GetLimit())
	}
}

func UpdateWorklog(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("worklog_id")
		actorID := c.GetString(middleware.CtxUserID)
		var req entity.UpdateWorklogReq
		if err := c.ShouldBindJSON(&req); err != nil {
			hs.BadRequest(c, err.Error())
			return
		}
		w, err := h.Worklog.Update(c.Request.Context(), id, actorID, &req)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Success(c, w)
	}
}

func DeleteWorklog(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("worklog_id")
		actorID := c.GetString(middleware.CtxUserID)
		if err := h.Worklog.Delete(c.Request.Context(), id, actorID); err != nil {
			hs.Error(c, err)
			return
		}
		hs.NoContent(c)
	}
}

func GetTimeSpentSummary(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		issueID := c.Param("id")
		summary, err := h.Worklog.GetTimeSummary(c.Request.Context(), issueID)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Success(c, summary)
	}
}

func UpdateEstimates(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		issueID := c.Param("id")
		var req struct {
			OriginalEstimate  *int `json:"original_estimate"`
			RemainingEstimate *int `json:"remaining_estimate"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			hs.BadRequest(c, err.Error())
			return
		}
		if err := h.Worklog.UpdateEstimates(c.Request.Context(), issueID, req.OriginalEstimate, req.RemainingEstimate); err != nil {
			hs.Error(c, err)
			return
		}
		hs.NoContent(c)
	}
}
