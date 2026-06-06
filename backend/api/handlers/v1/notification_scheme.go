package v1

import (
	"github.com/gin-gonic/gin"

	hs "github.com/jira-backend/jiraflow-backend/api/http_status"
	"github.com/jira-backend/jiraflow-backend/api/handlers"
	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

func ListNotificationSchemes(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		list, err := h.NotificationScheme.List(c.Request.Context())
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Success(c, list)
	}
}

func GetNotificationScheme(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		s, err := h.NotificationScheme.GetByID(c.Request.Context(), c.Param("id"))
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Success(c, s)
	}
}

func CreateNotificationScheme(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req entity.CreateNotificationSchemeReq
		if err := c.ShouldBindJSON(&req); err != nil {
			hs.BadRequest(c, err.Error())
			return
		}
		s, err := h.NotificationScheme.Create(c.Request.Context(), &req)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Created(c, s)
	}
}

func DeleteNotificationScheme(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := h.NotificationScheme.Delete(c.Request.Context(), c.Param("id")); err != nil {
			hs.Error(c, err)
			return
		}
		hs.NoContent(c)
	}
}
