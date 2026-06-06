package api

import (
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/jira-backend/jiraflow-backend/api/handlers"
	v1 "github.com/jira-backend/jiraflow-backend/api/handlers/v1"
	"github.com/jira-backend/jiraflow-backend/api/middleware"
	"github.com/jira-backend/jiraflow-backend/internal/pkg/logger"
	"github.com/jira-backend/jiraflow-backend/internal/pkg/token"
)

func NewRouter(h *handlers.Handler, tokenMaker token.Maker, enforcer *casbin.Enforcer, log logger.Logger) *gin.Engine {
	r := gin.New()

	r.Use(
		middleware.RequestID(),
		gin.Logger(),
		middleware.Logger(log),
		middleware.Recover(log),
		middleware.CORS(),
	)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.GET("/health", v1.HealthCheck())

	auth := middleware.Auth(tokenMaker)

	api := r.Group("/api/v1")

	// Auth (public)
	authGroup := api.Group("/auth")
	{
		authGroup.POST("/register", v1.Register(h))
		authGroup.POST("/login", v1.Login(h))
		authGroup.POST("/refresh", v1.Refresh(h))
		authGroup.POST("/forgot-password", v1.ForgotPassword(h))
		authGroup.POST("/reset-password", v1.ResetPassword(h))
	}

	// Public invite accept
	api.POST("/invites/accept", v1.AcceptInvite(h))

	// Public Telegram webhook
	api.POST("/telegram/webhook", v1.TelegramWebhook(h))

	// Public GitHub webhook
	api.POST("/github/webhook", v1.GitHubWebhook(h))

	// Protected routes
	protected := api.Group("/", auth, middleware.RateLimitByUser(30, 60))

	// Auth me/logout
	protected.POST("/auth/logout", v1.Logout(h))
	protected.GET("/auth/me", v1.Me(h))

	// Telegram
	protected.GET("/auth/telegram/status", v1.GetTelegramStatus(h))
	protected.POST("/auth/telegram/connect", v1.GenerateTelegramCode(h))
	protected.DELETE("/auth/telegram/disconnect", v1.DisconnectTelegram(h))

	// Users
	users := protected.Group("/users")
	{
		users.GET("", v1.ListUsers(h))
		users.POST("", v1.CreateUser(h))
		users.GET("/:id", v1.GetUser(h))
		users.PUT("/:id", v1.UpdateUser(h))
		users.POST("/:id/deactivate", v1.DeactivateUser(h))
		users.POST("/:id/activate", v1.ActivateUser(h))
		users.PUT("/:id/password", v1.ChangePassword(h))
	}

	// Invites
	invites := protected.Group("/invites")
	{
		invites.POST("", v1.CreateInvite(h))
		invites.GET("", v1.ListPendingInvites(h))
		invites.DELETE("/:id", v1.RevokeInvite(h))
	}

	// Projects
	projects := protected.Group("/projects")
	{
		projects.POST("", v1.CreateProject(h))
		projects.GET("", v1.ListProjects(h))
		projects.GET("/:id", v1.GetProject(h))
		projects.PUT("/:id", v1.UpdateProject(h))
		projects.DELETE("/:id", v1.DeleteProject(h))
		projects.POST("/:id/archive", v1.ArchiveProject(h))

		// Project members
		projects.POST("/:id/members", v1.AddProjectMember(h))
		projects.GET("/:id/members", v1.ListProjectMembers(h))
		projects.PUT("/:id/members/:user_id", v1.UpdateProjectMemberRole(h))
		projects.DELETE("/:id/members/:user_id", v1.RemoveProjectMember(h))

		// Sprints
		projects.POST("/:id/sprints", v1.CreateSprint(h))
		projects.GET("/:id/sprints", v1.ListSprints(h))

		// Components
		projects.POST("/:id/components", v1.CreateComponent(h))
		projects.GET("/:id/components", v1.ListComponents(h))

		// Versions
		projects.POST("/:id/versions", v1.CreateVersion(h))
		projects.GET("/:id/versions", v1.ListVersions(h))

		// GitHub integration
		projects.POST("/:id/github", v1.ConnectGitHubRepo(h))
		projects.GET("/:id/github", v1.GetGitHubRepo(h))
		projects.DELETE("/:id/github", v1.DisconnectGitHubRepo(h))

		// Project-level reports
		projects.GET("/:id/velocity", v1.GetProjectVelocity(h))
		projects.GET("/:id/cfd", v1.GetProjectCFD(h))
		projects.GET("/:id/roadmap", v1.GetRoadmap(h))
		projects.GET("/:id/backlog", v1.GetBacklog(h))
		projects.GET("/:id/space", v1.GetProjectLinkedSpace(h))
		projects.GET("/:id/dashboard", v1.GetProjectDashboard(h))

		// Labels
		projects.POST("/:id/labels", v1.CreateLabel(h))
		projects.GET("/:id/labels", v1.ListLabels(h))

		// Custom fields
		projects.POST("/:id/custom-fields", v1.CreateCustomField(h))
		projects.GET("/:id/custom-fields", v1.ListCustomFields(h))
		projects.PUT("/:id/custom-fields/reorder", v1.ReorderCustomFields(h))

		// Automation Rules
		projects.POST("/:id/automation-rules", v1.CreateAutomationRule(h))
		projects.GET("/:id/automation-rules", v1.ListAutomationRules(h))

		// Boards
		projects.POST("/:id/boards", v1.CreateBoard(h))
		projects.GET("/:id/boards", v1.ListBoards(h))

		// Scrum planning
		projects.GET("/:id/sprint-planning", v1.GetSprintPlanning(h))
		projects.POST("/:id/sprint-planning", v1.BulkAssignToSprint(h))
		projects.GET("/:id/release-plan", v1.GetReleasePlan(h))

		// Gantt
		projects.GET("/:id/gantt", v1.GetGanttData(h))
	}

	// Workflows
	workflows := protected.Group("/workflows")
	{
		workflows.POST("", v1.CreateWorkflow(h))
		workflows.GET("", v1.ListWorkflows(h))
		workflows.GET("/:id", v1.GetWorkflow(h))
		workflows.PUT("/:id", v1.UpdateWorkflow(h))
		workflows.DELETE("/:id", v1.DeleteWorkflow(h))

		workflows.POST("/:id/statuses", v1.CreateWorkflowStatus(h))
		workflows.PUT("/statuses/:id", v1.UpdateWorkflowStatus(h))
		workflows.DELETE("/statuses/:id", v1.DeleteWorkflowStatus(h))

		workflows.POST("/:id/transitions", v1.CreateWorkflowTransition(h))
		workflows.DELETE("/transitions/:id", v1.DeleteWorkflowTransition(h))
	}

	// Sprints (individual)
	sprints := protected.Group("/sprints")
	{
		sprints.GET("/:id", v1.GetSprint(h))
		sprints.PUT("/:id", v1.UpdateSprint(h))
		sprints.DELETE("/:id", v1.DeleteSprint(h))
		sprints.POST("/:id/start", v1.StartSprint(h))
		sprints.POST("/:id/complete", v1.CompleteSprint(h))
		sprints.GET("/:id/report", v1.GetSprintReport(h))
		sprints.GET("/:id/burndown", v1.GetSprintBurndown(h))
		sprints.GET("/:id/burnup", v1.GetSprintBurnup(h))
		sprints.POST("/:id/issues", v1.AddIssueToSprint(h))
		sprints.DELETE("/:id/issues/:issue_id", v1.RemoveIssueFromSprint(h))
		sprints.GET("/:id/capacity", v1.GetSprintCapacity(h))
		sprints.POST("/:id/goal", v1.UpdateSprintGoal(h))
		sprints.GET("/:id/impediments", v1.GetSprintImpediments(h))
	}

	// Components (individual)
	components := protected.Group("/components")
	{
		components.GET("/:id", v1.GetComponent(h))
		components.PUT("/:id", v1.UpdateComponent(h))
		components.DELETE("/:id", v1.DeleteComponent(h))
	}

	// Versions (individual)
	versions := protected.Group("/versions")
	{
		versions.GET("/:id", v1.GetVersion(h))
		versions.PUT("/:id", v1.UpdateVersion(h))
		versions.DELETE("/:id", v1.DeleteVersion(h))
		versions.POST("/:id/release", v1.ReleaseVersion(h))
		versions.POST("/:id/archive", v1.ArchiveVersion(h))
		versions.GET("/:id/release-notes", v1.GetVersionReleaseNotes(h))
	}

	// Issues
	issues := protected.Group("/issues")
	{
		issues.PUT("/reorder", v1.ReorderIssues(h))
		issues.PUT("/bulk", v1.BulkUpdateIssues(h))
		issues.DELETE("/bulk", v1.BulkDeleteIssues(h))
		issues.POST("", v1.CreateIssue(h))
		issues.GET("", v1.ListIssues(h))
		issues.GET("/key/:key", v1.GetIssueByKey(h))
		issues.GET("/:id", v1.GetIssue(h))
		issues.PUT("/:id", v1.UpdateIssue(h))
		issues.DELETE("/:id", v1.DeleteIssue(h))
		issues.POST("/:id/transition", v1.TransitionIssue(h))
		issues.POST("/:id/clone", v1.CloneIssue(h))
		issues.PUT("/:id/rank", v1.RankIssue(h))
		issues.PUT("/:id/move", v1.MoveIssueOnBoard(h))
		issues.GET("/:id/history", v1.ListIssueHistory(h))

		// Issue links
		issues.POST("/:id/links", v1.CreateIssueLink(h))
		issues.GET("/:id/links", v1.ListIssueLinks(h))
		issues.DELETE("/links/:id", v1.DeleteIssueLink(h))

		// Issue watchers
		issues.POST("/:id/watchers", v1.AddIssueWatcher(h))
		issues.GET("/:id/watchers", v1.ListIssueWatchers(h))
		issues.DELETE("/:id/watchers", v1.RemoveIssueWatcher(h))

		// Time tracking
		issues.POST("/:id/worklogs", v1.CreateWorklog(h))
		issues.GET("/:id/worklogs", v1.ListWorklogs(h))
		issues.PUT("/:id/worklogs/:worklog_id", v1.UpdateWorklog(h))
		issues.DELETE("/:id/worklogs/:worklog_id", v1.DeleteWorklog(h))
		issues.GET("/:id/time-summary", v1.GetTimeSpentSummary(h))
		issues.PUT("/:id/estimates", v1.UpdateEstimates(h))

		// Epic
		issues.GET("/:id/epic-progress", v1.GetEpicProgress(h))

		// Multiple assignees
		issues.PUT("/:id/assignees", v1.SetIssueAssignees(h))
		issues.GET("/:id/assignees", v1.ListIssueAssignees(h))
		issues.DELETE("/:id/assignees/:user_id", v1.RemoveIssueAssignee(h))

		// Issue ↔ Page links
		issues.POST("/:id/page-links", v1.LinkIssuePage(h))
		issues.GET("/:id/page-links", v1.ListIssuePageLinks(h))
		issues.DELETE("/:id/page-links/:page_id", v1.UnlinkIssuePage(h))

		// Issue votes
		issues.POST("/:id/votes", v1.ToggleIssueVote(h))
		issues.GET("/:id/votes", v1.GetIssueVoteSummary(h))

		// GitHub
		issues.GET("/:id/commits", v1.ListIssueCommits(h))
		issues.GET("/:id/pull-requests", v1.ListIssuePRs(h))
	}

	// Comments: /api/v1/:parent_type/:parent_id/comments
	// parent_type = "issues" | "pages"
	protected.POST("/:parent_type/:parent_id/comments", v1.CreateComment(h))
	protected.GET("/:parent_type/:parent_id/comments", v1.ListComments(h))

	// Attachments: /api/v1/:parent_type/:parent_id/attachments
	protected.POST("/:parent_type/:parent_id/attachments", v1.UploadAttachment(h))
	protected.GET("/:parent_type/:parent_id/attachments", v1.ListAttachments(h))

	// Comments (individual)
	protected.GET("/comments/:id", v1.GetComment(h))
	protected.PUT("/comments/:id", v1.UpdateComment(h))
	protected.DELETE("/comments/:id", v1.DeleteComment(h))
	protected.POST("/comments/:id/reactions", v1.ToggleCommentReaction(h))
	protected.GET("/comments/:id/reactions", v1.ListCommentReactions(h))

	// Attachments (individual)
	protected.GET("/attachments/:id", v1.GetAttachment(h))
	protected.GET("/attachments/:id/url", v1.GetAttachmentURL(h))
	protected.DELETE("/attachments/:id", v1.DeleteAttachment(h))

	// Boards
	boards := protected.Group("/boards")
	{
		boards.GET("/:id", v1.GetBoard(h))
		boards.PUT("/:id", v1.UpdateBoard(h))
		boards.DELETE("/:id", v1.DeleteBoard(h))
		boards.POST("/:id/columns", v1.CreateBoardColumn(h))
		boards.PUT("/:id/columns/reorder", v1.ReorderBoardColumns(h))
		boards.GET("/:id/swimlanes", v1.GetBoardSwimlanes(h))
		boards.PUT("/:id/swimlane-type", v1.SetBoardSwimlaneType(h))
	}

	// Board columns
	boardColumns := protected.Group("/board-columns")
	{
		boardColumns.PUT("/:id", v1.UpdateBoardColumn(h))
		boardColumns.DELETE("/:id", v1.DeleteBoardColumn(h))
	}

	// Labels (individual)
	protected.GET("/labels/:id", v1.GetLabel(h))
	protected.PUT("/labels/:id", v1.UpdateLabel(h))
	protected.DELETE("/labels/:id", v1.DeleteLabel(h))

	// Custom fields (individual)
	protected.GET("/custom-fields/:id", v1.GetCustomField(h))
	protected.PUT("/custom-fields/:id", v1.UpdateCustomField(h))
	protected.DELETE("/custom-fields/:id", v1.DeleteCustomField(h))

	// Files
	files := protected.Group("/files")
	{
		files.POST("/upload", v1.UploadFile(h))
		files.GET("/presign", v1.GetFilePresignedURL(h))
	}

	// Spaces
	spaces := protected.Group("/spaces")
	{
		spaces.POST("", v1.CreateSpace(h))
		spaces.GET("", v1.ListSpaces(h))
		spaces.GET("/:id", v1.GetSpace(h))
		spaces.PUT("/:id", v1.UpdateSpace(h))
		spaces.DELETE("/:id", v1.DeleteSpace(h))

		spaces.POST("/:id/members", v1.AddSpaceMember(h))
		spaces.GET("/:id/members", v1.ListSpaceMembers(h))
		spaces.PUT("/:id/members/:user_id", v1.UpdateSpaceMemberRole(h))
		spaces.DELETE("/:id/members/:user_id", v1.RemoveSpaceMember(h))

		spaces.POST("/:id/archive", v1.ArchiveSpace(h))
		spaces.POST("/:id/restore", v1.RestoreSpace(h))
		spaces.GET("/:id/statistics", v1.GetSpaceStatistics(h))

		spaces.POST("/:id/pages", v1.CreatePage(h))
		spaces.GET("/:id/pages/tree", v1.GetPageTree(h))

		// Space-level page tags
		spaces.POST("/:id/page-tags", v1.CreatePageTag(h))
		spaces.GET("/:id/page-tags", v1.ListPageTags(h))

		// Blog posts (space-scoped)
		spaces.POST("/:id/blog-posts", v1.CreateBlogPost(h))
		spaces.GET("/:id/blog-posts", v1.ListBlogPosts(h))

		// Space export (async ZIP)
		spaces.POST("/:id/export", v1.RequestSpaceExport(h))
		spaces.GET("/:id/export", v1.ListSpaceExports(h))
		spaces.GET("/:id/export/:export_id", v1.GetSpaceExport(h))
	}

	// Blog posts (individual)
	blogPosts := protected.Group("/blog-posts")
	{
		blogPosts.GET("/:id", v1.GetBlogPost(h))
		blogPosts.PUT("/:id", v1.UpdateBlogPost(h))
		blogPosts.DELETE("/:id", v1.DeleteBlogPost(h))
		blogPosts.POST("/:id/publish", v1.PublishBlogPost(h))
		blogPosts.POST("/:id/unpublish", v1.UnpublishBlogPost(h))
	}

	// Page tags (individual)
	pageTags := protected.Group("/page-tags")
	{
		pageTags.GET("/:id", v1.GetPageTag(h))
		pageTags.PUT("/:id", v1.UpdatePageTag(h))
		pageTags.DELETE("/:id", v1.DeletePageTag(h))
		pageTags.GET("/:id/pages", v1.GetPagesByTag(h))
	}

	// Pages
	pages := protected.Group("/pages")
	{
		pages.GET("", v1.ListPages(h))
		pages.GET("/:id", v1.GetPage(h))
		pages.PUT("/:id", v1.UpdatePage(h))
		pages.DELETE("/:id", v1.DeletePage(h))
		pages.PUT("/:id/move", v1.MovePage(h))
		pages.POST("/:id/watch", v1.WatchPage(h))
		pages.DELETE("/:id/watch", v1.UnwatchPage(h))
		pages.GET("/:id/watchers", v1.ListPageWatchers(h))

		pages.GET("/:id/versions", v1.ListPageVersions(h))
		pages.GET("/:id/versions/:version", v1.GetPageVersionByNumber(h))
		pages.GET("/:id/versions/:version/diff/:v2", v1.DiffPageVersions(h))

		// Page tags
		pages.PUT("/:id/tags", v1.SetPageTags(h))
		pages.GET("/:id/tags", v1.GetPageTagsForPage(h))

		// Page analytics
		pages.POST("/:id/view", v1.RecordPageView(h))
		pages.GET("/:id/analytics", v1.GetPageAnalytics(h))

		// Inline comments
		pages.POST("/:id/inline-comments", v1.CreateInlineComment(h))
		pages.GET("/:id/inline-comments", v1.ListInlineComments(h))

		// Restrictions
		pages.PUT("/:id/restrictions", v1.SetPageRestrictions(h))
		pages.GET("/:id/restrictions", v1.ListPageRestrictions(h))
		pages.DELETE("/:id/restrictions", v1.ClearPageRestrictions(h))
		pages.GET("/:id/access", v1.CheckPageAccess(h))

		// Export
		pages.GET("/:id/export/html", v1.ExportPageHTML(h))
		pages.GET("/:id/export/pdf", v1.ExportPagePDF(h))
		pages.GET("/:id/export/md", v1.ExportPageMarkdown(h))
		pages.GET("/:id/export/docx", v1.ExportPageDOCX(h))

		// Copy
		pages.POST("/:id/copy", v1.CopyPage(h))

		// Reactions
		pages.POST("/:id/reactions", v1.TogglePageReaction(h))
		pages.GET("/:id/reactions", v1.ListPageReactions(h))
		pages.GET("/:id/reactions/:emoji/users", v1.ListPageReactionUsers(h))

		// Collaborative lock
		pages.POST("/:id/lock", v1.AcquirePageLock(h))
		pages.DELETE("/:id/lock", v1.ReleasePageLock(h))
		pages.GET("/:id/lock", v1.GetPageLock(h))
		pages.PUT("/:id/lock", v1.ExtendPageLock(h))

		// Macros
		pages.POST("/:id/macros", v1.UpsertPageMacro(h))
		pages.GET("/:id/macros", v1.ListPageMacros(h))

		// Page ↔ Issue links
		pages.GET("/:id/issue-links", v1.ListPageIssueLinks(h))
	}

	// Inline comments (individual)
	inlineComments := protected.Group("/inline-comments")
	{
		inlineComments.PUT("/:id", v1.UpdateInlineComment(h))
		inlineComments.POST("/:id/resolve", v1.ResolveInlineComment(h))
		inlineComments.POST("/:id/unresolve", v1.UnresolveInlineComment(h))
		inlineComments.DELETE("/:id", v1.DeleteInlineComment(h))
	}

	// Page templates
	pageTemplates := protected.Group("/page-templates")
	{
		pageTemplates.POST("", v1.CreatePageTemplate(h))
		pageTemplates.GET("", v1.ListPageTemplates(h))
		pageTemplates.GET("/:id", v1.GetPageTemplate(h))
		pageTemplates.PUT("/:id", v1.UpdatePageTemplate(h))
		pageTemplates.DELETE("/:id", v1.DeletePageTemplate(h))
	}

	// Favorites
	favs := protected.Group("/favorites")
	{
		favs.POST("", v1.AddFavorite(h))
		favs.DELETE("", v1.RemoveFavorite(h))
		favs.GET("", v1.ListFavorites(h))
		favs.GET("/check", v1.IsFavorite(h))
	}

	// Recently visited
	protected.GET("/recently-visited", v1.ListRecentPages(h))

	// Macros (individual)
	protected.DELETE("/page-macros/:id", v1.DeletePageMacro(h))

	// Webhooks
	webhooks := protected.Group("/webhooks")
	{
		webhooks.POST("", v1.CreateWebhook(h))
		webhooks.GET("/:id", v1.GetWebhook(h))
		webhooks.PUT("/:id", v1.UpdateWebhook(h))
		webhooks.DELETE("/:id", v1.DeleteWebhook(h))
		webhooks.GET("/:id/deliveries", v1.ListWebhookDeliveries(h))
	}
	projects.GET("/:id/webhooks", v1.ListWebhooksByProject(h))
	spaces.GET("/:id/webhooks", v1.ListWebhooksBySpace(h))

	// Page versions (individual)
	protected.GET("/page-versions/:id", v1.GetPageVersion(h))

	// Notifications
	notifs := protected.Group("/notifications")
	{
		notifs.GET("", v1.ListNotifications(h))
		notifs.GET("/unread-count", v1.CountUnreadNotifications(h))
		notifs.POST("/mark-read", v1.MarkNotificationsRead(h))
		notifs.POST("/mark-all-read", v1.MarkAllNotificationsRead(h))
		notifs.DELETE("/:id", v1.DeleteNotification(h))
		notifs.GET("/preferences", v1.GetNotificationPreference(h))
		notifs.PUT("/preferences", v1.UpdateNotificationPreference(h))
	}

	// Activity feed
	protected.GET("/activity", v1.ListActivity(h))

	// Search
	protected.GET("/search", v1.Search(h))
	protected.GET("/search/suggestions", v1.SearchSuggestions(h))

	// Audit logs
	protected.GET("/audit-logs", v1.ListAuditLogs(h))
	protected.GET("/audit-logs/export", v1.ExportAuditLogs(h))

	// WebSocket
	protected.GET("/ws", v1.ServeWS(h))

	// ── Faza 7: Enterprise ─────────────────────────────────────────────────────

	// Google OAuth2 (public endpoints — no JWT required)
	authGroup.GET("/google", v1.GoogleLogin(h))
	authGroup.GET("/google/callback", v1.GoogleCallback(h))
	authGroup.GET("/google/link", auth, v1.GoogleLink(h))
	authGroup.GET("/google/link/callback", auth, v1.GoogleLinkCallback(h))
	protected.GET("/auth/providers", v1.ListLinkedProviders(h))
	protected.DELETE("/auth/providers/:provider", v1.UnlinkOAuthProvider(h))

	// API Keys
	apiKeys := protected.Group("/api-keys")
	{
		apiKeys.POST("", v1.CreateAPIKey(h))
		apiKeys.GET("", v1.ListAPIKeys(h))
		apiKeys.DELETE("/:id", v1.RevokeAPIKey(h))
	}

	// Permission Schemes
	schemes := protected.Group("/permission-schemes")
	{
		schemes.POST("", v1.CreatePermissionScheme(h))
		schemes.GET("", v1.ListPermissionSchemes(h))
		schemes.GET("/:id", v1.GetPermissionScheme(h))
		schemes.PUT("/:id", v1.UpdatePermissionScheme(h))
		schemes.DELETE("/:id", v1.DeletePermissionScheme(h))
		schemes.POST("/:id/grants", v1.AddSchemeGrant(h))
		schemes.DELETE("/:id/grants/:grant_id", v1.RemoveSchemeGrant(h))
	}
	projects.PUT("/:id/permission-scheme", v1.AssignPermissionScheme(h))
	projects.GET("/:id/permission-scheme", v1.GetProjectPermissionScheme(h))

	// Data Import
	importGroup := protected.Group("/import")
	{
		importGroup.POST("/jira", v1.ImportJira(h))
		importGroup.POST("/trello", v1.ImportTrello(h))
		importGroup.POST("/linear", v1.ImportLinear(h))
		importGroup.GET("/:id", v1.GetImportStatus(h))
	}

	// Saved filters
	savedFilters := protected.Group("/saved-filters")
	{
		savedFilters.POST("", v1.CreateSavedFilter(h))
		savedFilters.GET("", v1.ListSavedFilters(h))
		savedFilters.GET("/:id", v1.GetSavedFilter(h))
		savedFilters.PUT("/:id", v1.UpdateSavedFilter(h))
		savedFilters.DELETE("/:id", v1.DeleteSavedFilter(h))
	}

	projectTemplates := protected.Group("/project-templates")
	{
		projectTemplates.GET("", v1.ListProjectTemplates(h))
		projectTemplates.GET("/:id", v1.GetProjectTemplate(h))
	}

	notificationSchemes := protected.Group("/notification-schemes")
	{
		notificationSchemes.GET("", v1.ListNotificationSchemes(h))
		notificationSchemes.POST("", v1.CreateNotificationScheme(h))
		notificationSchemes.GET("/:id", v1.GetNotificationScheme(h))
		notificationSchemes.DELETE("/:id", v1.DeleteNotificationScheme(h))
	}

	fieldConfigurations := protected.Group("/field-configurations")
	{
		fieldConfigurations.GET("", v1.ListFieldConfigurations(h))
		fieldConfigurations.POST("", v1.CreateFieldConfiguration(h))
		fieldConfigurations.GET("/:id", v1.GetFieldConfiguration(h))
		fieldConfigurations.DELETE("/:id", v1.DeleteFieldConfiguration(h))
	}

	issueTypes := protected.Group("/issue-types")
	{
		issueTypes.GET("", v1.ListIssueTypes(h))
		issueTypes.POST("", v1.CreateIssueType(h))
		issueTypes.GET("/:id", v1.GetIssueType(h))
		issueTypes.DELETE("/:id", v1.DeleteIssueType(h))
	}

	issueTypeSchemes := protected.Group("/issue-type-schemes")
	{
		issueTypeSchemes.GET("", v1.ListIssueTypeSchemes(h))
		issueTypeSchemes.POST("", v1.CreateIssueTypeScheme(h))
		issueTypeSchemes.GET("/:id", v1.GetIssueTypeScheme(h))
		issueTypeSchemes.DELETE("/:id", v1.DeleteIssueTypeScheme(h))
	}

	projects.GET("/:id/issue-type-scheme", v1.GetProjectIssueTypeScheme(h))

	blueprints := protected.Group("/blueprints")
	{
		blueprints.GET("", v1.ListBlueprints(h))
		blueprints.GET("/:id", v1.GetBlueprint(h))
		blueprints.POST("", v1.CreateBlueprint(h))
		blueprints.DELETE("/:id", v1.DeleteBlueprint(h))
		blueprints.POST("/:id/create", v1.CreatePageFromBlueprint(h))
	}

	spaceCategories := protected.Group("/space-categories")
	{
		spaceCategories.POST("", v1.CreateSpaceCategory(h))
		spaceCategories.GET("", v1.ListSpaceCategories(h))
		spaceCategories.GET("/:id", v1.GetSpaceCategory(h))
		spaceCategories.PUT("/:id", v1.UpdateSpaceCategory(h))
		spaceCategories.DELETE("/:id", v1.DeleteSpaceCategory(h))
	}

	// Content Properties (pages)
	pages.GET("/:id/properties", v1.ListContentProperties("page")(h))
	pages.GET("/:id/properties/:key", v1.GetContentProperty("page")(h))
	pages.PUT("/:id/properties/:key", v1.SetContentProperty("page")(h))
	pages.DELETE("/:id/properties/:key", v1.DeleteContentProperty("page")(h))

	// Content Properties (issues)
	issues.GET("/:id/properties", v1.ListContentProperties("issue")(h))
	issues.GET("/:id/properties/:key", v1.GetContentProperty("issue")(h))
	issues.PUT("/:id/properties/:key", v1.SetContentProperty("issue")(h))
	issues.DELETE("/:id/properties/:key", v1.DeleteContentProperty("issue")(h))

	// Security Schemes & Levels
	securitySchemes := protected.Group("/security-schemes")
	{
		securitySchemes.POST("", v1.CreateSecurityScheme(h))
		securitySchemes.GET("", v1.ListSecuritySchemes(h))
		securitySchemes.GET("/:id", v1.GetSecurityScheme(h))
		securitySchemes.DELETE("/:id", v1.DeleteSecurityScheme(h))
		securitySchemes.POST("/:id/levels", v1.AddSecurityLevel(h))
		securitySchemes.GET("/:id/levels/:level_id", v1.GetSecurityLevel(h))
		securitySchemes.DELETE("/:id/levels/:level_id", v1.DeleteSecurityLevel(h))
		securitySchemes.POST("/:id/levels/:level_id/members", v1.AddSecurityLevelMember(h))
		securitySchemes.DELETE("/:id/levels/:level_id/members/:member_id", v1.DeleteSecurityLevelMember(h))
	}

	// Automation Rules (individual)
	automationRules := protected.Group("/automation-rules")
	{
		automationRules.GET("/:id", v1.GetAutomationRule(h))
		automationRules.PUT("/:id", v1.UpdateAutomationRule(h))
		automationRules.DELETE("/:id", v1.DeleteAutomationRule(h))
		automationRules.POST("/:id/enable", v1.EnableAutomationRule(h))
		automationRules.POST("/:id/disable", v1.DisableAutomationRule(h))
		automationRules.GET("/:id/logs", v1.ListAutomationLogs(h))
	}

	return r
}
