package repository

import (
	"context"
	"time"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

type VersionRepository interface {
	Create(ctx context.Context, v *entity.Version) error
	GetByID(ctx context.Context, id string) (*entity.Version, error)
	List(ctx context.Context, projectID string) ([]*entity.Version, error)
	Update(ctx context.Context, v *entity.Version) error
	Release(ctx context.Context, id string, releasedAt time.Time) error
	Archive(ctx context.Context, id string) error
	Delete(ctx context.Context, id string) error

	SetIssueVersions(ctx context.Context, issueID string, versionIDs []string) error
	SetIssueAffectsVersions(ctx context.Context, issueID string, versionIDs []string) error
	GetIssueVersions(ctx context.Context, issueID string) ([]*entity.Version, error)
	GetIssueAffectsVersions(ctx context.Context, issueID string) ([]*entity.Version, error)
	GetProgress(ctx context.Context, versionID string) (total, done int, err error)
}
