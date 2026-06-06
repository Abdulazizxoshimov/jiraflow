package audit

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
	"github.com/jira-backend/jiraflow-backend/internal/infrastructure/repository"
	"github.com/jira-backend/jiraflow-backend/internal/pkg/logger"
)

type useCase struct {
	auditRepo repository.AuditRepository
	log       logger.Logger
}

func New(auditRepo repository.AuditRepository, log logger.Logger) UseCase {
	return &useCase{auditRepo: auditRepo, log: log}
}

func (uc *useCase) Log(ctx context.Context, log *entity.AuditLog) error {
	if log.ID == "" {
		log.ID = uuid.NewString()
	}
	if log.CreatedAt.IsZero() {
		log.CreatedAt = time.Now().UTC()
	}
	if log.Details == nil {
		log.Details = map[string]any{}
	}
	if err := uc.auditRepo.Create(ctx, log); err != nil {
		uc.log.Error(ctx, "audit.Log: db error", logger.SafeString("err", err.Error()))
		return err
	}
	return nil
}

func (uc *useCase) List(ctx context.Context, filter *entity.AuditLogFilter) ([]*entity.AuditLog, int, error) {
	return uc.auditRepo.List(ctx, filter)
}
