package v1

import (
	"github.com/gin-gonic/gin"

	hs "github.com/jira-backend/jiraflow-backend/api/http_status"
	"github.com/jira-backend/jiraflow-backend/api/handlers"
	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

func CreateComponent(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		projectID := c.Param("id")
		var req entity.CreateComponentReq
		if err := c.ShouldBindJSON(&req); err != nil {
			hs.BadRequest(c, err.Error())
			return
		}
		comp, err := h.Component.Create(c.Request.Context(), projectID, &req)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Created(c, comp)
	}
}

func ListComponents(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		projectID := c.Param("id")
		comps, err := h.Component.List(c.Request.Context(), projectID)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Success(c, comps)
	}
}

func GetComponent(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		comp, err := h.Component.GetByID(c.Request.Context(), id)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Success(c, comp)
	}
}

func UpdateComponent(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var req entity.UpdateComponentReq
		if err := c.ShouldBindJSON(&req); err != nil {
			hs.BadRequest(c, err.Error())
			return
		}
		comp, err := h.Component.Update(c.Request.Context(), id, &req)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Success(c, comp)
	}
}

func DeleteComponent(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if err := h.Component.Delete(c.Request.Context(), id); err != nil {
			hs.Error(c, err)
			return
		}
		hs.NoContent(c)
	}
}
