package v1

import (
	"github.com/gin-gonic/gin"

	hs "github.com/jira-backend/jiraflow-backend/api/http_status"
	"github.com/jira-backend/jiraflow-backend/api/handlers"
	"github.com/jira-backend/jiraflow-backend/api/middleware"
	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

// CreateSprint godoc
// @Summary      Create sprint
// @Tags         sprints
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        project_id  path  string                 true  "Project ID"
// @Param        body        body  entity.CreateSprintReq  true  "Sprint data"
// @Success      201  {object}  object{data=entity.Sprint}
// @Failure      400  {object}  object{code=string,message=string}
// @Router       /api/v1/projects/{project_id}/sprints [post]
func CreateSprint(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		projectID := c.Param("id")
		actorID := c.GetString(middleware.CtxUserID)
		var req entity.CreateSprintReq
		if err := c.ShouldBindJSON(&req); err != nil {
			hs.BadRequest(c, err.Error())
			return
		}
		s := &entity.Sprint{
			Name:      req.Name,
			Goal:      req.Goal,
			StartDate: req.StartDate,
			EndDate:   req.EndDate,
		}
		isAdmin := c.GetString(middleware.CtxRole) == "admin"
		sprint, err := h.Sprint.Create(c.Request.Context(), projectID, actorID, isAdmin, s)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Created(c, sprint)
	}
}

// GetSprint godoc
// @Summary      Get sprint by ID
// @Tags         sprints
// @Produce      json
// @Security     BearerAuth
// @Param        id  path  string  true  "Sprint ID"
// @Success      200  {object}  object{data=entity.Sprint}
// @Failure      404  {object}  object{code=string,message=string}
// @Router       /api/v1/sprints/{id} [get]
func GetSprint(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		sprint, err := h.Sprint.GetByID(c.Request.Context(), id)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Success(c, sprint)
	}
}

// ListSprints godoc
// @Summary      List sprints by project
// @Tags         sprints
// @Produce      json
// @Security     BearerAuth
// @Param        project_id  path   string  true   "Project ID"
// @Param        page        query  int     false  "Page number"
// @Param        limit       query  int     false  "Page size"
// @Success      200  {object}  object{data=[]entity.Sprint,total=int,page=int,limit=int,total_pages=int}
// @Router       /api/v1/projects/{project_id}/sprints [get]
func ListSprints(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		projectID := c.Param("id")
		var filter entity.SprintFilter
		if err := c.ShouldBindQuery(&filter); err != nil {
			hs.BadRequest(c, err.Error())
			return
		}
		sprints, total, err := h.Sprint.List(c.Request.Context(), projectID, &filter)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.List(c, sprints, total, filter.Page, filter.GetLimit())
	}
}

// UpdateSprint godoc
// @Summary      Update sprint
// @Tags         sprints
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path  string                 true  "Sprint ID"
// @Param        body  body  entity.UpdateSprintReq  true  "Update data"
// @Success      200  {object}  object{data=entity.Sprint}
// @Failure      404  {object}  object{code=string,message=string}
// @Router       /api/v1/sprints/{id} [put]
func UpdateSprint(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var req entity.UpdateSprintReq
		if err := c.ShouldBindJSON(&req); err != nil {
			hs.BadRequest(c, err.Error())
			return
		}
		s := &entity.Sprint{}
		if req.Name != nil {
			s.Name = *req.Name
		}
		if req.Goal != nil {
			s.Goal = req.Goal
		}
		if req.StartDate != nil {
			s.StartDate = req.StartDate
		}
		if req.EndDate != nil {
			s.EndDate = req.EndDate
		}
		sprint, err := h.Sprint.Update(c.Request.Context(), id, s)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Success(c, sprint)
	}
}

// DeleteSprint godoc
// @Summary      Delete sprint
// @Tags         sprints
// @Security     BearerAuth
// @Param        id  path  string  true  "Sprint ID"
// @Success      204
// @Failure      404  {object}  object{code=string,message=string}
// @Router       /api/v1/sprints/{id} [delete]
func DeleteSprint(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if err := h.Sprint.Delete(c.Request.Context(), id); err != nil {
			hs.Error(c, err)
			return
		}
		hs.NoContent(c)
	}
}

// StartSprint godoc
// @Summary      Start sprint
// @Tags         sprints
// @Produce      json
// @Security     BearerAuth
// @Param        id  path  string  true  "Sprint ID"
// @Success      200  {object}  object{data=entity.Sprint}
// @Failure      404  {object}  object{code=string,message=string}
// @Router       /api/v1/sprints/{id}/start [post]
func StartSprint(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		actorID := c.GetString(middleware.CtxUserID)
		sprint, err := h.Sprint.Start(c.Request.Context(), id, actorID)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Success(c, sprint)
	}
}

// CompleteSprint godoc
// @Summary      Complete sprint
// @Tags         sprints
// @Produce      json
// @Security     BearerAuth
// @Param        id  path  string  true  "Sprint ID"
// @Success      200  {object}  object{data=entity.Sprint}
// @Failure      404  {object}  object{code=string,message=string}
// @Router       /api/v1/sprints/{id}/complete [post]
func CompleteSprint(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		actorID := c.GetString(middleware.CtxUserID)
		sprint, err := h.Sprint.Complete(c.Request.Context(), id, actorID)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Success(c, sprint)
	}
}

// AddIssueToSprint godoc
// @Summary      Add issue to sprint
// @Tags         sprints
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path  string  true  "Sprint ID"
// @Param        body  body  object{issue_id=string}  true  "Issue ID"
// @Success      204
// @Failure      400  {object}  object{code=string,message=string}
// @Failure      404  {object}  object{code=string,message=string}
// @Router       /api/v1/sprints/{id}/issues [post]
func AddIssueToSprint(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		sprintID := c.Param("id")
		actorID := c.GetString(middleware.CtxUserID)
		var req struct {
			IssueID string `json:"issue_id" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			hs.BadRequest(c, err.Error())
			return
		}
		if err := h.Sprint.AddIssue(c.Request.Context(), sprintID, req.IssueID, actorID); err != nil {
			hs.Error(c, err)
			return
		}
		hs.NoContent(c)
	}
}

// RemoveIssueFromSprint godoc
// @Summary      Remove issue from sprint
// @Tags         sprints
// @Produce      json
// @Security     BearerAuth
// @Param        id        path  string  true  "Sprint ID"
// @Param        issue_id  path  string  true  "Issue ID"
// @Success      204
// @Failure      400  {object}  object{code=string,message=string}
// @Failure      404  {object}  object{code=string,message=string}
// @Router       /api/v1/sprints/{id}/issues/{issue_id} [delete]
func RemoveIssueFromSprint(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		sprintID := c.Param("id")
		issueID := c.Param("issue_id")
		actorID := c.GetString(middleware.CtxUserID)
		if err := h.Sprint.RemoveIssue(c.Request.Context(), sprintID, issueID, actorID); err != nil {
			hs.Error(c, err)
			return
		}
		hs.NoContent(c)
	}
}
