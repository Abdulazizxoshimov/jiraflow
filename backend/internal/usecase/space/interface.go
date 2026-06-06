package space

import (
	"context"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

type UseCase interface {
	Create(ctx context.Context, s *entity.Space, leadID string) (*entity.Space, error)
	GetByID(ctx context.Context, id string) (*entity.Space, error)
	GetByKey(ctx context.Context, key string) (*entity.Space, error)
	List(ctx context.Context, filter *entity.Filter) ([]*entity.Space, int, error)
	Update(ctx context.Context, id string, s *entity.Space, actorID string) (*entity.Space, error)
	Delete(ctx context.Context, id, actorID string) error
	Archive(ctx context.Context, id, actorID string) error
	Restore(ctx context.Context, id, actorID string) error
	GetStatistics(ctx context.Context, id string) (*entity.SpaceStatistics, error)

	AddMember(ctx context.Context, spaceID string, m *entity.SpaceMember) error
	UpdateMemberRole(ctx context.Context, spaceID, userID, role string) error
	RemoveMember(ctx context.Context, spaceID, userID string) error
	ListMembers(ctx context.Context, spaceID string, filter *entity.Filter) ([]*entity.SpaceMember, int, error)
}
