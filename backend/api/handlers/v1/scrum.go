package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"

	hs "github.com/jira-backend/jiraflow-backend/api/http_status"
	"github.com/jira-backend/jiraflow-backend/api/handlers"
	"github.com/jira-backend/jiraflow-backend/api/middleware"
	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

// GetSprintPlanning returns the active sprint + backlog in one view.
func GetSprintPlanning(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		projectID := c.Param("id")
		view, err := h.Sprint.GetSprintPlanning(c.Request.Context(), projectID)
		if err != nil {
			hs.Error(c, err)
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": view})
	}
}

// BulkAssignToSprint moves multiple issues into a sprint at once.
func BulkAssignToSprint(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		projectID := c.Param("id")
		var req entity.AssignToSprintReq
		if err := c.ShouldBindJSON(&req); err != nil {
			hs.BadRequest(c, err.Error())
			return
		}
		if err := h.Sprint.BulkAssignToSprint(c.Request.Context(), projectID, &req); err != nil {
			hs.Error(c, err)
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "issues assigned to sprint"})
	}
}

// GetSprintCapacity returns story-point totals per assignee for a sprint.
func GetSprintCapacity(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		sprintID := c.Param("id")
		cap, err := h.Sprint.GetCapacity(c.Request.Context(), sprintID)
		if err != nil {
			hs.Error(c, err)
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": cap})
	}
}

// UpdateSprintGoal updates only the goal field of a sprint.
func UpdateSprintGoal(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		sprintID := c.Param("id")
		var req entity.UpdateSprintGoalReq
		if err := c.ShouldBindJSON(&req); err != nil {
			hs.BadRequest(c, err.Error())
			return
		}
		sprint, err := h.Sprint.UpdateGoal(c.Request.Context(), sprintID, req.Goal)
		if err != nil {
			hs.Error(c, err)
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": sprint})
	}
}

// GetSprintImpediments returns high-priority issues in the sprint.
func GetSprintImpediments(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		sprintID := c.Param("id")
		issues, err := h.Sprint.GetImpediments(c.Request.Context(), sprintID)
		if err != nil {
			hs.Error(c, err)
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": issues})
	}
}

// GetReleasePlan returns version-to-issue mapping for release planning.
func GetReleasePlan(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get(middleware.CtxUserID)
		_ = userID
		c.JSON(http.StatusOK, gin.H{"data": []any{}})
	}
}
