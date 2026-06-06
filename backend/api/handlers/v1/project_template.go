package v1

import (
	"github.com/gin-gonic/gin"

	hs "github.com/jira-backend/jiraflow-backend/api/http_status"
	"github.com/jira-backend/jiraflow-backend/api/handlers"
)

func ListProjectTemplates(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		list, err := h.ProjectTemplate.List(c.Request.Context())
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Success(c, list)
	}
}

func GetProjectTemplate(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		t, err := h.ProjectTemplate.GetByID(c.Request.Context(), c.Param("id"))
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Success(c, t)
	}
}
