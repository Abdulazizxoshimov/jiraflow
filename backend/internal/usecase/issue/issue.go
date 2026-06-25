package issue

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
	"github.com/jira-backend/jiraflow-backend/internal/infrastructure/repository"
	ws "github.com/jira-backend/jiraflow-backend/internal/infrastructure/websocket"
	apperr "github.com/jira-backend/jiraflow-backend/internal/pkg/errors"
	"github.com/jira-backend/jiraflow-backend/internal/pkg/lexorank"
	"github.com/jira-backend/jiraflow-backend/internal/pkg/logger"
	"github.com/jira-backend/jiraflow-backend/internal/usecase/notification"
)

type useCase struct {
	issueRepo    repository.IssueRepository
	projectRepo  repository.ProjectRepository
	workflowRepo repository.WorkflowRepository
	versionRepo  repository.VersionRepository
	memberRepo   repository.ProjectMemberRepository
	dispatcher   notification.Dispatcher
	hub          *ws.Hub
	log          logger.Logger
}

func New(
	issueRepo repository.IssueRepository,
	projectRepo repository.ProjectRepository,
	workflowRepo repository.WorkflowRepository,
	versionRepo repository.VersionRepository,
	memberRepo repository.ProjectMemberRepository,
	dispatcher notification.Dispatcher,
	hub *ws.Hub,
	log logger.Logger,
) UseCase {
	return &useCase{
		issueRepo:    issueRepo,
		projectRepo:  projectRepo,
		workflowRepo: workflowRepo,
		versionRepo:  versionRepo,
		memberRepo:   memberRepo,
		dispatcher:   dispatcher,
		hub:          hub,
		log:          log,
	}
}

func (uc *useCase) Create(ctx context.Context, projectID string, req *entity.CreateIssueReq, reporterID string, isAdmin bool) (*entity.Issue, error) {
	project, err := uc.projectRepo.GetByID(ctx, projectID)
	if err != nil {
		return nil, err
	}

	if !isAdmin {
		isMember, err := uc.memberRepo.IsMember(ctx, projectID, reporterID)
		if err != nil {
			return nil, fmt.Errorf("issue.Create membership check: %w", err)
		}
		if !isMember {
			return nil, apperr.Forbidden("you are not a member of this project")
		}
	}

	wf, err := uc.workflowRepo.GetWithDetails(ctx, project.WorkflowID)
	if err != nil {
		uc.log.Error(ctx, "issue.Create: get workflow failed", logger.String("project_id", projectID), logger.SafeString("err", err.Error()))
		return nil, fmt.Errorf("issue.Create get workflow: %w", err)
	}

	var initialStatusID string
	for _, s := range wf.Statuses {
		if s.IsInitial {
			initialStatusID = s.ID
			break
		}
	}
	if initialStatusID == "" && len(wf.Statuses) > 0 {
		initialStatusID = wf.Statuses[0].ID
	}

	counter, err := uc.projectRepo.IncrementIssueCounter(ctx, projectID)
	if err != nil {
		uc.log.Error(ctx, "issue.Create: increment counter failed", logger.String("project_id", projectID), logger.SafeString("err", err.Error()))
		return nil, fmt.Errorf("issue.Create increment counter: %w", err)
	}

	priority := req.Priority
	if priority == "" {
		priority = "medium"
	}

	// omitempty doesn't skip empty strings for *string — coerce "" to nil so DB never sees an invalid UUID.
	if req.ParentID != nil && *req.ParentID == "" {
		req.ParentID = nil
	}
	if req.AssigneeID != nil && *req.AssigneeID == "" {
		req.AssigneeID = nil
	}
	if req.SprintID != nil && *req.SprintID == "" {
		req.SprintID = nil
	}

	now := time.Now().UTC()
	issue := &entity.Issue{
		ID:           uuid.NewString(),
		ProjectID:    projectID,
		IssueNumber:  int(counter),
		Title:        req.Title,
		Description:  req.Description,
		Type:         req.Type,
		StatusID:     initialStatusID,
		Priority:     priority,
		AssigneeID:   req.AssigneeID,
		ReporterID:   reporterID,
		ParentID:     req.ParentID,
		SprintID:     req.SprintID,
		StoryPoints:  req.StoryPoints,
		DueDate:      req.DueDate,
		CustomFields: req.CustomFields,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	if issue.CustomFields == nil {
		issue.CustomFields = map[string]any{}
	}

	if err := uc.issueRepo.Create(ctx, issue); err != nil {
		uc.log.Error(ctx, "issue.Create: db error", logger.String("project_id", projectID), logger.SafeString("err", err.Error()))
		return nil, err
	}

	if len(req.LabelIDs) > 0 {
		if err := uc.issueRepo.SetLabels(ctx, issue.ID, req.LabelIDs); err != nil {
			uc.log.Warn(ctx, "issue.Create: set labels failed", logger.String("id", issue.ID), logger.SafeString("err", err.Error()))
		}
	}
	if len(req.FixVersionIDs) > 0 {
		if err := uc.versionRepo.SetIssueVersions(ctx, issue.ID, req.FixVersionIDs); err != nil {
			uc.log.Warn(ctx, "issue.Create: set fix versions failed", logger.String("id", issue.ID), logger.SafeString("err", err.Error()))
		}
	}
	if len(req.AffectsVersionIDs) > 0 {
		if err := uc.versionRepo.SetIssueAffectsVersions(ctx, issue.ID, req.AffectsVersionIDs); err != nil {
			uc.log.Warn(ctx, "issue.Create: set affects versions failed", logger.String("id", issue.ID), logger.SafeString("err", err.Error()))
		}
	}

	watcher := &entity.IssueWatcher{IssueID: issue.ID, UserID: reporterID, CreatedAt: now}
	if err := uc.issueRepo.AddWatcher(ctx, watcher); err != nil {
		uc.log.Warn(ctx, "issue.Create: add watcher failed", logger.String("id", issue.ID), logger.SafeString("err", err.Error()))
	}

	if issue.AssigneeID != nil {
		go uc.dispatcher.IssueAssigned(context.Background(), issue, *issue.AssigneeID, reporterID)
	}
	go uc.dispatcher.IssueCreated(context.Background(), issue, reporterID)

	uc.log.Info(ctx, "issue created", logger.String("id", issue.ID), logger.String("project_id", projectID))
	return issue, nil
}

func (uc *useCase) GetByID(ctx context.Context, id string) (*entity.Issue, error) {
	issue, err := uc.issueRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	uc.hydrate(ctx, issue)
	return issue, nil
}

func (uc *useCase) GetByKey(ctx context.Context, key string) (*entity.Issue, error) {
	issue, err := uc.issueRepo.GetByKey(ctx, key)
	if err != nil {
		return nil, err
	}
	uc.hydrate(ctx, issue)
	return issue, nil
}

// hydrate fetches labels, fix versions, and affects versions in parallel.
func (uc *useCase) hydrate(ctx context.Context, issue *entity.Issue) {
	var (
		mu      sync.Mutex
		labels  []*entity.Label
		fixVers []*entity.Version
		affVers []*entity.Version
	)
	g, gctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		l, err := uc.issueRepo.GetLabels(gctx, issue.ID)
		if err == nil {
			mu.Lock()
			labels = l
			mu.Unlock()
		}
		return nil
	})
	g.Go(func() error {
		v, err := uc.versionRepo.GetIssueVersions(gctx, issue.ID)
		if err == nil {
			mu.Lock()
			fixVers = v
			mu.Unlock()
		}
		return nil
	})
	g.Go(func() error {
		v, err := uc.versionRepo.GetIssueAffectsVersions(gctx, issue.ID)
		if err == nil {
			mu.Lock()
			affVers = v
			mu.Unlock()
		}
		return nil
	})
	_ = g.Wait()
	for _, l := range labels {
		issue.Labels = append(issue.Labels, *l)
	}
	for _, v := range fixVers {
		issue.Versions = append(issue.Versions, *v)
	}
	for _, v := range affVers {
		issue.AffectsVersions = append(issue.AffectsVersions, *v)
	}
}

func (uc *useCase) List(ctx context.Context, filter *entity.IssueFilter) ([]*entity.Issue, int, error) {
	return uc.issueRepo.List(ctx, filter)
}

func (uc *useCase) Update(ctx context.Context, id string, req *entity.UpdateIssueReq, actorID string) (*entity.Issue, error) {
	issue, err := uc.issueRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Title != nil {
		issue.Title = *req.Title
	}
	if req.Description != nil {
		issue.Description = req.Description
	}
	if req.Priority != nil {
		issue.Priority = *req.Priority
	}
	prevAssigneeID := issue.AssigneeID
	if req.AssigneeID != nil {
		issue.AssigneeID = req.AssigneeID
	}
	if req.SprintID != nil {
		issue.SprintID = req.SprintID
	}
	if req.StoryPoints != nil {
		issue.StoryPoints = req.StoryPoints
	}
	if req.DueDate != nil {
		issue.DueDate = req.DueDate
	}
	if req.CustomFields != nil {
		issue.CustomFields = req.CustomFields
	}
	if req.Resolution != nil {
		issue.Resolution = req.Resolution
	}

	if err := uc.issueRepo.Update(ctx, issue); err != nil {
		uc.log.Error(ctx, "issue.Update: db error", logger.String("id", id), logger.SafeString("err", err.Error()))
		return nil, err
	}
	if req.LabelIDs != nil {
		_ = uc.issueRepo.SetLabels(ctx, id, req.LabelIDs)
	}
	if req.FixVersionIDs != nil {
		_ = uc.versionRepo.SetIssueVersions(ctx, id, req.FixVersionIDs)
	}
	if req.AffectsVersionIDs != nil {
		_ = uc.versionRepo.SetIssueAffectsVersions(ctx, id, req.AffectsVersionIDs)
	}

	h := &entity.IssueHistory{
		ID:        uuid.NewString(),
		IssueID:   id,
		UserID:    &actorID,
		Field:     "updated",
		NewValue:  map[string]any{"updated_by": actorID},
		CreatedAt: time.Now().UTC(),
	}
	_ = uc.issueRepo.CreateHistory(ctx, h)

	// Assignee o'zgargan bo'lsa notification yuborish
	if req.AssigneeID != nil && issue.AssigneeID != nil {
		newID := *issue.AssigneeID
		if prevAssigneeID == nil || *prevAssigneeID != newID {
			go uc.dispatcher.IssueAssigned(context.Background(), issue, newID, actorID)
		}
	}

	go func() {
		bg := context.Background()
		watchers, _ := uc.issueRepo.ListWatchers(bg, id)
		ids := make([]string, 0, len(watchers))
		for _, w := range watchers {
			ids = append(ids, w.UserID)
		}
		uc.dispatcher.IssueUpdated(bg, issue, ids, actorID)
	}()

	uc.log.Info(ctx, "issue updated", logger.String("id", id))
	return issue, nil
}

func (uc *useCase) Transition(ctx context.Context, id, statusID, actorID string) (*entity.Issue, error) {
	issue, err := uc.issueRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	project, err := uc.projectRepo.GetByID(ctx, issue.ProjectID)
	if err != nil {
		return nil, err
	}

	allowed, err := uc.workflowRepo.IsTransitionAllowed(ctx, project.WorkflowID, issue.StatusID, statusID)
	if err != nil {
		uc.log.Error(ctx, "issue.Transition: check failed", logger.String("id", id), logger.SafeString("err", err.Error()))
		return nil, fmt.Errorf("issue.Transition check: %w", err)
	}
	if !allowed {
		return nil, apperr.BadRequest("transition not allowed by workflow")
	}

	newStatus, err := uc.workflowRepo.GetStatusByID(ctx, statusID)
	if err != nil {
		return nil, fmt.Errorf("issue.Transition get status: %w", err)
	}

	oldStatusID := issue.StatusID
	if err := uc.issueRepo.UpdateStatus(ctx, id, statusID); err != nil {
		uc.log.Error(ctx, "issue.Transition: db error", logger.String("id", id), logger.SafeString("err", err.Error()))
		return nil, err
	}
	issue.StatusID = statusID

	// Auto-set resolution when moving to "done" category; clear when moving away.
	if newStatus.Category == "done" && issue.Resolution == nil {
		r := "fixed"
		_ = uc.issueRepo.UpdateResolution(ctx, id, &r)
		issue.Resolution = &r
	} else if newStatus.Category != "done" && issue.Resolution != nil {
		_ = uc.issueRepo.UpdateResolution(ctx, id, nil)
		issue.Resolution = nil
	}

	h := &entity.IssueHistory{
		ID:        uuid.NewString(),
		IssueID:   id,
		UserID:    &actorID,
		Field:     "status",
		OldValue:  map[string]any{"status_id": oldStatusID},
		NewValue:  map[string]any{"status_id": statusID},
		CreatedAt: time.Now().UTC(),
	}
	_ = uc.issueRepo.CreateHistory(ctx, h)

	go func() {
		bg := context.Background()
		watchers, _ := uc.issueRepo.ListWatchers(bg, id)
		ids := make([]string, 0, len(watchers))
		for _, w := range watchers {
			ids = append(ids, w.UserID)
		}
		uc.dispatcher.IssueStatusChanged(bg, issue, ids, actorID)
	}()

	if uc.hub != nil {
		uc.hub.BroadcastToRoom(ws.NewIssueUpdatedMsg(issue.ProjectID, map[string]any{
			"id": issue.ID, "status_id": statusID, "project_id": issue.ProjectID,
		}))
	}

	uc.log.Info(ctx, "issue transitioned", logger.String("id", id), logger.String("status_id", statusID))
	return issue, nil
}

func (uc *useCase) Delete(ctx context.Context, id, actorID string) error {
	issue, err := uc.issueRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if issue.ReporterID != actorID {
		uc.log.Warn(ctx, "issue.Delete: forbidden", logger.String("id", id), logger.String("actor_id", actorID))
		return apperr.Forbidden("only the reporter can delete this issue")
	}
	if err := uc.issueRepo.SoftDelete(ctx, id); err != nil {
		uc.log.Error(ctx, "issue.Delete: db error", logger.String("id", id), logger.SafeString("err", err.Error()))
		return err
	}
	uc.log.Info(ctx, "issue deleted", logger.String("id", id))
	return nil
}

func (uc *useCase) AddWatcher(ctx context.Context, issueID, userID string) error {
	w := &entity.IssueWatcher{IssueID: issueID, UserID: userID, CreatedAt: time.Now().UTC()}
	return uc.issueRepo.AddWatcher(ctx, w)
}

func (uc *useCase) RemoveWatcher(ctx context.Context, issueID, userID string) error {
	return uc.issueRepo.RemoveWatcher(ctx, issueID, userID)
}

func (uc *useCase) ListWatchers(ctx context.Context, issueID string) ([]*entity.IssueWatcher, error) {
	return uc.issueRepo.ListWatchers(ctx, issueID)
}

func (uc *useCase) ListHistory(ctx context.Context, issueID string, filter *entity.Filter) ([]*entity.IssueHistory, int, error) {
	return uc.issueRepo.ListHistory(ctx, issueID, filter)
}

func (uc *useCase) ReorderIssues(ctx context.Context, req *entity.ReorderIssuesReq) error {
	if err := uc.issueRepo.BulkUpdatePositions(ctx, req.Items); err != nil {
		uc.log.Error(ctx, "issue.ReorderIssues: db error", logger.SafeString("err", err.Error()))
		return fmt.Errorf("reorder issues: %w", err)
	}
	return nil
}

func (uc *useCase) BulkUpdate(ctx context.Context, req *entity.BulkUpdateIssueReq, actorID string) (*entity.BulkResult, error) {
	updated, err := uc.issueRepo.BulkUpdate(ctx, req)
	if err != nil {
		uc.log.Error(ctx, "issue.BulkUpdate: db error", logger.SafeString("err", err.Error()))
		return nil, fmt.Errorf("bulk update: %w", err)
	}

	failed := make([]string, 0)
	for _, id := range req.IssueIDs {
		found := false
		for _, uid := range updated {
			if uid == id {
				found = true
				break
			}
		}
		if !found {
			failed = append(failed, id)
		}
	}

	return &entity.BulkResult{
		Updated: updated,
		Failed:  failed,
		Total:   len(req.IssueIDs),
	}, nil
}

func (uc *useCase) BulkCreate(ctx context.Context, projectID string, req *entity.BulkCreateIssueReq, reporterID string) (*entity.BulkCreateResult, error) {
	project, err := uc.projectRepo.GetByID(ctx, projectID)
	if err != nil {
		return nil, err
	}

	wf, err := uc.workflowRepo.GetWithDetails(ctx, project.WorkflowID)
	if err != nil {
		return nil, fmt.Errorf("issue.BulkCreate get workflow: %w", err)
	}

	var initialStatusID string
	for _, s := range wf.Statuses {
		if s.IsInitial {
			initialStatusID = s.ID
			break
		}
	}
	if initialStatusID == "" && len(wf.Statuses) > 0 {
		initialStatusID = wf.Statuses[0].ID
	}

	// Allocate all counters in one DB round-trip.
	startCounter, err := uc.projectRepo.AllocateIssueCounters(ctx, projectID, len(req.Issues))
	if err != nil {
		return nil, fmt.Errorf("issue.BulkCreate allocate counters: %w", err)
	}

	now := time.Now().UTC()
	issues := make([]*entity.Issue, 0, len(req.Issues))
	result := &entity.BulkCreateResult{Total: len(req.Issues)}

	for i, r := range req.Issues {
		if r.ParentID != nil && *r.ParentID == "" {
			r.ParentID = nil
		}
		if r.AssigneeID != nil && *r.AssigneeID == "" {
			r.AssigneeID = nil
		}
		if r.SprintID != nil && *r.SprintID == "" {
			r.SprintID = nil
		}

		priority := r.Priority
		if priority == "" {
			priority = "medium"
		}
		cf := r.CustomFields
		if cf == nil {
			cf = map[string]any{}
		}

		issues = append(issues, &entity.Issue{
			ID:           uuid.NewString(),
			ProjectID:    projectID,
			IssueNumber:  int(startCounter) + i,
			Title:        r.Title,
			Description:  r.Description,
			Type:         r.Type,
			StatusID:     initialStatusID,
			Priority:     priority,
			AssigneeID:   r.AssigneeID,
			ReporterID:   reporterID,
			ParentID:     r.ParentID,
			SprintID:     r.SprintID,
			StoryPoints:  r.StoryPoints,
			DueDate:      r.DueDate,
			CustomFields: cf,
			CreatedAt:    now,
			UpdatedAt:    now,
		})
	}

	if err := uc.issueRepo.BulkCreate(ctx, issues); err != nil {
		uc.log.Error(ctx, "issue.BulkCreate: db error", logger.String("project_id", projectID), logger.SafeString("err", err.Error()))
		return nil, fmt.Errorf("issue.BulkCreate: %w", err)
	}

	result.Created = issues
	uc.log.Info(ctx, "issues bulk created", logger.String("project_id", projectID), logger.Int("count", len(issues)))
	return result, nil
}

func (uc *useCase) BulkDelete(ctx context.Context, req *entity.BulkDeleteIssueReq, actorID string) error {
	if err := uc.issueRepo.BulkDelete(ctx, req.IssueIDs); err != nil {
		uc.log.Error(ctx, "issue.BulkDelete: db error", logger.SafeString("err", err.Error()))
		return fmt.Errorf("bulk delete: %w", err)
	}
	return nil
}

func (uc *useCase) GetEpicProgress(ctx context.Context, epicID string) (*entity.EpicProgress, error) {
	return uc.issueRepo.GetEpicProgress(ctx, epicID)
}

func (uc *useCase) GetRoadmap(ctx context.Context, projectID string) ([]*entity.RoadmapItem, error) {
	return uc.issueRepo.GetRoadmap(ctx, projectID)
}

func (uc *useCase) GetBacklog(ctx context.Context, projectID string, filter *entity.IssueFilter) ([]*entity.Issue, int, error) {
	return uc.issueRepo.GetBacklog(ctx, projectID, filter)
}

func (uc *useCase) Clone(ctx context.Context, id, reporterID string, req *entity.CloneIssueReq) (*entity.Issue, error) {
	src, err := uc.issueRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	counter, err := uc.projectRepo.IncrementIssueCounter(ctx, src.ProjectID)
	if err != nil {
		uc.log.Error(ctx, "issue.Clone: increment counter failed", logger.SafeString("err", err.Error()))
		return nil, fmt.Errorf("issue.Clone increment counter: %w", err)
	}

	title := src.Title + " (copy)"
	if req.Title != nil && *req.Title != "" {
		title = *req.Title
	}

	sprintID := src.SprintID
	if req.SprintID != nil {
		sprintID = req.SprintID
	}

	now := time.Now().UTC()
	clone := &entity.Issue{
		ID:               uuid.NewString(),
		ProjectID:        src.ProjectID,
		IssueNumber:      int(counter),
		Title:            title,
		Description:      src.Description,
		Type:             src.Type,
		StatusID:         src.StatusID,
		Priority:         src.Priority,
		ReporterID:       reporterID,
		ParentID:         src.ParentID,
		SprintID:         sprintID,
		StoryPoints:      src.StoryPoints,
		DueDate:          src.DueDate,
		OriginalEstimate: src.OriginalEstimate,
		CustomFields:     src.CustomFields,
		CreatedAt:        now,
		UpdatedAt:        now,
	}
	if clone.CustomFields == nil {
		clone.CustomFields = map[string]any{}
	}

	if err := uc.issueRepo.Create(ctx, clone); err != nil {
		uc.log.Error(ctx, "issue.Clone: db error", logger.SafeString("err", err.Error()))
		return nil, err
	}

	if req.IncludeLinks {
		// Labels kopi
		labels, _ := uc.issueRepo.GetLabels(ctx, src.ID)
		if len(labels) > 0 {
			ids := make([]string, len(labels))
			for i, l := range labels {
				ids[i] = l.ID
			}
			_ = uc.issueRepo.SetLabels(ctx, clone.ID, ids)
		}
	}

	watcher := &entity.IssueWatcher{IssueID: clone.ID, UserID: reporterID, CreatedAt: now}
	_ = uc.issueRepo.AddWatcher(ctx, watcher)

	go uc.dispatcher.IssueCreated(context.Background(), clone, reporterID)

	uc.log.Info(ctx, "issue cloned", logger.String("src", id), logger.String("clone", clone.ID))
	return clone, nil
}

func (uc *useCase) RankBetween(ctx context.Context, id string, req *entity.RankIssueReq) error {
	issue, err := uc.issueRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	loRank, hiRank, err := uc.issueRepo.GetRankNeighbors(ctx, issue.ProjectID, req.Before, req.After)
	if err != nil {
		return fmt.Errorf("issue.RankBetween: get neighbors: %w", err)
	}

	newRank := lexorank.Between(loRank, hiRank)
	if err := uc.issueRepo.UpdateRank(ctx, id, newRank); err != nil {
		uc.log.Error(ctx, "issue.RankBetween: db error", logger.String("id", id), logger.SafeString("err", err.Error()))
		return fmt.Errorf("issue.RankBetween: %w", err)
	}
	return nil
}

func (uc *useCase) MoveOnBoard(ctx context.Context, issueID string, req *entity.MoveIssueReq, actorID string) (*entity.Issue, error) {
	issue, err := uc.issueRepo.GetByID(ctx, issueID)
	if err != nil {
		return nil, err
	}

	project, err := uc.projectRepo.GetByID(ctx, issue.ProjectID)
	if err != nil {
		return nil, err
	}

	wf, err := uc.workflowRepo.GetWithDetails(ctx, project.WorkflowID)
	if err != nil {
		return nil, fmt.Errorf("issue.MoveOnBoard get workflow: %w", err)
	}

	validStatus := false
	for _, s := range wf.Statuses {
		if s.ID == req.StatusID {
			validStatus = true
			break
		}
	}
	if !validStatus {
		return nil, apperr.BadRequest("status does not belong to this project's workflow")
	}

	oldStatusID := issue.StatusID

	if issue.StatusID != req.StatusID {
		if err := uc.issueRepo.UpdateStatus(ctx, issueID, req.StatusID); err != nil {
			uc.log.Error(ctx, "issue.MoveOnBoard: update status error",
				logger.String("id", issueID), logger.SafeString("err", err.Error()))
			return nil, err
		}
		issue.StatusID = req.StatusID

		h := &entity.IssueHistory{
			ID:        uuid.NewString(),
			IssueID:   issueID,
			UserID:    &actorID,
			Field:     "status",
			OldValue:  map[string]any{"status_id": oldStatusID},
			NewValue:  map[string]any{"status_id": req.StatusID},
			CreatedAt: time.Now().UTC(),
		}
		_ = uc.issueRepo.CreateHistory(ctx, h)

		go func() {
			bg := context.Background()
			watchers, _ := uc.issueRepo.ListWatchers(bg, issueID)
			ids := make([]string, 0, len(watchers))
			for _, w := range watchers {
				ids = append(ids, w.UserID)
			}
			uc.dispatcher.IssueStatusChanged(bg, issue, ids, actorID)
		}()
	}

	if req.Position >= 0 {
		_ = uc.issueRepo.BulkUpdatePositions(ctx, []entity.IssuePositionItem{
			{IssueID: issueID, Position: req.Position},
		})
		issue.Position = req.Position
	}

	if uc.hub != nil {
		uc.hub.BroadcastToRoom(ws.NewIssueMovedMsg(issue.ProjectID, map[string]any{
			"id": issue.ID, "status_id": issue.StatusID, "position": issue.Position, "project_id": issue.ProjectID,
		}))
	}

	uc.log.Info(ctx, "issue moved on board",
		logger.String("id", issueID), logger.String("status_id", req.StatusID))
	return issue, nil
}
