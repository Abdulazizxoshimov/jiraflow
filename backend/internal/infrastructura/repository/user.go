package repository

import (
	"context"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	GetByID(ctx context.Context, id string) (*entity.User, error)
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
	List(ctx context.Context, filter *entity.UserFilter) ([]*entity.User, int, error)
	Update(ctx context.Context, user *entity.User) error
	UpdatePassword(ctx context.Context, userID, passwordHash string) error
	UpdateLastLogin(ctx context.Context, userID string) error
	SoftDelete(ctx context.Context, id string) error
	ExistsByEmail(ctx context.Context, email string) (bool, error)
}
