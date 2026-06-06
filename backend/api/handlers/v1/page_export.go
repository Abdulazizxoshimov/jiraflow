package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"

	hs "github.com/jira-backend/jiraflow-backend/api/http_status"
	"github.com/jira-backend/jiraflow-backend/api/handlers"
)

func ExportPageHTML(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		content, filename, err := h.PageExport.ExportHTML(c.Request.Context(), c.Param("id"))
		if err != nil {
			hs.Error(c, err)
			return
		}
		c.Header("Content-Disposition", "attachment; filename=\""+filename+"\"")
		c.Header("Content-Type", "text/html; charset=UTF-8")
		c.Data(http.StatusOK, "text/html; charset=UTF-8", content)
	}
}

func ExportPagePDF(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		content, filename, err := h.PageExport.ExportPDF(c.Request.Context(), c.Param("id"))
		if err != nil {
			hs.Error(c, err)
			return
		}
		c.Header("Content-Disposition", "attachment; filename=\""+filename+"\"")
		c.Data(http.StatusOK, "application/pdf", content)
	}
}

func ExportPageMarkdown(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		content, filename, err := h.PageExport.ExportMarkdown(c.Request.Context(), c.Param("id"))
		if err != nil {
			hs.Error(c, err)
			return
		}
		c.Header("Content-Disposition", "attachment; filename=\""+filename+"\"")
		c.Data(http.StatusOK, "text/markdown; charset=UTF-8", content)
	}
}

func ExportPageDOCX(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		content, filename, err := h.PageExport.ExportDOCX(c.Request.Context(), c.Param("id"))
		if err != nil {
			hs.Error(c, err)
			return
		}
		c.Header("Content-Disposition", "attachment; filename=\""+filename+"\"")
		c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.wordprocessingml.document", content)
	}
}
