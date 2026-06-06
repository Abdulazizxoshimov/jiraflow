package oauth

import (
	"context"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

type UseCase interface {
	// GenerateAuthURL creates a state token and returns the Google OAuth2 URL.
	GenerateAuthURL(ctx context.Context, redirectURL string) (string, error)
	// HandleCallback exchanges the code for tokens, upserts the user, and returns JWT pair.
	HandleCallback(ctx context.Context, state, code string) (*entity.OAuthCallbackResp, error)
	// LinkAccount links a Google account to an already-authenticated user.
	LinkAccount(ctx context.Context, userID, state, code string) error
	// UnlinkAccount removes the social login link for a provider.
	UnlinkAccount(ctx context.Context, userID, provider string) error
	// ListLinkedAccounts returns all OAuth accounts connected to the user.
	ListLinkedAccounts(ctx context.Context, userID string) ([]*entity.OAuthAccount, error)
}
