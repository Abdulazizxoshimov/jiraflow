package automation

import (
	"context"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

type UseCase interface {
	Create(ctx context.Context, projectID, createdBy string, req *entity.CreateAutomationRuleReq) (*entity.AutomationRule, error)
	GetByID(ctx context.Context, id string) (*entity.AutomationRule, error)
	List(ctx context.Context, filter *entity.AutomationFilter) ([]*entity.AutomationRule, int, error)
	Update(ctx context.Context, id string, req *entity.UpdateAutomationRuleReq) (*entity.AutomationRule, error)
	Delete(ctx context.Context, id string) error
	Enable(ctx context.Context, id string) error
	Disable(ctx context.Context, id string) error

	// TriggerEvent — dispatcher tomonidan chaqiriladi; mos qoidalarni topib, bajaradi.
	TriggerEvent(ctx context.Context, event *entity.AutomationEvent) error

	ListLogs(ctx context.Context, ruleID string, filter *entity.Filter) ([]*entity.AutomationLog, int, error)
}
