package v1

import (
	"strconv"

	"github.com/gin-gonic/gin"

	hs "github.com/jira-backend/jiraflow-backend/api/http_status"
	"github.com/jira-backend/jiraflow-backend/api/handlers"
	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

// SearchSuggestions godoc
// @Summary      Autocomplete suggestions
// @Tags         search
// @Produce      json
// @Security     BearerAuth
// @Param        q      query  string  true   "Partial query (min 2 chars)"
// @Param        limit  query  int     false  "Max results (default 10)"
// @Success      200  {object}  object{data=[]entity.SearchSuggestion}
// @Router       /api/v1/search/suggestions [get]
func SearchSuggestions(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		q := c.Query("q")
		limit := 10
		if l := c.Query("limit"); l != "" {
			if n, err := strconv.Atoi(l); err == nil && n > 0 {
				limit = n
			}
		}
		suggestions, err := h.Search.Suggest(c.Request.Context(), q, limit)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Success(c, suggestions)
	}
}

// Search godoc
// @Summary      Full-text search across issues and pages
// @Tags         search
// @Produce      json
// @Security     BearerAuth
// @Param        q           query  string  true   "Search query"
// @Param        project_id  query  string  false  "Limit to project"
// @Param        type        query  string  false  "Entity type (issue|page)"
// @Param        page        query  int     false  "Page number"
// @Param        limit       query  int     false  "Page size"
// @Success      200  {object}  object{data=[]entity.SearchResult,total=int,page=int,limit=int,total_pages=int}
// @Router       /api/v1/search [get]
func Search(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		var filter entity.SearchFilter
		if err := c.ShouldBindQuery(&filter); err != nil {
			hs.BadRequest(c, err.Error())
			return
		}
		results, total, err := h.Search.Search(c.Request.Context(), &filter)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.List(c, results, total, filter.Page, filter.Limit)
	}
}
