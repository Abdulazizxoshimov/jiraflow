package v1

import (
	"github.com/gin-gonic/gin"

	hs "github.com/jira-backend/jiraflow-backend/api/http_status"
	"github.com/jira-backend/jiraflow-backend/api/handlers"
	"github.com/jira-backend/jiraflow-backend/api/middleware"
)

// AddIssueWatcher godoc
// @Summary      Watch an issue
// @Tags         issue-watchers
// @Security     BearerAuth
// @Param        issue_id  path  string  true  "Issue ID"
// @Success      204
// @Failure      404  {object}  object{code=string,message=string}
// @Router       /api/v1/issues/{issue_id}/watchers [post]
func AddIssueWatcher(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		issueID := c.Param("id")
		userID := c.GetString(middleware.CtxUserID)
		if err := h.Issue.AddWatcher(c.Request.Context(), issueID, userID); err != nil {
			hs.Error(c, err)
			return
		}
		hs.NoContent(c)
	}
}

// RemoveIssueWatcher godoc
// @Summary      Unwatch an issue
// @Tags         issue-watchers
// @Security     BearerAuth
// @Param        issue_id  path  string  true  "Issue ID"
// @Success      204
// @Failure      404  {object}  object{code=string,message=string}
// @Router       /api/v1/issues/{issue_id}/watchers [delete]
func RemoveIssueWatcher(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		issueID := c.Param("id")
		userID := c.GetString(middleware.CtxUserID)
		if err := h.Issue.RemoveWatcher(c.Request.Context(), issueID, userID); err != nil {
			hs.Error(c, err)
			return
		}
		hs.NoContent(c)
	}
}

// ListIssueWatchers godoc
// @Summary      List issue watchers
// @Tags         issue-watchers
// @Produce      json
// @Security     BearerAuth
// @Param        issue_id  path  string  true  "Issue ID"
// @Success      200  {object}  object{data=[]entity.IssueWatcher}
// @Failure      404  {object}  object{code=string,message=string}
// @Router       /api/v1/issues/{issue_id}/watchers [get]
func ListIssueWatchers(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		issueID := c.Param("id")
		watchers, err := h.Issue.ListWatchers(c.Request.Context(), issueID)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Success(c, watchers)
	}
}
