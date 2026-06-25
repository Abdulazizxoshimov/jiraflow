package page_lock

import (
	"context"
	"fmt"
	"time"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
	"github.com/jira-backend/jiraflow-backend/internal/infrastructure/repository"
	ws "github.com/jira-backend/jiraflow-backend/internal/infrastructure/websocket"
)

const (
	defaultTTL = 60
	maxTTL     = 3600
)

type useCase struct {
	repo repository.PageLockRepository
	hub  *ws.Hub
}

func New(repo repository.PageLockRepository, hub *ws.Hub) UseCase {
	return &useCase{repo: repo, hub: hub}
}

func (uc *useCase) Acquire(ctx context.Context, pageID, userID string, ttlSeconds int) (*entity.PageLock, error) {
	if ttlSeconds <= 0 {
		ttlSeconds = defaultTTL
	}
	if ttlSeconds > maxTTL {
		ttlSeconds = maxTTL
	}

	lock := &entity.PageLock{
		PageID:    pageID,
		UserID:    userID,
		ExpiresAt: time.Now().UTC().Add(time.Duration(ttlSeconds) * time.Second),
	}
	// Acquire is atomic in the DB — no separate Get() needed (eliminates TOCTOU).
	if err := uc.repo.Acquire(ctx, lock); err != nil {
		return nil, err
	}

	if uc.hub != nil {
		uc.hub.BroadcastToRoom(ws.NewPageLockedMsg(pageID, map[string]any{
			"page_id": pageID, "user_id": userID, "expires_at": lock.ExpiresAt,
		}))
	}
	return lock, nil
}

func (uc *useCase) Release(ctx context.Context, pageID, userID string) error {
	err := uc.repo.Release(ctx, pageID, userID)
	if err == nil && uc.hub != nil {
		uc.hub.BroadcastToRoom(ws.NewPageUnlockedMsg(pageID, map[string]any{
			"page_id": pageID, "user_id": userID,
		}))
	}
	return err
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
