package v1

import (
	"encoding/csv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/jira-backend/jiraflow-backend/api/handlers"
	hs "github.com/jira-backend/jiraflow-backend/api/http_status"
	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

// ListAuditLogs godoc
// @Summary      List audit logs
// @Tags         audit
// @Produce      json
// @Security     BearerAuth
// @Param        actor_id     query  string  false  "Filter by actor"
// @Param        entity_type  query  string  false  "Filter by entity type"
// @Param        entity_id    query  string  false  "Filter by entity ID"
// @Param        action       query  string  false  "Filter by action"
// @Param        page         query  int     false  "Page number"
// @Param        limit        query  int     false  "Page size"
// @Success      200  {object}  object{data=[]entity.AuditLog,total=int,page=int,limit=int,total_pages=int}
// @Router       /api/v1/audit-logs [get]
func ListAuditLogs(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		var filter entity.AuditLogFilter
		if err := c.ShouldBindQuery(&filter); err != nil {
			hs.BadRequest(c, err.Error())
			return
		}
		logs, total, err := h.Audit.List(c.Request.Context(), &filter)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.List(c, logs, total, filter.Page, filter.GetLimit())
	}
}

// ExportAuditLogs godoc
// @Summary      Export audit logs as CSV
// @Tags         audit
// @Produce      text/csv
// @Security     BearerAuth
// @Param        user_id      query  string  false  "Filter by user"
// @Param        entity_type  query  string  false  "Filter by entity type"
// @Param        action       query  string  false  "Filter by action"
// @Success      200
// @Router       /api/v1/audit-logs/export [get]
func ExportAuditLogs(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		var filter entity.AuditLogFilter
		if err := c.ShouldBindQuery(&filter); err != nil {
			hs.BadRequest(c, err.Error())
			return
		}

		c.Header("Content-Type", "text/csv; charset=utf-8")
		c.Header("Content-Disposition", `attachment; filename="audit_logs.csv"`)

		w := csv.NewWriter(c.Writer)
		_ = w.Write([]string{"id", "user_id", "action", "entity_type", "entity_id", "ip_address", "created_at"})

		// Stream in pages of 500 to avoid loading the entire table into memory.
		const pageSize = 500
		filter.Limit = pageSize
		for page := 1; ; page++ {
			filter.Page = page
			logs, total, err := h.Audit.List(c.Request.Context(), &filter)
			if err != nil {
				break
			}
			for _, l := range logs {
				_ = w.Write([]string{
					l.ID,
					derefStr(l.UserID),
					l.Action,
					derefStr(l.EntityType),
					derefStr(l.EntityID),
					derefStr(l.IPAddress),
					l.CreatedAt.Format(time.RFC3339),
				})
			}
			w.Flush()
			if page*pageSize >= total || len(logs) < pageSize {
				break
			}
		}
	}
}

func derefStr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
