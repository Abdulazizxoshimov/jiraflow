package v1

import (
	"github.com/gin-gonic/gin"

	hs "github.com/jira-backend/jiraflow-backend/api/http_status"
	"github.com/jira-backend/jiraflow-backend/api/handlers"
	"github.com/jira-backend/jiraflow-backend/api/middleware"
	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

func AcquirePageLock(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req entity.AcquireLockReq
		_ = c.ShouldBindJSON(&req)
		userID := c.GetString(middleware.CtxUserID)
		lock, err := h.PageLock.Acquire(c.Request.Context(), c.Param("id"), userID, req.TTLSeconds)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Success(c, lock)
	}
}

func ReleasePageLock(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString(middleware.CtxUserID)
		if err := h.PageLock.Release(c.Request.Context(), c.Param("id"), userID); err != nil {
			hs.Error(c, err)
			return
		}
		hs.NoContent(c)
	}
}

func GetPageLock(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		lock, err := h.PageLock.Get(c.Request.Context(), c.Param("id"))
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Success(c, lock)
	}
}

func ExtendPageLock(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req entity.AcquireLockReq
		_ = c.ShouldBindJSON(&req)
		userID := c.GetString(middleware.CtxUserID)
		lock, err := h.PageLock.Extend(c.Request.Context(), c.Param("id"), userID, req.TTLSeconds)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Success(c, lock)
	}
}
