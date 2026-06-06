package v1

import (
	"github.com/gin-gonic/gin"

	hs "github.com/jira-backend/jiraflow-backend/api/http_status"
	"github.com/jira-backend/jiraflow-backend/api/handlers"
	"github.com/jira-backend/jiraflow-backend/api/middleware"
	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

// CreateSpace godoc
// @Summary      Create space
// @Tags         spaces
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body  entity.CreateSpaceReq  true  "Space data"
// @Success      201  {object}  object{data=entity.Space}
// @Failure      400  {object}  object{code=string,message=string}
// @Router       /api/v1/spaces [post]
func CreateSpace(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		leadID := c.GetString(middleware.CtxUserID)
		var req entity.CreateSpaceReq
		if err := c.ShouldBindJSON(&req); err != nil {
			hs.BadRequest(c, err.Error())
			return
		}
		s := &entity.Space{
			Key:         req.Key,
			Name:        req.Name,
			Description: req.Description,
			IconURL:     req.IconURL,
			Type:        req.Type,
			ProjectID:   req.ProjectID,
		}
		space, err := h.Space.Create(c.Request.Context(), s, leadID)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Created(c, space)
	}
}

// GetSpace godoc
// @Summary      Get space by ID
// @Tags         spaces
// @Produce      json
// @Security     BearerAuth
// @Param        id  path  string  true  "Space ID"
// @Success      200  {object}  object{data=entity.Space}
// @Failure      404  {object}  object{code=string,message=string}
// @Router       /api/v1/spaces/{id} [get]
func GetSpace(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		space, err := h.Space.GetByID(c.Request.Context(), id)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Success(c, space)
	}
}

// ListSpaces godoc
// @Summary      List spaces
// @Tags         spaces
// @Produce      json
// @Security     BearerAuth
// @Param        page   query  int  false  "Page number"
// @Param        limit  query  int  false  "Page size"
// @Success      200  {object}  object{data=[]entity.Space,total=int,page=int,limit=int,total_pages=int}
// @Router       /api/v1/spaces [get]
func ListSpaces(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		var filter entity.Filter
		if err := c.ShouldBindQuery(&filter); err != nil {
			hs.BadRequest(c, err.Error())
			return
		}
		spaces, total, err := h.Space.List(c.Request.Context(), &filter)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.List(c, spaces, total, filter.Page, filter.GetLimit())
	}
}

// UpdateSpace godoc
// @Summary      Update space
// @Tags         spaces
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path  string                true  "Space ID"
// @Param        body  body  entity.UpdateSpaceReq  true  "Update data"
// @Success      200  {object}  object{data=entity.Space}
// @Failure      404  {object}  object{code=string,message=string}
// @Router       /api/v1/spaces/{id} [put]
func UpdateSpace(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		actorID := c.GetString(middleware.CtxUserID)
		var req entity.UpdateSpaceReq
		if err := c.ShouldBindJSON(&req); err != nil {
			hs.BadRequest(c, err.Error())
			return
		}
		s := &entity.Space{}
		if req.Name != nil {
			s.Name = *req.Name
		}
		if req.Description != nil {
			s.Description = req.Description
		}
		if req.IconURL != nil {
			s.IconURL = req.IconURL
		}
		if req.IsArchived != nil {
			s.IsArchived = *req.IsArchived
		}
		space, err := h.Space.Update(c.Request.Context(), id, s, actorID)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Success(c, space)
	}
}

// DeleteSpace godoc
// @Summary      Delete space
// @Tags         spaces
// @Security     BearerAuth
// @Param        id  path  string  true  "Space ID"
// @Success      204
// @Failure      403  {object}  object{code=string,message=string}
// @Failure      404  {object}  object{code=string,message=string}
// @Router       /api/v1/spaces/{id} [delete]
func DeleteSpace(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		actorID := c.GetString(middleware.CtxUserID)
		if err := h.Space.Delete(c.Request.Context(), id, actorID); err != nil {
			hs.Error(c, err)
			return
		}
		hs.NoContent(c)
	}
}

// ArchiveSpace godoc
// @Summary      Archive space
// @Tags         spaces
// @Security     BearerAuth
// @Param        id  path  string  true  "Space ID"
// @Success      204
// @Failure      400  {object}  object{code=string,message=string}
// @Failure      403  {object}  object{code=string,message=string}
// @Failure      404  {object}  object{code=string,message=string}
// @Router       /api/v1/spaces/{id}/archive [post]
func ArchiveSpace(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		actorID := c.GetString(middleware.CtxUserID)
		if err := h.Space.Archive(c.Request.Context(), id, actorID); err != nil {
			hs.Error(c, err)
			return
		}
		hs.NoContent(c)
	}
}

// RestoreSpace godoc
// @Summary      Restore archived space
// @Tags         spaces
// @Security     BearerAuth
// @Param        id  path  string  true  "Space ID"
// @Success      204
// @Failure      400  {object}  object{code=string,message=string}
// @Failure      403  {object}  object{code=string,message=string}
// @Failure      404  {object}  object{code=string,message=string}
// @Router       /api/v1/spaces/{id}/restore [post]
func RestoreSpace(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		actorID := c.GetString(middleware.CtxUserID)
		if err := h.Space.Restore(c.Request.Context(), id, actorID); err != nil {
			hs.Error(c, err)
			return
		}
		hs.NoContent(c)
	}
}

// GetSpaceStatistics godoc
// @Summary      Get space statistics
// @Tags         spaces
// @Produce      json
// @Security     BearerAuth
// @Param        id  path  string  true  "Space ID"
// @Success      200  {object}  object{data=entity.SpaceStatistics}
// @Failure      404  {object}  object{code=string,message=string}
// @Router       /api/v1/spaces/{id}/statistics [get]
func GetSpaceStatistics(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		stats, err := h.Space.GetStatistics(c.Request.Context(), id)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Success(c, stats)
	}
}
