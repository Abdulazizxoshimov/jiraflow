package v1

import (
	"strconv"

	"github.com/gin-gonic/gin"

	hs "github.com/jira-backend/jiraflow-backend/api/http_status"
	"github.com/jira-backend/jiraflow-backend/api/handlers"
)

func GetSprintReport(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		sprintID := c.Param("id")
		report, err := h.Sprint.GetReport(c.Request.Context(), sprintID)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Success(c, report)
	}
}

func GetSprintBurndown(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		sprintID := c.Param("id")
		chart, err := h.Sprint.GetBurndown(c.Request.Context(), sprintID)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Success(c, chart)
	}
}

func GetSprintBurnup(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		sprintID := c.Param("id")
		chart, err := h.Sprint.GetBurnup(c.Request.Context(), sprintID)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Success(c, chart)
	}
}

func GetProjectCFD(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		projectID := c.Param("id")
		from := c.Query("from")
		to := c.Query("to")
		var fromPtr, toPtr *string
		if from != "" {
			fromPtr = &from
		}
		if to != "" {
			toPtr = &to
		}
		chart, err := h.Sprint.GetCFD(c.Request.Context(), projectID, fromPtr, toPtr)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Success(c, chart)
	}
}

func GetProjectVelocity(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		projectID := c.Param("id")
		limit := 10
		if l := c.Query("limit"); l != "" {
			if v, err := strconv.Atoi(l); err == nil && v > 0 {
				limit = v
			}
		}
		report, err := h.Sprint.GetVelocity(c.Request.Context(), projectID, limit)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Success(c, report)
	}
}
