package v1

import (
	"github.com/gin-gonic/gin"

	hs "github.com/jira-backend/jiraflow-backend/api/http_status"
	"github.com/jira-backend/jiraflow-backend/api/handlers"
	"github.com/jira-backend/jiraflow-backend/api/middleware"
	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

// LinkIssuePage godoc
// @Summary      Link issue to a Confluence page
// @Tags         issue-page-links
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path  string                       true  "Issue ID"
// @Param        body  body  entity.CreateIssuePageLinkReq  true  "Page ID"
// @Success      201  {object}  object{data=entity.IssuePageLink}
// @Failure      400,404,409  {object}  object{code=string,message=string}
// @Router       /api/v1/issues/{id}/page-links [post]
func LinkIssuePage(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		issueID := c.Param("id")
		actorID := c.GetString(middleware.CtxUserID)
		var req entity.CreateIssuePageLinkReq
		if err := c.ShouldBindJSON(&req); err != nil {
			hs.BadRequest(c, err.Error())
			return
		}
		link, err := h.IssuePageLink.Link(c.Request.Context(), issueID, req.PageID, actorID)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Created(c, link)
	}
}

// UnlinkIssuePage godoc
// @Summary      Remove issue-page link
// @Tags         issue-page-links
// @Security     BearerAuth
// @Param        id       path  string  true  "Issue ID"
// @Param        page_id  path  string  true  "Page ID"
// @Success      204
// @Failure      404  {object}  object{code=string,message=string}
// @Router       /api/v1/issues/{id}/page-links/{page_id} [delete]
func UnlinkIssuePage(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		issueID := c.Param("id")
		pageID := c.Param("page_id")
		if err := h.IssuePageLink.Unlink(c.Request.Context(), issueID, pageID); err != nil {
			hs.Error(c, err)
			return
		}
		hs.NoContent(c)
	}
}

// ListIssuePageLinks godoc
// @Summary      List pages linked to an issue
// @Tags         issue-page-links
// @Produce      json
// @Security     BearerAuth
// @Param        id  path  string  true  "Issue ID"
// @Success      200  {object}  object{data=[]entity.IssuePageLink}
// @Router       /api/v1/issues/{id}/page-links [get]
func ListIssuePageLinks(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		issueID := c.Param("id")
		links, err := h.IssuePageLink.ListByIssue(c.Request.Context(), issueID)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Success(c, links)
	}
}

// ListPageIssueLinks godoc
// @Summary      List issues linked to a page
// @Tags         issue-page-links
// @Produce      json
// @Security     BearerAuth
// @Param        id  path  string  true  "Page ID"
// @Success      200  {object}  object{data=[]entity.IssuePageLink}
// @Router       /api/v1/pages/{id}/issue-links [get]
func ListPageIssueLinks(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		pageID := c.Param("id")
		links, err := h.IssuePageLink.ListByPage(c.Request.Context(), pageID)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Success(c, links)
	}
}
