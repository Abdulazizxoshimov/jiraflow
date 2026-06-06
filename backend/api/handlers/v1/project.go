package v1

import (
	"github.com/gin-gonic/gin"

	hs "github.com/jira-backend/jiraflow-backend/api/http_status"
	"github.com/jira-backend/jiraflow-backend/api/handlers"
	"github.com/jira-backend/jiraflow-backend/api/middleware"
	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

// CreateProject godoc
// @Summary      Create project
// @Tags         projects
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body  entity.CreateProjectReq  true  "Project data"
// @Success      201  {object}  object{data=entity.Project}
// @Failure      400  {object}  object{code=string,message=string}
// @Router       /api/v1/projects [post]
func CreateProject(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		ownerID := c.GetString(middleware.CtxUserID)
		var req entity.CreateProjectReq
		if err := c.ShouldBindJSON(&req); err != nil {
			hs.BadRequest(c, err.Error())
			return
		}
		p := &entity.Project{
			Key:         req.Key,
			Name:        req.Name,
			Description: req.Description,
			IconURL:     req.IconURL,
			WorkflowID:  req.WorkflowID,
			LeadID:      ownerID,
		}
		proj, err := h.Project.Create(c.Request.Context(), p, ownerID)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Created(c, proj)
	}
}

// GetProject godoc
// @Summary      Get project by ID
// @Tags         projects
// @Produce      json
// @Security     BearerAuth
// @Param        id  path  string  true  "Project ID"
// @Success      200  {object}  object{data=entity.Project}
// @Failure      404  {object}  object{code=string,message=string}
// @Router       /api/v1/projects/{id} [get]
func GetProject(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		proj, err := h.Project.GetByID(c.Request.Context(), id)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Success(c, proj)
	}
}

// ListProjects godoc
// @Summary      List projects
// @Tags         projects
// @Produce      json
// @Security     BearerAuth
// @Param        page   query  int     false  "Page number"
// @Param        limit  query  int     false  "Page size"
// @Param        q      query  string  false  "Search query"
// @Success      200  {object}  object{data=[]entity.Project,total=int,page=int,limit=int,total_pages=int}
// @Router       /api/v1/projects [get]
func ListProjects(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		var filter entity.ProjectFilter
		if err := c.ShouldBindQuery(&filter); err != nil {
			hs.BadRequest(c, err.Error())
			return
		}
		projects, total, err := h.Project.List(c.Request.Context(), &filter)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.List(c, projects, total, filter.Page, filter.GetLimit())
	}
}

// UpdateProject godoc
// @Summary      Update project
// @Tags         projects
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path  string                true  "Project ID"
// @Param        body  body  entity.UpdateProjectReq  true  "Update data"
// @Success      200  {object}  object{data=entity.Project}
// @Failure      400  {object}  object{code=string,message=string}
// @Failure      404  {object}  object{code=string,message=string}
// @Router       /api/v1/projects/{id} [put]
func UpdateProject(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		actorID := c.GetString(middleware.CtxUserID)
		var req entity.UpdateProjectReq
		if err := c.ShouldBindJSON(&req); err != nil {
			hs.BadRequest(c, err.Error())
			return
		}
		p := &entity.Project{}
		if req.Name != nil {
			p.Name = *req.Name
		}
		if req.Description != nil {
			p.Description = req.Description
		}
		if req.IconURL != nil {
			p.IconURL = req.IconURL
		}
		if req.LeadID != nil {
			p.LeadID = *req.LeadID
		}
		proj, err := h.Project.Update(c.Request.Context(), id, p, actorID)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Success(c, proj)
	}
}

// DeleteProject godoc
// @Summary      Delete project
// @Tags         projects
// @Security     BearerAuth
// @Param        id  path  string  true  "Project ID"
// @Success      204
// @Failure      403  {object}  object{code=string,message=string}
// @Failure      404  {object}  object{code=string,message=string}
// @Router       /api/v1/projects/{id} [delete]
func DeleteProject(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		actorID := c.GetString(middleware.CtxUserID)
		if err := h.Project.Delete(c.Request.Context(), id, actorID); err != nil {
			hs.Error(c, err)
			return
		}
		hs.NoContent(c)
	}
}

// ArchiveProject godoc
// @Summary      Archive project
// @Tags         projects
// @Security     BearerAuth
// @Param        id  path  string  true  "Project ID"
// @Success      204
// @Failure      403  {object}  object{code=string,message=string}
// @Failure      404  {object}  object{code=string,message=string}
// @Router       /api/v1/projects/{id}/archive [post]
func ArchiveProject(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		actorID := c.GetString(middleware.CtxUserID)
		if err := h.Project.Archive(c.Request.Context(), id, actorID); err != nil {
			hs.Error(c, err)
			return
		}
		hs.NoContent(c)
	}
}

// GetProjectDashboard godoc
// @Summary      Get project dashboard statistics
// @Tags         projects
// @Produce      json
// @Security     BearerAuth
// @Param        id  path  string  true  "Project ID"
// @Success      200  {object}  object{data=entity.ProjectDashboard}
// @Failure      404  {object}  object{code=string,message=string}
// @Router       /api/v1/projects/{id}/dashboard [get]
func GetProjectDashboard(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		dash, err := h.Project.GetDashboard(c.Request.Context(), id)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Success(c, dash)
	}
}
