package v1

import (
	"strconv"

	"github.com/gin-gonic/gin"

	hs "github.com/jira-backend/jiraflow-backend/api/http_status"
	"github.com/jira-backend/jiraflow-backend/api/handlers"
	"github.com/jira-backend/jiraflow-backend/api/middleware"
)

func RecordPageView(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		pageID := c.Param("id")
		userID := c.GetString(middleware.CtxUserID)
		ip := c.ClientIP()

		var uid *string
		if userID != "" {
			uid = &userID
		}
		var ipAddr *string
		if ip != "" && ip != "::1" {
			ipAddr = &ip
		}

		_ = h.PageView.RecordView(c.Request.Context(), pageID, uid, ipAddr)
		hs.NoContent(c)
	}
}

func GetPageAnalytics(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		analytics, err := h.PageView.GetAnalytics(c.Request.Context(), c.Param("id"))
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Success(c, analytics)
	}
}

func ListRecentPages(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString(middleware.CtxUserID)
		limit := 20
		if l := c.Query("limit"); l != "" {
			if n, err := strconv.Atoi(l); err == nil && n > 0 && n <= 50 {
				limit = n
			}
		}
		pages, err := h.PageView.ListRecentByUser(c.Request.Context(), userID, limit)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Success(c, pages)
	}
}
