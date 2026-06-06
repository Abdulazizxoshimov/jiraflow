package v1

import (
	"github.com/gin-gonic/gin"

	hs "github.com/jira-backend/jiraflow-backend/api/http_status"
	"github.com/jira-backend/jiraflow-backend/api/handlers"
	"github.com/jira-backend/jiraflow-backend/api/middleware"
	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

// CreatePermissionScheme godoc
// @Summary      Create a permission scheme
// @Tags         permission-schemes
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body  entity.CreatePermissionSchemeReq  true  "Scheme"
// @Success      201  {object}  object{data=entity.PermissionScheme}
// @Router       /api/v1/permission-schemes [post]
func CreatePermissionScheme(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString(middleware.CtxUserID)
		var req entity.CreatePermissionSchemeReq
		if err := c.ShouldBindJSON(&req); err != nil {
			hs.BadRequest(c, err.Error())
			return
		}
		s, err := h.PermissionScheme.Create(c.Request.Context(), &req, userID)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Created(c, s)
	}
}

// ListPermissionSchemes godoc
// @Summary      List all permission schemes
// @Tags         permission-schemes
// @Security     BearerAuth
// @Success      200  {object}  object{data=[]entity.PermissionScheme}
// @Router       /api/v1/permission-schemes [get]
func ListPermissionSchemes(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		list, err := h.PermissionScheme.List(c.Request.Context())
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Success(c, list)
	}
}

// GetPermissionScheme godoc
// @Summary      Get a permission scheme by ID
// @Tags         permission-schemes
// @Security     BearerAuth
// @Param        id  path  string  true  "Scheme ID"
// @Success      200  {object}  object{data=entity.PermissionScheme}
// @Router       /api/v1/permission-schemes/{id} [get]
func GetPermissionScheme(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		s, err := h.PermissionScheme.Get(c.Request.Context(), c.Param("id"))
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Success(c, s)
	}
}

// UpdatePermissionScheme godoc
// @Summary      Update a permission scheme
// @Tags         permission-schemes
// @Security     BearerAuth
// @Param        id    path  string                          true  "Scheme ID"
// @Param        body  body  entity.UpdatePermissionSchemeReq  true  "Updates"
// @Success      200  {object}  object{data=entity.PermissionScheme}
// @Router       /api/v1/permission-schemes/{id} [put]
func UpdatePermissionScheme(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req entity.UpdatePermissionSchemeReq
		if err := c.ShouldBindJSON(&req); err != nil {
			hs.BadRequest(c, err.Error())
			return
		}
		s, err := h.PermissionScheme.Update(c.Request.Context(), c.Param("id"), &req)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Success(c, s)
	}
}

// DeletePermissionScheme godoc
// @Summary      Delete a permission scheme
// @Tags         permission-schemes
// @Security     BearerAuth
// @Param        id  path  string  true  "Scheme ID"
// @Success      204
// @Router       /api/v1/permission-schemes/{id} [delete]
func DeletePermissionScheme(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := h.PermissionScheme.Delete(c.Request.Context(), c.Param("id")); err != nil {
			hs.Error(c, err)
			return
		}
		hs.NoContent(c)
	}
}

// AddSchemeGrant godoc
// @Summary      Add a grant to a permission scheme
// @Tags         permission-schemes
// @Security     BearerAuth
// @Param        id    path  string             true  "Scheme ID"
// @Param        body  body  entity.AddGrantReq  true  "Grant"
// @Success      201  {object}  object{data=entity.PermissionSchemeGrant}
// @Router       /api/v1/permission-schemes/{id}/grants [post]
func AddSchemeGrant(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req entity.AddGrantReq
		if err := c.ShouldBindJSON(&req); err != nil {
			hs.BadRequest(c, err.Error())
			return
		}
		g, err := h.PermissionScheme.AddGrant(c.Request.Context(), c.Param("id"), &req)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Created(c, g)
	}
}

// RemoveSchemeGrant godoc
// @Summary      Remove a grant from a permission scheme
// @Tags         permission-schemes
// @Security     BearerAuth
// @Param        id        path  string  true  "Scheme ID"
// @Param        grant_id  path  string  true  "Grant ID"
// @Success      204
// @Router       /api/v1/permission-schemes/{id}/grants/{grant_id} [delete]
func RemoveSchemeGrant(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := h.PermissionScheme.RemoveGrant(c.Request.Context(), c.Param("id"), c.Param("grant_id")); err != nil {
			hs.Error(c, err)
			return
		}
		hs.NoContent(c)
	}
}

// AssignPermissionScheme godoc
// @Summary      Assign a permission scheme to a project
// @Tags         permission-schemes
// @Security     BearerAuth
// @Param        id    path  string                   true  "Project ID"
// @Param        body  body  entity.AssignSchemeReq    true  "Scheme"
// @Success      204
// @Router       /api/v1/projects/{id}/permission-scheme [put]
func AssignPermissionScheme(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req entity.AssignSchemeReq
		if err := c.ShouldBindJSON(&req); err != nil {
			hs.BadRequest(c, err.Error())
			return
		}
		if err := h.PermissionScheme.AssignToProject(c.Request.Context(), c.Param("id"), req.SchemeID); err != nil {
			hs.Error(c, err)
			return
		}
		hs.NoContent(c)
	}
}

// GetProjectPermissionScheme godoc
// @Summary      Get the permission scheme assigned to a project
// @Tags         permission-schemes
// @Security     BearerAuth
// @Param        id  path  string  true  "Project ID"
// @Success      200  {object}  object{data=entity.PermissionScheme}
// @Router       /api/v1/projects/{id}/permission-scheme [get]
func GetProjectPermissionScheme(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		s, err := h.PermissionScheme.GetByProject(c.Request.Context(), c.Param("id"))
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Success(c, s)
	}
}
