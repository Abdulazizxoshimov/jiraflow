package v1

import (
	"github.com/gin-gonic/gin"

	hs "github.com/jira-backend/jiraflow-backend/api/http_status"
	"github.com/jira-backend/jiraflow-backend/api/handlers"
	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

func ListIssueTypes(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		list, err := h.IssueType.ListTypes(c.Request.Context())
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Success(c, list)
	}
}

func CreateIssueType(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req entity.CreateIssueTypeReq
		if err := c.ShouldBindJSON(&req); err != nil {
			hs.BadRequest(c, err.Error())
			return
		}
		t, err := h.IssueType.CreateType(c.Request.Context(), &req)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Created(c, t)
	}
}

func GetIssueType(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		t, err := h.IssueType.GetTypeByID(c.Request.Context(), c.Param("id"))
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Success(c, t)
	}
}

func DeleteIssueType(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := h.IssueType.DeleteType(c.Request.Context(), c.Param("id")); err != nil {
			hs.Error(c, err)
			return
		}
		hs.NoContent(c)
	}
}

func ListIssueTypeSchemes(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		list, err := h.IssueType.ListSchemes(c.Request.Context())
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Success(c, list)
	}
}

func CreateIssueTypeScheme(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req entity.CreateIssueTypeSchemeReq
		if err := c.ShouldBindJSON(&req); err != nil {
			hs.BadRequest(c, err.Error())
			return
		}
		s, err := h.IssueType.CreateScheme(c.Request.Context(), &req)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Created(c, s)
	}
}

func GetIssueTypeScheme(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		s, err := h.IssueType.GetSchemeByID(c.Request.Context(), c.Param("id"))
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Success(c, s)
	}
}

func DeleteIssueTypeScheme(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := h.IssueType.DeleteScheme(c.Request.Context(), c.Param("id")); err != nil {
			hs.Error(c, err)
			return
		}
		hs.NoContent(c)
	}
}

func GetProjectIssueTypeScheme(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		s, err := h.IssueType.GetSchemeByProject(c.Request.Context(), c.Param("id"))
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Success(c, s)
	}
}
