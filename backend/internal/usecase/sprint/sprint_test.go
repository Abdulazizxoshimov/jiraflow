package sprint_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
	apperr "github.com/jira-backend/jiraflow-backend/internal/pkg/errors"
	"github.com/jira-backend/jiraflow-backend/internal/testutil"
	"github.com/jira-backend/jiraflow-backend/internal/usecase/sprint"
)

func newUC(repo *testutil.SprintRepoMock) sprint.UseCase {
	return sprint.New(repo, nil, nil, nil, nil, testutil.ProjectMemberRepoMock{}, testutil.DispatcherMock{}, testutil.NopLogger{})
}

// ─── Create ──────────────────────────────────────────────────────────────────

func TestCreate_SetsPlannedStatus(t *testing.T) {
	ctx := context.Background()

	var saved *entity.Sprint
	repo := &testutil.SprintRepoMock{
		CreateFn: func(_ context.Context, s *entity.Sprint) error {
			saved = s
			return nil
		},
	}

	uc := newUC(repo)
	result, err := uc.Create(ctx, "proj-1", "user-1", false, &entity.Sprint{Name: "Sprint 1"})

	require.NoError(t, err)
	assert.Equal(t, "planned", result.Status)
	assert.Equal(t, "Sprint 1", result.Name)
	assert.NotEmpty(t, result.ID)
	assert.NotNil(t, saved)
}

// ─── AddIssue ────────────────────────────────────────────────────────────────

func TestAddIssue_Success(t *testing.T) {
	ctx := context.Background()

	var addedSprintID, addedIssueID string
	repo := &testutil.SprintRepoMock{
		GetByIDFn: func(_ context.Context, id string) (*entity.Sprint, error) {
			return &entity.Sprint{ID: id, Status: "active"}, nil
		},
		AddIssueFn: func(_ context.Context, sprintID, issueID string) error {
			addedSprintID = sprintID
			addedIssueID = issueID
			return nil
		},
	}

	uc := newUC(repo)
	err := uc.AddIssue(ctx, "sprint-1", "issue-1", "actor-1")

	require.NoError(t, err)
	assert.Equal(t, "sprint-1", addedSprintID)
	assert.Equal(t, "issue-1", addedIssueID)
}

func TestAddIssue_CompletedSprint(t *testing.T) {
	ctx := context.Background()

	repo := &testutil.SprintRepoMock{
		GetByIDFn: func(_ context.Context, id string) (*entity.Sprint, error) {
			return &entity.Sprint{ID: id, Status: "completed"}, nil
		},
	}

	uc := newUC(repo)
	err := uc.AddIssue(ctx, "sprint-1", "issue-1", "actor-1")

	require.Error(t, err)
	appErr, ok := err.(*apperr.AppError)
	require.True(t, ok)
	assert.Equal(t, apperr.CodeBadRequest, appErr.Code)
}

func TestAddIssue_SprintNotFound(t *testing.T) {
	ctx := context.Background()

	repo := &testutil.SprintRepoMock{
		GetByIDFn: func(_ context.Context, _ string) (*entity.Sprint, error) {
			return nil, apperr.NotFound("sprint not found")
		},
	}

	uc := newUC(repo)
	err := uc.AddIssue(ctx, "sprint-x", "issue-1", "actor-1")

	require.Error(t, err)
	appErr, ok := err.(*apperr.AppError)
	require.True(t, ok)
	assert.Equal(t, apperr.CodeNotFound, appErr.Code)
}

// ─── RemoveIssue ─────────────────────────────────────────────────────────────

func TestRemoveIssue_Success(t *testing.T) {
	ctx := context.Background()

	var removedSprintID, removedIssueID string
	repo := &testutil.SprintRepoMock{
		GetByIDFn: func(_ context.Context, id string) (*entity.Sprint, error) {
			return &entity.Sprint{ID: id, Status: "active"}, nil
		},
		RemoveIssueFn: func(_ context.Context, sprintID, issueID string) error {
			removedSprintID = sprintID
			removedIssueID = issueID
			return nil
		},
	}

	uc := newUC(repo)
	err := uc.RemoveIssue(ctx, "sprint-1", "issue-1", "actor-1")

	require.NoError(t, err)
	assert.Equal(t, "sprint-1", removedSprintID)
	assert.Equal(t, "issue-1", removedIssueID)
}

func TestRemoveIssue_CompletedSprint(t *testing.T) {
	ctx := context.Background()

	repo := &testutil.SprintRepoMock{
		GetByIDFn: func(_ context.Context, id string) (*entity.Sprint, error) {
			return &entity.Sprint{ID: id, Status: "completed"}, nil
		},
	}

	uc := newUC(repo)
	err := uc.RemoveIssue(ctx, "sprint-1", "issue-1", "actor-1")

	require.Error(t, err)
	appErr, ok := err.(*apperr.AppError)
	require.True(t, ok)
	assert.Equal(t, apperr.CodeBadRequest, appErr.Code)
}

// ─── Start ───────────────────────────────────────────────────────────────────

func TestStart_Success(t *testing.T) {
	ctx := context.Background()

	repo := &testutil.SprintRepoMock{
		GetByIDFn: func(_ context.Context, id string) (*entity.Sprint, error) {
			return &entity.Sprint{ID: id, ProjectID: "proj-1", Status: "planned"}, nil
		},
		GetActiveFn: func(_ context.Context, _ string) (*entity.Sprint, error) {
			return nil, apperr.NotFound("no active sprint")
		},
	}

	uc := newUC(repo)
	result, err := uc.Start(ctx, "sprint-1", "actor-1")

	require.NoError(t, err)
	assert.Equal(t, "active", result.Status)
	assert.NotNil(t, result.StartedAt)
}

func TestStart_AlreadyActiveSprintExists(t *testing.T) {
	ctx := context.Background()

	repo := &testutil.SprintRepoMock{
		GetByIDFn: func(_ context.Context, id string) (*entity.Sprint, error) {
			return &entity.Sprint{ID: id, ProjectID: "proj-1", Status: "planned"}, nil
		},
		GetActiveFn: func(_ context.Context, _ string) (*entity.Sprint, error) {
			return &entity.Sprint{ID: "other-sprint", Status: "active"}, nil
		},
	}

	uc := newUC(repo)
	_, err := uc.Start(ctx, "sprint-1", "actor-1")

	require.Error(t, err)
	appErr, ok := err.(*apperr.AppError)
	require.True(t, ok)
	assert.Equal(t, apperr.CodeConflict, appErr.Code)
}

func TestStart_NotPlanned(t *testing.T) {
	ctx := context.Background()

	repo := &testutil.SprintRepoMock{
		GetByIDFn: func(_ context.Context, id string) (*entity.Sprint, error) {
			return &entity.Sprint{ID: id, Status: "completed"}, nil
		},
	}

	uc := newUC(repo)
	_, err := uc.Start(ctx, "sprint-1", "actor-1")

	require.Error(t, err)
	appErr, ok := err.(*apperr.AppError)
	require.True(t, ok)
	assert.Equal(t, apperr.CodeBadRequest, appErr.Code)
}

// ─── Complete ────────────────────────────────────────────────────────────────

func TestComplete_Success(t *testing.T) {
	ctx := context.Background()

	repo := &testutil.SprintRepoMock{
		GetByIDFn: func(_ context.Context, id string) (*entity.Sprint, error) {
			return &entity.Sprint{ID: id, ProjectID: "proj-1", Status: "active", Name: "Sprint 1"}, nil
		},
	}

	uc := newUC(repo)
	result, err := uc.Complete(ctx, "sprint-1", "actor-1")

	require.NoError(t, err)
	assert.Equal(t, "completed", result.Status)
	assert.NotNil(t, result.CompletedAt)
}

func TestComplete_NotActive(t *testing.T) {
	ctx := context.Background()

	repo := &testutil.SprintRepoMock{
		GetByIDFn: func(_ context.Context, id string) (*entity.Sprint, error) {
			return &entity.Sprint{ID: id, Status: "planned"}, nil
		},
	}

	uc := newUC(repo)
	_, err := uc.Complete(ctx, "sprint-1", "actor-1")

	require.Error(t, err)
	appErr, ok := err.(*apperr.AppError)
	require.True(t, ok)
	assert.Equal(t, apperr.CodeBadRequest, appErr.Code)
}
