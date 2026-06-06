package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
	"github.com/jira-backend/jiraflow-backend/internal/infrastructura/repository"
	apperr "github.com/jira-backend/jiraflow-backend/internal/pkg/errors"
	"github.com/jira-backend/jiraflow-backend/internal/pkg/jql"
	pg "github.com/jira-backend/jiraflow-backend/internal/pkg/postgres"
)

type issueRepo struct {
	db      *pgxpool.Pool
	builder sq.StatementBuilderType
}

func NewIssueRepo(p *pg.Postgres) repository.IssueRepository {
	return &issueRepo{db: p.DB, builder: p.Builder}
}

const issueCols = `id, project_id, issue_number, title, description, type, status_id, priority,
	assignee_id, reporter_id, parent_id, sprint_id, story_points, due_date,
	original_estimate, remaining_estimate, custom_fields, resolution, position, created_at, updated_at, deleted_at`

func scanIssue(row pgx.Row) (*entity.Issue, error) {
	i := &entity.Issue{}
	var cfJSON []byte
	err := row.Scan(
		&i.ID, &i.ProjectID, &i.IssueNumber, &i.Title, &i.Description,
		&i.Type, &i.StatusID, &i.Priority, &i.AssigneeID, &i.ReporterID,
		&i.ParentID, &i.SprintID, &i.StoryPoints, &i.DueDate,
		&i.OriginalEstimate, &i.RemainingEstimate,
		&cfJSON, &i.Resolution, &i.Position, &i.CreatedAt, &i.UpdatedAt, &i.DeletedAt,
	)
	if err != nil {
		return nil, err
	}
	if len(cfJSON) > 0 {
		_ = json.Unmarshal(cfJSON, &i.CustomFields)
	}
	return i, nil
}

func (r *issueRepo) Create(ctx context.Context, issue *entity.Issue) error {
	cfJSON, err := json.Marshal(issue.CustomFields)
	if err != nil {
		return fmt.Errorf("issueRepo.Create marshal custom_fields: %w", err)
	}
	sql, args, err := r.builder.
		Insert("issues").
		Columns("id", "project_id", "issue_number", "title", "description", "type", "status_id", "priority",
			"assignee_id", "reporter_id", "parent_id", "sprint_id", "story_points", "due_date",
			"original_estimate", "remaining_estimate", "custom_fields", "created_at", "updated_at").
		Values(issue.ID, issue.ProjectID, issue.IssueNumber, issue.Title, issue.Description, issue.Type, issue.StatusID, issue.Priority,
			issue.AssigneeID, issue.ReporterID, issue.ParentID, issue.SprintID, issue.StoryPoints, issue.DueDate,
			issue.OriginalEstimate, issue.RemainingEstimate, cfJSON, issue.CreatedAt, issue.UpdatedAt).
		ToSql()
	if err != nil {
		return fmt.Errorf("issueRepo.Create: %w", err)
	}
	_, err = r.db.Exec(ctx, sql, args...)
	return err
}

func (r *issueRepo) GetByID(ctx context.Context, id string) (*entity.Issue, error) {
	sql, args, err := r.builder.
		Select(issueCols).From("issues").
		Where(sq.And{sq.Eq{"id": id}, sq.Eq{"deleted_at": nil}}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("issueRepo.GetByID: %w", err)
	}
	issue, err := scanIssue(r.db.QueryRow(ctx, sql, args...))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, apperr.NotFound("issue")
	}
	return issue, err
}

func (r *issueRepo) GetByKey(ctx context.Context, key string) (*entity.Issue, error) {
	// key = "PROJ-42" — proekt kaliti va issue raqamini join orqali qidiradi
	sql, args, err := r.builder.
		Select("i."+issueCols).
		From("issues i").
		Join("projects p ON p.id = i.project_id").
		Where(sq.And{
			sq.Expr("UPPER(p.key) || '-' || i.issue_number = ?", key),
			sq.Eq{"i.deleted_at": nil},
		}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("issueRepo.GetByKey build: %w", err)
	}
	issue, err := scanIssue(r.db.QueryRow(ctx, sql, args...))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, apperr.NotFound("issue")
	}
	return issue, err
}

func (r *issueRepo) applyIssueFilter(where sq.And, filter *entity.IssueFilter) sq.And {
	if filter.ProjectID != "" {
		where = append(where, sq.Eq{"project_id": filter.ProjectID})
	}
	if filter.SprintID != "" {
		where = append(where, sq.Eq{"sprint_id": filter.SprintID})
	}
	if filter.NoSprint {
		where = append(where, sq.Eq{"sprint_id": nil})
	}
	if filter.AssigneeID != "" {
		where = append(where, sq.Eq{"assignee_id": filter.AssigneeID})
	}
	if len(filter.AssigneeIDs) > 0 {
		where = append(where, sq.Eq{"assignee_id": filter.AssigneeIDs})
	}
	if filter.ReporterID != "" {
		where = append(where, sq.Eq{"reporter_id": filter.ReporterID})
	}
	if filter.StatusID != "" {
		where = append(where, sq.Eq{"status_id": filter.StatusID})
	}
	if len(filter.StatusIDs) > 0 {
		where = append(where, sq.Eq{"status_id": filter.StatusIDs})
	}
	if filter.Type != "" {
		where = append(where, sq.Eq{"type": filter.Type})
	}
	if len(filter.Types) > 0 {
		where = append(where, sq.Eq{"type": filter.Types})
	}
	if filter.Priority != "" {
		where = append(where, sq.Eq{"priority": filter.Priority})
	}
	if len(filter.Priorities) > 0 {
		where = append(where, sq.Eq{"priority": filter.Priorities})
	}
	if filter.ParentID != "" {
		where = append(where, sq.Eq{"parent_id": filter.ParentID})
	}
	if filter.EpicID != "" {
		where = append(where, sq.Eq{"parent_id": filter.EpicID})
	}
	if len(filter.LabelIDs) > 0 {
		where = append(where, sq.Expr(
			"id IN (SELECT issue_id FROM issue_labels WHERE label_id = ANY(?))",
			filter.LabelIDs,
		))
	}
	if len(filter.ComponentIDs) > 0 {
		where = append(where, sq.Expr(
			"id IN (SELECT issue_id FROM issue_components WHERE component_id = ANY(?))",
			filter.ComponentIDs,
		))
	}
	if len(filter.VersionIDs) > 0 {
		where = append(where, sq.Expr(
			"id IN (SELECT issue_id FROM issue_versions WHERE version_id = ANY(?) AND version_type = 'fix')",
			filter.VersionIDs,
		))
	}
	if len(filter.AffectsVersionIDs) > 0 {
		where = append(where, sq.Expr(
			"id IN (SELECT issue_id FROM issue_versions WHERE version_id = ANY(?) AND version_type = 'affects')",
			filter.AffectsVersionIDs,
		))
	}
	if filter.DueDateFrom != nil {
		where = append(where, sq.GtOrEq{"due_date": filter.DueDateFrom})
	}
	if filter.DueDateTo != nil {
		where = append(where, sq.LtOrEq{"due_date": filter.DueDateTo})
	}
	if filter.CreatedAfter != nil {
		where = append(where, sq.GtOrEq{"created_at": filter.CreatedAfter})
	}
	if filter.CreatedBefore != nil {
		where = append(where, sq.LtOrEq{"created_at": filter.CreatedBefore})
	}
	if filter.TextSearch != "" {
		where = append(where, sq.Expr(
			"search_vector @@ plainto_tsquery('simple', ?)",
			filter.TextSearch,
		))
	}
	if filter.Search != "" {
		where = append(where, sq.ILike{"title": "%" + filter.Search + "%"})
	}
	if filter.JQL != "" {
		if q, err := jql.Parse(filter.JQL); err == nil {
			if s := jql.ToSqlizer(q, filter.CurrentUserID); s != nil {
				where = append(where, s)
			}
		}
	}
	return where
}

func (r *issueRepo) List(ctx context.Context, filter *entity.IssueFilter) ([]*entity.Issue, int, error) {
	where := r.applyIssueFilter(sq.And{sq.Eq{"deleted_at": nil}}, filter)

	var total int
	cntSQL, cntArgs, _ := r.builder.Select("COUNT(*)").From("issues").Where(where).ToSql()
	if err := r.db.QueryRow(ctx, cntSQL, cntArgs...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("issueRepo.List count: %w", err)
	}

	q := r.builder.
		Select(issueCols).From("issues").Where(where).
		Limit(uint64(filter.GetLimit())).Offset(uint64(filter.Offset()))

	if filter.TextSearch != "" {
		q = q.OrderBy("ts_rank(search_vector, plainto_tsquery('simple', ?)) DESC", filter.TextSearch)
	} else {
		q = q.OrderBy("position ASC, created_at DESC")
	}

	dataSQL, dataArgs, _ := q.ToSql()
	rows, err := r.db.Query(ctx, dataSQL, dataArgs...)
	if err != nil {
		return nil, 0, fmt.Errorf("issueRepo.List query: %w", err)
	}
	defer rows.Close()

	var issues []*entity.Issue
	for rows.Next() {
		issue, err := scanIssue(rows)
		if err != nil {
			return nil, 0, err
		}
		issues = append(issues, issue)
	}
	return issues, total, rows.Err()
}

func (r *issueRepo) Update(ctx context.Context, issue *entity.Issue) error {
	cfJSON, err := json.Marshal(issue.CustomFields)
	if err != nil {
		return fmt.Errorf("issueRepo.Update marshal custom_fields: %w", err)
	}
	sql, args, err := r.builder.
		Update("issues").
		Set("title", issue.Title).Set("description", issue.Description).
		Set("priority", issue.Priority).Set("assignee_id", issue.AssigneeID).
		Set("sprint_id", issue.SprintID).Set("story_points", issue.StoryPoints).
		Set("due_date", issue.DueDate).
		Set("original_estimate", issue.OriginalEstimate).
		Set("remaining_estimate", issue.RemainingEstimate).
		Set("custom_fields", cfJSON).
		Set("resolution", issue.Resolution).
		Set("updated_at", sq.Expr("NOW()")).
		Where(sq.And{sq.Eq{"id": issue.ID}, sq.Eq{"deleted_at": nil}}).ToSql()
	if err != nil {
		return fmt.Errorf("issueRepo.Update: %w", err)
	}
	_, err = r.db.Exec(ctx, sql, args...)
	return err
}

func (r *issueRepo) UpdateStatus(ctx context.Context, id, statusID string) error {
	sql, args, err := r.builder.
		Update("issues").
		Set("status_id", statusID).Set("updated_at", sq.Expr("NOW()")).
		Where(sq.And{sq.Eq{"id": id}, sq.Eq{"deleted_at": nil}}).ToSql()
	if err != nil {
		return fmt.Errorf("issueRepo.UpdateStatus: %w", err)
	}
	_, err = r.db.Exec(ctx, sql, args...)
	return err
}

func (r *issueRepo) UpdateResolution(ctx context.Context, id string, resolution *string) error {
	sql, args, err := r.builder.
		Update("issues").
		Set("resolution", resolution).Set("updated_at", sq.Expr("NOW()")).
		Where(sq.And{sq.Eq{"id": id}, sq.Eq{"deleted_at": nil}}).ToSql()
	if err != nil {
		return fmt.Errorf("issueRepo.UpdateResolution: %w", err)
	}
	_, err = r.db.Exec(ctx, sql, args...)
	return err
}

func (r *issueRepo) SoftDelete(ctx context.Context, id string) error {
	sql, args, err := r.builder.
		Update("issues").Set("deleted_at", sq.Expr("NOW()")).
		Where(sq.And{sq.Eq{"id": id}, sq.Eq{"deleted_at": nil}}).ToSql()
	if err != nil {
		return fmt.Errorf("issueRepo.SoftDelete: %w", err)
	}
	_, err = r.db.Exec(ctx, sql, args...)
	return err
}

func (r *issueRepo) CountByProject(ctx context.Context, projectID string) (int, error) {
	var count int
	err := r.db.QueryRow(ctx,
		`SELECT COUNT(*) FROM issues WHERE project_id=$1 AND deleted_at IS NULL`, projectID,
	).Scan(&count)
	return count, err
}

// ─── Labels ───────────────────────────────────────────────────────────────────

func (r *issueRepo) SetLabels(ctx context.Context, issueID string, labelIDs []string) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("issueRepo.SetLabels begin: %w", err)
	}
	defer tx.Rollback(ctx)

	if _, err = tx.Exec(ctx, `DELETE FROM issue_labels WHERE issue_id=$1`, issueID); err != nil {
		return fmt.Errorf("issueRepo.SetLabels delete: %w", err)
	}
	for _, lid := range labelIDs {
		if _, err = tx.Exec(ctx,
			`INSERT INTO issue_labels(issue_id, label_id) VALUES($1,$2) ON CONFLICT DO NOTHING`,
			issueID, lid,
		); err != nil {
			return fmt.Errorf("issueRepo.SetLabels insert: %w", err)
		}
	}
	return tx.Commit(ctx)
}

func (r *issueRepo) GetLabels(ctx context.Context, issueID string) ([]*entity.Label, error) {
	rows, err := r.db.Query(ctx, `
		SELECT l.id, l.project_id, l.name, l.color, l.created_at
		FROM labels l
		JOIN issue_labels il ON il.label_id = l.id
		WHERE il.issue_id = $1
		ORDER BY l.name ASC
	`, issueID)
	if err != nil {
		return nil, fmt.Errorf("issueRepo.GetLabels: %w", err)
	}
	defer rows.Close()

	var labels []*entity.Label
	for rows.Next() {
		l := &entity.Label{}
		if err := rows.Scan(&l.ID, &l.ProjectID, &l.Name, &l.Color, &l.CreatedAt); err != nil {
			return nil, err
		}
		labels = append(labels, l)
	}
	return labels, rows.Err()
}

// ─── Watchers ─────────────────────────────────────────────────────────────────

func (r *issueRepo) AddWatcher(ctx context.Context, w *entity.IssueWatcher) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO issue_watchers(issue_id, user_id, created_at) VALUES($1,$2,$3) ON CONFLICT DO NOTHING`,
		w.IssueID, w.UserID, w.CreatedAt,
	)
	return err
}

func (r *issueRepo) RemoveWatcher(ctx context.Context, issueID, userID string) error {
	_, err := r.db.Exec(ctx,
		`DELETE FROM issue_watchers WHERE issue_id=$1 AND user_id=$2`, issueID, userID,
	)
	return err
}

func (r *issueRepo) ListWatchers(ctx context.Context, issueID string) ([]*entity.IssueWatcher, error) {
	rows, err := r.db.Query(ctx, `
		SELECT iw.issue_id, iw.user_id, iw.created_at,
		       u.id, u.full_name, u.email, u.avatar_url, u.color
		FROM issue_watchers iw
		JOIN users u ON u.id = iw.user_id
		WHERE iw.issue_id = $1
		ORDER BY iw.created_at ASC
	`, issueID)
	if err != nil {
		return nil, fmt.Errorf("issueRepo.ListWatchers: %w", err)
	}
	defer rows.Close()

	var watchers []*entity.IssueWatcher
	for rows.Next() {
		w := &entity.IssueWatcher{User: &entity.UserShort{}}
		if err := rows.Scan(
			&w.IssueID, &w.UserID, &w.CreatedAt,
			&w.User.ID, &w.User.FullName, &w.User.Email, &w.User.AvatarURL, &w.User.Color,
		); err != nil {
			return nil, err
		}
		watchers = append(watchers, w)
	}
	return watchers, rows.Err()
}

func (r *issueRepo) IsWatcher(ctx context.Context, issueID, userID string) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx,
		`SELECT EXISTS(SELECT 1 FROM issue_watchers WHERE issue_id=$1 AND user_id=$2)`,
		issueID, userID,
	).Scan(&exists)
	return exists, err
}

// ─── History ──────────────────────────────────────────────────────────────────

func (r *issueRepo) CreateHistory(ctx context.Context, h *entity.IssueHistory) error {
	oldJSON, err := json.Marshal(h.OldValue)
	if err != nil {
		return fmt.Errorf("issueRepo.CreateHistory marshal old_value: %w", err)
	}
	newJSON, err := json.Marshal(h.NewValue)
	if err != nil {
		return fmt.Errorf("issueRepo.CreateHistory marshal new_value: %w", err)
	}
	_, err = r.db.Exec(ctx,
		`INSERT INTO issue_history(id, issue_id, user_id, field, old_value, new_value, created_at)
		 VALUES($1,$2,$3,$4,$5,$6,$7)`,
		h.ID, h.IssueID, h.UserID, h.Field, oldJSON, newJSON, h.CreatedAt,
	)
	return err
}

func (r *issueRepo) ListHistory(ctx context.Context, issueID string, filter *entity.Filter) ([]*entity.IssueHistory, int, error) {
	var total int
	if err := r.db.QueryRow(ctx,
		`SELECT COUNT(*) FROM issue_history WHERE issue_id=$1`, issueID,
	).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("issueRepo.ListHistory count: %w", err)
	}

	rows, err := r.db.Query(ctx, `
		SELECT h.id, h.issue_id, h.user_id, h.field, h.old_value, h.new_value, h.created_at,
		       u.id, u.full_name, u.email, u.avatar_url, u.color
		FROM issue_history h
		LEFT JOIN users u ON u.id = h.user_id
		WHERE h.issue_id = $1
		ORDER BY h.created_at DESC
		LIMIT $2 OFFSET $3
	`, issueID, filter.GetLimit(), filter.Offset())
	if err != nil {
		return nil, 0, fmt.Errorf("issueRepo.ListHistory query: %w", err)
	}
	defer rows.Close()

	var history []*entity.IssueHistory
	for rows.Next() {
		h := &entity.IssueHistory{}
		var oldJSON, newJSON []byte
		u := &entity.UserShort{}
		if err := rows.Scan(
			&h.ID, &h.IssueID, &h.UserID, &h.Field, &oldJSON, &newJSON, &h.CreatedAt,
			&u.ID, &u.FullName, &u.Email, &u.AvatarURL, &u.Color,
		); err != nil {
			return nil, 0, err
		}
		_ = json.Unmarshal(oldJSON, &h.OldValue)
		_ = json.Unmarshal(newJSON, &h.NewValue)
		if u.ID != "" {
			h.User = u
		}
		history = append(history, h)
	}
	return history, total, rows.Err()
}

func (r *issueRepo) BulkUpdatePositions(ctx context.Context, items []entity.IssuePositionItem) error {
	if len(items) == 0 {
		return nil
	}
	batch := &pgx.Batch{}
	for _, item := range items {
		batch.Queue(
			"UPDATE issues SET position=$1, updated_at=NOW() WHERE id=$2 AND deleted_at IS NULL",
			item.Position, item.IssueID,
		)
	}
	br := r.db.SendBatch(ctx, batch)
	defer br.Close()
	for range items {
		if _, err := br.Exec(); err != nil {
			return fmt.Errorf("issueRepo.BulkUpdatePositions: %w", err)
		}
	}
	return nil
}

func (r *issueRepo) BulkUpdate(ctx context.Context, req *entity.BulkUpdateIssueReq) ([]string, error) {
	if len(req.IssueIDs) == 0 {
		return nil, nil
	}
	upd := r.builder.Update("issues").Set("updated_at", sq.Expr("NOW()")).
		Where(sq.And{sq.Eq{"id": req.IssueIDs}, sq.Eq{"deleted_at": nil}})

	if req.AssigneeID != nil {
		upd = upd.Set("assignee_id", req.AssigneeID)
	}
	if req.StatusID != nil {
		upd = upd.Set("status_id", req.StatusID)
	}
	if req.Priority != nil {
		upd = upd.Set("priority", req.Priority)
	}
	if req.SprintID != nil {
		upd = upd.Set("sprint_id", req.SprintID)
	}

	sql, args, err := upd.Suffix("RETURNING id").ToSql()
	if err != nil {
		return nil, fmt.Errorf("issueRepo.BulkUpdate build: %w", err)
	}
	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("issueRepo.BulkUpdate exec: %w", err)
	}
	defer rows.Close()

	var updated []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		updated = append(updated, id)
	}

	// Labels va components batch update
	if len(req.LabelIDs) > 0 {
		batch := &pgx.Batch{}
		for _, id := range updated {
			batch.Queue(`DELETE FROM issue_labels WHERE issue_id=$1`, id)
			for _, lid := range req.LabelIDs {
				batch.Queue(`INSERT INTO issue_labels(issue_id,label_id) VALUES($1,$2) ON CONFLICT DO NOTHING`, id, lid)
			}
		}
		br := r.db.SendBatch(ctx, batch)
		br.Close()
	}
	if len(req.ComponentIDs) > 0 {
		batch := &pgx.Batch{}
		for _, id := range updated {
			batch.Queue(`DELETE FROM issue_components WHERE issue_id=$1`, id)
			for _, cid := range req.ComponentIDs {
				batch.Queue(`INSERT INTO issue_components(issue_id,component_id) VALUES($1,$2) ON CONFLICT DO NOTHING`, id, cid)
			}
		}
		br := r.db.SendBatch(ctx, batch)
		br.Close()
	}

	return updated, rows.Err()
}

func (r *issueRepo) BulkDelete(ctx context.Context, issueIDs []string) error {
	if len(issueIDs) == 0 {
		return nil
	}
	sql, args, err := r.builder.
		Update("issues").Set("deleted_at", sq.Expr("NOW()")).
		Where(sq.And{sq.Eq{"id": issueIDs}, sq.Eq{"deleted_at": nil}}).ToSql()
	if err != nil {
		return fmt.Errorf("issueRepo.BulkDelete: %w", err)
	}
	_, err = r.db.Exec(ctx, sql, args...)
	return err
}

func (r *issueRepo) UpdateEstimates(ctx context.Context, issueID string, original, remaining *int) error {
	upd := r.builder.Update("issues").Set("updated_at", sq.Expr("NOW()")).
		Where(sq.And{sq.Eq{"id": issueID}, sq.Eq{"deleted_at": nil}})
	if original != nil {
		upd = upd.Set("original_estimate", original)
	}
	if remaining != nil {
		upd = upd.Set("remaining_estimate", remaining)
	}
	sql, args, err := upd.ToSql()
	if err != nil {
		return fmt.Errorf("issueRepo.UpdateEstimates: %w", err)
	}
	_, err = r.db.Exec(ctx, sql, args...)
	return err
}

func (r *issueRepo) GetEpicProgress(ctx context.Context, epicID string) (*entity.EpicProgress, error) {
	var total, done int
	err := r.db.QueryRow(ctx, `
		SELECT
			COUNT(*)                                              AS total,
			COUNT(*) FILTER (WHERE ws.category = 'done')         AS done
		FROM issues i
		JOIN workflow_statuses ws ON ws.id = i.status_id
		WHERE i.parent_id = $1 AND i.deleted_at IS NULL`, epicID,
	).Scan(&total, &done)
	if err != nil {
		return nil, fmt.Errorf("issueRepo.GetEpicProgress: %w", err)
	}
	progress := 0.0
	if total > 0 {
		progress = float64(done) / float64(total) * 100
	}
	return &entity.EpicProgress{Total: total, Done: done, Progress: progress}, nil
}

func (r *issueRepo) GetRoadmap(ctx context.Context, projectID string) ([]*entity.RoadmapItem, error) {
	rows, err := r.db.Query(ctx, `
		SELECT
			i.id, i.issue_number, i.title, i.status_id, i.priority,
			i.assignee_id, i.due_date,
			s.start_date, s.end_date,
			COALESCE(
				(SELECT COUNT(*) FILTER (WHERE ws2.category='done')
				 FROM issues sub
				 JOIN workflow_statuses ws2 ON ws2.id=sub.status_id
				 WHERE sub.parent_id=i.id AND sub.deleted_at IS NULL)::float
				/ NULLIF((SELECT COUNT(*) FROM issues sub WHERE sub.parent_id=i.id AND sub.deleted_at IS NULL),0) * 100,
			0) AS progress
		FROM issues i
		LEFT JOIN sprints s ON s.id = i.sprint_id AND s.deleted_at IS NULL
		WHERE i.project_id=$1 AND i.type='epic' AND i.deleted_at IS NULL
		ORDER BY COALESCE(s.start_date, i.due_date) ASC NULLS LAST`, projectID,
	)
	if err != nil {
		return nil, fmt.Errorf("issueRepo.GetRoadmap: %w", err)
	}
	defer rows.Close()

	var items []*entity.RoadmapItem
	for rows.Next() {
		item := &entity.RoadmapItem{}
		var dueDate, sprintEnd *time.Time
		if err := rows.Scan(
			&item.ID, &item.IssueNumber, &item.Title, &item.StatusID, &item.Priority,
			&item.AssigneeID, &dueDate, &item.StartDate, &sprintEnd, &item.Progress,
		); err != nil {
			return nil, err
		}
		if sprintEnd != nil {
			item.EndDate = sprintEnd
		} else {
			item.EndDate = dueDate
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *issueRepo) GetBacklog(ctx context.Context, projectID string, filter *entity.IssueFilter) ([]*entity.Issue, int, error) {
	filter.ProjectID = projectID
	filter.NoSprint = true
	return r.List(ctx, filter)
}

func (r *issueRepo) UpdateRank(ctx context.Context, id, rank string) error {
	_, err := r.db.Exec(ctx,
		`UPDATE issues SET rank=$1, updated_at=NOW() WHERE id=$2 AND deleted_at IS NULL`,
		rank, id,
	)
	return err
}

func (r *issueRepo) GetRankNeighbors(ctx context.Context, projectID string, beforeID, afterID *string) (string, string, error) {
	var loRank, hiRank string

	if afterID != nil {
		if err := r.db.QueryRow(ctx,
			`SELECT COALESCE(rank,'') FROM issues WHERE id=$1 AND deleted_at IS NULL`,
			*afterID,
		).Scan(&loRank); err != nil {
			return "", "", fmt.Errorf("issueRepo.GetRankNeighbors after: %w", err)
		}
	}
	if beforeID != nil {
		if err := r.db.QueryRow(ctx,
			`SELECT COALESCE(rank,'') FROM issues WHERE id=$1 AND deleted_at IS NULL`,
			*beforeID,
		).Scan(&hiRank); err != nil {
			return "", "", fmt.Errorf("issueRepo.GetRankNeighbors before: %w", err)
		}
	}
	return loRank, hiRank, nil
}
