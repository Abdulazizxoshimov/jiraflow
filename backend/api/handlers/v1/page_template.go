package v1

import (
	"github.com/gin-gonic/gin"

	hs "github.com/jira-backend/jiraflow-backend/api/http_status"
	"github.com/jira-backend/jiraflow-backend/api/handlers"
	"github.com/jira-backend/jiraflow-backend/api/middleware"
	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

func CreatePageTemplate(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		createdBy := c.GetString(middleware.CtxUserID)
		var req entity.CreatePageTemplateReq
		if err := c.ShouldBindJSON(&req); err != nil {
			hs.BadRequest(c, err.Error())
			return
		}

		var spaceID *string
		if sid := c.Query("space_id"); sid != "" {
			spaceID = &sid
		}

		t, err := h.PageTemplate.Create(c.Request.Context(), spaceID, createdBy, &req)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Created(c, t)
	}
}

func ListPageTemplates(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		var filter entity.PageTemplateFilter
		if err := c.ShouldBindQuery(&filter); err != nil {
			hs.BadRequest(c, err.Error())
			return
		}
		filter.SpaceID = c.Query("space_id")
		filter.Category = c.Query("category")

		templates, total, err := h.PageTemplate.List(c.Request.Context(), &filter)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.List(c, templates, total, filter.Page, filter.GetLimit())
	}
}

func GetPageTemplate(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		t, err := h.PageTemplate.GetByID(c.Request.Context(), c.Param("id"))
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Success(c, t)
	}
}

func UpdatePageTemplate(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		actorID := c.GetString(middleware.CtxUserID)
		var req entity.UpdatePageTemplateReq
		if err := c.ShouldBindJSON(&req); err != nil {
			hs.BadRequest(c, err.Error())
			return
		}
		t, err := h.PageTemplate.Update(c.Request.Context(), c.Param("id"), actorID, &req)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Success(c, t)
	}
}

func DeletePageTemplate(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		actorID := c.GetString(middleware.CtxUserID)
		if err := h.PageTemplate.Delete(c.Request.Context(), c.Param("id"), actorID); err != nil {
			hs.Error(c, err)
			return
		}
		hs.NoContent(c)
	}
}
