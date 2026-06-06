package v1

import (
	"github.com/gin-gonic/gin"

	hs "github.com/jira-backend/jiraflow-backend/api/http_status"
	"github.com/jira-backend/jiraflow-backend/api/handlers"
	"github.com/jira-backend/jiraflow-backend/api/middleware"
	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

// ListNotifications godoc
// @Summary      List notifications for current user
// @Tags         notifications
// @Produce      json
// @Security     BearerAuth
// @Param        is_read  query  bool  false  "Filter by read status"
// @Param        page     query  int   false  "Page number"
// @Param        limit    query  int   false  "Page size"
// @Success      200  {object}  object{data=[]entity.Notification,total=int,page=int,limit=int,total_pages=int}
// @Router       /api/v1/notifications [get]
func ListNotifications(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString(middleware.CtxUserID)
		var filter entity.NotificationFilter
		if err := c.ShouldBindQuery(&filter); err != nil {
			hs.BadRequest(c, err.Error())
			return
		}
		notifs, total, err := h.Notification.ListByUser(c.Request.Context(), userID, &filter)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.List(c, notifs, total, filter.Page, filter.GetLimit())
	}
}

// CountUnreadNotifications godoc
// @Summary      Count unread notifications
// @Tags         notifications
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  object{data=object{count=int}}
// @Router       /api/v1/notifications/unread-count [get]
func CountUnreadNotifications(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString(middleware.CtxUserID)
		count, err := h.Notification.CountUnread(c.Request.Context(), userID)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Success(c, gin.H{"count": count})
	}
}

// MarkNotificationsRead godoc
// @Summary      Mark specific notifications as read
// @Tags         notifications
// @Accept       json
// @Security     BearerAuth
// @Param        body  body  entity.MarkReadReq  true  "Notification IDs"
// @Success      204
// @Router       /api/v1/notifications/mark-read [post]
func MarkNotificationsRead(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString(middleware.CtxUserID)
		var req entity.MarkReadReq
		if err := c.ShouldBindJSON(&req); err != nil {
			hs.BadRequest(c, err.Error())
			return
		}
		if err := h.Notification.MarkRead(c.Request.Context(), userID, req.IDs); err != nil {
			hs.Error(c, err)
			return
		}
		hs.NoContent(c)
	}
}

// MarkAllNotificationsRead godoc
// @Summary      Mark all notifications as read
// @Tags         notifications
// @Security     BearerAuth
// @Success      204
// @Router       /api/v1/notifications/mark-all-read [post]
func MarkAllNotificationsRead(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString(middleware.CtxUserID)
		if err := h.Notification.MarkAllRead(c.Request.Context(), userID); err != nil {
			hs.Error(c, err)
			return
		}
		hs.NoContent(c)
	}
}

// DeleteNotification godoc
// @Summary      Delete notification
// @Tags         notifications
// @Security     BearerAuth
// @Param        id  path  string  true  "Notification ID"
// @Success      204
// @Failure      404  {object}  object{code=string,message=string}
// @Router       /api/v1/notifications/{id} [delete]
func DeleteNotification(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		userID := c.GetString(middleware.CtxUserID)
		if err := h.Notification.Delete(c.Request.Context(), id, userID); err != nil {
			hs.Error(c, err)
			return
		}
		hs.NoContent(c)
	}
}

// GetNotificationPreference godoc
// @Summary      Get notification preferences
// @Tags         notifications
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  object{data=entity.NotificationPreference}
// @Router       /api/v1/notifications/preferences [get]
func GetNotificationPreference(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString(middleware.CtxUserID)
		pref, err := h.Notification.GetPreference(c.Request.Context(), userID)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Success(c, pref)
	}
}

// UpdateNotificationPreference godoc
// @Summary      Update notification preferences
// @Tags         notifications
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body  entity.UpdateNotificationPreferenceReq  true  "Preferences"
// @Success      200  {object}  object{data=entity.NotificationPreference}
// @Router       /api/v1/notifications/preferences [put]
func UpdateNotificationPreference(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString(middleware.CtxUserID)
		var req entity.UpdateNotificationPreferenceReq
		if err := c.ShouldBindJSON(&req); err != nil {
			hs.BadRequest(c, err.Error())
			return
		}
		pref, err := h.Notification.UpdatePreference(c.Request.Context(), userID, &req)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Success(c, pref)
	}
}
