package repository

import (
	"context"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

type CommentRepository interface {
	Create(ctx context.Context, c *entity.Comment) error
	GetByID(ctx context.Context, id string) (*entity.Comment, error)
	ListByParent(ctx context.Context, parentType, parentID string, filter *entity.Filter) ([]*entity.Comment, int, error)
	Update(ctx context.Context, c *entity.Comment) error
	SoftDelete(ctx context.Context, id string) error

	AddMention(ctx context.Context, m *entity.CommentMention) error
	ListMentions(ctx context.Context, commentID string) ([]*entity.CommentMention, error)
	DeleteMentions(ctx context.Context, commentID string) error

	ToggleReaction(ctx context.Context, commentID, userID, emoji string) error
	ListReactions(ctx context.Context, commentID, viewerID string) ([]entity.CommentReactionSummary, error)
}
