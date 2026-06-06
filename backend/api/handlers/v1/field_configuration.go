package v1

import (
	"github.com/gin-gonic/gin"

	hs "github.com/jira-backend/jiraflow-backend/api/http_status"
	"github.com/jira-backend/jiraflow-backend/api/handlers"
	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

func ListFieldConfigurations(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		var projectID *string
		if pid := c.Query("project_id"); pid != "" {
			projectID = &pid
		}
		list, err := h.FieldConfiguration.List(c.Request.Context(), projectID)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Success(c, list)
	}
}

func GetFieldConfiguration(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		fc, err := h.FieldConfiguration.GetByID(c.Request.Context(), c.Param("id"))
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Success(c, fc)
	}
}

func CreateFieldConfiguration(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req entity.CreateFieldConfigurationReq
		if err := c.ShouldBindJSON(&req); err != nil {
			hs.BadRequest(c, err.Error())
			return
		}
		fc, err := h.FieldConfiguration.Create(c.Request.Context(), &req)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Created(c, fc)
	}
}

func DeleteFieldConfiguration(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := h.FieldConfiguration.Delete(c.Request.Context(), c.Param("id")); err != nil {
			hs.Error(c, err)
			return
		}
		hs.NoContent(c)
	}
}
