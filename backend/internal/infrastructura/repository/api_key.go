package repository

import (
	"context"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

type APIKeyRepository interface {
	Create(ctx context.Context, key *entity.APIKey, keyHash string) error
	ListByUser(ctx context.Context, userID string) ([]*entity.APIKey, error)
	GetByID(ctx context.Context, id string) (*entity.APIKey, error)
	GetByHash(ctx context.Context, keyHash string) (*entity.APIKey, error)
	Revoke(ctx context.Context, id, userID string) error
	UpdateLastUsed(ctx context.Context, id string) error
}
