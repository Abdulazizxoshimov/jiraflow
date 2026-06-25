package v1

import (
	"github.com/gin-gonic/gin"

	hs "github.com/jira-backend/jiraflow-backend/api/http_status"
	"github.com/jira-backend/jiraflow-backend/api/handlers"
	"github.com/jira-backend/jiraflow-backend/api/middleware"
	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

// CreateIssue godoc
// @Summary      Create issue
// @Tags         issues
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body  entity.CreateIssueReq  true  "Issue data"
// @Success      201  {object}  object{data=entity.Issue}
// @Failure      400  {object}  object{code=string,message=string}
// @Router       /api/v1/issues [post]
func CreateIssue(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		reporterID := c.GetString(middleware.CtxUserID)
		var req entity.CreateIssueReq
		if err := c.ShouldBindJSON(&req); err != nil {
			hs.BadRequest(c, err.Error())
			return
		}
		isAdmin := c.GetString(middleware.CtxRole) == "admin"
		issue, err := h.Issue.Create(c.Request.Context(), req.ProjectID, &req, reporterID, isAdmin)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Created(c, issue)
	}
}

// GetIssue godoc
// @Summary      Get issue by ID
// @Tags         issues
// @Produce      json
// @Security     BearerAuth
// @Param        id  path  string  true  "Issue ID"
// @Success      200  {object}  object{data=entity.Issue}
// @Failure      404  {object}  object{code=string,message=string}
// @Router       /api/v1/issues/{id} [get]
func GetIssue(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		issue, err := h.Issue.GetByID(c.Request.Context(), id)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Success(c, issue)
	}
}

// GetIssueByKey godoc
// @Summary      Get issue by project key (e.g. PROJ-42)
// @Tags         issues
// @Produce      json
// @Security     BearerAuth
// @Param        key  path  string  true  "Issue key"
// @Success      200  {object}  object{data=entity.Issue}
// @Failure      404  {object}  object{code=string,message=string}
// @Router       /api/v1/issues/key/{key} [get]
func GetIssueByKey(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.Param("key")
		issue, err := h.Issue.GetByKey(c.Request.Context(), key)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Success(c, issue)
	}
}

// ListIssues godoc
// @Summary      List issues
// @Tags         issues
// @Produce      json
// @Security     BearerAuth
// @Param        project_id   query  string  false  "Filter by project"
// @Param        sprint_id    query  string  false  "Filter by sprint"
// @Param        assignee_id  query  string  false  "Filter by assignee"
// @Param        status_id    query  string  false  "Filter by status"
// @Param        page         query  int     false  "Page number"
// @Param        limit        query  int     false  "Page size"
// @Success      200  {object}  object{data=[]entity.Issue,total=int,page=int,limit=int,total_pages=int}
// @Router       /api/v1/issues [get]
func ListIssues(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		var filter entity.IssueFilter
		if err := c.ShouldBindQuery(&filter); err != nil {
			hs.BadRequest(c, err.Error())
			return
		}
		filter.CurrentUserID = c.GetString(middleware.CtxUserID)
		issues, total, err := h.Issue.List(c.Request.Context(), &filter)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.List(c, issues, total, filter.Page, filter.GetLimit())
	}
}

// UpdateIssue godoc
// @Summary      Update issue
// @Tags         issues
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path  string                true  "Issue ID"
// @Param        body  body  entity.UpdateIssueReq  true  "Update data"
// @Success      200  {object}  object{data=entity.Issue}
// @Failure      400  {object}  object{code=string,message=string}
// @Failure      404  {object}  object{code=string,message=string}
// @Router       /api/v1/issues/{id} [put]
func UpdateIssue(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		actorID := c.GetString(middleware.CtxUserID)
		var req entity.UpdateIssueReq
		if err := c.ShouldBindJSON(&req); err != nil {
			hs.BadRequest(c, err.Error())
			return
		}
		issue, err := h.Issue.Update(c.Request.Context(), id, &req, actorID)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Success(c, issue)
	}
}

// DeleteIssue godoc
// @Summary      Delete issue
// @Tags         issues
// @Security     BearerAuth
// @Param        id  path  string  true  "Issue ID"
// @Success      204
// @Failure      403  {object}  object{code=string,message=string}
// @Failure      404  {object}  object{code=string,message=string}
// @Router       /api/v1/issues/{id} [delete]
func DeleteIssue(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		actorID := c.GetString(middleware.CtxUserID)
		if err := h.Issue.Delete(c.Request.Context(), id, actorID); err != nil {
			hs.Error(c, err)
			return
		}
		hs.NoContent(c)
	}
}

// ReorderIssues godoc
// @Summary      Reorder issues (backlog drag-and-drop)
// @Tags         issues
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body  entity.ReorderIssuesReq  true  "New positions"
// @Success      204
// @Failure      400  {object}  object{code=string,message=string}
// @Router       /api/v1/issues/reorder [put]
func ReorderIssues(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req entity.ReorderIssuesReq
		if err := c.ShouldBindJSON(&req); err != nil {
			hs.BadRequest(c, err.Error())
			return
		}
		if err := h.Issue.ReorderIssues(c.Request.Context(), &req); err != nil {
			hs.Error(c, err)
			return
		}
		hs.NoContent(c)
	}
}

func CloneIssue(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		reporterID := c.GetString(middleware.CtxUserID)
		var req entity.CloneIssueReq
		if err := c.ShouldBindJSON(&req); err != nil {
			hs.BadRequest(c, err.Error())
			return
		}
		clone, err := h.Issue.Clone(c.Request.Context(), c.Param("id"), reporterID, &req)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Created(c, clone)
	}
}

func RankIssue(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req entity.RankIssueReq
		if err := c.ShouldBindJSON(&req); err != nil {
			hs.BadRequest(c, err.Error())
			return
		}
		if err := h.Issue.RankBetween(c.Request.Context(), c.Param("id"), &req); err != nil {
			hs.Error(c, err)
			return
		}
		hs.NoContent(c)
	}
}

// MoveIssueOnBoard godoc
// @Summary      Move issue on board (change column/status and position)
// @Tags         issues
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path  string               true  "Issue ID"
// @Param        body  body  entity.MoveIssueReq  true  "Move data"
// @Success      200  {object}  object{data=entity.Issue}
// @Failure      400  {object}  object{code=string,message=string}
// @Failure      404  {object}  object{code=string,message=string}
// @Router       /api/v1/issues/{id}/move [put]
func MoveIssueOnBoard(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		actorID := c.GetString(middleware.CtxUserID)
		var req entity.MoveIssueReq
		if err := c.ShouldBindJSON(&req); err != nil {
			hs.BadRequest(c, err.Error())
			return
		}
		issue, err := h.Issue.MoveOnBoard(c.Request.Context(), id, &req, actorID)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Success(c, issue)
	}
}
