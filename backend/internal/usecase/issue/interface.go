package issue

import (
	"context"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

type UseCase interface {
	Create(ctx context.Context, projectID string, req *entity.CreateIssueReq, reporterID string) (*entity.Issue, error)
	GetByID(ctx context.Context, id string) (*entity.Issue, error)
	GetByKey(ctx context.Context, key string) (*entity.Issue, error)
	List(ctx context.Context, filter *entity.IssueFilter) ([]*entity.Issue, int, error)
	Update(ctx context.Context, id string, req *entity.UpdateIssueReq, actorID string) (*entity.Issue, error)
	Transition(ctx context.Context, id, statusID, actorID string) (*entity.Issue, error)
	Delete(ctx context.Context, id, actorID string) error
	AddWatcher(ctx context.Context, issueID, userID string) error
	RemoveWatcher(ctx context.Context, issueID, userID string) error
	ListWatchers(ctx context.Context, issueID string) ([]*entity.IssueWatcher, error)
	ListHistory(ctx context.Context, issueID string, filter *entity.Filter) ([]*entity.IssueHistory, int, error)
	ReorderIssues(ctx context.Context, req *entity.ReorderIssuesReq) error

	BulkUpdate(ctx context.Context, req *entity.BulkUpdateIssueReq, actorID string) (*entity.BulkResult, error)
	BulkDelete(ctx context.Context, req *entity.BulkDeleteIssueReq, actorID string) error
	GetEpicProgress(ctx context.Context, epicID string) (*entity.EpicProgress, error)
	GetRoadmap(ctx context.Context, projectID string) ([]*entity.RoadmapItem, error)
	GetBacklog(ctx context.Context, projectID string, filter *entity.IssueFilter) ([]*entity.Issue, int, error)
	Clone(ctx context.Context, id, reporterID string, req *entity.CloneIssueReq) (*entity.Issue, error)
	RankBetween(ctx context.Context, id string, req *entity.RankIssueReq) error
	MoveOnBoard(ctx context.Context, issueID string, req *entity.MoveIssueReq, actorID string) (*entity.Issue, error)
}
