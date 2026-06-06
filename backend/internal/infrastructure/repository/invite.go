package repository

import (
	"context"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

type InviteRepository interface {
	Create(ctx context.Context, invite *entity.Invite) error
	GetByID(ctx context.Context, id string) (*entity.Invite, error)
	GetByTokenHash(ctx context.Context, hash string) (*entity.Invite, error)
	GetPendingByEmail(ctx context.Context, email string) (*entity.Invite, error)
	List(ctx context.Context, filter *entity.Filter) ([]*entity.Invite, int, error)
	MarkAccepted(ctx context.Context, id, acceptedByUserID string) error
	Delete(ctx context.Context, id string) error
	DeleteExpired(ctx context.Context) error
}
