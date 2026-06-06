package repository

import (
	"context"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

type ProjectMemberRepository interface {
	Add(ctx context.Context, m *entity.ProjectMember) error
	GetMember(ctx context.Context, projectID, userID string) (*entity.ProjectMember, error)
	ListByProject(ctx context.Context, projectID string, filter *entity.Filter) ([]*entity.ProjectMember, int, error)
	ListByUser(ctx context.Context, userID string) ([]*entity.ProjectMember, error)
	UpdateRole(ctx context.Context, projectID, userID, role string) error
	Remove(ctx context.Context, projectID, userID string) error
	IsMember(ctx context.Context, projectID, userID string) (bool, error)
}
