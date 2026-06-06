package repository

import (
	"context"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

type AutomationRepository interface {
	Create(ctx context.Context, rule *entity.AutomationRule) error
	GetByID(ctx context.Context, id string) (*entity.AutomationRule, error)
	List(ctx context.Context, filter *entity.AutomationFilter) ([]*entity.AutomationRule, int, error)
	Update(ctx context.Context, rule *entity.AutomationRule) error
	Delete(ctx context.Context, id string) error
	SetActive(ctx context.Context, id string, isActive bool) error

	// Trigger engine uchun: event trigger_type bo'yicha active rule'larni topadi
	FindByTrigger(ctx context.Context, projectID, triggerType string) ([]*entity.AutomationRule, error)

	// Log
	SaveLog(ctx context.Context, log *entity.AutomationLog) error
	ListLogs(ctx context.Context, ruleID string, filter *entity.Filter) ([]*entity.AutomationLog, int, error)
}
