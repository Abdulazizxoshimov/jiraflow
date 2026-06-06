package v1

import (
	"github.com/gin-gonic/gin"

	hs "github.com/jira-backend/jiraflow-backend/api/http_status"
	"github.com/jira-backend/jiraflow-backend/api/handlers"
	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

func CreateSpaceCategory(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req entity.CreateSpaceCategoryReq
		if err := c.ShouldBindJSON(&req); err != nil {
			hs.BadRequest(c, err.Error())
			return
		}
		cat, err := h.SpaceCategory.Create(c.Request.Context(), &req)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Created(c, cat)
	}
}

func ListSpaceCategories(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		cats, err := h.SpaceCategory.List(c.Request.Context())
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Success(c, cats)
	}
}

func GetSpaceCategory(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		cat, err := h.SpaceCategory.GetByID(c.Request.Context(), c.Param("id"))
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Success(c, cat)
	}
}

func UpdateSpaceCategory(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req entity.UpdateSpaceCategoryReq
		if err := c.ShouldBindJSON(&req); err != nil {
			hs.BadRequest(c, err.Error())
			return
		}
		cat, err := h.SpaceCategory.Update(c.Request.Context(), c.Param("id"), &req)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Success(c, cat)
	}
}

func DeleteSpaceCategory(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := h.SpaceCategory.Delete(c.Request.Context(), c.Param("id")); err != nil {
			hs.Error(c, err)
			return
		}
		hs.NoContent(c)
	}
}
