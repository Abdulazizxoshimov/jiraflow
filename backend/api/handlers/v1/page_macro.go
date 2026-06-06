package v1

import (
	"github.com/gin-gonic/gin"

	hs "github.com/jira-backend/jiraflow-backend/api/http_status"
	"github.com/jira-backend/jiraflow-backend/api/handlers"
	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

func UpsertPageMacro(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req entity.UpsertPageMacroReq
		if err := c.ShouldBindJSON(&req); err != nil {
			hs.BadRequest(c, err.Error())
			return
		}
		macro, err := h.PageMacro.Upsert(c.Request.Context(), c.Param("id"), &req)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Created(c, macro)
	}
}

func ListPageMacros(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		macros, err := h.PageMacro.ListByPage(c.Request.Context(), c.Param("id"))
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Success(c, macros)
	}
}

func DeletePageMacro(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := h.PageMacro.Delete(c.Request.Context(), c.Param("id")); err != nil {
			hs.Error(c, err)
			return
		}
		hs.NoContent(c)
	}
}
