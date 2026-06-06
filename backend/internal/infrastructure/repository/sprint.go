package repository

import (
	"context"
	"time"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

type SprintRepository interface {
	Create(ctx context.Context, s *entity.Sprint) error
	GetByID(ctx context.Context, id string) (*entity.Sprint, error)
	List(ctx context.Context, projectID string, filter *entity.SprintFilter) ([]*entity.Sprint, int, error)
	Update(ctx context.Context, s *entity.Sprint) error
	SoftDelete(ctx context.Context, id string) error
	GetActive(ctx context.Context, projectID string) (*entity.Sprint, error)
	Start(ctx context.Context, id string, startedAt time.Time) error
	Complete(ctx context.Context, id string, completedAt time.Time) error

	AddIssue(ctx context.Context, sprintID, issueID string) error
	RemoveIssue(ctx context.Context, sprintID, issueID string) error

	GetReport(ctx context.Context, sprintID string) (*entity.SprintReport, error)
	GetBurndown(ctx context.Context, sprintID string) (*entity.BurndownChart, error)
	GetBurnup(ctx context.Context, sprintID string) (*entity.BurnupChart, error)
	GetCFD(ctx context.Context, projectID string, from, to *string) (*entity.CFDChart, error)
	GetVelocity(ctx context.Context, projectID string, limit int) (*entity.VelocityReport, error)
}
