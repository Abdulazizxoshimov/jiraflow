package repository

import (
	"context"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

type IssueTypeRepository interface {
	CreateType(ctx context.Context, t *entity.IssueType) error
	ListTypes(ctx context.Context) ([]*entity.IssueType, error)
	GetTypeByID(ctx context.Context, id string) (*entity.IssueType, error)
	DeleteType(ctx context.Context, id string) error

	CreateScheme(ctx context.Context, s *entity.IssueTypeScheme, issueTypeIDs []string) error
	GetSchemeByID(ctx context.Context, id string) (*entity.IssueTypeScheme, error)
	GetSchemeByProject(ctx context.Context, projectID string) (*entity.IssueTypeScheme, error)
	ListSchemes(ctx context.Context) ([]*entity.IssueTypeScheme, error)
	DeleteScheme(ctx context.Context, id string) error
}
