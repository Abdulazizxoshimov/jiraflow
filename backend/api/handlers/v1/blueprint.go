package v1

import (
	"github.com/gin-gonic/gin"

	hs "github.com/jira-backend/jiraflow-backend/api/http_status"
	"github.com/jira-backend/jiraflow-backend/api/handlers"
	"github.com/jira-backend/jiraflow-backend/api/middleware"
	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

func ListBlueprints(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		list, err := h.Blueprint.List(c.Request.Context())
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Success(c, list)
	}
}

func GetBlueprint(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		b, err := h.Blueprint.GetByID(c.Request.Context(), c.Param("id"))
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Success(c, b)
	}
}

func CreateBlueprint(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req entity.CreateBlueprintReq
		if err := c.ShouldBindJSON(&req); err != nil {
			hs.BadRequest(c, err.Error())
			return
		}
		b, err := h.Blueprint.Create(c.Request.Context(), &req)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Created(c, b)
	}
}

func DeleteBlueprint(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := h.Blueprint.Delete(c.Request.Context(), c.Param("id")); err != nil {
			hs.Error(c, err)
			return
		}
		hs.NoContent(c)
	}
}

func CreatePageFromBlueprint(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		actorID := c.GetString(middleware.CtxUserID)
		var req entity.CreatePageFromBlueprintReq
		if err := c.ShouldBindJSON(&req); err != nil {
			hs.BadRequest(c, err.Error())
			return
		}
		p, err := h.Blueprint.CreatePage(c.Request.Context(), c.Param("id"), actorID, &req)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Created(c, p)
	}
}
