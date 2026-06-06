package audit

import (
	"context"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

type UseCase interface {
	Log(ctx context.Context, log *entity.AuditLog) error
	List(ctx context.Context, filter *entity.AuditLogFilter) ([]*entity.AuditLog, int, error)
}
