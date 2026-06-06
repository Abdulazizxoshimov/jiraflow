package v1

import (
	"github.com/gin-gonic/gin"

	hs "github.com/jira-backend/jiraflow-backend/api/http_status"
	"github.com/jira-backend/jiraflow-backend/api/handlers"
	"github.com/jira-backend/jiraflow-backend/api/middleware"
	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

func SetPageRestrictions(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		pageID := c.Param("id")
		actorID := c.GetString(middleware.CtxUserID)
		var req entity.SetPageRestrictionsReq
		if err := c.ShouldBindJSON(&req); err != nil {
			hs.BadRequest(c, err.Error())
			return
		}
		if err := h.PageRestriction.Set(c.Request.Context(), pageID, actorID, &req); err != nil {
			hs.Error(c, err)
			return
		}
		hs.NoContent(c)
	}
}

func ListPageRestrictions(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		restrictions, err := h.PageRestriction.List(c.Request.Context(), c.Param("id"))
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Success(c, restrictions)
	}
}

func ClearPageRestrictions(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		actorID := c.GetString(middleware.CtxUserID)
		if err := h.PageRestriction.Clear(c.Request.Context(), c.Param("id"), actorID); err != nil {
			hs.Error(c, err)
			return
		}
		hs.NoContent(c)
	}
}

func CheckPageAccess(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString(middleware.CtxUserID)
		accessType := c.DefaultQuery("type", "view")
		info, err := h.PageRestriction.CheckAccess(c.Request.Context(), c.Param("id"), userID, accessType)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Success(c, info)
	}
}
