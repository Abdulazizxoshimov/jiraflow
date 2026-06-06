package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"

	hs "github.com/jira-backend/jiraflow-backend/api/http_status"
	"github.com/jira-backend/jiraflow-backend/api/handlers"
	"github.com/jira-backend/jiraflow-backend/api/middleware"
	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

// GetTelegramStatus godoc
// @Summary      Get Telegram connection status
// @Tags         telegram
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  entity.TelegramStatusResp
// @Router       /api/v1/auth/telegram/status [get]
func GetTelegramStatus(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString(middleware.CtxUserID)
		conn, err := h.Telegram.GetStatus(c.Request.Context(), userID)
		if err != nil {
			hs.Error(c, err)
			return
		}
		resp := entity.TelegramStatusResp{Connected: false}
		if conn != nil && conn.VerifiedAt != nil {
			resp.Connected = true
			resp.Username = conn.Username
			resp.VerifiedAt = conn.VerifiedAt
		}
		hs.Success(c, resp)
	}
}

// GenerateTelegramCode godoc
// @Summary      Generate Telegram verification code
// @Tags         telegram
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  object{code=string}
// @Router       /api/v1/auth/telegram/connect [post]
func GenerateTelegramCode(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString(middleware.CtxUserID)
		code, err := h.Telegram.GenerateCode(c.Request.Context(), userID)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Success(c, gin.H{"code": code})
	}
}

// DisconnectTelegram godoc
// @Summary      Disconnect Telegram account
// @Tags         telegram
// @Security     BearerAuth
// @Success      204
// @Router       /api/v1/auth/telegram/disconnect [delete]
func DisconnectTelegram(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString(middleware.CtxUserID)
		if err := h.Telegram.Disconnect(c.Request.Context(), userID); err != nil {
			hs.Error(c, err)
			return
		}
		c.Status(http.StatusNoContent)
	}
}

// TelegramWebhook godoc
// @Summary      Receive Telegram webhook updates
// @Tags         telegram
// @Accept       json
// @Success      200
// @Router       /api/v1/telegram/webhook [post]
func TelegramWebhook(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		if h.TelegramWebhookSecret != "" {
			token := c.GetHeader("X-Telegram-Bot-Api-Secret-Token")
			if token != h.TelegramWebhookSecret {
				c.Status(http.StatusUnauthorized)
				return
			}
		}
		var update entity.TelegramUpdate
		if err := c.ShouldBindJSON(&update); err != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		_ = h.Telegram.HandleUpdate(c.Request.Context(), &update)
		c.Status(http.StatusOK)
	}
}
