package favorite

import (
	"context"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

type UseCase interface {
	Add(ctx context.Context, userID string, req *entity.AddFavoriteReq) (*entity.Favorite, error)
	Remove(ctx context.Context, userID, entityType, entityID string) error
	List(ctx context.Context, userID string, filter *entity.FavoriteFilter) ([]*entity.Favorite, int, error)
	IsFavorite(ctx context.Context, userID, entityType, entityID string) (bool, error)
}
