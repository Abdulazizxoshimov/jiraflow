package v1

import (
	"github.com/gin-gonic/gin"

	hs "github.com/jira-backend/jiraflow-backend/api/http_status"
	"github.com/jira-backend/jiraflow-backend/api/handlers"
	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

func CreatePageTag(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		spaceID := c.Param("id")
		var req entity.CreatePageTagReq
		if err := c.ShouldBindJSON(&req); err != nil {
			hs.BadRequest(c, err.Error())
			return
		}
		tag, err := h.PageTag.Create(c.Request.Context(), spaceID, &req)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Created(c, tag)
	}
}

func ListPageTags(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		tags, err := h.PageTag.List(c.Request.Context(), c.Param("id"))
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Success(c, tags)
	}
}

func GetPageTag(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		tag, err := h.PageTag.GetByID(c.Request.Context(), c.Param("id"))
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Success(c, tag)
	}
}

func UpdatePageTag(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req entity.UpdatePageTagReq
		if err := c.ShouldBindJSON(&req); err != nil {
			hs.BadRequest(c, err.Error())
			return
		}
		tag, err := h.PageTag.Update(c.Request.Context(), c.Param("id"), &req)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Success(c, tag)
	}
}

func DeletePageTag(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := h.PageTag.Delete(c.Request.Context(), c.Param("id")); err != nil {
			hs.Error(c, err)
			return
		}
		hs.NoContent(c)
	}
}

func SetPageTags(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req entity.SetPageTagsReq
		if err := c.ShouldBindJSON(&req); err != nil {
			hs.BadRequest(c, err.Error())
			return
		}
		if err := h.PageTag.SetPageTags(c.Request.Context(), c.Param("id"), &req); err != nil {
			hs.Error(c, err)
			return
		}
		hs.NoContent(c)
	}
}

func GetPageTagsForPage(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		tags, err := h.PageTag.GetPageTags(c.Request.Context(), c.Param("id"))
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Success(c, tags)
	}
}

func GetPagesByTag(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		var filter entity.Filter
		if err := c.ShouldBindQuery(&filter); err != nil {
			hs.BadRequest(c, err.Error())
			return
		}
		pages, total, err := h.PageTag.GetPagesByTag(c.Request.Context(), c.Param("id"), &filter)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.List(c, pages, total, filter.Page, filter.GetLimit())
	}
}
