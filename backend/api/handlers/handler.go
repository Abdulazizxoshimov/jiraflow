package handlers

import (
	"github.com/jira-backend/jiraflow-backend/internal/infrastructura/websocket"
	"github.com/jira-backend/jiraflow-backend/internal/usecase/attachment"
	"github.com/jira-backend/jiraflow-backend/internal/usecase/audit"
	"github.com/jira-backend/jiraflow-backend/internal/usecase/auth"
	"github.com/jira-backend/jiraflow-backend/internal/usecase/board"
	"github.com/jira-backend/jiraflow-backend/internal/usecase/comment"
	"github.com/jira-backend/jiraflow-backend/internal/usecase/component"
	"github.com/jira-backend/jiraflow-backend/internal/usecase/custom_field"
	"github.com/jira-backend/jiraflow-backend/internal/usecase/file"
	"github.com/jira-backend/jiraflow-backend/internal/usecase/invite"
	"github.com/jira-backend/jiraflow-backend/internal/usecase/issue"
	"github.com/jira-backend/jiraflow-backend/internal/usecase/activity_feed"
	"github.com/jira-backend/jiraflow-backend/internal/usecase/issue_link"
	"github.com/jira-backend/jiraflow-backend/internal/usecase/issue_page_link"
	"github.com/jira-backend/jiraflow-backend/internal/usecase/label"
	"github.com/jira-backend/jiraflow-backend/internal/usecase/notification"
	"github.com/jira-backend/jiraflow-backend/internal/usecase/favorite"
	"github.com/jira-backend/jiraflow-backend/internal/usecase/inline_comment"
	"github.com/jira-backend/jiraflow-backend/internal/usecase/page"
	"github.com/jira-backend/jiraflow-backend/internal/usecase/page_export"
	"github.com/jira-backend/jiraflow-backend/internal/usecase/page_restriction"
	"github.com/jira-backend/jiraflow-backend/internal/usecase/page_tag"
	"github.com/jira-backend/jiraflow-backend/internal/usecase/page_template"
	"github.com/jira-backend/jiraflow-backend/internal/usecase/page_version"
	"github.com/jira-backend/jiraflow-backend/internal/usecase/page_view"
	"github.com/jira-backend/jiraflow-backend/internal/usecase/project"
	"github.com/jira-backend/jiraflow-backend/internal/usecase/project_member"
	"github.com/jira-backend/jiraflow-backend/internal/usecase/search"
	"github.com/jira-backend/jiraflow-backend/internal/usecase/space"
	"github.com/jira-backend/jiraflow-backend/internal/usecase/sprint"
	"github.com/jira-backend/jiraflow-backend/internal/usecase/user"
	"github.com/jira-backend/jiraflow-backend/internal/usecase/version"
	"github.com/jira-backend/jiraflow-backend/internal/usecase/workflow"
	"github.com/jira-backend/jiraflow-backend/internal/usecase/worklog"
	"github.com/jira-backend/jiraflow-backend/internal/usecase/issue_assignee"
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
	security_scheme "github.com/jira-backend/jiraflow-backend/internal/usecase/security_scheme"
	"github.com/jira-backend/jiraflow-backend/internal/usecase/automation"
	"github.com/jira-backend/jiraflow-backend/internal/usecase/telegram"
	github_uc "github.com/jira-backend/jiraflow-backend/internal/usecase/github"
)

type Handler struct {

	Auth          auth.UseCase
	User          user.UseCase
	Project       project.UseCase
	ProjectMember project_member.UseCase
	Invite        invite.UseCase
	Workflow      workflow.UseCase
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
	PageTemplate    page_template.UseCase
	PageRestriction page_restriction.UseCase
	PageView        page_view.UseCase
	PageExport      page_export.UseCase
	Favorite        favorite.UseCase
	InlineComment   inline_comment.UseCase
	Notification    notification.UseCase
	Search          search.UseCase
	Audit           audit.UseCase
	Hub             *websocket.Hub
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
	SecurityScheme          security_scheme.UseCase
	Automation              automation.UseCase
	Telegram                telegram.UseCase
	GitHub                  github_uc.UseCase
	TelegramWebhookSecret   string
	GitHubWebhookSecret     string
}
