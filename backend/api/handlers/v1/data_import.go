package v1

import (
	"io"

	"github.com/gin-gonic/gin"

	hs "github.com/jira-backend/jiraflow-backend/api/http_status"
	"github.com/jira-backend/jiraflow-backend/api/handlers"
	"github.com/jira-backend/jiraflow-backend/api/middleware"
)

// ImportJira godoc
// @Summary      Import issues from Jira XML export
// @Tags         import
// @Accept       application/xml
// @Produce      json
// @Security     BearerAuth
// @Success      202  {object}  object{data=entity.DataImport}
// @Router       /api/v1/import/jira [post]
func ImportJira(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString(middleware.CtxUserID)
		data, err := io.ReadAll(c.Request.Body)
		if err != nil || len(data) == 0 {
			hs.BadRequest(c, "request body required")
			return
		}
		imp, err := h.DataImport.ImportJira(c.Request.Context(), userID, data)
		if err != nil {
			hs.Error(c, err)
			return
		}
		c.JSON(202, gin.H{"data": imp})
	}
}

// ImportTrello godoc
// @Summary      Import issues from Trello JSON export
// @Tags         import
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      202  {object}  object{data=entity.DataImport}
// @Router       /api/v1/import/trello [post]
func ImportTrello(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString(middleware.CtxUserID)
		data, err := io.ReadAll(c.Request.Body)
		if err != nil || len(data) == 0 {
			hs.BadRequest(c, "request body required")
			return
		}
		imp, err := h.DataImport.ImportTrello(c.Request.Context(), userID, data)
		if err != nil {
			hs.Error(c, err)
			return
		}
		c.JSON(202, gin.H{"data": imp})
	}
}

// ImportLinear godoc
// @Summary      Import issues from Linear CSV export
// @Tags         import
// @Accept       text/csv
// @Produce      json
// @Security     BearerAuth
// @Success      202  {object}  object{data=entity.DataImport}
// @Router       /api/v1/import/linear [post]
func ImportLinear(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString(middleware.CtxUserID)
		data, err := io.ReadAll(c.Request.Body)
		if err != nil || len(data) == 0 {
			hs.BadRequest(c, "request body required")
			return
		}
		imp, err := h.DataImport.ImportLinear(c.Request.Context(), userID, data)
		if err != nil {
			hs.Error(c, err)
			return
		}
		c.JSON(202, gin.H{"data": imp})
	}
}

// GetImportStatus godoc
// @Summary      Get the status of an import job
// @Tags         import
// @Security     BearerAuth
// @Param        id  path  string  true  "Import job ID"
// @Success      200  {object}  object{data=entity.DataImport}
// @Router       /api/v1/import/{id} [get]
func GetImportStatus(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		imp, err := h.DataImport.GetStatus(c.Request.Context(), c.Param("id"))
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Success(c, imp)
	}
}
