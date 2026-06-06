package repository

import (
	"context"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

type SpaceRepository interface {
	Create(ctx context.Context, s *entity.Space) error
	GetByID(ctx context.Context, id string) (*entity.Space, error)
	GetByKey(ctx context.Context, key string) (*entity.Space, error)
	GetByProjectID(ctx context.Context, projectID string) (*entity.Space, error)
	List(ctx context.Context, filter *entity.Filter) ([]*entity.Space, int, error)
	Update(ctx context.Context, s *entity.Space) error
	SoftDelete(ctx context.Context, id string) error
	ExistsByKey(ctx context.Context, key string) (bool, error)

	Archive(ctx context.Context, id string) error
	Restore(ctx context.Context, id string) error
	GetStatistics(ctx context.Context, spaceID string) (*entity.SpaceStatistics, error)

	AddMember(ctx context.Context, m *entity.SpaceMember) error
	GetMember(ctx context.Context, spaceID, userID string) (*entity.SpaceMember, error)
	ListMembers(ctx context.Context, spaceID string, filter *entity.Filter) ([]*entity.SpaceMember, int, error)
	UpdateMemberRole(ctx context.Context, spaceID, userID, role string) error
	RemoveMember(ctx context.Context, spaceID, userID string) error
	IsMember(ctx context.Context, spaceID, userID string) (bool, error)
}
