package repository

import (
	"context"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

type PageReactionRepository interface {
	Toggle(ctx context.Context, pageID, userID, emoji string) (added bool, err error)
	ListByPage(ctx context.Context, pageID, viewerUserID string) ([]*entity.PageReactionSummary, error)
	ListUsers(ctx context.Context, pageID, emoji string) ([]*entity.PageReaction, error)
}
