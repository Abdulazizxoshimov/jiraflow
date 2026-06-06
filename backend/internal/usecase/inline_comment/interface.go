package inline_comment

import (
	"context"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

type UseCase interface {
	Create(ctx context.Context, pageID, authorID string, req *entity.CreateInlineCommentReq) (*entity.InlineComment, error)
	GetByID(ctx context.Context, id string) (*entity.InlineComment, error)
	ListByPage(ctx context.Context, pageID string) ([]*entity.InlineComment, error)
	ListByAnchor(ctx context.Context, pageID, anchorID string) ([]*entity.InlineComment, error)
	Update(ctx context.Context, id, actorID string, req *entity.UpdateInlineCommentReq) (*entity.InlineComment, error)
	Resolve(ctx context.Context, id, resolverID string) error
	Unresolve(ctx context.Context, id, actorID string) error
	Delete(ctx context.Context, id, actorID string) error
}
