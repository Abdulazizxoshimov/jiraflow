package v1

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/jira-backend/jiraflow-backend/api/handlers"
	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

// GanttItem represents one row on the Gantt chart.
type GanttItem struct {
	ID           string     `json:"id"`
	Title        string     `json:"title"`
	Type         string     `json:"type"` // "epic" | "version"
	StartDate    *time.Time `json:"start_date,omitempty"`
	DueDate      *time.Time `json:"due_date,omitempty"`
	Progress     int        `json:"progress"` // 0–100
	Dependencies []string   `json:"dependencies"`
}

// GetGanttData returns epics and versions for a project as Gantt items.
func GetGanttData(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		projectID := c.Param("id")
		ctx := c.Request.Context()

		// Epics
		epics, _, _ := h.Issue.List(ctx, &entity.IssueFilter{
			Filter:    entity.Filter{Limit: 200},
			ProjectID: projectID,
			Type:      "epic",
		})

		// Versions
		versions, _ := h.Version.List(ctx, projectID)

		items := make([]GanttItem, 0, len(epics)+len(versions))

		for _, e := range epics {
			progress := 0
			if e.EpicProgress != nil {
				progress = int(e.EpicProgress.Progress)
			}
			items = append(items, GanttItem{
				ID:           e.ID,
				Title:        e.Title,
				Type:         "epic",
				DueDate:      e.DueDate,
				Progress:     progress,
				Dependencies: []string{},
			})
		}

		for _, v := range versions {
			items = append(items, GanttItem{
				ID:           v.ID,
				Title:        v.Name,
				Type:         "version",
				StartDate:    v.StartDate,
				DueDate:      v.ReleaseDate,
				Progress:     0,
				Dependencies: []string{},
			})
		}

		c.JSON(http.StatusOK, gin.H{"data": gin.H{"items": items}})
	}
}
