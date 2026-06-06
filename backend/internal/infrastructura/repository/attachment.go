package repository

import (
	"context"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

type AttachmentRepository interface {
	Create(ctx context.Context, a *entity.Attachment) error
	GetByID(ctx context.Context, id string) (*entity.Attachment, error)
	ListByParent(ctx context.Context, parentType, parentID string) ([]*entity.Attachment, error)
	SoftDelete(ctx context.Context, id string) error
	GetTotalSizeByParent(ctx context.Context, parentType, parentID string) (int64, error)
}
