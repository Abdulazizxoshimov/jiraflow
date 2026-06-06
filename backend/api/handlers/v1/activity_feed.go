package v1

import (
	"github.com/gin-gonic/gin"

	hs "github.com/jira-backend/jiraflow-backend/api/http_status"
	"github.com/jira-backend/jiraflow-backend/api/handlers"
	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

// ListActivity godoc
// @Summary      Unified activity feed
// @Tags         activity
// @Produce      json
// @Security     BearerAuth
// @Param        actor_id     query  string  false  "Filter by actor"
// @Param        project_id   query  string  false  "Filter by project"
// @Param        space_id     query  string  false  "Filter by space"
// @Param        entity_type  query  string  false  "issue|page|comment|sprint|space|project"
// @Param        page         query  int     false  "Page number"
// @Param        limit        query  int     false  "Page size"
// @Success      200  {object}  object{data=[]entity.ActivityEvent,total=int}
// @Router       /api/v1/activity [get]
func ListActivity(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		var filter entity.ActivityFilter
		if err := c.ShouldBindQuery(&filter); err != nil {
			hs.BadRequest(c, err.Error())
			return
		}
		events, total, err := h.ActivityFeed.List(c.Request.Context(), &filter)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.List(c, events, total, filter.Page, filter.GetLimit())
	}
}

// GetProjectLinkedSpace godoc
// @Summary      Get Confluence space linked to a project
// @Tags         projects
// @Produce      json
// @Security     BearerAuth
// @Param        id  path  string  true  "Project ID"
// @Success      200  {object}  object{data=entity.Space}
// @Failure      404  {object}  object{code=string,message=string}
// @Router       /api/v1/projects/{id}/space [get]
func GetProjectLinkedSpace(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		projectID := c.Param("id")
		space, err := h.Project.GetLinkedSpace(c.Request.Context(), projectID)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Success(c, space)
	}
}
