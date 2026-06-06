package v1

import (
	"github.com/gin-gonic/gin"

	hs "github.com/jira-backend/jiraflow-backend/api/http_status"
	"github.com/jira-backend/jiraflow-backend/api/handlers"
	"github.com/jira-backend/jiraflow-backend/api/middleware"
	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

// AddProjectMember godoc
// @Summary      Add member to project
// @Tags         project-members
// @Accept       json
// @Security     BearerAuth
// @Param        project_id  path  string                    true  "Project ID"
// @Param        body        body  entity.AddProjectMemberReq  true  "Member data"
// @Success      204
// @Failure      400  {object}  object{code=string,message=string}
// @Failure      404  {object}  object{code=string,message=string}
// @Router       /api/v1/projects/{project_id}/members [post]
func AddProjectMember(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		projectID := c.Param("id")
		actorID := c.GetString(middleware.CtxUserID)
		var req entity.AddProjectMemberReq
		if err := c.ShouldBindJSON(&req); err != nil {
			hs.BadRequest(c, err.Error())
			return
		}
		if err := h.ProjectMember.Add(c.Request.Context(), projectID, &req, actorID); err != nil {
			hs.Error(c, err)
			return
		}
		hs.NoContent(c)
	}
}

// ListProjectMembers godoc
// @Summary      List project members
// @Tags         project-members
// @Produce      json
// @Security     BearerAuth
// @Param        project_id  path  string  true  "Project ID"
// @Param        page        query int     false "Page number"
// @Param        limit       query int     false "Page size"
// @Success      200  {object}  object{data=[]entity.ProjectMember,total=int,page=int,limit=int,total_pages=int}
// @Router       /api/v1/projects/{project_id}/members [get]
func ListProjectMembers(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		projectID := c.Param("id")
		var filter entity.Filter
		if err := c.ShouldBindQuery(&filter); err != nil {
			hs.BadRequest(c, err.Error())
			return
		}
		members, total, err := h.ProjectMember.ListByProject(c.Request.Context(), projectID, &filter)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.List(c, members, total, filter.Page, filter.GetLimit())
	}
}

// UpdateProjectMemberRole godoc
// @Summary      Update project member role
// @Tags         project-members
// @Accept       json
// @Security     BearerAuth
// @Param        project_id  path  string                          true  "Project ID"
// @Param        user_id     path  string                          true  "User ID"
// @Param        body        body  entity.UpdateProjectMemberRoleReq  true  "Role"
// @Success      204
// @Failure      403  {object}  object{code=string,message=string}
// @Router       /api/v1/projects/{project_id}/members/{user_id} [put]
func UpdateProjectMemberRole(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		projectID := c.Param("id")
		userID := c.Param("user_id")
		actorID := c.GetString(middleware.CtxUserID)
		var req entity.UpdateProjectMemberRoleReq
		if err := c.ShouldBindJSON(&req); err != nil {
			hs.BadRequest(c, err.Error())
			return
		}
		if err := h.ProjectMember.UpdateRole(c.Request.Context(), projectID, userID, &req, actorID); err != nil {
			hs.Error(c, err)
			return
		}
		hs.NoContent(c)
	}
}

// RemoveProjectMember godoc
// @Summary      Remove project member
// @Tags         project-members
// @Security     BearerAuth
// @Param        project_id  path  string  true  "Project ID"
// @Param        user_id     path  string  true  "User ID"
// @Success      204
// @Failure      403  {object}  object{code=string,message=string}
// @Failure      404  {object}  object{code=string,message=string}
// @Router       /api/v1/projects/{project_id}/members/{user_id} [delete]
func RemoveProjectMember(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		projectID := c.Param("id")
		userID := c.Param("user_id")
		actorID := c.GetString(middleware.CtxUserID)
		if err := h.ProjectMember.Remove(c.Request.Context(), projectID, userID, actorID); err != nil {
			hs.Error(c, err)
			return
		}
		hs.NoContent(c)
	}
}
