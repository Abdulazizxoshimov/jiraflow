package v1

import (
	"github.com/gin-gonic/gin"

	hs "github.com/jira-backend/jiraflow-backend/api/http_status"
	"github.com/jira-backend/jiraflow-backend/api/handlers"
	"github.com/jira-backend/jiraflow-backend/api/middleware"
	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

func AddFavorite(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString(middleware.CtxUserID)
		var req entity.AddFavoriteReq
		if err := c.ShouldBindJSON(&req); err != nil {
			hs.BadRequest(c, err.Error())
			return
		}
		fav, err := h.Favorite.Add(c.Request.Context(), userID, &req)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Created(c, fav)
	}
}

func RemoveFavorite(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString(middleware.CtxUserID)
		entityType := c.Query("entity_type")
		entityID := c.Query("entity_id")
		if entityType == "" || entityID == "" {
			hs.BadRequest(c, "entity_type and entity_id are required")
			return
		}
		if err := h.Favorite.Remove(c.Request.Context(), userID, entityType, entityID); err != nil {
			hs.Error(c, err)
			return
		}
		hs.NoContent(c)
	}
}

func ListFavorites(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString(middleware.CtxUserID)
		var filter entity.FavoriteFilter
		if err := c.ShouldBindQuery(&filter); err != nil {
			hs.BadRequest(c, err.Error())
			return
		}
		filter.EntityType = c.Query("entity_type")
		favs, total, err := h.Favorite.List(c.Request.Context(), userID, &filter)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.List(c, favs, total, filter.Page, filter.GetLimit())
	}
}

func IsFavorite(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString(middleware.CtxUserID)
		entityType := c.Query("entity_type")
		entityID := c.Query("entity_id")
		ok, err := h.Favorite.IsFavorite(c.Request.Context(), userID, entityType, entityID)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Success(c, gin.H{"is_favorite": ok})
	}
}
