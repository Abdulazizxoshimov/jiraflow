package usecase

import (
	"time"

	emailpkg "github.com/jira-backend/jiraflow-backend/internal/infrastructure/email"
	"github.com/jira-backend/jiraflow-backend/internal/infrastructure/minio"
	ws "github.com/jira-backend/jiraflow-backend/internal/infrastructure/websocket"
	"github.com/jira-backend/jiraflow-backend/internal/storage"
	"github.com/jira-backend/jiraflow-backend/internal/pkg/hasher"
	"github.com/jira-backend/jiraflow-backend/internal/pkg/logger"
	"github.com/jira-backend/jiraflow-backend/internal/pkg/token"
	"github.com/jira-backend/jiraflow-backend/internal/usecase/attachment"
	"github.com/jira-backend/jiraflow-backend/internal/usecase/audit"
	"github.com/jira-backend/jiraflow-backend/internal/usecase/auth"
	"github.com/jira-backend/jiraflow-backend/internal/usecase/board"
	"github.com/jira-backend/jiraflow-backend/internal/usecase/comment"
	"github.com/jira-backend/jiraflow-backend/internal/usecase/component"
	custom_field "github.com/jira-backend/jiraflow-backend/internal/usecase/custom_field"
	"github.com/jira-backend/jiraflow-backend/internal/usecase/file"
	"github.com/jira-backend/jiraflow-backend/internal/usecase/invite"
	"github.com/jira-backend/jiraflow-backend/internal/usecase/issue"
	issue_link "github.com/jira-backend/jiraflow-backend/internal/usecase/issue_link"
	"github.com/jira-backend/jiraflow-backend/internal/usecase/label"
	"github.com/jira-backend/jiraflow-backend/internal/usecase/notification"
	"github.com/jira-backend/jiraflow-backend/internal/usecase/favorite"
	"github.com/jira-backend/jiraflow-backend/internal/usecase/inline_comment"
	"github.com/jira-backend/jiraflow-backend/internal/usecase/page"
	"github.com/jira-backend/jiraflow-backend/internal/usecase/page_export"
	"github.com/jira-backend/jiraflow-backend/internal/usecase/page_restriction"
	page_tag "github.com/jira-backend/jiraflow-backend/internal/usecase/page_tag"
	page_template "github.com/jira-backend/jiraflow-backend/internal/usecase/page_template"
	page_version "github.com/jira-backend/jiraflow-backend/internal/usecase/page_version"
	page_view "github.com/jira-backend/jiraflow-backend/internal/usecase/page_view"
	"github.com/jira-backend/jiraflow-backend/internal/usecase/project"
	project_member "github.com/jira-backend/jiraflow-backend/internal/usecase/project_member"
	"github.com/jira-backend/jiraflow-backend/internal/usecase/search"
	"github.com/jira-backend/jiraflow-backend/internal/usecase/space"
	"github.com/jira-backend/jiraflow-backend/internal/usecase/sprint"
	"github.com/jira-backend/jiraflow-backend/internal/usecase/user"
	"github.com/jira-backend/jiraflow-backend/internal/usecase/version"
	"github.com/jira-backend/jiraflow-backend/internal/usecase/workflow"
	"github.com/jira-backend/jiraflow-backend/internal/usecase/worklog"
	"github.com/jira-backend/jiraflow-backend/internal/usecase/activity_feed"
	"github.com/jira-backend/jiraflow-backend/internal/usecase/issue_assignee"
	issue_page_link "github.com/jira-backend/jiraflow-backend/internal/usecase/issue_page_link"
	"github.com/jira-backend/jiraflow-backend/internal/usecase/page_lock"
	"github.com/jira-backend/jiraflow-backend/internal/usecase/page_macro"
	"github.com/jira-backend/jiraflow-backend/internal/usecase/page_reaction"
	"github.com/jira-backend/jiraflow-backend/internal/usecase/webhook"
	api_key "github.com/jira-backend/jiraflow-backend/internal/usecase/api_key"
	data_import "github.com/jira-backend/jiraflow-backend/internal/usecase/data_import"
	"github.com/jira-backend/jiraflow-backend/internal/usecase/oauth"
	permission_scheme "github.com/jira-backend/jiraflow-backend/internal/usecase/permission_scheme"
	issue_vote "github.com/jira-backend/jiraflow-backend/internal/usecase/issue_vote"
	blog_post "github.com/jira-backend/jiraflow-backend/internal/usecase/blog_post"
	saved_filter "github.com/jira-backend/jiraflow-backend/internal/usecase/saved_filter"
	space_export "github.com/jira-backend/jiraflow-backend/internal/usecase/space_export"
	space_category "github.com/jira-backend/jiraflow-backend/internal/usecase/space_category"
	content_property "github.com/jira-backend/jiraflow-backend/internal/usecase/content_property"
	"github.com/jira-backend/jiraflow-backend/internal/usecase/blueprint"
	issue_type "github.com/jira-backend/jiraflow-backend/internal/usecase/issue_type"
	notification_scheme "github.com/jira-backend/jiraflow-backend/internal/usecase/notification_scheme"
	project_template "github.com/jira-backend/jiraflow-backend/internal/usecase/project_template"
	field_configuration "github.com/jira-backend/jiraflow-backend/internal/usecase/field_configuration"
	automation_uc "github.com/jira-backend/jiraflow-backend/internal/usecase/automation"
	security_scheme "github.com/jira-backend/jiraflow-backend/internal/usecase/security_scheme"
	telegram_uc "github.com/jira-backend/jiraflow-backend/internal/usecase/telegram"
	github_uc "github.com/jira-backend/jiraflow-backend/internal/usecase/github"
	tgclient "github.com/jira-backend/jiraflow-backend/internal/infrastructure/telegram"
)

// UseCases barcha usecase interfeyslari.
type UseCases struct {
	SecurityScheme      security_scheme.UseCase
	Auth          auth.UseCase
	User          user.UseCase
	Workflow      workflow.UseCase
	Project       project.UseCase
	ProjectMember project_member.UseCase
	Invite        invite.UseCase
	Sprint        sprint.UseCase
	Issue         issue.UseCase
	IssueLink     issue_link.UseCase
	Worklog       worklog.UseCase
	Component     component.UseCase
	Version       version.UseCase
	Label         label.UseCase
	CustomField   custom_field.UseCase
	Board         board.UseCase
	Comment       comment.UseCase
	Attachment    attachment.UseCase
	File          file.UseCase
	Space           space.UseCase
	Page            page.UseCase
	PageVersion     page_version.UseCase
	PageTag         page_tag.UseCase
	PageView        page_view.UseCase
	PageTemplate    page_template.UseCase
	PageRestriction page_restriction.UseCase
	PageExport      page_export.UseCase
	InlineComment   inline_comment.UseCase
	Favorite        favorite.UseCase
	Notification    notification.UseCase
	Search          search.UseCase
	Audit           audit.UseCase
	IssueAssignee   issue_assignee.UseCase
	PageReaction    page_reaction.UseCase
	Webhook         webhook.UseCase
	PageLock        page_lock.UseCase
	PageMacro       page_macro.UseCase
	IssuePageLink    issue_page_link.UseCase
	ActivityFeed     activity_feed.UseCase
	OAuth            oauth.UseCase
	APIKey           api_key.UseCase
	PermissionScheme permission_scheme.UseCase
	DataImport       data_import.UseCase
	IssueVote           issue_vote.UseCase
	BlogPost            blog_post.UseCase
	SavedFilter         saved_filter.UseCase
	SpaceExport         space_export.UseCase
	SpaceCategory       space_category.UseCase
	ContentProperty     content_property.UseCase
	Blueprint           blueprint.UseCase
	IssueType           issue_type.UseCase
	NotificationScheme  notification_scheme.UseCase
	ProjectTemplate     project_template.UseCase
	FieldConfiguration  field_configuration.UseCase
	Automation          automation_uc.UseCase
	Telegram            telegram_uc.UseCase
	GitHub              github_uc.UseCase
}

// Deps — UseCases uchun tashqi bog'liqliklar.
type Deps struct {
	Store               *storage.Storage
	TokenMaker          token.Maker
	Hasher              hasher.Hasher
	Minio               minio.Client
	Log                 logger.Logger
	Hub                 *ws.Hub
	EmailSender         emailpkg.Sender
	FrontendBaseURL     string
	GoogleClientID      string
	GoogleClientSecret  string
	GoogleRedirectURL   string
	TelegramBotToken      string
	TelegramWebhookURL    string
	TelegramWebhookSecret string
}

func New(d Deps) *UseCases {
	notifUC := notification.New(d.Store.Notification, d.Log)
	disp := notification.NewDispatcher(notifUC, d.Store.User, d.Hub, d.EmailSender)

	automationUC := automation_uc.New(d.Store.Automation, d.Store.Issue, d.Store.Project, d.Log)
	notification.SetAutomation(disp, automationUC)

	tgClient := tgclient.New(d.TelegramBotToken)
	telegramUC := telegram_uc.New(d.Store.Telegram, tgClient, d.TelegramWebhookURL, d.TelegramWebhookSecret)
	notification.SetTelegram(disp, telegramUC)

	return &UseCases{
		Auth:          auth.New(d.Store.User, d.Store.Auth, d.TokenMaker, d.Hasher, 24*time.Hour, d.EmailSender, d.FrontendBaseURL, d.Log),
		User:          user.New(d.Store.User, d.Store.Space, d.Hasher, d.Log),
		Workflow:      workflow.New(d.Store.Workflow, d.Log),
		Project:       project.New(d.Store.Project, d.Store.Workflow, d.Store.Space, d.Log),
		ProjectMember: project_member.New(d.Store.ProjectMember, d.Store.Project, d.Log),
		Invite:        invite.New(d.Store.Invite, d.Store.User, d.TokenMaker, d.Hasher, d.Log),
		Sprint:        sprint.New(d.Store.Sprint, d.Store.Issue, d.Store.Space, d.Store.Page, d.Store.PageVersion, disp, d.Log),
		Issue:         issue.New(d.Store.Issue, d.Store.Project, d.Store.Workflow, d.Store.Version, d.Store.ProjectMember, disp, d.Log),
		IssueLink:     issue_link.New(d.Store.IssueLink, d.Store.Issue, d.Log),
		Worklog:       worklog.New(d.Store.Worklog, d.Store.Issue, d.Log),
		Component:     component.New(d.Store.Component, d.Log),
		Version:       version.New(d.Store.Version, d.Log),
		Label:         label.New(d.Store.Label, d.Log),
		CustomField:   custom_field.New(d.Store.CustomField, d.Log),
		Board:         board.New(d.Store.Board, d.Log),
		Comment:       comment.New(d.Store.Comment, d.Store.Issue, d.Store.Page, disp, d.Log),
		Attachment:    attachment.New(d.Store.Attachment, d.Minio, d.Log),
		File:          file.New(d.Minio, d.Log),
		Space:           space.New(d.Store.Space, d.Log),
		PageVersion:     page_version.New(d.Store.PageVersion, d.Log),
		Page:            page.New(d.Store.Page, d.Store.PageVersion, d.Store.Space, d.Store.IssuePageLink, d.Store.Issue, d.Log),
		PageTag:         page_tag.New(d.Store.PageTag, d.Log),
		PageView:        page_view.New(d.Store.PageView, d.Log),
		PageTemplate:    page_template.New(d.Store.PageTemplate, d.Log),
		PageRestriction: page_restriction.New(d.Store.PageRestriction, d.Store.Page, d.Store.Space, d.Log),
		PageExport:      page_export.New(d.Store.Page),
		InlineComment:   inline_comment.New(d.Store.InlineComment, d.Log),
		Favorite:        favorite.New(d.Store.Favorite, d.Store.Page, d.Store.Space, d.Log),
		Notification:    notifUC,
		Search:          search.New(d.Store.Search, d.Log),
		Audit:           audit.New(d.Store.Audit, d.Log),
		IssueAssignee:   issue_assignee.New(d.Store.IssueAssignee, d.Store.Issue, d.Log),
		PageReaction:    page_reaction.New(d.Store.PageReaction, d.Log),
		Webhook:         webhook.New(d.Store.Webhook, d.Log),
		PageLock:        page_lock.New(d.Store.PageLock),
		PageMacro:       page_macro.New(d.Store.PageMacro),
		IssuePageLink:    issue_page_link.New(d.Store.IssuePageLink, d.Store.Issue, d.Store.Page, d.Log),
		ActivityFeed:     activity_feed.New(d.Store.ActivityFeed, d.Log),
		OAuth:            oauth.New(d.Store.OAuth, d.Store.User, d.Store.Auth, d.TokenMaker, d.GoogleClientID, d.GoogleClientSecret, d.GoogleRedirectURL),
		APIKey:           api_key.New(d.Store.APIKey),
		PermissionScheme: permission_scheme.New(d.Store.PermissionScheme),
		DataImport:       data_import.New(d.Store.DataImport, d.Store.Issue, d.Store.Project, d.Store.Workflow),
		IssueVote:        issue_vote.New(d.Store.IssueVote, d.Store.Issue),
		BlogPost:         blog_post.New(d.Store.BlogPost),
		SavedFilter:        saved_filter.New(d.Store.SavedFilter),
		SpaceExport:        space_export.New(d.Store.SpaceExport, d.Store.Page, d.Minio, d.Log),
		SpaceCategory:      space_category.New(d.Store.SpaceCategory),
		ContentProperty:    content_property.New(d.Store.ContentProperty),
		Blueprint:          blueprint.New(d.Store.Blueprint, d.Store.Page),
		IssueType:          issue_type.New(d.Store.IssueType),
		NotificationScheme: notification_scheme.New(d.Store.NotificationScheme),
		ProjectTemplate:    project_template.New(d.Store.ProjectTemplate),
		FieldConfiguration: field_configuration.New(d.Store.FieldConfiguration),
		SecurityScheme:     security_scheme.New(d.Store.SecurityScheme),
		Automation:         automationUC,
		Telegram:           telegramUC,
		GitHub:             github_uc.New(d.Store.GitHubRepo, d.Store.IssueCommit, d.Store.IssuePR, d.Store.Issue, d.Store.Project, d.Store.Workflow),
	}
}
