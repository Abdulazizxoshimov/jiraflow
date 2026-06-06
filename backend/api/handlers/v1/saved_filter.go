package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"

	hs "github.com/jira-backend/jiraflow-backend/api/http_status"
	"github.com/jira-backend/jiraflow-backend/api/handlers"
	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

func CreateSavedFilter(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req entity.CreateSavedFilterReq
		if err := c.ShouldBindJSON(&req); err != nil {
			hs.BadRequest(c, err.Error())
			return
		}
		sf, err := h.SavedFilter.Create(c.Request.Context(), c.GetString("user_id"), &req)
		if err != nil {
			hs.Error(c, err)
			return
		}
		c.JSON(http.StatusCreated, sf)
	}
}

func ListSavedFilters(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		filterType := c.Query("type") // issue | page | ""
		filters, err := h.SavedFilter.List(c.Request.Context(), c.GetString("user_id"), filterType)
		if err != nil {
			hs.Error(c, err)
			return
		}
		c.JSON(http.StatusOK, filters)
	}
}

func GetSavedFilter(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		sf, err := h.SavedFilter.GetByID(c.Request.Context(), c.Param("id"), c.GetString("user_id"))
		if err != nil {
			hs.Error(c, err)
			return
		}
		c.JSON(http.StatusOK, sf)
	}
}

func UpdateSavedFilter(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req entity.UpdateSavedFilterReq
		if err := c.ShouldBindJSON(&req); err != nil {
			hs.BadRequest(c, err.Error())
			return
		}
		sf, err := h.SavedFilter.Update(c.Request.Context(), c.Param("id"), c.GetString("user_id"), &req)
		if err != nil {
			hs.Error(c, err)
			return
		}
		c.JSON(http.StatusOK, sf)
	}
}

func DeleteSavedFilter(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := h.SavedFilter.Delete(c.Request.Context(), c.Param("id"), c.GetString("user_id")); err != nil {
			hs.Error(c, err)
			return
		}
		c.Status(http.StatusNoContent)
	}
}
