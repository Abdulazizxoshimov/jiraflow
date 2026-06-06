package automation

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
	"github.com/jira-backend/jiraflow-backend/internal/infrastructura/repository"
	apperr "github.com/jira-backend/jiraflow-backend/internal/pkg/errors"
	"github.com/jira-backend/jiraflow-backend/internal/pkg/logger"
)

type useCase struct {
	repo        repository.AutomationRepository
	issueRepo   repository.IssueRepository
	projectRepo repository.ProjectRepository
	log         logger.Logger
}

func New(
	repo repository.AutomationRepository,
	issueRepo repository.IssueRepository,
	projectRepo repository.ProjectRepository,
	log logger.Logger,
) UseCase {
	return &useCase{
		repo:        repo,
		issueRepo:   issueRepo,
		projectRepo: projectRepo,
		log:         log,
	}
}

func (uc *useCase) Create(ctx context.Context, projectID, createdBy string, req *entity.CreateAutomationRuleReq) (*entity.AutomationRule, error) {
	if _, err := uc.projectRepo.GetByID(ctx, projectID); err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	rule := &entity.AutomationRule{
		ID:            uuid.NewString(),
		ProjectID:     projectID,
		Name:          req.Name,
		Description:   req.Description,
		TriggerType:   req.TriggerType,
		TriggerConfig: req.TriggerConfig,
		Conditions:    req.Conditions,
		Actions:       req.Actions,
		IsActive:      true,
		CreatedBy:     createdBy,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
	if rule.TriggerConfig == nil {
		rule.TriggerConfig = map[string]any{}
	}
	if rule.Conditions == nil {
		rule.Conditions = []entity.AutomationCondition{}
	}

	if err := uc.repo.Create(ctx, rule); err != nil {
		uc.log.Error(ctx, "automation.Create: db error", logger.SafeString("err", err.Error()))
		return nil, err
	}
	uc.log.Info(ctx, "automation rule created", logger.String("id", rule.ID))
	return rule, nil
}

func (uc *useCase) GetByID(ctx context.Context, id string) (*entity.AutomationRule, error) {
	return uc.repo.GetByID(ctx, id)
}

func (uc *useCase) List(ctx context.Context, filter *entity.AutomationFilter) ([]*entity.AutomationRule, int, error) {
	return uc.repo.List(ctx, filter)
}

func (uc *useCase) Update(ctx context.Context, id string, req *entity.UpdateAutomationRuleReq) (*entity.AutomationRule, error) {
	rule, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		rule.Name = *req.Name
	}
	if req.Description != nil {
		rule.Description = req.Description
	}
	if req.TriggerType != nil {
		rule.TriggerType = *req.TriggerType
	}
	if req.TriggerConfig != nil {
		rule.TriggerConfig = req.TriggerConfig
	}
	if req.Conditions != nil {
		rule.Conditions = req.Conditions
	}
	if req.Actions != nil {
		rule.Actions = req.Actions
	}
	if req.IsActive != nil {
		rule.IsActive = *req.IsActive
	}
	rule.UpdatedAt = time.Now().UTC()

	if err := uc.repo.Update(ctx, rule); err != nil {
		uc.log.Error(ctx, "automation.Update: db error", logger.String("id", id), logger.SafeString("err", err.Error()))
		return nil, err
	}
	return rule, nil
}

func (uc *useCase) Delete(ctx context.Context, id string) error {
	if _, err := uc.repo.GetByID(ctx, id); err != nil {
		return err
	}
	return uc.repo.Delete(ctx, id)
}

func (uc *useCase) Enable(ctx context.Context, id string) error {
	if _, err := uc.repo.GetByID(ctx, id); err != nil {
		return err
	}
	return uc.repo.SetActive(ctx, id, true)
}

func (uc *useCase) Disable(ctx context.Context, id string) error {
	if _, err := uc.repo.GetByID(ctx, id); err != nil {
		return err
	}
	return uc.repo.SetActive(ctx, id, false)
}

func (uc *useCase) TriggerEvent(ctx context.Context, event *entity.AutomationEvent) error {
	rules, err := uc.repo.FindByTrigger(ctx, event.ProjectID, event.Type)
	if err != nil {
		return fmt.Errorf("automation.TriggerEvent FindByTrigger: %w", err)
	}

	for _, rule := range rules {
		go uc.executeRule(context.Background(), rule, event)
	}
	return nil
}

func (uc *useCase) executeRule(ctx context.Context, rule *entity.AutomationRule, event *entity.AutomationEvent) {
	logEntry := &entity.AutomationLog{
		ID:         uuid.NewString(),
		RuleID:     rule.ID,
		EntityID:   &event.EntityID,
		EntityType: &event.EntityType,
		ExecutedAt: time.Now().UTC(),
		Status:     "success",
	}

	if !uc.checkConditions(rule.Conditions, event.Payload) {
		logEntry.Status = "skipped"
		_ = uc.repo.SaveLog(ctx, logEntry)
		return
	}

	var execErr error
	for _, action := range rule.Actions {
		if err := uc.executeAction(ctx, action, event); err != nil {
			execErr = err
			uc.log.Error(ctx, "automation.executeAction failed",
				logger.String("rule_id", rule.ID),
				logger.String("action_type", action.Type),
				logger.SafeString("err", err.Error()),
			)
			break
		}
	}

	if execErr != nil {
		logEntry.Status = "failed"
		errMsg := execErr.Error()
		logEntry.ErrorMsg = &errMsg
	}

	_ = uc.repo.SaveLog(ctx, logEntry)
}

// checkConditions — barcha shartlar to'g'ri bo'lsa true qaytaradi.
func (uc *useCase) checkConditions(conditions []entity.AutomationCondition, payload map[string]any) bool {
	for _, c := range conditions {
		val, ok := payload[c.Field]
		if !ok {
			if c.Operator == "is_empty" {
				continue
			}
			return false
		}

		strVal := fmt.Sprintf("%v", val)
		condVal := fmt.Sprintf("%v", c.Value)

		switch c.Operator {
		case "=":
			if strVal != condVal {
				return false
			}
		case "!=":
			if strVal == condVal {
				return false
			}
		case "in":
			if vals, ok := c.Value.([]any); ok {
				found := false
				for _, v := range vals {
					if fmt.Sprintf("%v", v) == strVal {
						found = true
						break
					}
				}
				if !found {
					return false
				}
			}
		case "not_in":
			if vals, ok := c.Value.([]any); ok {
				for _, v := range vals {
					if fmt.Sprintf("%v", v) == strVal {
						return false
					}
				}
			}
		case "is_empty":
			if strVal != "" && strVal != "<nil>" {
				return false
			}
		case "is_not_empty":
			if strVal == "" || strVal == "<nil>" {
				return false
			}
		case "contains":
			if !strings.Contains(strings.ToLower(strVal), strings.ToLower(condVal)) {
				return false
			}
		}
	}
	return true
}

// executeAction — bitta action'ni bajaradi.
func (uc *useCase) executeAction(ctx context.Context, action entity.AutomationAction, event *entity.AutomationEvent) error {
	switch action.Type {
	case "transition_issue":
		statusID, _ := action.Config["status_id"].(string)
		if statusID == "" {
			return apperr.BadRequest("transition_issue action requires status_id")
		}
		return uc.issueRepo.UpdateStatus(ctx, event.EntityID, statusID)

	case "assign_issue":
		assigneeID, _ := action.Config["assignee_id"].(string)
		if assigneeID == "" {
			return apperr.BadRequest("assign_issue action requires assignee_id")
		}
		issue, err := uc.issueRepo.GetByID(ctx, event.EntityID)
		if err != nil {
			return err
		}
		issue.AssigneeID = &assigneeID
		return uc.issueRepo.Update(ctx, issue)

	case "add_label":
		labelID, _ := action.Config["label_id"].(string)
		if labelID == "" {
			return apperr.BadRequest("add_label action requires label_id")
		}
		existing, _ := uc.issueRepo.GetLabels(ctx, event.EntityID)
		ids := make([]string, 0, len(existing)+1)
		for _, l := range existing {
			ids = append(ids, l.ID)
		}
		ids = append(ids, labelID)
		return uc.issueRepo.SetLabels(ctx, event.EntityID, ids)

	case "set_field":
		// Generic field update — supported: priority, story_points
		field, _ := action.Config["field"].(string)
		value := action.Config["value"]
		issue, err := uc.issueRepo.GetByID(ctx, event.EntityID)
		if err != nil {
			return err
		}
		switch field {
		case "priority":
			if v, ok := value.(string); ok {
				issue.Priority = v
			}
		case "story_points":
			if v, ok := value.(float64); ok {
				sp := int(v)
				issue.StoryPoints = &sp
			}
		}
		return uc.issueRepo.Update(ctx, issue)

	case "send_notification":
		// Notification sending is handled via dispatcher — log it
		uc.log.Info(ctx, "automation: send_notification action triggered",
			logger.String("entity_id", event.EntityID))
		return nil

	case "create_issue":
		// Clone/create new issue — minimal implementation
		title, _ := action.Config["title"].(string)
		issueType, _ := action.Config["type"].(string)
		if title == "" {
			title = "Auto-created issue"
		}
		if issueType == "" {
			issueType = "task"
		}
		now := time.Now().UTC()
		newIssue := &entity.Issue{
			ID:           uuid.NewString(),
			ProjectID:    event.ProjectID,
			Title:        title,
			Type:         issueType,
			Priority:     "medium",
			ReporterID:   "automation",
			CustomFields: map[string]any{},
			CreatedAt:    now,
			UpdatedAt:    now,
		}
		return uc.issueRepo.Create(ctx, newIssue)

	default:
		uc.log.Warn(ctx, "automation: unknown action type", logger.String("type", action.Type))
		return nil
	}
}

func (uc *useCase) ListLogs(ctx context.Context, ruleID string, filter *entity.Filter) ([]*entity.AutomationLog, int, error) {
	return uc.repo.ListLogs(ctx, ruleID, filter)
}
