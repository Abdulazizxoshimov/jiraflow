package v1

import (
	"github.com/gin-gonic/gin"

	hs "github.com/jira-backend/jiraflow-backend/api/http_status"
	"github.com/jira-backend/jiraflow-backend/api/handlers"
	"github.com/jira-backend/jiraflow-backend/api/middleware"
)

type transitionReq struct {
	StatusID string `json:"status_id" binding:"required"`
}

// TransitionIssue godoc
// @Summary      Transition issue status
// @Tags         issues
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path  string          true  "Issue ID"
// @Param        body  body  transitionReq   true  "Target status ID"
// @Success      200  {object}  object{data=entity.Issue}
// @Failure      400  {object}  object{code=string,message=string}
// @Failure      404  {object}  object{code=string,message=string}
// @Router       /api/v1/issues/{id}/transition [post]
func TransitionIssue(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		actorID := c.GetString(middleware.CtxUserID)
		var req transitionReq
		if err := c.ShouldBindJSON(&req); err != nil {
			hs.BadRequest(c, err.Error())
			return
		}
		issue, err := h.Issue.Transition(c.Request.Context(), id, req.StatusID, actorID)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Success(c, issue)
	}
}
