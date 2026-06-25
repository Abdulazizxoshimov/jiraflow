package workflow_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
	apperr "github.com/jira-backend/jiraflow-backend/internal/pkg/errors"
	"github.com/jira-backend/jiraflow-backend/internal/testutil"
	"github.com/jira-backend/jiraflow-backend/internal/usecase/workflow"
)

func newUC(repo *testutil.WorkflowRepoMock) workflow.UseCase {
	return workflow.New(repo, nil, testutil.NopLogger{})
}

// ─── Create ──────────────────────────────────────────────────────────────────

func TestCreate_SetsID(t *testing.T) {
	ctx := context.Background()

	var saved *entity.Workflow
	repo := &testutil.WorkflowRepoMock{
		CreateFn: func(_ context.Context, wf *entity.Workflow) error {
			saved = wf
			return nil
		},
	}

	uc := newUC(repo)
	result, err := uc.Create(ctx, &entity.Workflow{Name: "Default"})

	require.NoError(t, err)
	assert.NotEmpty(t, result.ID)
	assert.Equal(t, "Default", result.Name)
	require.NotNil(t, saved)
	assert.Equal(t, result.ID, saved.ID)
}

func TestCreate_RepoError(t *testing.T) {
	ctx := context.Background()

	repo := &testutil.WorkflowRepoMock{
		CreateFn: func(_ context.Context, _ *entity.Workflow) error {
			return errors.New("db error")
		},
	}

	uc := newUC(repo)
	_, err := uc.Create(ctx, &entity.Workflow{Name: "Bad"})
	require.Error(t, err)
}

// ─── Update ──────────────────────────────────────────────────────────────────

func TestUpdate_Success(t *testing.T) {
	ctx := context.Background()

	repo := &testutil.WorkflowRepoMock{
		GetByIDFn: func(_ context.Context, id string) (*entity.Workflow, error) {
			return &entity.Workflow{ID: id, Name: "Old Name"}, nil
		},
	}

	uc := newUC(repo)
	result, err := uc.Update(ctx, "wf-1", &entity.Workflow{Name: "New Name"})

	require.NoError(t, err)
	assert.Equal(t, "New Name", result.Name)
}

func TestUpdate_NotFound(t *testing.T) {
	ctx := context.Background()

	repo := &testutil.WorkflowRepoMock{
		GetByIDFn: func(_ context.Context, _ string) (*entity.Workflow, error) {
			return nil, apperr.NotFound("workflow not found")
		},
	}

	uc := newUC(repo)
	_, err := uc.Update(ctx, "wf-x", &entity.Workflow{Name: "New"})

	require.Error(t, err)
	var appErr *apperr.AppError
	require.True(t, errors.As(err, &appErr))
	assert.Equal(t, apperr.CodeNotFound, appErr.Code)
}

// ─── Delete ──────────────────────────────────────────────────────────────────

func TestDelete_Success(t *testing.T) {
	ctx := context.Background()

	var deletedID string
	repo := &testutil.WorkflowRepoMock{
		SoftDeleteFn: func(_ context.Context, id string) error {
			deletedID = id
			return nil
		},
	}

	uc := newUC(repo)
	err := uc.Delete(ctx, "wf-1")

	require.NoError(t, err)
	assert.Equal(t, "wf-1", deletedID)
}

func TestDelete_RepoError(t *testing.T) {
	ctx := context.Background()

	repo := &testutil.WorkflowRepoMock{
		SoftDeleteFn: func(_ context.Context, _ string) error {
			return errors.New("db error")
		},
	}

	uc := newUC(repo)
	err := uc.Delete(ctx, "wf-bad")
	require.Error(t, err)
}

// ─── SetDefault ──────────────────────────────────────────────────────────────

func TestSetDefault_NotFound(t *testing.T) {
	ctx := context.Background()

	repo := &testutil.WorkflowRepoMock{
		GetByIDFn: func(_ context.Context, _ string) (*entity.Workflow, error) {
			return nil, apperr.NotFound("not found")
		},
	}

	uc := newUC(repo)
	err := uc.SetDefault(ctx, "wf-x")

	require.Error(t, err)
	var appErr *apperr.AppError
	require.True(t, errors.As(err, &appErr))
	assert.Equal(t, apperr.CodeNotFound, appErr.Code)
}
