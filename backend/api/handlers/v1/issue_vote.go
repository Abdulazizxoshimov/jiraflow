package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"

	hs "github.com/jira-backend/jiraflow-backend/api/http_status"
	"github.com/jira-backend/jiraflow-backend/api/handlers"
)

// ToggleIssueVote godoc
// @Summary  Toggle vote on an issue
// @Tags     issues
// @Produce  json
// @Param    id  path  string  true  "Issue ID"
// @Success  200 {object} map[string]interface{}
// @Router   /issues/{id}/votes [post]
func ToggleIssueVote(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		actorID := c.GetString("user_id")
		added, err := h.IssueVote.Toggle(c.Request.Context(), c.Param("id"), actorID)
		if err != nil {
			hs.Error(c, err)
			return
		}
		c.JSON(http.StatusOK, gin.H{"voted": added})
	}
}

// GetIssueVoteSummary godoc
// @Summary  Get vote summary for an issue
// @Tags     issues
// @Produce  json
// @Param    id  path  string  true  "Issue ID"
// @Success  200 {object} entity.IssueVoteSummary
// @Router   /issues/{id}/votes [get]
func GetIssueVoteSummary(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		actorID := c.GetString("user_id")
		summary, err := h.IssueVote.GetSummary(c.Request.Context(), c.Param("id"), actorID)
		if err != nil {
			hs.Error(c, err)
			return
		}
		c.JSON(http.StatusOK, summary)
	}
}
