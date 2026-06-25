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

// VerifyTelegramCode godoc
// @Summary      Verify a code received from the Telegram bot to link the account
// @Description  Flow: open bot → /start → bot sends code → enter code here
// @Tags         telegram
// @Accept       json
// @Security     BearerAuth
// @Param        body  body  object{code=string}  true  "6-digit code from the bot"
// @Success      200   {object}  object{message=string}
// @Router       /api/v1/auth/telegram/verify [post]
func VerifyTelegramCode(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Code string `json:"code" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "code is required"})
			return
		}
		userID := c.GetString(middleware.CtxUserID)
		if err := h.Telegram.VerifyCode(c.Request.Context(), userID, req.Code); err != nil {
			hs.Error(c, err)
			return
		}
		hs.Success(c, gin.H{"message": "Telegram account linked successfully"})
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
// @Router       /telegram/webhook [post]
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

// SetupTelegramWebhook godoc
// @Summary      Register Telegram webhook with Telegram servers (admin only)
// @Tags         telegram
// @Security     BearerAuth
// @Success      200  {object}  object{message=string}
// @Router       /api/v1/admin/telegram/webhook [post]
func SetupTelegramWebhook(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := h.Telegram.SetupWebhook(c.Request.Context()); err != nil {
			hs.Error(c, err)
			return
		}
		hs.Success(c, gin.H{"message": "webhook registered"})
	}
}

// DeleteTelegramWebhook godoc
// @Summary      Unregister Telegram webhook
// @Tags         telegram
// @Security     BearerAuth
// @Success      200  {object}  object{message=string}
// @Router       /api/v1/admin/telegram/webhook [delete]
func DeleteTelegramWebhook(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := h.Telegram.DeleteWebhook(c.Request.Context()); err != nil {
			hs.Error(c, err)
			return
		}
		hs.Success(c, gin.H{"message": "webhook removed"})
	}
}

// GetTelegramBotInfo godoc
// @Summary      Get Telegram bot configuration info
// @Tags         telegram
// @Security     BearerAuth
// @Success      200  {object}  entity.TelegramBotInfo
// @Router       /api/v1/admin/telegram/info [get]
func GetTelegramBotInfo(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		hs.Success(c, h.Telegram.BotInfo(c.Request.Context()))
	}
}
