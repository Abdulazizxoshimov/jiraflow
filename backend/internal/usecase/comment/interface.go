package comment

import (
	"context"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

type UseCase interface {
	Create(ctx context.Context, parentType, parentID, authorID string, req *entity.CreateCommentReq) (*entity.Comment, error)
	GetByID(ctx context.Context, id string) (*entity.Comment, error)
	ListByParent(ctx context.Context, parentType, parentID string, filter *entity.Filter) ([]*entity.Comment, int, error)
	Update(ctx context.Context, id, actorID string, req *entity.UpdateCommentReq) (*entity.Comment, error)
	Delete(ctx context.Context, id, actorID string) error
	ToggleReaction(ctx context.Context, commentID, userID, emoji string) error
	ListReactions(ctx context.Context, commentID, viewerID string) ([]entity.CommentReactionSummary, error)
}
