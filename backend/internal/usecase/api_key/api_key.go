package api_key

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
	"github.com/jira-backend/jiraflow-backend/internal/infrastructure/repository"
	apperr "github.com/jira-backend/jiraflow-backend/internal/pkg/errors"
)

type useCase struct {
	repo repository.APIKeyRepository
}

func New(repo repository.APIKeyRepository) UseCase {
	return &useCase{repo: repo}
}

func (uc *useCase) Create(ctx context.Context, userID string, req *entity.CreateAPIKeyReq) (*entity.CreateAPIKeyResp, error) {
	rawBytes := make([]byte, 24)
	if _, err := rand.Read(rawBytes); err != nil {
		return nil, fmt.Errorf("api_key.Create rand: %w", err)
	}
	plainKey := "jfk_" + hex.EncodeToString(rawBytes) // jfk = jiraflow key

	keyHash := hashKey(plainKey)
	prefix := plainKey[:12]

	now := time.Now().UTC()
	key := &entity.APIKey{
		ID:        uuid.NewString(),
		UserID:    userID,
		Name:      req.Name,
		KeyPrefix: prefix,
		Scopes:    req.Scopes,
		ExpiresAt: req.ExpiresAt,
		CreatedAt: now,
	}
	if key.Scopes == nil {
		key.Scopes = []string{}
	}

	if err := uc.repo.Create(ctx, key, keyHash); err != nil {
		return nil, fmt.Errorf("api_key.Create: %w", err)
	}

	return &entity.CreateAPIKeyResp{APIKey: key, PlainKey: plainKey}, nil
}

func (uc *useCase) List(ctx context.Context, userID string) ([]*entity.APIKey, error) {
	return uc.repo.ListByUser(ctx, userID)
}

func (uc *useCase) Revoke(ctx context.Context, id, userID string) error {
	return uc.repo.Revoke(ctx, id, userID)
}

func (uc *useCase) ValidateKey(ctx context.Context, plainKey string) (*entity.APIKey, error) {
	if len(plainKey) < 4 || plainKey[:4] != "jfk_" {
		return nil, apperr.Unauthorized("invalid api key format")
	}

	keyHash := hashKey(plainKey)
	key, err := uc.repo.GetByHash(ctx, keyHash)
	if err != nil {
		return nil, apperr.Unauthorized("invalid api key")
	}
	if key.RevokedAt != nil {
		return nil, apperr.Unauthorized("api key has been revoked")
	}
	if key.ExpiresAt != nil && time.Now().UTC().After(*key.ExpiresAt) {
		return nil, apperr.Unauthorized("api key has expired")
	}

	// fire-and-forget
	go func() {
		_ = uc.repo.UpdateLastUsed(context.Background(), key.ID)
	}()

	return key, nil
}

func hashKey(plainKey string) string {
	h := sha256.Sum256([]byte(plainKey))
	return hex.EncodeToString(h[:])
}
