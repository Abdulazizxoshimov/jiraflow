package repository

import (
	"context"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

type ComponentRepository interface {
	Create(ctx context.Context, c *entity.Component) error
	GetByID(ctx context.Context, id string) (*entity.Component, error)
	List(ctx context.Context, projectID string) ([]*entity.Component, error)
	Update(ctx context.Context, c *entity.Component) error
	Delete(ctx context.Context, id string) error

	SetIssueComponents(ctx context.Context, issueID string, componentIDs []string) error
	GetIssueComponents(ctx context.Context, issueID string) ([]*entity.Component, error)
}
