package v1

import (
	"github.com/gin-gonic/gin"

	hs "github.com/jira-backend/jiraflow-backend/api/http_status"
	"github.com/jira-backend/jiraflow-backend/api/handlers"
	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

// ListIssueHistory godoc
// @Summary      List issue change history
// @Tags         issues
// @Produce      json
// @Security     BearerAuth
// @Param        id     path   string  true   "Issue ID"
// @Param        page   query  int     false  "Page number"
// @Param        limit  query  int     false  "Page size"
// @Success      200  {object}  object{data=[]entity.IssueHistory,total=int,page=int,limit=int,total_pages=int}
// @Failure      404  {object}  object{code=string,message=string}
// @Router       /api/v1/issues/{id}/history [get]
func ListIssueHistory(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		issueID := c.Param("id")
		var filter entity.Filter
		if err := c.ShouldBindQuery(&filter); err != nil {
			hs.BadRequest(c, err.Error())
			return
		}
		history, total, err := h.Issue.ListHistory(c.Request.Context(), issueID, &filter)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.List(c, history, total, filter.Page, filter.GetLimit())
	}
}
