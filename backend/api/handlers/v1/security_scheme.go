package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jira-backend/jiraflow-backend/api/handlers"
	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

func CreateSecurityScheme(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req entity.CreateSecuritySchemeReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		scheme, err := h.SecurityScheme.Create(c.Request.Context(), &req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, scheme)
	}
}

func ListSecuritySchemes(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		projectID := c.Query("project_id")
		schemes, err := h.SecurityScheme.List(c.Request.Context(), projectID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"schemes": schemes})
	}
}

func GetSecurityScheme(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		scheme, err := h.SecurityScheme.GetByID(c.Request.Context(), c.Param("id"))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, scheme)
	}
}

func DeleteSecurityScheme(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := h.SecurityScheme.Delete(c.Request.Context(), c.Param("id")); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.Status(http.StatusNoContent)
	}
}

func AddSecurityLevel(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req entity.CreateSecurityLevelReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		level, err := h.SecurityScheme.AddLevel(c.Request.Context(), c.Param("id"), &req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, level)
	}
}

func GetSecurityLevel(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		level, err := h.SecurityScheme.GetLevel(c.Request.Context(), c.Param("level_id"))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, level)
	}
}

func DeleteSecurityLevel(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := h.SecurityScheme.DeleteLevel(c.Request.Context(), c.Param("level_id")); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.Status(http.StatusNoContent)
	}
}

func AddSecurityLevelMember(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req entity.CreateSecurityLevelMemberReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		member, err := h.SecurityScheme.AddMember(c.Request.Context(), c.Param("level_id"), &req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, member)
	}
}

func DeleteSecurityLevelMember(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := h.SecurityScheme.DeleteMember(c.Request.Context(), c.Param("member_id")); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.Status(http.StatusNoContent)
	}
}
