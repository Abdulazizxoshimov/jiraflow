package v1

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	hs "github.com/jira-backend/jiraflow-backend/api/http_status"
	"github.com/jira-backend/jiraflow-backend/api/handlers"
	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

// GetPageVersion godoc
// @Summary      Get page version by ID
// @Tags         page-versions
// @Produce      json
// @Security     BearerAuth
// @Param        id  path  string  true  "Page version ID"
// @Success      200  {object}  object{data=entity.PageVersion}
// @Failure      404  {object}  object{code=string,message=string}
// @Router       /api/v1/page-versions/{id} [get]
func GetPageVersion(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		version, err := h.PageVersion.GetByID(c.Request.Context(), id)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Success(c, version)
	}
}

// GetPageVersionByNumber godoc
// @Summary      Get specific version of a page
// @Tags         page-versions
// @Produce      json
// @Security     BearerAuth
// @Param        page_id  path  string  true  "Page ID"
// @Param        version  path  int     true  "Version number"
// @Success      200  {object}  object{data=entity.PageVersion}
// @Failure      404  {object}  object{code=string,message=string}
// @Router       /api/v1/pages/{page_id}/versions/{version} [get]
func GetPageVersionByNumber(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		pageID := c.Param("id")
		versionStr := c.Param("version")
		versionNum, err := strconv.Atoi(versionStr)
		if err != nil {
			hs.BadRequest(c, "invalid version number")
			return
		}
		version, err := h.PageVersion.GetByVersion(c.Request.Context(), pageID, versionNum)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Success(c, version)
	}
}

// ListPageVersions godoc
// @Summary      List all versions of a page
// @Tags         page-versions
// @Produce      json
// @Security     BearerAuth
// @Param        page_id  path   string  true   "Page ID"
// @Param        page     query  int     false  "Page number"
// @Param        limit    query  int     false  "Page size"
// @Success      200  {object}  object{data=[]entity.PageVersion,total=int,page=int,limit=int,total_pages=int}
// @Router       /api/v1/pages/{page_id}/versions [get]
func ListPageVersions(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		pageID := c.Param("id")
		var filter entity.Filter
		if err := c.ShouldBindQuery(&filter); err != nil {
			hs.BadRequest(c, err.Error())
			return
		}
		versions, total, err := h.PageVersion.ListByPage(c.Request.Context(), pageID, &filter)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.List(c, versions, total, filter.Page, filter.GetLimit())
	}
}

// DiffPageVersions returns a line-level diff between two versions of a page.
func DiffPageVersions(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		pageID := c.Param("id")
		v1Str := c.Param("version")
		v2Str := c.Param("v2")

		v1, err := strconv.Atoi(v1Str)
		if err != nil {
			hs.BadRequest(c, "v1 must be an integer")
			return
		}
		v2, err2 := strconv.Atoi(v2Str)
		if err2 != nil {
			hs.BadRequest(c, "v2 must be an integer")
			return
		}

		diff, err := h.PageVersion.Diff(c.Request.Context(), pageID, v1, v2)
		if err != nil {
			hs.Error(c, err)
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": diff})
	}
}
