// Package testutil provides mock implementations of repository and service interfaces for unit tests.
package testutil

import (
	"context"
	"time"

	"go.uber.org/zap/zapcore"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
	"github.com/jira-backend/jiraflow-backend/internal/infrastructure/repository"
	"github.com/jira-backend/jiraflow-backend/internal/pkg/token"
	"github.com/jira-backend/jiraflow-backend/internal/usecase/notification"
)

// ─── Nop Logger ──────────────────────────────────────────────────────────────

type NopLogger struct{}

func (NopLogger) Debug(_ context.Context, _ string, _ ...zapcore.Field) {}
func (NopLogger) Info(_ context.Context, _ string, _ ...zapcore.Field)  {}
func (NopLogger) Warn(_ context.Context, _ string, _ ...zapcore.Field)  {}
func (NopLogger) Error(_ context.Context, _ string, _ ...zapcore.Field) {}
func (NopLogger) Fatal(_ context.Context, _ string, _ ...zapcore.Field) {}

// ─── IssueRepository ─────────────────────────────────────────────────────────

type IssueRepoMock struct {
	CreateFn              func(ctx context.Context, issue *entity.Issue) error
	GetByIDFn             func(ctx context.Context, id string) (*entity.Issue, error)
	GetByKeyFn            func(ctx context.Context, key string) (*entity.Issue, error)
	ListFn                func(ctx context.Context, filter *entity.IssueFilter) ([]*entity.Issue, int, error)
	UpdateFn              func(ctx context.Context, issue *entity.Issue) error
	UpdateStatusFn        func(ctx context.Context, id, statusID string) error
	UpdateResolutionFn    func(ctx context.Context, id string, resolution *string) error
	SoftDeleteFn          func(ctx context.Context, id string) error
	CountByProjectFn      func(ctx context.Context, projectID string) (int, error)
	SetLabelsFn           func(ctx context.Context, issueID string, labelIDs []string) error
	GetLabelsFn           func(ctx context.Context, issueID string) ([]*entity.Label, error)
	AddWatcherFn          func(ctx context.Context, w *entity.IssueWatcher) error
	RemoveWatcherFn       func(ctx context.Context, issueID, userID string) error
	ListWatchersFn        func(ctx context.Context, issueID string) ([]*entity.IssueWatcher, error)
	IsWatcherFn           func(ctx context.Context, issueID, userID string) (bool, error)
	CreateHistoryFn       func(ctx context.Context, h *entity.IssueHistory) error
	ListHistoryFn         func(ctx context.Context, issueID string, filter *entity.Filter) ([]*entity.IssueHistory, int, error)
	BulkCreateFn          func(ctx context.Context, issues []*entity.Issue) error
	BulkUpdatePositionsFn func(ctx context.Context, items []entity.IssuePositionItem) error
	BulkUpdateFn          func(ctx context.Context, req *entity.BulkUpdateIssueReq) ([]string, error)
	BulkDeleteFn          func(ctx context.Context, issueIDs []string) error
	UpdateEstimatesFn     func(ctx context.Context, issueID string, original, remaining *int) error
	GetEpicProgressFn     func(ctx context.Context, epicID string) (*entity.EpicProgress, error)
	GetRoadmapFn          func(ctx context.Context, projectID string) ([]*entity.RoadmapItem, error)
	GetBacklogFn          func(ctx context.Context, projectID string, filter *entity.IssueFilter) ([]*entity.Issue, int, error)
	UpdateRankFn          func(ctx context.Context, id, rank string) error
	GetRankNeighborsFn    func(ctx context.Context, projectID string, beforeID, afterID *string) (string, string, error)
}

var _ repository.IssueRepository = (*IssueRepoMock)(nil)

func (m *IssueRepoMock) Create(ctx context.Context, issue *entity.Issue) error {
	if m.CreateFn != nil {
		return m.CreateFn(ctx, issue)
	}
	return nil
}
func (m *IssueRepoMock) GetByID(ctx context.Context, id string) (*entity.Issue, error) {
	if m.GetByIDFn != nil {
		return m.GetByIDFn(ctx, id)
	}
	return nil, nil
}
func (m *IssueRepoMock) GetByKey(ctx context.Context, key string) (*entity.Issue, error) {
	if m.GetByKeyFn != nil {
		return m.GetByKeyFn(ctx, key)
	}
	return nil, nil
}
func (m *IssueRepoMock) List(ctx context.Context, filter *entity.IssueFilter) ([]*entity.Issue, int, error) {
	if m.ListFn != nil {
		return m.ListFn(ctx, filter)
	}
	return nil, 0, nil
}
func (m *IssueRepoMock) Update(ctx context.Context, issue *entity.Issue) error {
	if m.UpdateFn != nil {
		return m.UpdateFn(ctx, issue)
	}
	return nil
}
func (m *IssueRepoMock) UpdateStatus(ctx context.Context, id, statusID string) error {
	if m.UpdateStatusFn != nil {
		return m.UpdateStatusFn(ctx, id, statusID)
	}
	return nil
}
func (m *IssueRepoMock) UpdateResolution(ctx context.Context, id string, resolution *string) error {
	if m.UpdateResolutionFn != nil {
		return m.UpdateResolutionFn(ctx, id, resolution)
	}
	return nil
}
func (m *IssueRepoMock) SoftDelete(ctx context.Context, id string) error {
	if m.SoftDeleteFn != nil {
		return m.SoftDeleteFn(ctx, id)
	}
	return nil
}
func (m *IssueRepoMock) CountByProject(ctx context.Context, projectID string) (int, error) {
	if m.CountByProjectFn != nil {
		return m.CountByProjectFn(ctx, projectID)
	}
	return 0, nil
}
func (m *IssueRepoMock) SetLabels(ctx context.Context, issueID string, labelIDs []string) error {
	if m.SetLabelsFn != nil {
		return m.SetLabelsFn(ctx, issueID, labelIDs)
	}
	return nil
}
func (m *IssueRepoMock) GetLabels(ctx context.Context, issueID string) ([]*entity.Label, error) {
	if m.GetLabelsFn != nil {
		return m.GetLabelsFn(ctx, issueID)
	}
	return nil, nil
}
func (m *IssueRepoMock) AddWatcher(ctx context.Context, w *entity.IssueWatcher) error {
	if m.AddWatcherFn != nil {
		return m.AddWatcherFn(ctx, w)
	}
	return nil
}
func (m *IssueRepoMock) RemoveWatcher(ctx context.Context, issueID, userID string) error {
	if m.RemoveWatcherFn != nil {
		return m.RemoveWatcherFn(ctx, issueID, userID)
	}
	return nil
}
func (m *IssueRepoMock) ListWatchers(ctx context.Context, issueID string) ([]*entity.IssueWatcher, error) {
	if m.ListWatchersFn != nil {
		return m.ListWatchersFn(ctx, issueID)
	}
	return nil, nil
}
func (m *IssueRepoMock) IsWatcher(ctx context.Context, issueID, userID string) (bool, error) {
	if m.IsWatcherFn != nil {
		return m.IsWatcherFn(ctx, issueID, userID)
	}
	return false, nil
}
func (m *IssueRepoMock) CreateHistory(ctx context.Context, h *entity.IssueHistory) error {
	if m.CreateHistoryFn != nil {
		return m.CreateHistoryFn(ctx, h)
	}
	return nil
}
func (m *IssueRepoMock) ListHistory(ctx context.Context, issueID string, filter *entity.Filter) ([]*entity.IssueHistory, int, error) {
	if m.ListHistoryFn != nil {
		return m.ListHistoryFn(ctx, issueID, filter)
	}
	return nil, 0, nil
}
func (m *IssueRepoMock) BulkUpdatePositions(ctx context.Context, items []entity.IssuePositionItem) error {
	if m.BulkUpdatePositionsFn != nil {
		return m.BulkUpdatePositionsFn(ctx, items)
	}
	return nil
}
func (m *IssueRepoMock) BulkCreate(ctx context.Context, issues []*entity.Issue) error {
	if m.BulkCreateFn != nil {
		return m.BulkCreateFn(ctx, issues)
	}
	return nil
}
func (m *IssueRepoMock) BulkUpdate(ctx context.Context, req *entity.BulkUpdateIssueReq) ([]string, error) {
	if m.BulkUpdateFn != nil {
		return m.BulkUpdateFn(ctx, req)
	}
	return nil, nil
}
func (m *IssueRepoMock) BulkDelete(ctx context.Context, issueIDs []string) error {
	if m.BulkDeleteFn != nil {
		return m.BulkDeleteFn(ctx, issueIDs)
	}
	return nil
}
func (m *IssueRepoMock) UpdateEstimates(ctx context.Context, issueID string, original, remaining *int) error {
	if m.UpdateEstimatesFn != nil {
		return m.UpdateEstimatesFn(ctx, issueID, original, remaining)
	}
	return nil
}
func (m *IssueRepoMock) GetEpicProgress(ctx context.Context, epicID string) (*entity.EpicProgress, error) {
	if m.GetEpicProgressFn != nil {
		return m.GetEpicProgressFn(ctx, epicID)
	}
	return nil, nil
}
func (m *IssueRepoMock) GetRoadmap(ctx context.Context, projectID string) ([]*entity.RoadmapItem, error) {
	if m.GetRoadmapFn != nil {
		return m.GetRoadmapFn(ctx, projectID)
	}
	return nil, nil
}
func (m *IssueRepoMock) GetBacklog(ctx context.Context, projectID string, filter *entity.IssueFilter) ([]*entity.Issue, int, error) {
	if m.GetBacklogFn != nil {
		return m.GetBacklogFn(ctx, projectID, filter)
	}
	return nil, 0, nil
}
func (m *IssueRepoMock) UpdateRank(ctx context.Context, id, rank string) error {
	if m.UpdateRankFn != nil {
		return m.UpdateRankFn(ctx, id, rank)
	}
	return nil
}
func (m *IssueRepoMock) GetRankNeighbors(ctx context.Context, projectID string, beforeID, afterID *string) (string, string, error) {
	if m.GetRankNeighborsFn != nil {
		return m.GetRankNeighborsFn(ctx, projectID, beforeID, afterID)
	}
	return "", "", nil
}

// ─── ProjectRepository ───────────────────────────────────────────────────────

type ProjectRepoMock struct {
	GetByIDFn               func(ctx context.Context, id string) (*entity.Project, error)
	GetByKeyFn              func(ctx context.Context, key string) (*entity.Project, error)
	CreateFn                func(ctx context.Context, p *entity.Project) error
	ListFn                  func(ctx context.Context, filter *entity.ProjectFilter) ([]*entity.Project, int, error)
	UpdateFn                func(ctx context.Context, p *entity.Project) error
	SoftDeleteFn            func(ctx context.Context, id string) error
	ExistsByKeyFn           func(ctx context.Context, key string) (bool, error)
	IncrementIssueCounterFn func(ctx context.Context, id string) (int64, error)
	AllocateIssueCountersFn func(ctx context.Context, id string, n int) (int64, error)
	GetDashboardFn          func(ctx context.Context, projectID string) (*entity.ProjectDashboard, error)
}

var _ repository.ProjectRepository = (*ProjectRepoMock)(nil)

func (m *ProjectRepoMock) GetByID(ctx context.Context, id string) (*entity.Project, error) {
	if m.GetByIDFn != nil {
		return m.GetByIDFn(ctx, id)
	}
	return nil, nil
}
func (m *ProjectRepoMock) GetByKey(ctx context.Context, key string) (*entity.Project, error) {
	if m.GetByKeyFn != nil {
		return m.GetByKeyFn(ctx, key)
	}
	return nil, nil
}
func (m *ProjectRepoMock) Create(ctx context.Context, p *entity.Project) error {
	if m.CreateFn != nil {
		return m.CreateFn(ctx, p)
	}
	return nil
}
func (m *ProjectRepoMock) List(ctx context.Context, filter *entity.ProjectFilter) ([]*entity.Project, int, error) {
	if m.ListFn != nil {
		return m.ListFn(ctx, filter)
	}
	return nil, 0, nil
}
func (m *ProjectRepoMock) Update(ctx context.Context, p *entity.Project) error {
	if m.UpdateFn != nil {
		return m.UpdateFn(ctx, p)
	}
	return nil
}
func (m *ProjectRepoMock) SoftDelete(ctx context.Context, id string) error {
	if m.SoftDeleteFn != nil {
		return m.SoftDeleteFn(ctx, id)
	}
	return nil
}
func (m *ProjectRepoMock) ExistsByKey(ctx context.Context, key string) (bool, error) {
	if m.ExistsByKeyFn != nil {
		return m.ExistsByKeyFn(ctx, key)
	}
	return false, nil
}
func (m *ProjectRepoMock) IncrementIssueCounter(ctx context.Context, id string) (int64, error) {
	if m.IncrementIssueCounterFn != nil {
		return m.IncrementIssueCounterFn(ctx, id)
	}
	return 1, nil
}
func (m *ProjectRepoMock) AllocateIssueCounters(ctx context.Context, id string, n int) (int64, error) {
	if m.AllocateIssueCountersFn != nil {
		return m.AllocateIssueCountersFn(ctx, id, n)
	}
	return 1, nil
}
func (m *ProjectRepoMock) GetDashboard(ctx context.Context, projectID string) (*entity.ProjectDashboard, error) {
	if m.GetDashboardFn != nil {
		return m.GetDashboardFn(ctx, projectID)
	}
	return nil, nil
}

// ─── WorkflowRepository ──────────────────────────────────────────────────────

type WorkflowRepoMock struct {
	GetWithDetailsFn      func(ctx context.Context, id string) (*entity.Workflow, error)
	GetByIDFn             func(ctx context.Context, id string) (*entity.Workflow, error)
	CreateFn              func(ctx context.Context, wf *entity.Workflow) error
	ListFn                func(ctx context.Context, filter *entity.Filter) ([]*entity.Workflow, int, error)
	UpdateFn              func(ctx context.Context, wf *entity.Workflow) error
	SoftDeleteFn          func(ctx context.Context, id string) error
	SetDefaultFn          func(ctx context.Context, id string) error
	GetDefaultFn          func(ctx context.Context) (*entity.Workflow, error)
	CreateStatusFn        func(ctx context.Context, s *entity.WorkflowStatus) error
	GetStatusByIDFn       func(ctx context.Context, id string) (*entity.WorkflowStatus, error)
	ListStatusesFn        func(ctx context.Context, workflowID string) ([]*entity.WorkflowStatus, error)
	UpdateStatusFn        func(ctx context.Context, s *entity.WorkflowStatus) error
	DeleteStatusFn        func(ctx context.Context, id string) error
	CreateTransitionFn    func(ctx context.Context, t *entity.WorkflowTransition) error
	GetTransitionByIDFn   func(ctx context.Context, id string) (*entity.WorkflowTransition, error)
	ListTransitionsFn     func(ctx context.Context, workflowID string) ([]*entity.WorkflowTransition, error)
	DeleteTransitionFn    func(ctx context.Context, id string) error
	IsTransitionAllowedFn func(ctx context.Context, workflowID, fromStatusID, toStatusID string) (bool, error)
}

var _ repository.WorkflowRepository = (*WorkflowRepoMock)(nil)

func (m *WorkflowRepoMock) GetWithDetails(ctx context.Context, id string) (*entity.Workflow, error) {
	if m.GetWithDetailsFn != nil {
		return m.GetWithDetailsFn(ctx, id)
	}
	return nil, nil
}
func (m *WorkflowRepoMock) GetByID(ctx context.Context, id string) (*entity.Workflow, error) {
	if m.GetByIDFn != nil {
		return m.GetByIDFn(ctx, id)
	}
	return nil, nil
}
func (m *WorkflowRepoMock) Create(ctx context.Context, wf *entity.Workflow) error {
	if m.CreateFn != nil {
		return m.CreateFn(ctx, wf)
	}
	return nil
}
func (m *WorkflowRepoMock) List(ctx context.Context, filter *entity.Filter) ([]*entity.Workflow, int, error) {
	if m.ListFn != nil {
		return m.ListFn(ctx, filter)
	}
	return nil, 0, nil
}
func (m *WorkflowRepoMock) Update(ctx context.Context, wf *entity.Workflow) error {
	if m.UpdateFn != nil {
		return m.UpdateFn(ctx, wf)
	}
	return nil
}
func (m *WorkflowRepoMock) SoftDelete(ctx context.Context, id string) error {
	if m.SoftDeleteFn != nil {
		return m.SoftDeleteFn(ctx, id)
	}
	return nil
}
func (m *WorkflowRepoMock) SetDefault(ctx context.Context, id string) error {
	if m.SetDefaultFn != nil {
		return m.SetDefaultFn(ctx, id)
	}
	return nil
}
func (m *WorkflowRepoMock) GetDefault(ctx context.Context) (*entity.Workflow, error) {
	if m.GetDefaultFn != nil {
		return m.GetDefaultFn(ctx)
	}
	return nil, nil
}
func (m *WorkflowRepoMock) CreateStatus(ctx context.Context, s *entity.WorkflowStatus) error {
	if m.CreateStatusFn != nil {
		return m.CreateStatusFn(ctx, s)
	}
	return nil
}
func (m *WorkflowRepoMock) GetStatusByID(ctx context.Context, id string) (*entity.WorkflowStatus, error) {
	if m.GetStatusByIDFn != nil {
		return m.GetStatusByIDFn(ctx, id)
	}
	return nil, nil
}
func (m *WorkflowRepoMock) ListStatuses(ctx context.Context, workflowID string) ([]*entity.WorkflowStatus, error) {
	if m.ListStatusesFn != nil {
		return m.ListStatusesFn(ctx, workflowID)
	}
	return nil, nil
}
func (m *WorkflowRepoMock) UpdateStatus(ctx context.Context, s *entity.WorkflowStatus) error {
	if m.UpdateStatusFn != nil {
		return m.UpdateStatusFn(ctx, s)
	}
	return nil
}
func (m *WorkflowRepoMock) DeleteStatus(ctx context.Context, id string) error {
	if m.DeleteStatusFn != nil {
		return m.DeleteStatusFn(ctx, id)
	}
	return nil
}
func (m *WorkflowRepoMock) CreateTransition(ctx context.Context, t *entity.WorkflowTransition) error {
	if m.CreateTransitionFn != nil {
		return m.CreateTransitionFn(ctx, t)
	}
	return nil
}
func (m *WorkflowRepoMock) GetTransitionByID(ctx context.Context, id string) (*entity.WorkflowTransition, error) {
	if m.GetTransitionByIDFn != nil {
		return m.GetTransitionByIDFn(ctx, id)
	}
	return nil, nil
}
func (m *WorkflowRepoMock) ListTransitions(ctx context.Context, workflowID string) ([]*entity.WorkflowTransition, error) {
	if m.ListTransitionsFn != nil {
		return m.ListTransitionsFn(ctx, workflowID)
	}
	return nil, nil
}
func (m *WorkflowRepoMock) DeleteTransition(ctx context.Context, id string) error {
	if m.DeleteTransitionFn != nil {
		return m.DeleteTransitionFn(ctx, id)
	}
	return nil
}
func (m *WorkflowRepoMock) IsTransitionAllowed(ctx context.Context, workflowID, fromStatusID, toStatusID string) (bool, error) {
	if m.IsTransitionAllowedFn != nil {
		return m.IsTransitionAllowedFn(ctx, workflowID, fromStatusID, toStatusID)
	}
	return true, nil
}

// ─── SprintRepository ────────────────────────────────────────────────────────

type SprintRepoMock struct {
	CreateFn      func(ctx context.Context, s *entity.Sprint) error
	GetByIDFn     func(ctx context.Context, id string) (*entity.Sprint, error)
	ListFn        func(ctx context.Context, projectID string, filter *entity.SprintFilter) ([]*entity.Sprint, int, error)
	UpdateFn      func(ctx context.Context, s *entity.Sprint) error
	SoftDeleteFn  func(ctx context.Context, id string) error
	GetActiveFn   func(ctx context.Context, projectID string) (*entity.Sprint, error)
	StartFn       func(ctx context.Context, id string, startedAt time.Time) error
	CompleteFn    func(ctx context.Context, id string, completedAt time.Time) error
	AddIssueFn    func(ctx context.Context, sprintID, issueID string) error
	RemoveIssueFn func(ctx context.Context, sprintID, issueID string) error
	GetReportFn   func(ctx context.Context, sprintID string) (*entity.SprintReport, error)
	GetBurndownFn func(ctx context.Context, sprintID string) (*entity.BurndownChart, error)
	GetBurnupFn   func(ctx context.Context, sprintID string) (*entity.BurnupChart, error)
	GetVelocityFn func(ctx context.Context, projectID string, limit int) (*entity.VelocityReport, error)
	GetCFDFn      func(ctx context.Context, projectID string, from, to *string) (*entity.CFDChart, error)
}

var _ repository.SprintRepository = (*SprintRepoMock)(nil)

func (m *SprintRepoMock) Create(ctx context.Context, s *entity.Sprint) error {
	if m.CreateFn != nil {
		return m.CreateFn(ctx, s)
	}
	return nil
}
func (m *SprintRepoMock) GetByID(ctx context.Context, id string) (*entity.Sprint, error) {
	if m.GetByIDFn != nil {
		return m.GetByIDFn(ctx, id)
	}
	return nil, nil
}
func (m *SprintRepoMock) List(ctx context.Context, projectID string, filter *entity.SprintFilter) ([]*entity.Sprint, int, error) {
	if m.ListFn != nil {
		return m.ListFn(ctx, projectID, filter)
	}
	return nil, 0, nil
}
func (m *SprintRepoMock) Update(ctx context.Context, s *entity.Sprint) error {
	if m.UpdateFn != nil {
		return m.UpdateFn(ctx, s)
	}
	return nil
}
func (m *SprintRepoMock) SoftDelete(ctx context.Context, id string) error {
	if m.SoftDeleteFn != nil {
		return m.SoftDeleteFn(ctx, id)
	}
	return nil
}
func (m *SprintRepoMock) GetActive(ctx context.Context, projectID string) (*entity.Sprint, error) {
	if m.GetActiveFn != nil {
		return m.GetActiveFn(ctx, projectID)
	}
	return nil, nil
}
func (m *SprintRepoMock) Start(ctx context.Context, id string, startedAt time.Time) error {
	if m.StartFn != nil {
		return m.StartFn(ctx, id, startedAt)
	}
	return nil
}
func (m *SprintRepoMock) Complete(ctx context.Context, id string, completedAt time.Time) error {
	if m.CompleteFn != nil {
		return m.CompleteFn(ctx, id, completedAt)
	}
	return nil
}
func (m *SprintRepoMock) AddIssue(ctx context.Context, sprintID, issueID string) error {
	if m.AddIssueFn != nil {
		return m.AddIssueFn(ctx, sprintID, issueID)
	}
	return nil
}
func (m *SprintRepoMock) BulkAddIssues(ctx context.Context, sprintID string, issueIDs []string) error {
	return nil
}
func (m *SprintRepoMock) RemoveIssue(ctx context.Context, sprintID, issueID string) error {
	if m.RemoveIssueFn != nil {
		return m.RemoveIssueFn(ctx, sprintID, issueID)
	}
	return nil
}
func (m *SprintRepoMock) GetReport(ctx context.Context, sprintID string) (*entity.SprintReport, error) {
	if m.GetReportFn != nil {
		return m.GetReportFn(ctx, sprintID)
	}
	return nil, nil
}
func (m *SprintRepoMock) GetBurndown(ctx context.Context, sprintID string) (*entity.BurndownChart, error) {
	if m.GetBurndownFn != nil {
		return m.GetBurndownFn(ctx, sprintID)
	}
	return nil, nil
}
func (m *SprintRepoMock) GetBurnup(ctx context.Context, sprintID string) (*entity.BurnupChart, error) {
	if m.GetBurnupFn != nil {
		return m.GetBurnupFn(ctx, sprintID)
	}
	return nil, nil
}
func (m *SprintRepoMock) GetVelocity(ctx context.Context, projectID string, limit int) (*entity.VelocityReport, error) {
	if m.GetVelocityFn != nil {
		return m.GetVelocityFn(ctx, projectID, limit)
	}
	return nil, nil
}
func (m *SprintRepoMock) GetCFD(ctx context.Context, projectID string, from, to *string) (*entity.CFDChart, error) {
	if m.GetCFDFn != nil {
		return m.GetCFDFn(ctx, projectID, from, to)
	}
	return nil, nil
}

// ─── VersionRepository ───────────────────────────────────────────────────────

type VersionRepoMock struct{}

var _ repository.VersionRepository = (*VersionRepoMock)(nil)

func (VersionRepoMock) Create(_ context.Context, _ *entity.Version) error              { return nil }
func (VersionRepoMock) GetByID(_ context.Context, _ string) (*entity.Version, error)   { return nil, nil }
func (VersionRepoMock) List(_ context.Context, _ string) ([]*entity.Version, error)    { return nil, nil }
func (VersionRepoMock) Update(_ context.Context, _ *entity.Version) error              { return nil }
func (VersionRepoMock) Release(_ context.Context, _ string, _ time.Time) error         { return nil }
func (VersionRepoMock) Archive(_ context.Context, _ string) error                      { return nil }
func (VersionRepoMock) Delete(_ context.Context, _ string) error                       { return nil }
func (VersionRepoMock) SetIssueVersions(_ context.Context, _ string, _ []string) error { return nil }
func (VersionRepoMock) SetIssueAffectsVersions(_ context.Context, _ string, _ []string) error {
	return nil
}
func (VersionRepoMock) GetIssueVersions(_ context.Context, _ string) ([]*entity.Version, error) {
	return nil, nil
}
func (VersionRepoMock) GetIssueAffectsVersions(_ context.Context, _ string) ([]*entity.Version, error) {
	return nil, nil
}
func (VersionRepoMock) GetProgress(_ context.Context, _ string) (int, int, error) { return 0, 0, nil }

// ─── ProjectMemberRepository ─────────────────────────────────────────────────

type ProjectMemberRepoMock struct {
	GetMemberFn func(ctx context.Context, projectID, userID string) (*entity.ProjectMember, error)
}

func (m ProjectMemberRepoMock) Add(_ context.Context, _ *entity.ProjectMember) error { return nil }
func (m ProjectMemberRepoMock) GetMember(ctx context.Context, projectID, userID string) (*entity.ProjectMember, error) {
	if m.GetMemberFn != nil {
		return m.GetMemberFn(ctx, projectID, userID)
	}
	return &entity.ProjectMember{Role: "admin"}, nil
}
func (m ProjectMemberRepoMock) ListByProject(_ context.Context, _ string, _ *entity.Filter) ([]*entity.ProjectMember, int, error) {
	return nil, 0, nil
}
func (m ProjectMemberRepoMock) ListByUser(_ context.Context, _ string) ([]*entity.ProjectMember, error) {
	return nil, nil
}
func (m ProjectMemberRepoMock) UpdateRole(_ context.Context, _, _, _ string) error { return nil }
func (m ProjectMemberRepoMock) Remove(_ context.Context, _, _ string) error        { return nil }
func (m ProjectMemberRepoMock) IsMember(_ context.Context, _, _ string) (bool, error) {
	return true, nil
}

// ─── UserRepository ──────────────────────────────────────────────────────────

type UserRepoMock struct {
	CreateFn          func(ctx context.Context, user *entity.User) error
	GetByIDFn         func(ctx context.Context, id string) (*entity.User, error)
	GetByEmailFn      func(ctx context.Context, email string) (*entity.User, error)
	UpdateLastLoginFn func(ctx context.Context, userID string) error
	UpdatePasswordFn  func(ctx context.Context, userID, hash string) error
}

var _ repository.UserRepository = (*UserRepoMock)(nil)

func (m *UserRepoMock) Create(ctx context.Context, user *entity.User) error {
	if m.CreateFn != nil {
		return m.CreateFn(ctx, user)
	}
	return nil
}
func (m *UserRepoMock) GetByID(ctx context.Context, id string) (*entity.User, error) {
	if m.GetByIDFn != nil {
		return m.GetByIDFn(ctx, id)
	}
	return &entity.User{ID: id, IsActive: true, Role: "member"}, nil
}
func (m *UserRepoMock) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	if m.GetByEmailFn != nil {
		return m.GetByEmailFn(ctx, email)
	}
	return &entity.User{ID: "user-1", Email: email, IsActive: true, Role: "member"}, nil
}
func (m *UserRepoMock) List(_ context.Context, _ *entity.UserFilter) ([]*entity.User, int, error) {
	return nil, 0, nil
}
func (m *UserRepoMock) Update(_ context.Context, _ *entity.User) error { return nil }
func (m *UserRepoMock) UpdatePassword(ctx context.Context, userID, hash string) error {
	if m.UpdatePasswordFn != nil {
		return m.UpdatePasswordFn(ctx, userID, hash)
	}
	return nil
}
func (m *UserRepoMock) UpdateLastLogin(ctx context.Context, userID string) error {
	if m.UpdateLastLoginFn != nil {
		return m.UpdateLastLoginFn(ctx, userID)
	}
	return nil
}
func (m *UserRepoMock) SoftDelete(_ context.Context, _ string) error            { return nil }
func (m *UserRepoMock) ExistsByEmail(_ context.Context, _ string) (bool, error) { return false, nil }

// ─── AuthRepository ───────────────────────────────────────────────────────────

type AuthRepoMock struct {
	CreateRefreshTokenFn     func(ctx context.Context, rt *entity.RefreshToken) error
	GetRefreshTokenByHashFn  func(ctx context.Context, hash string) (*entity.RefreshToken, error)
	RevokeRefreshTokenFn     func(ctx context.Context, id string) error
	CreatePasswordResetFn    func(ctx context.Context, pr *entity.PasswordReset) error
	GetPasswordResetByHashFn func(ctx context.Context, hash string) (*entity.PasswordReset, error)
	MarkPasswordResetUsedFn  func(ctx context.Context, id string) error
}

var _ repository.AuthRepository = (*AuthRepoMock)(nil)

func (m *AuthRepoMock) CreateRefreshToken(ctx context.Context, rt *entity.RefreshToken) error {
	if m.CreateRefreshTokenFn != nil {
		return m.CreateRefreshTokenFn(ctx, rt)
	}
	return nil
}
func (m *AuthRepoMock) GetRefreshTokenByHash(ctx context.Context, hash string) (*entity.RefreshToken, error) {
	if m.GetRefreshTokenByHashFn != nil {
		return m.GetRefreshTokenByHashFn(ctx, hash)
	}
	return &entity.RefreshToken{ID: "rt-1", UserID: "user-1"}, nil
}
func (m *AuthRepoMock) RevokeRefreshToken(ctx context.Context, id string) error {
	if m.RevokeRefreshTokenFn != nil {
		return m.RevokeRefreshTokenFn(ctx, id)
	}
	return nil
}
func (m *AuthRepoMock) RevokeAllUserTokens(_ context.Context, _ string) error { return nil }
func (m *AuthRepoMock) DeleteExpiredTokens(_ context.Context) error           { return nil }
func (m *AuthRepoMock) CreatePasswordReset(ctx context.Context, pr *entity.PasswordReset) error {
	if m.CreatePasswordResetFn != nil {
		return m.CreatePasswordResetFn(ctx, pr)
	}
	return nil
}
func (m *AuthRepoMock) GetPasswordResetByHash(ctx context.Context, hash string) (*entity.PasswordReset, error) {
	if m.GetPasswordResetByHashFn != nil {
		return m.GetPasswordResetByHashFn(ctx, hash)
	}
	return nil, nil
}
func (m *AuthRepoMock) MarkPasswordResetUsed(ctx context.Context, id string) error {
	if m.MarkPasswordResetUsedFn != nil {
		return m.MarkPasswordResetUsedFn(ctx, id)
	}
	return nil
}
func (m *AuthRepoMock) DeleteExpiredPasswordResets(_ context.Context) error { return nil }

// ─── TokenMakerMock ───────────────────────────────────────────────────────────

type TokenMakerMock struct {
	GenerateFn      func(ctx context.Context, sub, sessionID, role string) (string, string, error)
	RotateFn        func(ctx context.Context, refreshToken string) (string, string, error)
	RevokeSessionFn func(ctx context.Context, sessionID string) error
}

var _ token.Maker = (*TokenMakerMock)(nil)

func (m *TokenMakerMock) Generate(ctx context.Context, sub, sessionID, role string) (string, string, error) {
	if m.GenerateFn != nil {
		return m.GenerateFn(ctx, sub, sessionID, role)
	}
	return "access-token", "refresh-token", nil
}
func (m *TokenMakerMock) ValidateAccess(_ context.Context, _ string) (*token.Claims, error) {
	return &token.Claims{}, nil
}
func (m *TokenMakerMock) Rotate(ctx context.Context, refreshToken string) (string, string, error) {
	if m.RotateFn != nil {
		return m.RotateFn(ctx, refreshToken)
	}
	return "new-access", "new-refresh", nil
}
func (m *TokenMakerMock) Revoke(_ context.Context, _ string) error { return nil }
func (m *TokenMakerMock) StoreSession(_ context.Context, _, _ string, _ time.Duration) error {
	return nil
}
func (m *TokenMakerMock) RevokeSession(ctx context.Context, sessionID string) error {
	if m.RevokeSessionFn != nil {
		return m.RevokeSessionFn(ctx, sessionID)
	}
	return nil
}

// ─── HasherMock ───────────────────────────────────────────────────────────────

type HasherMock struct {
	HashFn  func(password string) (string, error)
	CheckFn func(plain, hashed string) bool
}

func (m *HasherMock) Hash(password string) (string, error) {
	if m.HashFn != nil {
		return m.HashFn(password)
	}
	return "hashed:" + password, nil
}
func (m *HasherMock) Check(plain, hashed string) bool {
	if m.CheckFn != nil {
		return m.CheckFn(plain, hashed)
	}
	return hashed == "hashed:"+plain
}

// ─── EmailSenderMock ──────────────────────────────────────────────────────────

type EmailSenderMock struct {
	SendFn    func(ctx context.Context, to []string, subject, templateName string, data any) error
	SendRawFn func(ctx context.Context, to []string, subject, body string) error
}

func (m *EmailSenderMock) Send(ctx context.Context, to []string, subject, templateName string, data any) error {
	if m.SendFn != nil {
		return m.SendFn(ctx, to, subject, templateName, data)
	}
	return nil
}

func (m *EmailSenderMock) SendRaw(ctx context.Context, to []string, subject, body string) error {
	if m.SendRawFn != nil {
		return m.SendRawFn(ctx, to, subject, body)
	}
	return nil
}

// ─── Dispatcher ──────────────────────────────────────────────────────────────

type DispatcherMock struct{}

var _ notification.Dispatcher = (*DispatcherMock)(nil)

func (DispatcherMock) IssueAssigned(_ context.Context, _ *entity.Issue, _, _ string)               {}
func (DispatcherMock) IssueCreated(_ context.Context, _ *entity.Issue, _ string)                   {}
func (DispatcherMock) IssueUpdated(_ context.Context, _ *entity.Issue, _ []string, _ string)       {}
func (DispatcherMock) IssueCommented(_ context.Context, _ string, _ []string, _ string)            {}
func (DispatcherMock) IssueMentioned(_ context.Context, _ string, _ []string, _ string)            {}
func (DispatcherMock) IssueStatusChanged(_ context.Context, _ *entity.Issue, _ []string, _ string) {}
func (DispatcherMock) PageCommented(_ context.Context, _ string, _ []string, _ string)             {}
func (DispatcherMock) PageMentioned(_ context.Context, _ string, _ []string, _ string)             {}
func (DispatcherMock) SprintStarted(_ context.Context, _, _, _ string)                             {}
func (DispatcherMock) SprintCompleted(_ context.Context, _, _, _ string)                           {}
