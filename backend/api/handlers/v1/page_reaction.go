package v1

import (
	"github.com/gin-gonic/gin"

	hs "github.com/jira-backend/jiraflow-backend/api/http_status"
	"github.com/jira-backend/jiraflow-backend/api/handlers"
	"github.com/jira-backend/jiraflow-backend/api/middleware"
	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

func TogglePageReaction(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString(middleware.CtxUserID)
		var req entity.ToggleReactionReq
		if err := c.ShouldBindJSON(&req); err != nil {
			hs.BadRequest(c, err.Error())
			return
		}
		added, err := h.PageReaction.Toggle(c.Request.Context(), c.Param("id"), userID, req.Emoji)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Success(c, gin.H{"added": added, "emoji": req.Emoji})
	}
}

func ListPageReactions(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		viewerID := c.GetString(middleware.CtxUserID)
		summaries, err := h.PageReaction.ListByPage(c.Request.Context(), c.Param("id"), viewerID)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Success(c, summaries)
	}
}

func ListPageReactionUsers(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		emoji := c.Query("emoji")
		if emoji == "" {
			hs.BadRequest(c, "emoji query param is required")
			return
		}
		users, err := h.PageReaction.ListUsers(c.Request.Context(), c.Param("id"), emoji)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Success(c, users)
	}
}
