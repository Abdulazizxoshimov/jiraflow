package project_member

import (
	"context"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

type UseCase interface {
	Add(ctx context.Context, projectID string, req *entity.AddProjectMemberReq, actorID string) error
	UpdateRole(ctx context.Context, projectID, userID string, req *entity.UpdateProjectMemberRoleReq, actorID string) error
	Remove(ctx context.Context, projectID, userID, actorID string) error
	ListByProject(ctx context.Context, projectID string, filter *entity.Filter) ([]*entity.ProjectMember, int, error)
	IsMember(ctx context.Context, projectID, userID string) (bool, error)
}
