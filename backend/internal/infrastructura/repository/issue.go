package repository

import (
	"context"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

type IssueRepository interface {
	Create(ctx context.Context, issue *entity.Issue) error
	GetByID(ctx context.Context, id string) (*entity.Issue, error)
	GetByKey(ctx context.Context, key string) (*entity.Issue, error) // key = "PROJ-42"
	List(ctx context.Context, filter *entity.IssueFilter) ([]*entity.Issue, int, error)
	Update(ctx context.Context, issue *entity.Issue) error
	UpdateStatus(ctx context.Context, id, statusID string) error
	UpdateResolution(ctx context.Context, id string, resolution *string) error
	SoftDelete(ctx context.Context, id string) error
	CountByProject(ctx context.Context, projectID string) (int, error)

	SetLabels(ctx context.Context, issueID string, labelIDs []string) error
	GetLabels(ctx context.Context, issueID string) ([]*entity.Label, error)

	AddWatcher(ctx context.Context, w *entity.IssueWatcher) error
	RemoveWatcher(ctx context.Context, issueID, userID string) error
	ListWatchers(ctx context.Context, issueID string) ([]*entity.IssueWatcher, error)
	IsWatcher(ctx context.Context, issueID, userID string) (bool, error)

	CreateHistory(ctx context.Context, h *entity.IssueHistory) error
	ListHistory(ctx context.Context, issueID string, filter *entity.Filter) ([]*entity.IssueHistory, int, error)

	BulkUpdatePositions(ctx context.Context, items []entity.IssuePositionItem) error
	BulkUpdate(ctx context.Context, req *entity.BulkUpdateIssueReq) ([]string, error)
	BulkDelete(ctx context.Context, issueIDs []string) error

	UpdateEstimates(ctx context.Context, issueID string, original, remaining *int) error
	GetEpicProgress(ctx context.Context, epicID string) (*entity.EpicProgress, error)
	GetRoadmap(ctx context.Context, projectID string) ([]*entity.RoadmapItem, error)
	GetBacklog(ctx context.Context, projectID string, filter *entity.IssueFilter) ([]*entity.Issue, int, error)

	UpdateRank(ctx context.Context, id, rank string) error
	GetRankNeighbors(ctx context.Context, projectID string, beforeID, afterID *string) (loRank, hiRank string, err error)
}
