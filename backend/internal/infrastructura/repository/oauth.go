package repository

import (
	"context"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

type OAuthRepository interface {
	SaveState(ctx context.Context, s *entity.OAuthState) error
	GetState(ctx context.Context, state string) (*entity.OAuthState, error)
	DeleteState(ctx context.Context, state string) error

	UpsertAccount(ctx context.Context, acc *entity.OAuthAccount) error
	GetAccountByProvider(ctx context.Context, provider, providerUserID string) (*entity.OAuthAccount, error)
	ListByUser(ctx context.Context, userID string) ([]*entity.OAuthAccount, error)
	DeleteAccount(ctx context.Context, userID, provider string) error
}
