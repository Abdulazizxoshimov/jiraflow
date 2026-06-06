package v1

import (
	"github.com/gin-gonic/gin"

	hs "github.com/jira-backend/jiraflow-backend/api/http_status"
	"github.com/jira-backend/jiraflow-backend/api/handlers"
	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

func ListContentProperties(entityType string) func(h *handlers.Handler) gin.HandlerFunc {
	return func(h *handlers.Handler) gin.HandlerFunc {
		return func(c *gin.Context) {
			props, err := h.ContentProperty.List(c.Request.Context(), entityType, c.Param("id"))
			if err != nil {
				hs.Error(c, err)
				return
			}
			hs.Success(c, props)
		}
	}
}

func GetContentProperty(entityType string) func(h *handlers.Handler) gin.HandlerFunc {
	return func(h *handlers.Handler) gin.HandlerFunc {
		return func(c *gin.Context) {
			prop, err := h.ContentProperty.Get(c.Request.Context(), entityType, c.Param("id"), c.Param("key"))
			if err != nil {
				hs.Error(c, err)
				return
			}
			hs.Success(c, prop)
		}
	}
}

func SetContentProperty(entityType string) func(h *handlers.Handler) gin.HandlerFunc {
	return func(h *handlers.Handler) gin.HandlerFunc {
		return func(c *gin.Context) {
			var req entity.SetContentPropertyReq
			if err := c.ShouldBindJSON(&req); err != nil {
				hs.BadRequest(c, err.Error())
				return
			}
			prop, err := h.ContentProperty.Set(c.Request.Context(), entityType, c.Param("id"), c.Param("key"), req.Value)
			if err != nil {
				hs.Error(c, err)
				return
			}
			hs.Success(c, prop)
		}
	}
}

func DeleteContentProperty(entityType string) func(h *handlers.Handler) gin.HandlerFunc {
	return func(h *handlers.Handler) gin.HandlerFunc {
		return func(c *gin.Context) {
			if err := h.ContentProperty.Delete(c.Request.Context(), entityType, c.Param("id"), c.Param("key")); err != nil {
				hs.Error(c, err)
				return
			}
			hs.NoContent(c)
		}
	}
}
