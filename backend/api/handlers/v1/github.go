package v1

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"

	hs "github.com/jira-backend/jiraflow-backend/api/http_status"
	"github.com/jira-backend/jiraflow-backend/api/handlers"
	"github.com/jira-backend/jiraflow-backend/api/middleware"
	"github.com/jira-backend/jiraflow-backend/internal/entity"
	ghinfra "github.com/jira-backend/jiraflow-backend/internal/infrastructure/github"
)

// ConnectGitHubRepo godoc
// @Summary      Connect GitHub repository to a project
// @Tags         github
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path  string                  true  "Project ID"
// @Param        body  body  entity.ConnectRepoReq   true  "Repo data"
// @Success      201  {object}  object{data=entity.GitHubRepo}
// @Router       /api/v1/projects/{id}/github [post]
func ConnectGitHubRepo(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		projectID := c.Param("id")
		userID := c.GetString(middleware.CtxUserID)
		var req entity.ConnectRepoReq
		if err := c.ShouldBindJSON(&req); err != nil {
			hs.BadRequest(c, err.Error())
			return
		}
		repo, err := h.GitHub.ConnectRepo(c.Request.Context(), projectID, userID, &req)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Created(c, repo)
	}
}

// DisconnectGitHubRepo godoc
// @Summary      Disconnect GitHub repository from a project
// @Tags         github
// @Security     BearerAuth
// @Param        id  path  string  true  "Project ID"
// @Success      204
// @Router       /api/v1/projects/{id}/github [delete]
func DisconnectGitHubRepo(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		projectID := c.Param("id")
		if err := h.GitHub.DisconnectRepo(c.Request.Context(), projectID); err != nil {
			hs.Error(c, err)
			return
		}
		hs.NoContent(c)
	}
}

// GetGitHubRepo godoc
// @Summary      Get connected GitHub repository for a project
// @Tags         github
// @Produce      json
// @Security     BearerAuth
// @Param        id  path  string  true  "Project ID"
// @Success      200  {object}  object{data=entity.GitHubRepo}
// @Router       /api/v1/projects/{id}/github [get]
func GetGitHubRepo(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		projectID := c.Param("id")
		repo, err := h.GitHub.GetRepo(c.Request.Context(), projectID)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Success(c, repo)
	}
}

// ListIssueCommits godoc
// @Summary      List commits linked to an issue
// @Tags         github
// @Produce      json
// @Security     BearerAuth
// @Param        id  path  string  true  "Issue ID"
// @Success      200  {object}  object{data=[]entity.IssueCommit}
// @Router       /api/v1/issues/{id}/commits [get]
func ListIssueCommits(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		issueID := c.Param("id")
		commits, err := h.GitHub.ListCommits(c.Request.Context(), issueID)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Success(c, commits)
	}
}

// ListIssuePRs godoc
// @Summary      List pull requests linked to an issue
// @Tags         github
// @Produce      json
// @Security     BearerAuth
// @Param        id  path  string  true  "Issue ID"
// @Success      200  {object}  object{data=[]entity.IssuePullRequest}
// @Router       /api/v1/issues/{id}/pull-requests [get]
func ListIssuePRs(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		issueID := c.Param("id")
		prs, err := h.GitHub.ListPRs(c.Request.Context(), issueID)
		if err != nil {
			hs.Error(c, err)
			return
		}
		hs.Success(c, prs)
	}
}

// GitHubWebhook godoc
// @Summary      Receive GitHub webhook events
// @Tags         github
// @Accept       json
// @Success      200
// @Router       /api/v1/github/webhook [post]
func GitHubWebhook(h *handlers.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}

		// Signature verification
		if h.GitHubWebhookSecret != "" {
			sig := c.GetHeader("X-Hub-Signature-256")
			if !ghinfra.VerifyWebhookSignature([]byte(h.GitHubWebhookSecret), body, sig) {
				c.Status(http.StatusUnauthorized)
				return
			}
		}

		event := c.GetHeader("X-GitHub-Event")
		ctx := c.Request.Context()

		switch event {
		case "push":
			ev, err := ghinfra.ParsePushEvent(body)
			if err == nil {
				_ = h.GitHub.HandlePushEvent(ctx, ev.Repository.FullName, ev.Commits)
			}
		case "pull_request":
			ev, err := ghinfra.ParsePREvent(body)
			if err == nil {
				_ = h.GitHub.HandlePREvent(ctx, ev.Repository.FullName, ev)
			}
		}

		c.Status(http.StatusOK)
	}
}
