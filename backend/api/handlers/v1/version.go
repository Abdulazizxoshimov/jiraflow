package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"

	hs "github.com/jira-backend/jiraflow-backend/api/http_status"
	"github.com/jira-backend/jiraflow-backend/api/handlers"
	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

func CreateVersion(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		projectID := c.Param("id")
		var req entity.CreateVersionReq
		if err := c.ShouldBindJSON(&req); err != nil {
			hs.BadRequest(c, err.Error())
			return
		}
		v, err := h.Version.Create(c.Request.Context(), projectID, &req)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Created(c, v)
	}
}

func ListVersions(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		projectID := c.Param("id")
		versions, err := h.Version.List(c.Request.Context(), projectID)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Success(c, versions)
	}
}

func GetVersion(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		v, err := h.Version.GetByID(c.Request.Context(), id)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Success(c, v)
	}
}

func UpdateVersion(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var req entity.UpdateVersionReq
		if err := c.ShouldBindJSON(&req); err != nil {
			hs.BadRequest(c, err.Error())
			return
		}
		v, err := h.Version.Update(c.Request.Context(), id, &req)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Success(c, v)
	}
}

func ReleaseVersion(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var req entity.ReleaseVersionReq
		if err := c.ShouldBindJSON(&req); err != nil {
			hs.BadRequest(c, err.Error())
			return
		}
		v, err := h.Version.Release(c.Request.Context(), id, &req)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Success(c, v)
	}
}

func ArchiveVersion(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if err := h.Version.Archive(c.Request.Context(), id); err != nil {
			hs.Error(c, err)
			return
		}
		hs.NoContent(c)
	}
}

func DeleteVersion(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if err := h.Version.Delete(c.Request.Context(), id); err != nil {
			hs.Error(c, err)
			return
		}
		hs.NoContent(c)
	}
}

// GetVersionReleaseNotes returns issues fixed in a version grouped by type.
func GetVersionReleaseNotes(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		versionID := c.Param("id")
		ctx := c.Request.Context()

		version, err := h.Version.GetByID(ctx, versionID)
		if err != nil {
			hs.Error(c, err)
			return
		}

		issues, _, err := h.Issue.List(ctx, &entity.IssueFilter{
			Filter:     entity.Filter{Limit: 500},
			ProjectID:  version.ProjectID,
			VersionIDs: []string{versionID},
		})
		if err != nil {
			hs.Error(c, err)
			return
		}

		grouped := map[string][]*entity.Issue{
			"bug":         {},
			"improvement": {},
			"new_feature": {},
			"other":       {},
		}
		for _, issue := range issues {
			switch issue.Type {
			case "bug":
				grouped["bug"] = append(grouped["bug"], issue)
			case "story":
				grouped["new_feature"] = append(grouped["new_feature"], issue)
			case "task":
				grouped["improvement"] = append(grouped["improvement"], issue)
			default:
				grouped["other"] = append(grouped["other"], issue)
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"data": gin.H{
				"version": version,
				"notes":   grouped,
			},
		})
	}
}
