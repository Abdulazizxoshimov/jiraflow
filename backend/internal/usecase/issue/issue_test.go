package issue_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
	apperr "github.com/jira-backend/jiraflow-backend/internal/pkg/errors"
	"github.com/jira-backend/jiraflow-backend/internal/testutil"
	"github.com/jira-backend/jiraflow-backend/internal/usecase/issue"
)

func newUC(issueRepo *testutil.IssueRepoMock, projRepo *testutil.ProjectRepoMock, wfRepo *testutil.WorkflowRepoMock) issue.UseCase {
	return issue.New(issueRepo, projRepo, wfRepo, testutil.VersionRepoMock{}, testutil.ProjectMemberRepoMock{}, testutil.DispatcherMock{}, testutil.NopLogger{})
}

// ─── Create ──────────────────────────────────────────────────────────────────

func TestCreate_Success(t *testing.T) {
	ctx := context.Background()

	issueRepo := &testutil.IssueRepoMock{}
	projRepo := &testutil.ProjectRepoMock{
		GetByIDFn: func(_ context.Context, id string) (*entity.Project, error) {
			return &entity.Project{ID: id, WorkflowID: "wf-1"}, nil
		},
		IncrementIssueCounterFn: func(_ context.Context, _ string) (int64, error) {
			return 42, nil
		},
	}
	wfRepo := &testutil.WorkflowRepoMock{
		GetWithDetailsFn: func(_ context.Context, _ string) (*entity.Workflow, error) {
			return &entity.Workflow{
				ID:       "wf-1",
				Statuses: []entity.WorkflowStatus{{ID: "status-1", IsInitial: true}},
			}, nil
		},
	}

	uc := newUC(issueRepo, projRepo, wfRepo)
	result, err := uc.Create(ctx, "proj-1", &entity.CreateIssueReq{
		Title: "Fix login bug",
		Type:  "bug",
	}, "user-1")

	require.NoError(t, err)
	assert.Equal(t, "Fix login bug", result.Title)
	assert.Equal(t, "bug", result.Type)
	assert.Equal(t, "proj-1", result.ProjectID)
	assert.Equal(t, "status-1", result.StatusID)
	assert.Equal(t, 42, result.IssueNumber)
}

func TestCreate_DefaultPriority(t *testing.T) {
	ctx := context.Background()

	projRepo := &testutil.ProjectRepoMock{
		GetByIDFn: func(_ context.Context, id string) (*entity.Project, error) {
			return &entity.Project{ID: id, WorkflowID: "wf-1"}, nil
		},
	}
	wfRepo := &testutil.WorkflowRepoMock{
		GetWithDetailsFn: func(_ context.Context, _ string) (*entity.Workflow, error) {
			return &entity.Workflow{
				ID:       "wf-1",
				Statuses: []entity.WorkflowStatus{{ID: "status-1", IsInitial: true}},
			}, nil
		},
	}

	uc := newUC(&testutil.IssueRepoMock{}, projRepo, wfRepo)
	result, err := uc.Create(ctx, "proj-1", &entity.CreateIssueReq{Title: "Task"}, "user-1")

	require.NoError(t, err)
	assert.Equal(t, "medium", result.Priority)
}

func TestCreate_ProjectNotFound(t *testing.T) {
	ctx := context.Background()

	projRepo := &testutil.ProjectRepoMock{
		GetByIDFn: func(_ context.Context, _ string) (*entity.Project, error) {
			return nil, apperr.NotFound("project not found")
		},
	}

	uc := newUC(&testutil.IssueRepoMock{}, projRepo, &testutil.WorkflowRepoMock{})
	_, err := uc.Create(ctx, "proj-x", &entity.CreateIssueReq{Title: "Task"}, "user-1")

	require.Error(t, err)
	var appErr *apperr.AppError
	assert.True(t, errors.As(err, &appErr))
	assert.Equal(t, apperr.CodeNotFound, appErr.Code)
}

func TestCreate_WorkflowError(t *testing.T) {
	ctx := context.Background()

	projRepo := &testutil.ProjectRepoMock{
		GetByIDFn: func(_ context.Context, id string) (*entity.Project, error) {
			return &entity.Project{ID: id, WorkflowID: "wf-bad"}, nil
		},
	}
	wfRepo := &testutil.WorkflowRepoMock{
		GetWithDetailsFn: func(_ context.Context, _ string) (*entity.Workflow, error) {
			return nil, errors.New("db error")
		},
	}

	uc := newUC(&testutil.IssueRepoMock{}, projRepo, wfRepo)
	_, err := uc.Create(ctx, "proj-1", &entity.CreateIssueReq{Title: "Task"}, "user-1")

	require.Error(t, err)
}

// ─── Clone ───────────────────────────────────────────────────────────────────

func TestClone_Success(t *testing.T) {
	ctx := context.Background()

	src := &entity.Issue{
		ID:        "issue-src",
		ProjectID: "proj-1",
		Title:     "Original issue",
		Type:      "task",
		Priority:  "high",
		StatusID:  "status-1",
	}

	issueRepo := &testutil.IssueRepoMock{
		GetByIDFn: func(_ context.Context, id string) (*entity.Issue, error) {
			if id == "issue-src" {
				return src, nil
			}
			return &entity.Issue{ID: id}, nil
		},
	}
	projRepo := &testutil.ProjectRepoMock{
		IncrementIssueCounterFn: func(_ context.Context, _ string) (int64, error) {
			return 5, nil
		},
	}

	uc := newUC(issueRepo, projRepo, &testutil.WorkflowRepoMock{})
	clone, err := uc.Clone(ctx, "issue-src", "user-1", &entity.CloneIssueReq{})

	require.NoError(t, err)
	assert.NotEqual(t, src.ID, clone.ID)
	assert.Equal(t, "Original issue (copy)", clone.Title)
	assert.Equal(t, src.ProjectID, clone.ProjectID)
	assert.Equal(t, src.Type, clone.Type)
	assert.Equal(t, src.Priority, clone.Priority)
}

func TestClone_CustomTitle(t *testing.T) {
	ctx := context.Background()

	src := &entity.Issue{ID: "issue-src", ProjectID: "proj-1", Title: "Original"}
	issueRepo := &testutil.IssueRepoMock{
		GetByIDFn: func(_ context.Context, _ string) (*entity.Issue, error) { return src, nil },
	}
	projRepo := &testutil.ProjectRepoMock{
		IncrementIssueCounterFn: func(_ context.Context, _ string) (int64, error) { return 2, nil },
	}

	customTitle := "My custom clone"
	uc := newUC(issueRepo, projRepo, &testutil.WorkflowRepoMock{})
	clone, err := uc.Clone(ctx, "issue-src", "user-1", &entity.CloneIssueReq{Title: &customTitle})

	require.NoError(t, err)
	assert.Equal(t, "My custom clone", clone.Title)
}

// ─── GetByID ─────────────────────────────────────────────────────────────────

func TestGetByID_NotFound(t *testing.T) {
	ctx := context.Background()

	issueRepo := &testutil.IssueRepoMock{
		GetByIDFn: func(_ context.Context, _ string) (*entity.Issue, error) {
			return nil, apperr.NotFound("issue not found")
		},
	}

	uc := newUC(issueRepo, &testutil.ProjectRepoMock{}, &testutil.WorkflowRepoMock{})
	_, err := uc.GetByID(ctx, "nonexistent")

	require.Error(t, err)
	var appErr *apperr.AppError
	assert.True(t, errors.As(err, &appErr))
	assert.Equal(t, apperr.CodeNotFound, appErr.Code)
}

// ─── BulkUpdate ──────────────────────────────────────────────────────────────

func TestBulkUpdate_PartialFailure(t *testing.T) {
	ctx := context.Background()

	issueRepo := &testutil.IssueRepoMock{
		BulkUpdateFn: func(_ context.Context, req *entity.BulkUpdateIssueReq) ([]string, error) {
			// Simulate only first ID was updated
			return []string{req.IssueIDs[0]}, nil
		},
	}

	uc := newUC(issueRepo, &testutil.ProjectRepoMock{}, &testutil.WorkflowRepoMock{})
	result, err := uc.BulkUpdate(ctx, &entity.BulkUpdateIssueReq{
		IssueIDs: []string{"id-1", "id-2", "id-3"},
	}, "actor-1")

	require.NoError(t, err)
	assert.Len(t, result.Updated, 1)
	assert.Len(t, result.Failed, 2)
	assert.Equal(t, 3, result.Total)
}

func TestBulkUpdate_RepoError(t *testing.T) {
	ctx := context.Background()

	issueRepo := &testutil.IssueRepoMock{
		BulkUpdateFn: func(_ context.Context, _ *entity.BulkUpdateIssueReq) ([]string, error) {
			return nil, errors.New("db timeout")
		},
	}

	uc := newUC(issueRepo, &testutil.ProjectRepoMock{}, &testutil.WorkflowRepoMock{})
	_, err := uc.BulkUpdate(ctx, &entity.BulkUpdateIssueReq{IssueIDs: []string{"id-1"}}, "actor-1")

	require.Error(t, err)
}
