package repository

import (
	"context"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

type InlineCommentRepository interface {
	Create(ctx context.Context, c *entity.InlineComment) error
	GetByID(ctx context.Context, id string) (*entity.InlineComment, error)
	ListByPage(ctx context.Context, pageID string) ([]*entity.InlineComment, error)
	ListByAnchor(ctx context.Context, pageID, anchorID string) ([]*entity.InlineComment, error)
	Update(ctx context.Context, c *entity.InlineComment) error
	Resolve(ctx context.Context, id, resolverID string) error
	Unresolve(ctx context.Context, id string) error
	SoftDelete(ctx context.Context, id string) error
}
