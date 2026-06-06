package comment

import (
	"context"
	"time"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

func (uc *useCase) AddMention(ctx context.Context, commentID, userID string) error {
	m := &entity.CommentMention{
		CommentID: commentID,
		UserID:    userID,
		CreatedAt: time.Now().UTC(),
	}
	return uc.repo.AddMention(ctx, m)
}

func (uc *useCase) ListMentions(ctx context.Context, commentID string) ([]*entity.CommentMention, error) {
	return uc.repo.ListMentions(ctx, commentID)
}
