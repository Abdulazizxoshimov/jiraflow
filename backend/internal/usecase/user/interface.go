package user

import (
	"context"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

type UseCase interface {
	Create(ctx context.Context, req *entity.CreateUserReq) (*entity.User, error)
	GetByID(ctx context.Context, id string) (*entity.User, error)
	List(ctx context.Context, filter *entity.UserFilter) ([]*entity.User, int, error)
	Update(ctx context.Context, id string, req *entity.UpdateUserReq) (*entity.User, error)
	ChangePassword(ctx context.Context, id string, req *entity.ChangePasswordReq) error
	Deactivate(ctx context.Context, id string) error
	Activate(ctx context.Context, id string) error
	Delete(ctx context.Context, id string) error
}
