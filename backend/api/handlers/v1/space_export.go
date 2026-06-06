package v1

import (
	"github.com/gin-gonic/gin"

	hs "github.com/jira-backend/jiraflow-backend/api/http_status"
	"github.com/jira-backend/jiraflow-backend/api/handlers"
	"github.com/jira-backend/jiraflow-backend/api/middleware"
)

func RequestSpaceExport(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		spaceID := c.Param("id")
		actorID := c.GetString(middleware.CtxUserID)
		export, err := h.SpaceExport.RequestExport(c.Request.Context(), spaceID, actorID)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Created(c, export)
	}
}

func GetSpaceExport(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		export, err := h.SpaceExport.GetExport(c.Request.Context(), c.Param("export_id"))
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Success(c, export)
	}
}

func ListSpaceExports(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		exports, err := h.SpaceExport.ListExports(c.Request.Context(), c.Param("id"))
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Success(c, exports)
	}
}
