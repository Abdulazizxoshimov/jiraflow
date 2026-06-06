package api_key

import (
	"context"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

type UseCase interface {
	Create(ctx context.Context, userID string, req *entity.CreateAPIKeyReq) (*entity.CreateAPIKeyResp, error)
	List(ctx context.Context, userID string) ([]*entity.APIKey, error)
	Revoke(ctx context.Context, id, userID string) error
	// ValidateKey checks the plain-text key, updates last_used_at, and returns the key record.
	ValidateKey(ctx context.Context, plainKey string) (*entity.APIKey, error)
}
