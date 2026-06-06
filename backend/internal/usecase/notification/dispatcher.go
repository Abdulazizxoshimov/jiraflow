package notification

import (
	"context"
	"fmt"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
	"github.com/jira-backend/jiraflow-backend/internal/infrastructure/email"
	"github.com/jira-backend/jiraflow-backend/internal/infrastructure/repository"
	"github.com/jira-backend/jiraflow-backend/internal/infrastructure/websocket"
	tguc "github.com/jira-backend/jiraflow-backend/internal/usecase/telegram"
)

// AutomationTrigger — automation engine'ga event yuboradi.
type AutomationTrigger interface {
	TriggerEvent(ctx context.Context, event *entity.AutomationEvent) error
}

// Dispatcher sends notifications to watchers of an entity.
type Dispatcher interface {
	IssueAssigned(ctx context.Context, issue *entity.Issue, assigneeID, actorID string)
	IssueCreated(ctx context.Context, issue *entity.Issue, actorID string)
	IssueUpdated(ctx context.Context, issue *entity.Issue, watcherIDs []string, actorID string)
	IssueCommented(ctx context.Context, issueID string, watcherIDs []string, actorID string)
	IssueMentioned(ctx context.Context, issueID string, mentionedUserIDs []string, actorID string)
	IssueStatusChanged(ctx context.Context, issue *entity.Issue, watcherIDs []string, actorID string)
	PageCommented(ctx context.Context, pageID string, watcherIDs []string, actorID string)
	PageMentioned(ctx context.Context, pageID string, mentionedUserIDs []string, actorID string)
	SprintStarted(ctx context.Context, sprintID, projectID, name string)
	SprintCompleted(ctx context.Context, sprintID, projectID, name string)
}

type dispatcher struct {
	uc          UseCase
	userRepo    repository.UserRepository
	hub         *websocket.Hub
	email       email.Sender
	automation  AutomationTrigger
	telegramUC  tguc.UseCase
}

func NewDispatcher(uc UseCase, userRepo repository.UserRepository, hub *websocket.Hub, emailSender email.Sender) Dispatcher {
	return &dispatcher{uc: uc, userRepo: userRepo, hub: hub, email: emailSender}
}

// SetAutomation — automation engine'ni dispatcher'ga ulaydi (circular import oldini olish uchun setter).
func SetAutomation(d Dispatcher, at AutomationTrigger) {
	if dp, ok := d.(*dispatcher); ok {
		dp.automation = at
	}
}

// SetTelegram — telegram usecase'ni dispatcher'ga ulaydi.
func SetTelegram(d Dispatcher, tg tguc.UseCase) {
	if dp, ok := d.(*dispatcher); ok {
		dp.telegramUC = tg
	}
}

func (d *dispatcher) notify(ctx context.Context, n *entity.Notification, emailSubject, emailTemplate string, shouldEmail bool) {
	_ = d.uc.Notify(ctx, n)

	// WebSocket — foydalanuvchi onlayn bo'lsa real-time yuborish
	if d.hub != nil {
		d.hub.Send(n.UserID, websocket.NewNotificationMsg(n))
	}

	// Email — preference'ga qarab
	if shouldEmail && d.email != nil && n.UserID != "" {
		go func() {
			pref, err := d.uc.GetPreference(ctx, n.UserID)
			if err != nil {
				return
			}
			if !d.prefAllows(pref, n.Type) {
				return
			}
			user, err := d.userRepo.GetByID(ctx, n.UserID)
			if err != nil || user.Email == "" {
				return
			}
			_ = d.email.Send(ctx, []string{user.Email}, emailSubject, "notification", map[string]any{
				"Title":     emailSubject,
				"Body":      emailTemplate,
				"ActionURL": "",
			})
		}()
	}

	// Telegram — ulangan bo'lsa xabar yuborish
	if d.telegramUC != nil && n.UserID != "" && emailTemplate != "" {
		go func() {
			_ = d.telegramUC.SendNotification(ctx, n.UserID, emailTemplate)
		}()
	}
}

func (d *dispatcher) prefAllows(pref *entity.NotificationPreference, notifType string) bool {
	switch notifType {
	case "issue_assigned":
		return pref.EmailAssigned
	case "mentioned":
		return pref.EmailMentioned
	case "issue_commented", "page_commented":
		return pref.EmailCommented
	case "issue_status_changed":
		return pref.EmailStatus
	default:
		return pref.EmailWatcher
	}
}

func (d *dispatcher) IssueCreated(ctx context.Context, issue *entity.Issue, actorID string) {
	if d.automation != nil {
		go d.automation.TriggerEvent(ctx, &entity.AutomationEvent{
			Type:       "issue.created",
			ProjectID:  issue.ProjectID,
			EntityID:   issue.ID,
			EntityType: "issue",
			Payload:    map[string]any{"type": issue.Type, "priority": issue.Priority, "status_id": issue.StatusID, "reporter_id": issue.ReporterID},
		})
	}
}

func (d *dispatcher) IssueUpdated(ctx context.Context, issue *entity.Issue, watcherIDs []string, actorID string) {
	for _, uid := range watcherIDs {
		if uid == actorID {
			continue
		}
		n := &entity.Notification{
			UserID:     uid,
			Type:       "issue_updated",
			EntityType: strPtr("issue"),
			EntityID:   &issue.ID,
			ActorID:    &actorID,
			Payload:    map[string]any{"issue_id": issue.ID, "title": issue.Title},
		}
		d.notify(ctx, n, fmt.Sprintf("Issue updated: %s", issue.Title),
			fmt.Sprintf("Issue <strong>%s</strong> has been updated.", issue.Title), false)
	}

	if d.automation != nil {
		go d.automation.TriggerEvent(ctx, &entity.AutomationEvent{
			Type:       "issue.updated",
			ProjectID:  issue.ProjectID,
			EntityID:   issue.ID,
			EntityType: "issue",
			Payload:    map[string]any{"type": issue.Type, "priority": issue.Priority, "status_id": issue.StatusID},
		})
	}
}

func (d *dispatcher) SprintStarted(ctx context.Context, sprintID, projectID, name string) {
	if d.automation != nil {
		go d.automation.TriggerEvent(ctx, &entity.AutomationEvent{
			Type:       "sprint.started",
			ProjectID:  projectID,
			EntityID:   sprintID,
			EntityType: "sprint",
			Payload:    map[string]any{"sprint_id": sprintID, "name": name},
		})
	}
}

func (d *dispatcher) SprintCompleted(ctx context.Context, sprintID, projectID, name string) {
	if d.automation != nil {
		go d.automation.TriggerEvent(ctx, &entity.AutomationEvent{
			Type:       "sprint.completed",
			ProjectID:  projectID,
			EntityID:   sprintID,
			EntityType: "sprint",
			Payload:    map[string]any{"sprint_id": sprintID, "name": name},
		})
	}
}

func (d *dispatcher) IssueAssigned(ctx context.Context, issue *entity.Issue, assigneeID, actorID string) {
	if assigneeID == actorID {
		return
	}
	n := &entity.Notification{
		UserID:     assigneeID,
		Type:       "issue_assigned",
		EntityType: strPtr("issue"),
		EntityID:   &issue.ID,
		ActorID:    &actorID,
		Payload:    map[string]any{"issue_id": issue.ID, "title": issue.Title},
	}
	d.notify(ctx, n, fmt.Sprintf("Issue assigned: %s", issue.Title),
		fmt.Sprintf("Issue <strong>%s</strong> has been assigned to you.", issue.Title), true)

	if d.automation != nil {
		go d.automation.TriggerEvent(ctx, &entity.AutomationEvent{
			Type:       "issue.assigned",
			ProjectID:  issue.ProjectID,
			EntityID:   issue.ID,
			EntityType: "issue",
			Payload:    map[string]any{"assignee_id": assigneeID, "type": issue.Type, "priority": issue.Priority, "status_id": issue.StatusID},
		})
	}
}

func (d *dispatcher) IssueCommented(ctx context.Context, issueID string, watcherIDs []string, actorID string) {
	for _, uid := range watcherIDs {
		if uid == actorID {
			continue
		}
		n := &entity.Notification{
			UserID:     uid,
			Type:       "issue_commented",
			EntityType: strPtr("issue"),
			EntityID:   &issueID,
			ActorID:    &actorID,
			Payload:    map[string]any{"issue_id": issueID},
		}
		d.notify(ctx, n, "New comment on an issue you're watching", "Someone added a new comment on an issue you are watching.", true)
	}
}

func (d *dispatcher) IssueMentioned(ctx context.Context, issueID string, mentionedUserIDs []string, actorID string) {
	for _, uid := range mentionedUserIDs {
		if uid == actorID {
			continue
		}
		n := &entity.Notification{
			UserID:     uid,
			Type:       "mentioned",
			EntityType: strPtr("issue"),
			EntityID:   &issueID,
			ActorID:    &actorID,
			Payload:    map[string]any{"issue_id": issueID},
		}
		d.notify(ctx, n, "You were mentioned in an issue", "You were mentioned in a comment on an issue.", true)
	}
}

func (d *dispatcher) IssueStatusChanged(ctx context.Context, issue *entity.Issue, watcherIDs []string, actorID string) {
	for _, uid := range watcherIDs {
		if uid == actorID {
			continue
		}
		n := &entity.Notification{
			UserID:     uid,
			Type:       "issue_status_changed",
			EntityType: strPtr("issue"),
			EntityID:   &issue.ID,
			ActorID:    &actorID,
			Payload:    map[string]any{"issue_id": issue.ID, "status_id": issue.StatusID},
		}
		d.notify(ctx, n, fmt.Sprintf("Issue status changed: %s", issue.Title),
			fmt.Sprintf("The status of issue <strong>%s</strong> has changed.", issue.Title), true)
	}

	if d.automation != nil {
		go d.automation.TriggerEvent(ctx, &entity.AutomationEvent{
			Type:       "issue.transition",
			ProjectID:  issue.ProjectID,
			EntityID:   issue.ID,
			EntityType: "issue",
			Payload:    map[string]any{"status_id": issue.StatusID, "type": issue.Type, "priority": issue.Priority},
		})
	}
}

func (d *dispatcher) PageCommented(ctx context.Context, pageID string, watcherIDs []string, actorID string) {
	for _, uid := range watcherIDs {
		if uid == actorID {
			continue
		}
		n := &entity.Notification{
			UserID:     uid,
			Type:       "page_commented",
			EntityType: strPtr("page"),
			EntityID:   &pageID,
			ActorID:    &actorID,
			Payload:    map[string]any{"page_id": pageID},
		}
		d.notify(ctx, n, "New comment on a page you're watching", "Someone added a new comment on a page you are watching.", true)
	}
}

func (d *dispatcher) PageMentioned(ctx context.Context, pageID string, mentionedUserIDs []string, actorID string) {
	for _, uid := range mentionedUserIDs {
		if uid == actorID {
			continue
		}
		n := &entity.Notification{
			UserID:     uid,
			Type:       "mentioned",
			EntityType: strPtr("page"),
			EntityID:   &pageID,
			ActorID:    &actorID,
			Payload:    map[string]any{"page_id": pageID},
		}
		d.notify(ctx, n, "You were mentioned in a page comment", "You were mentioned in a comment on a Confluence page.", true)
	}
}

func strPtr(s string) *string { return &s }
