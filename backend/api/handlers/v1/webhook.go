package v1

import (
	"strconv"

	"github.com/gin-gonic/gin"

	hs "github.com/jira-backend/jiraflow-backend/api/http_status"
	"github.com/jira-backend/jiraflow-backend/api/handlers"
	"github.com/jira-backend/jiraflow-backend/api/middleware"
	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

func CreateWebhook(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req entity.CreateWebhookReq
		if err := c.ShouldBindJSON(&req); err != nil {
			hs.BadRequest(c, err.Error())
			return
		}
		actorID := c.GetString(middleware.CtxUserID)
		wh, err := h.Webhook.Create(c.Request.Context(), actorID, &req)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Created(c, wh)
	}
}

func GetWebhook(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		wh, err := h.Webhook.GetByID(c.Request.Context(), c.Param("id"))
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Success(c, wh)
	}
}

func ListWebhooksByProject(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		whs, err := h.Webhook.ListByProject(c.Request.Context(), c.Param("id"))
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Success(c, whs)
	}
}

func ListWebhooksBySpace(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		whs, err := h.Webhook.ListBySpace(c.Request.Context(), c.Param("id"))
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Success(c, whs)
	}
}

func UpdateWebhook(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req entity.UpdateWebhookReq
		if err := c.ShouldBindJSON(&req); err != nil {
			hs.BadRequest(c, err.Error())
			return
		}
		wh, err := h.Webhook.Update(c.Request.Context(), c.Param("id"), &req)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Success(c, wh)
	}
}

func DeleteWebhook(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := h.Webhook.Delete(c.Request.Context(), c.Param("id")); err != nil {
			hs.Error(c, err)
			return
		}
		hs.NoContent(c)
	}
}

func ListWebhookDeliveries(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "25"))
		deliveries, err := h.Webhook.ListDeliveries(c.Request.Context(), c.Param("id"), limit)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Success(c, deliveries)
	}
}
