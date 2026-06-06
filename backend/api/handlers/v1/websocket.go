package v1

import (
	"github.com/gin-gonic/gin"

	"github.com/jira-backend/jiraflow-backend/api/handlers"
	"github.com/jira-backend/jiraflow-backend/api/middleware"
)

// ServeWS godoc
// @Summary      WebSocket connection
// @Description  Upgrades HTTP connection to WebSocket for real-time notifications.
//
//	Client should send JWT token as query param: /api/v1/ws?token=<jwt>
//
// @Tags         websocket
// @Produce      json
// @Security     BearerAuth
// @Success      101  "Switching Protocols"
// @Router       /api/v1/ws [get]
func ServeWS(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString(middleware.CtxUserID)
		h.Hub.ServeWS(c.Request.Context(), c.Writer, c.Request, userID)
	}
}
