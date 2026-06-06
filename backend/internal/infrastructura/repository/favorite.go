package repository

import (
	"context"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

type FavoriteRepository interface {
	Add(ctx context.Context, fav *entity.Favorite) error
	Remove(ctx context.Context, userID, entityType, entityID string) error
	List(ctx context.Context, userID string, filter *entity.FavoriteFilter) ([]*entity.Favorite, int, error)
	Exists(ctx context.Context, userID, entityType, entityID string) (bool, error)
}
