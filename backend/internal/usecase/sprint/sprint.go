package sprint

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
	"github.com/jira-backend/jiraflow-backend/internal/infrastructura/repository"
	apperr "github.com/jira-backend/jiraflow-backend/internal/pkg/errors"
	"github.com/jira-backend/jiraflow-backend/internal/pkg/logger"
	"github.com/jira-backend/jiraflow-backend/internal/usecase/notification"
)

type useCase struct {
	repo        repository.SprintRepository
	issueRepo   repository.IssueRepository
	spaceRepo   repository.SpaceRepository
	pageRepo    repository.PageRepository
	versionRepo repository.PageVersionRepository
	dispatcher  notification.Dispatcher
	log         logger.Logger
}

func New(
	repo repository.SprintRepository,
	issueRepo repository.IssueRepository,
	spaceRepo repository.SpaceRepository,
	pageRepo repository.PageRepository,
	versionRepo repository.PageVersionRepository,
	dispatcher notification.Dispatcher,
	log logger.Logger,
) UseCase {
	return &useCase{repo: repo, issueRepo: issueRepo, spaceRepo: spaceRepo, pageRepo: pageRepo, versionRepo: versionRepo, dispatcher: dispatcher, log: log}
}

func (uc *useCase) Create(ctx context.Context, projectID, createdBy string, s *entity.Sprint) (*entity.Sprint, error) {
	now := time.Now().UTC()
	s.ID = uuid.NewString()
	s.ProjectID = projectID
	s.CreatedBy = createdBy
	s.Status = "planned"
	s.CreatedAt = now
	s.UpdatedAt = now
	if err := uc.repo.Create(ctx, s); err != nil {
		uc.log.Error(ctx, "sprint.Create: db error", logger.String("project_id", projectID), logger.SafeString("err", err.Error()))
		return nil, err
	}
	uc.log.Info(ctx, "sprint created", logger.String("id", s.ID), logger.String("project_id", projectID))
	return s, nil
}

func (uc *useCase) GetByID(ctx context.Context, id string) (*entity.Sprint, error) {
	return uc.repo.GetByID(ctx, id)
}

func (uc *useCase) List(ctx context.Context, projectID string, filter *entity.SprintFilter) ([]*entity.Sprint, int, error) {
	return uc.repo.List(ctx, projectID, filter)
}

func (uc *useCase) Update(ctx context.Context, id string, s *entity.Sprint) (*entity.Sprint, error) {
	existing, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if s.Name != "" {
		existing.Name = s.Name
	}
	existing.Goal = s.Goal
	existing.StartDate = s.StartDate
	existing.EndDate = s.EndDate
	if err := uc.repo.Update(ctx, existing); err != nil {
		uc.log.Error(ctx, "sprint.Update: db error", logger.String("id", id), logger.SafeString("err", err.Error()))
		return nil, err
	}
	uc.log.Info(ctx, "sprint updated", logger.String("id", id))
	return existing, nil
}

func (uc *useCase) Start(ctx context.Context, id, actorID string) (*entity.Sprint, error) {
	s, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if s.Status != "planned" {
		return nil, apperr.BadRequest("only planned sprints can be started")
	}
	active, err := uc.repo.GetActive(ctx, s.ProjectID)
	if err == nil && active != nil {
		return nil, apperr.Conflict("project already has an active sprint")
	}
	now := time.Now().UTC()
	if err := uc.repo.Start(ctx, id, now); err != nil {
		uc.log.Error(ctx, "sprint.Start: db error", logger.String("id", id), logger.SafeString("err", err.Error()))
		return nil, err
	}
	s.Status = "active"
	s.StartedAt = &now
	uc.log.Info(ctx, "sprint started", logger.String("id", id), logger.String("actor_id", actorID))

	if uc.dispatcher != nil {
		go uc.dispatcher.SprintStarted(context.Background(), s.ID, s.ProjectID, s.Name)
	}
	return s, nil
}

func (uc *useCase) Complete(ctx context.Context, id, actorID string) (*entity.Sprint, error) {
	s, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if s.Status != "active" {
		return nil, apperr.BadRequest("only active sprints can be completed")
	}
	now := time.Now().UTC()
	if err := uc.repo.Complete(ctx, id, now); err != nil {
		uc.log.Error(ctx, "sprint.Complete: db error", logger.String("id", id), logger.SafeString("err", err.Error()))
		return nil, err
	}
	s.Status = "completed"
	s.CompletedAt = &now
	uc.log.Info(ctx, "sprint completed", logger.String("id", id), logger.String("actor_id", actorID))

	if uc.dispatcher != nil {
		go uc.dispatcher.SprintCompleted(context.Background(), s.ID, s.ProjectID, s.Name)
	}
	go uc.autoCreateRetroPage(context.Background(), s, actorID)

	return s, nil
}

func (uc *useCase) autoCreateRetroPage(ctx context.Context, s *entity.Sprint, actorID string) {
	if uc.spaceRepo == nil || uc.pageRepo == nil {
		return
	}
	space, err := uc.spaceRepo.GetByProjectID(ctx, s.ProjectID)
	if err != nil {
		uc.log.Warn(ctx, "sprint.autoCreateRetroPage: no linked space",
			logger.String("project_id", s.ProjectID), logger.SafeString("err", err.Error()))
		return
	}

	title := fmt.Sprintf("Retrospective: %s", s.Name)
	contentText := fmt.Sprintf("Sprint: %s\nGoal: %s\n\n## What went well\n\n## What could be improved\n\n## Action items\n",
		s.Name, stringOrEmpty(s.Goal))

	now := time.Now().UTC()
	p := &entity.Page{
		ID:             uuid.NewString(),
		SpaceID:        space.ID,
		Title:          title,
		Content:        map[string]any{"type": "doc", "content": []any{}},
		ContentText:    contentText,
		AuthorID:       actorID,
		LastEditorID:   actorID,
		CurrentVersion: 1,
		Status:         "draft",
		Position:       1,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
	if err := uc.pageRepo.Create(ctx, p); err != nil {
		uc.log.Warn(ctx, "sprint.autoCreateRetroPage: page create failed",
			logger.String("sprint_id", s.ID), logger.SafeString("err", err.Error()))
		return
	}
	if uc.versionRepo != nil {
		v := &entity.PageVersion{
			ID:          uuid.NewString(),
			PageID:      p.ID,
			Version:     1,
			Title:       p.Title,
			Content:     p.Content,
			ContentText: p.ContentText,
			AuthorID:    actorID,
			CreatedAt:   now,
		}
		_ = uc.versionRepo.Create(ctx, v)
	}
	uc.log.Info(ctx, "sprint retro page created",
		logger.String("sprint_id", s.ID), logger.String("page_id", p.ID))
}

func stringOrEmpty(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func (uc *useCase) Delete(ctx context.Context, id string) error {
	if err := uc.repo.SoftDelete(ctx, id); err != nil {
		uc.log.Error(ctx, "sprint.Delete: db error", logger.String("id", id), logger.SafeString("err", err.Error()))
		return err
	}
	uc.log.Info(ctx, "sprint deleted", logger.String("id", id))
	return nil
}

func (uc *useCase) AddIssue(ctx context.Context, sprintID, issueID, actorID string) error {
	s, err := uc.repo.GetByID(ctx, sprintID)
	if err != nil {
		return err
	}
	if s.Status == "completed" {
		return apperr.BadRequest("cannot add issue to a completed sprint")
	}
	if err := uc.repo.AddIssue(ctx, sprintID, issueID); err != nil {
		uc.log.Error(ctx, "sprint.AddIssue: db error",
			logger.String("sprint_id", sprintID),
			logger.String("issue_id", issueID),
			logger.SafeString("err", err.Error()))
		return err
	}
	uc.log.Info(ctx, "issue added to sprint",
		logger.String("sprint_id", sprintID),
		logger.String("issue_id", issueID))
	return nil
}

func (uc *useCase) RemoveIssue(ctx context.Context, sprintID, issueID, actorID string) error {
	s, err := uc.repo.GetByID(ctx, sprintID)
	if err != nil {
		return err
	}
	if s.Status == "completed" {
		return apperr.BadRequest("cannot remove issue from a completed sprint")
	}
	if err := uc.repo.RemoveIssue(ctx, sprintID, issueID); err != nil {
		uc.log.Error(ctx, "sprint.RemoveIssue: db error",
			logger.String("sprint_id", sprintID),
			logger.String("issue_id", issueID),
			logger.SafeString("err", err.Error()))
		return err
	}
	uc.log.Info(ctx, "issue removed from sprint",
		logger.String("sprint_id", sprintID),
		logger.String("issue_id", issueID))
	return nil
}

func (uc *useCase) GetReport(ctx context.Context, sprintID string) (*entity.SprintReport, error) {
	report, err := uc.repo.GetReport(ctx, sprintID)
	if err != nil {
		uc.log.Error(ctx, "sprint.GetReport: db error", logger.String("sprint_id", sprintID), logger.SafeString("err", err.Error()))
		return nil, err
	}
	return report, nil
}

func (uc *useCase) GetBurndown(ctx context.Context, sprintID string) (*entity.BurndownChart, error) {
	chart, err := uc.repo.GetBurndown(ctx, sprintID)
	if err != nil {
		uc.log.Error(ctx, "sprint.GetBurndown: db error", logger.String("sprint_id", sprintID), logger.SafeString("err", err.Error()))
		return nil, err
	}
	return chart, nil
}

func (uc *useCase) GetBurnup(ctx context.Context, sprintID string) (*entity.BurnupChart, error) {
	chart, err := uc.repo.GetBurnup(ctx, sprintID)
	if err != nil {
		uc.log.Error(ctx, "sprint.GetBurnup: db error", logger.String("sprint_id", sprintID), logger.SafeString("err", err.Error()))
		return nil, err
	}
	return chart, nil
}

func (uc *useCase) GetCFD(ctx context.Context, projectID string, from, to *string) (*entity.CFDChart, error) {
	chart, err := uc.repo.GetCFD(ctx, projectID, from, to)
	if err != nil {
		uc.log.Error(ctx, "sprint.GetCFD: db error", logger.String("project_id", projectID), logger.SafeString("err", err.Error()))
		return nil, err
	}
	return chart, nil
}

func (uc *useCase) GetVelocity(ctx context.Context, projectID string, limit int) (*entity.VelocityReport, error) {
	if limit <= 0 {
		limit = 10
	}
	report, err := uc.repo.GetVelocity(ctx, projectID, limit)
	if err != nil {
		uc.log.Error(ctx, "sprint.GetVelocity: db error", logger.String("project_id", projectID), logger.SafeString("err", err.Error()))
		return nil, err
	}
	return report, nil
}

func (uc *useCase) GetSprintPlanning(ctx context.Context, projectID string) (*entity.SprintPlanningView, error) {
	active, _ := uc.repo.GetActive(ctx, projectID)

	large := 200
	backlogFilter := &entity.IssueFilter{
		Filter:    entity.Filter{Limit: large},
		ProjectID: projectID,
		NoSprint:  true,
	}
	backlog, _, err := uc.issueRepo.List(ctx, backlogFilter)
	if err != nil {
		return nil, fmt.Errorf("sprint.GetSprintPlanning backlog: %w", err)
	}

	view := &entity.SprintPlanningView{
		ActiveSprint: active,
		BacklogItems: backlog,
	}

	if active != nil {
		sprintFilter := &entity.IssueFilter{
			Filter:   entity.Filter{Limit: large},
			ProjectID: projectID,
			SprintID: active.ID,
		}
		items, _, err := uc.issueRepo.List(ctx, sprintFilter)
		if err != nil {
			return nil, fmt.Errorf("sprint.GetSprintPlanning sprint items: %w", err)
		}
		view.SprintItems = items
	}
	return view, nil
}

func (uc *useCase) BulkAssignToSprint(ctx context.Context, projectID string, req *entity.AssignToSprintReq) error {
	sprint, err := uc.repo.GetByID(ctx, req.SprintID)
	if err != nil {
		return err
	}
	if sprint.ProjectID != projectID {
		return apperr.Forbidden("sprint does not belong to this project")
	}
	if sprint.Status == "completed" {
		return apperr.BadRequest("cannot assign issues to a completed sprint")
	}
	for _, issueID := range req.IssueIDs {
		if err := uc.repo.AddIssue(ctx, req.SprintID, issueID); err != nil {
			uc.log.Warn(ctx, "sprint.BulkAssign: add issue failed",
				logger.String("sprint_id", req.SprintID),
				logger.String("issue_id", issueID),
				logger.SafeString("err", err.Error()),
			)
		}
	}
	return nil
}

func (uc *useCase) GetCapacity(ctx context.Context, sprintID string) (*entity.SprintCapacity, error) {
	sprint, err := uc.repo.GetByID(ctx, sprintID)
	if err != nil {
		return nil, err
	}

	issues, _, err := uc.issueRepo.List(ctx, &entity.IssueFilter{
		Filter:   entity.Filter{Limit: 500},
		ProjectID: sprint.ProjectID,
		SprintID: sprintID,
	})
	if err != nil {
		return nil, fmt.Errorf("sprint.GetCapacity list issues: %w", err)
	}

	totals := map[string]int{}
	names := map[string]string{}
	total := 0
	for _, issue := range issues {
		if issue.StoryPoints == nil {
			continue
		}
		pts := *issue.StoryPoints
		total += pts
		if issue.Assignee != nil {
			totals[issue.Assignee.ID] += pts
			names[issue.Assignee.ID] = issue.Assignee.FullName
		}
	}

	cap := &entity.SprintCapacity{TotalPoints: total}
	for uid, pts := range totals {
		cap.ByAssignee = append(cap.ByAssignee, entity.AssigneeCapacity{
			UserID:      uid,
			DisplayName: names[uid],
			Points:      pts,
		})
	}
	return cap, nil
}

func (uc *useCase) UpdateGoal(ctx context.Context, sprintID string, goal string) (*entity.Sprint, error) {
	sprint, err := uc.repo.GetByID(ctx, sprintID)
	if err != nil {
		return nil, err
	}
	sprint.Goal = &goal
	sprint.UpdatedAt = time.Now().UTC()
	if err := uc.repo.Update(ctx, sprint); err != nil {
		return nil, fmt.Errorf("sprint.UpdateGoal: %w", err)
	}
	return sprint, nil
}

func (uc *useCase) GetImpediments(ctx context.Context, sprintID string) ([]*entity.Issue, error) {
	sprint, err := uc.repo.GetByID(ctx, sprintID)
	if err != nil {
		return nil, err
	}
	issues, _, err := uc.issueRepo.List(ctx, &entity.IssueFilter{
		Filter:    entity.Filter{Limit: 200},
		ProjectID: sprint.ProjectID,
		SprintID:  sprintID,
		Priority:  "highest",
	})
	if err != nil {
		return nil, fmt.Errorf("sprint.GetImpediments: %w", err)
	}
	return issues, nil
}
