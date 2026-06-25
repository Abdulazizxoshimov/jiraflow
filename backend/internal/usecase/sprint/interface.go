package sprint

import (
	"context"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

type UseCase interface {
	Create(ctx context.Context, projectID, createdBy string, isAdmin bool, s *entity.Sprint) (*entity.Sprint, error)
	GetByID(ctx context.Context, id string) (*entity.Sprint, error)
	List(ctx context.Context, projectID string, filter *entity.SprintFilter) ([]*entity.Sprint, int, error)
	Update(ctx context.Context, id string, s *entity.Sprint) (*entity.Sprint, error)
	Start(ctx context.Context, id, actorID string) (*entity.Sprint, error)
	Complete(ctx context.Context, id, actorID string) (*entity.Sprint, error)
	Delete(ctx context.Context, id string) error

	AddIssue(ctx context.Context, sprintID, issueID, actorID string) error
	RemoveIssue(ctx context.Context, sprintID, issueID, actorID string) error

	GetReport(ctx context.Context, sprintID string) (*entity.SprintReport, error)
	GetBurndown(ctx context.Context, sprintID string) (*entity.BurndownChart, error)
	GetBurnup(ctx context.Context, sprintID string) (*entity.BurnupChart, error)
	GetCFD(ctx context.Context, projectID string, from, to *string) (*entity.CFDChart, error)
	GetVelocity(ctx context.Context, projectID string, limit int) (*entity.VelocityReport, error)

	GetSprintPlanning(ctx context.Context, projectID string) (*entity.SprintPlanningView, error)
	BulkAssignToSprint(ctx context.Context, projectID string, req *entity.AssignToSprintReq) error
	GetCapacity(ctx context.Context, sprintID string) (*entity.SprintCapacity, error)
	UpdateGoal(ctx context.Context, sprintID string, goal string) (*entity.Sprint, error)
	GetImpediments(ctx context.Context, sprintID string) ([]*entity.Issue, error)
}
