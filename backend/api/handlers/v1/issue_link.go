package v1

import (
	"github.com/gin-gonic/gin"

	hs "github.com/jira-backend/jiraflow-backend/api/http_status"
	"github.com/jira-backend/jiraflow-backend/api/handlers"
	"github.com/jira-backend/jiraflow-backend/api/middleware"
	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

// CreateIssueLink godoc
// @Summary      Link two issues
// @Tags         issue-links
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        issue_id  path  string                   true  "Source Issue ID"
// @Param        body      body  entity.CreateIssueLinkReq  true  "Link data"
// @Success      201  {object}  object{data=entity.IssueLink}
// @Failure      400  {object}  object{code=string,message=string}
// @Router       /api/v1/issues/{issue_id}/links [post]
func CreateIssueLink(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		sourceID := c.Param("id")
		actorID := c.GetString(middleware.CtxUserID)
		var req entity.CreateIssueLinkReq
		if err := c.ShouldBindJSON(&req); err != nil {
			hs.BadRequest(c, err.Error())
			return
		}
		link := &entity.IssueLink{
			SourceID:  sourceID,
			TargetID:  req.TargetID,
			LinkType:  req.LinkType,
			CreatedBy: actorID,
		}
		created, err := h.IssueLink.Create(c.Request.Context(), link)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Created(c, created)
	}
}

// ListIssueLinks godoc
// @Summary      List links for an issue
// @Tags         issue-links
// @Produce      json
// @Security     BearerAuth
// @Param        issue_id  path  string  true  "Issue ID"
// @Success      200  {object}  object{data=[]entity.IssueLink}
// @Failure      404  {object}  object{code=string,message=string}
// @Router       /api/v1/issues/{issue_id}/links [get]
func ListIssueLinks(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		issueID := c.Param("id")
		links, err := h.IssueLink.ListByIssue(c.Request.Context(), issueID)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Success(c, links)
	}
}

// DeleteIssueLink godoc
// @Summary      Delete issue link
// @Tags         issue-links
// @Security     BearerAuth
// @Param        id  path  string  true  "Link ID"
// @Success      204
// @Failure      404  {object}  object{code=string,message=string}
// @Router       /api/v1/issues/links/{id} [delete]
func DeleteIssueLink(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		linkID := c.Param("id")
		actorID := c.GetString(middleware.CtxUserID)
		if err := h.IssueLink.Delete(c.Request.Context(), linkID, actorID); err != nil {
			hs.Error(c, err)
			return
		}
		hs.NoContent(c)
	}
}
