package page_lock

import (
	"context"
	"fmt"
	"time"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
	"github.com/jira-backend/jiraflow-backend/internal/infrastructure/repository"
	apperr "github.com/jira-backend/jiraflow-backend/internal/pkg/errors"
)

const (
	defaultTTL = 60
	maxTTL     = 3600
)

type useCase struct {
	repo repository.PageLockRepository
}

func New(repo repository.PageLockRepository) UseCase {
	return &useCase{repo: repo}
}

func (uc *useCase) Acquire(ctx context.Context, pageID, userID string, ttlSeconds int) (*entity.PageLock, error) {
	if ttlSeconds <= 0 {
		ttlSeconds = defaultTTL
	}
	if ttlSeconds > maxTTL {
		ttlSeconds = maxTTL
	}

	existing, err := uc.repo.Get(ctx, pageID)
	if err == nil && existing.UserID != userID {
		return nil, apperr.Conflict("page is locked by another user")
	}

	lock := &entity.PageLock{
		PageID:    pageID,
		UserID:    userID,
		ExpiresAt: time.Now().UTC().Add(time.Duration(ttlSeconds) * time.Second),
	}
	if err := uc.repo.Acquire(ctx, lock); err != nil {
		return nil, fmt.Errorf("pageLock.Acquire: %w", err)
	}
	return lock, nil
}

func (uc *useCase) Release(ctx context.Context, pageID, userID string) error {
	return uc.repo.Release(ctx, pageID, userID)
}

func (uc *useCase) Get(ctx context.Context, pageID string) (*entity.PageLock, error) {
	return uc.repo.Get(ctx, pageID)
}

func (uc *useCase) Extend(ctx context.Context, pageID, userID string, ttlSeconds int) (*entity.PageLock, error) {
	if ttlSeconds <= 0 {
		ttlSeconds = defaultTTL
	}
	expiresAt := time.Now().UTC().Add(time.Duration(ttlSeconds) * time.Second)
	if err := uc.repo.Extend(ctx, pageID, userID, expiresAt); err != nil {
		return nil, fmt.Errorf("pageLock.Extend: %w", err)
	}
	return uc.repo.Get(ctx, pageID)
}
