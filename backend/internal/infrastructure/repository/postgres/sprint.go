package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
	"github.com/jira-backend/jiraflow-backend/internal/infrastructure/repository"
	apperr "github.com/jira-backend/jiraflow-backend/internal/pkg/errors"
	pg "github.com/jira-backend/jiraflow-backend/internal/pkg/postgres"
)

type sprintRepo struct {
	db      *pgxpool.Pool
	builder sq.StatementBuilderType
}

func NewSprintRepo(p *pg.Postgres) repository.SprintRepository {
	return &sprintRepo{db: p.DB, builder: p.Builder}
}

const sprintCols = "id, project_id, name, goal, status, start_date, end_date, started_at, completed_at, created_by, created_at, updated_at, deleted_at"

func scanSprint(row pgx.Row) (*entity.Sprint, error) {
	s := &entity.Sprint{}
	err := row.Scan(
		&s.ID, &s.ProjectID, &s.Name, &s.Goal, &s.Status,
		&s.StartDate, &s.EndDate, &s.StartedAt, &s.CompletedAt,
		&s.CreatedBy, &s.CreatedAt, &s.UpdatedAt, &s.DeletedAt,
	)
	return s, err
}

func (r *sprintRepo) Create(ctx context.Context, s *entity.Sprint) error {
	sql, args, err := r.builder.
		Insert("sprints").
		Columns("id", "project_id", "name", "goal", "status", "start_date", "end_date", "created_by", "created_at", "updated_at").
		Values(s.ID, s.ProjectID, s.Name, s.Goal, s.Status, s.StartDate, s.EndDate, s.CreatedBy, s.CreatedAt, s.UpdatedAt).
		ToSql()
	if err != nil {
		return fmt.Errorf("sprintRepo.Create: %w", err)
	}
	_, err = r.db.Exec(ctx, sql, args...)
	return err
}

func (r *sprintRepo) GetByID(ctx context.Context, id string) (*entity.Sprint, error) {
	sql, args, err := r.builder.
		Select(sprintCols).From("sprints").
		Where(sq.And{sq.Eq{"id": id}, sq.Eq{"deleted_at": nil}}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("sprintRepo.GetByID: %w", err)
	}
	s, err := scanSprint(r.db.QueryRow(ctx, sql, args...))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, apperr.NotFound("sprint")
	}
	return s, err
}

func (r *sprintRepo) List(ctx context.Context, projectID string, filter *entity.SprintFilter) ([]*entity.Sprint, int, error) {
	where := sq.And{sq.Eq{"project_id": projectID}, sq.Eq{"deleted_at": nil}}
	if filter.Status != "" {
		where = append(where, sq.Eq{"status": filter.Status})
	}

	var total int
	cntSQL, cntArgs, _ := r.builder.Select("COUNT(*)").From("sprints").Where(where).ToSql()
	if err := r.db.QueryRow(ctx, cntSQL, cntArgs...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("sprintRepo.List count: %w", err)
	}

	dataSQL, dataArgs, _ := r.builder.
		Select(sprintCols).From("sprints").Where(where).
		OrderBy("created_at DESC").
		Limit(uint64(filter.GetLimit())).Offset(uint64(filter.Offset())).ToSql()

	rows, err := r.db.Query(ctx, dataSQL, dataArgs...)
	if err != nil {
		return nil, 0, fmt.Errorf("sprintRepo.List query: %w", err)
	}
	defer rows.Close()

	var sprints []*entity.Sprint
	for rows.Next() {
		s, err := scanSprint(rows)
		if err != nil {
			return nil, 0, err
		}
		sprints = append(sprints, s)
	}
	return sprints, total, rows.Err()
}

func (r *sprintRepo) Update(ctx context.Context, s *entity.Sprint) error {
	sql, args, err := r.builder.
		Update("sprints").
		Set("name", s.Name).Set("goal", s.Goal).
		Set("start_date", s.StartDate).Set("end_date", s.EndDate).
		Set("updated_at", sq.Expr("NOW()")).
		Where(sq.And{sq.Eq{"id": s.ID}, sq.Eq{"deleted_at": nil}}).ToSql()
	if err != nil {
		return fmt.Errorf("sprintRepo.Update: %w", err)
	}
	_, err = r.db.Exec(ctx, sql, args...)
	return err
}

func (r *sprintRepo) SoftDelete(ctx context.Context, id string) error {
	sql, args, err := r.builder.
		Update("sprints").Set("deleted_at", sq.Expr("NOW()")).
		Where(sq.And{sq.Eq{"id": id}, sq.Eq{"deleted_at": nil}}).ToSql()
	if err != nil {
		return fmt.Errorf("sprintRepo.SoftDelete: %w", err)
	}
	_, err = r.db.Exec(ctx, sql, args...)
	return err
}

func (r *sprintRepo) GetActive(ctx context.Context, projectID string) (*entity.Sprint, error) {
	sql, args, err := r.builder.
		Select(sprintCols).From("sprints").
		Where(sq.And{sq.Eq{"project_id": projectID}, sq.Eq{"status": "active"}, sq.Eq{"deleted_at": nil}}).
		Limit(1).ToSql()
	if err != nil {
		return nil, fmt.Errorf("sprintRepo.GetActive: %w", err)
	}
	s, err := scanSprint(r.db.QueryRow(ctx, sql, args...))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, apperr.NotFound("active sprint")
	}
	return s, err
}

func (r *sprintRepo) Start(ctx context.Context, id string, startedAt time.Time) error {
	sql, args, err := r.builder.
		Update("sprints").
		Set("status", "active").Set("started_at", startedAt).
		Set("updated_at", sq.Expr("NOW()")).
		Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		return fmt.Errorf("sprintRepo.Start: %w", err)
	}
	_, err = r.db.Exec(ctx, sql, args...)
	return err
}

func (r *sprintRepo) AddIssue(ctx context.Context, sprintID, issueID string) error {
	sql, args, err := r.builder.
		Update("issues").
		Set("sprint_id", sprintID).
		Set("updated_at", sq.Expr("NOW()")).
		Where(sq.And{sq.Eq{"id": issueID}, sq.Eq{"deleted_at": nil}}).
		ToSql()
	if err != nil {
		return fmt.Errorf("sprintRepo.AddIssue: %w", err)
	}
	_, err = r.db.Exec(ctx, sql, args...)
	return err
}

func (r *sprintRepo) BulkAddIssues(ctx context.Context, sprintID string, issueIDs []string) error {
	if len(issueIDs) == 0 {
		return nil
	}
	_, err := r.db.Exec(ctx,
		`UPDATE issues SET sprint_id = $1, updated_at = NOW()
		 WHERE id = ANY($2::uuid[]) AND deleted_at IS NULL`,
		sprintID, issueIDs,
	)
	return err
}

func (r *sprintRepo) RemoveIssue(ctx context.Context, sprintID, issueID string) error {
	sql, args, err := r.builder.
		Update("issues").
		Set("sprint_id", nil).
		Set("updated_at", sq.Expr("NOW()")).
		Where(sq.And{sq.Eq{"id": issueID}, sq.Eq{"sprint_id": sprintID}, sq.Eq{"deleted_at": nil}}).
		ToSql()
	if err != nil {
		return fmt.Errorf("sprintRepo.RemoveIssue: %w", err)
	}
	_, err = r.db.Exec(ctx, sql, args...)
	return err
}

func (r *sprintRepo) Complete(ctx context.Context, id string, completedAt time.Time) error {
	sql, args, err := r.builder.
		Update("sprints").
		Set("status", "completed").Set("completed_at", completedAt).
		Set("updated_at", sq.Expr("NOW()")).
		Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		return fmt.Errorf("sprintRepo.Complete: %w", err)
	}
	_, err = r.db.Exec(ctx, sql, args...)
	return err
}

func (r *sprintRepo) GetReport(ctx context.Context, sprintID string) (*entity.SprintReport, error) {
	sprint, err := r.GetByID(ctx, sprintID)
	if err != nil {
		return nil, err
	}

	var total, completed, totalSP, completedSP int
	err = r.db.QueryRow(ctx, `
		SELECT
			COUNT(*)                                                                 AS total,
			COUNT(*) FILTER (WHERE ws.category = 'done')                            AS completed,
			COALESCE(SUM(i.story_points), 0)                                        AS total_sp,
			COALESCE(SUM(i.story_points) FILTER (WHERE ws.category = 'done'), 0)   AS completed_sp
		FROM issues i
		JOIN workflow_statuses ws ON ws.id = i.status_id
		WHERE i.sprint_id = $1 AND i.deleted_at IS NULL`, sprintID,
	).Scan(&total, &completed, &totalSP, &completedSP)
	if err != nil {
		return nil, fmt.Errorf("sprintRepo.GetReport: %w", err)
	}

	completionRate := 0.0
	if total > 0 {
		completionRate = float64(completed) / float64(total)
	}

	return &entity.SprintReport{
		SprintID:             sprintID,
		SprintName:           sprint.Name,
		StartDate:            sprint.StartDate,
		EndDate:              sprint.EndDate,
		CompletedAt:          sprint.CompletedAt,
		TotalIssues:          total,
		CompletedIssues:      completed,
		IncompleteIssues:     total - completed,
		TotalStoryPoints:     totalSP,
		CompletedStoryPoints: completedSP,
		CompletionRate:       completionRate,
	}, nil
}

func (r *sprintRepo) GetBurndown(ctx context.Context, sprintID string) (*entity.BurndownChart, error) {
	sprint, err := r.GetByID(ctx, sprintID)
	if err != nil {
		return nil, err
	}
	if sprint.StartDate == nil || sprint.EndDate == nil {
		return &entity.BurndownChart{SprintID: sprintID}, nil
	}

	var totalPoints int
	r.db.QueryRow(ctx,
		`SELECT COALESCE(SUM(story_points),0) FROM issues WHERE sprint_id=$1 AND deleted_at IS NULL`,
		sprintID,
	).Scan(&totalPoints)

	rows, err := r.db.Query(ctx, `
		SELECT
			gs::date                                    AS day,
			COALESCE(SUM(i.story_points) FILTER (
				WHERE ws.category != 'done'
				  OR NOT EXISTS (
					SELECT 1 FROM issue_history ih
					WHERE ih.issue_id = i.id
					  AND ih.field = 'status'
					  AND (ih.new_value->>'status_id') IN (
						SELECT id::text FROM workflow_statuses WHERE category='done'
					  )
					  AND ih.created_at::date <= gs::date
				  )
			), $3) AS remaining
		FROM generate_series($1::date, $2::date, '1 day'::interval) gs
		LEFT JOIN issues i ON i.sprint_id = $4 AND i.deleted_at IS NULL
		LEFT JOIN workflow_statuses ws ON ws.id = i.status_id
		GROUP BY gs
		ORDER BY gs`, sprint.StartDate, sprint.EndDate, totalPoints, sprintID,
	)
	if err != nil {
		return nil, fmt.Errorf("sprintRepo.GetBurndown query: %w", err)
	}
	defer rows.Close()

	totalDays := sprint.EndDate.Sub(*sprint.StartDate).Hours() / 24
	var points []entity.BurndownPoint
	day := 0
	for rows.Next() {
		var p entity.BurndownPoint
		if err := rows.Scan(&p.Date, &p.RemainingPoints); err != nil {
			return nil, err
		}
		idealRemaining := float64(totalPoints)
		if totalDays > 0 {
			idealRemaining = float64(totalPoints) * (1 - float64(day)/totalDays)
		}
		p.IdealPoints = idealRemaining
		points = append(points, p)
		day++
	}

	return &entity.BurndownChart{
		SprintID:    sprintID,
		TotalPoints: totalPoints,
		Points:      points,
	}, rows.Err()
}

func (r *sprintRepo) GetBurnup(ctx context.Context, sprintID string) (*entity.BurnupChart, error) {
	sprint, err := r.GetByID(ctx, sprintID)
	if err != nil {
		return nil, err
	}
	if sprint.StartDate == nil || sprint.EndDate == nil {
		return &entity.BurnupChart{SprintID: sprintID}, nil
	}

	var totalPoints int
	r.db.QueryRow(ctx,
		`SELECT COALESCE(SUM(story_points),0) FROM issues WHERE sprint_id=$1 AND deleted_at IS NULL`,
		sprintID,
	).Scan(&totalPoints)

	rows, err := r.db.Query(ctx, `
		SELECT
			gs::date AS day,
			COALESCE(SUM(i.story_points) FILTER (
				WHERE EXISTS (
					SELECT 1 FROM issue_history ih
					WHERE ih.issue_id = i.id
					  AND ih.field = 'status'
					  AND (ih.new_value->>'status_id') IN (
						SELECT id::text FROM workflow_statuses WHERE category='done'
					  )
					  AND ih.created_at::date <= gs::date
				)
			), 0) AS completed,
			COALESCE(COUNT(i.id) FILTER (WHERE i.created_at::date <= gs::date) *
				COALESCE(AVG(i.story_points),0), $3) AS scope
		FROM generate_series($1::date, $2::date, '1 day'::interval) gs
		LEFT JOIN issues i ON i.sprint_id = $4 AND i.deleted_at IS NULL
		GROUP BY gs
		ORDER BY gs`, sprint.StartDate, sprint.EndDate, totalPoints, sprintID,
	)
	if err != nil {
		return nil, fmt.Errorf("sprintRepo.GetBurnup query: %w", err)
	}
	defer rows.Close()

	var points []entity.BurnupPoint
	for rows.Next() {
		var p entity.BurnupPoint
		var scope float64
		if err := rows.Scan(&p.Date, &p.CompletedPoints, &scope); err != nil {
			return nil, err
		}
		p.TotalScope = int(scope)
		points = append(points, p)
	}

	return &entity.BurnupChart{
		SprintID:    sprintID,
		TotalPoints: totalPoints,
		Points:      points,
	}, rows.Err()
}

func (r *sprintRepo) GetCFD(ctx context.Context, projectID string, from, to *string) (*entity.CFDChart, error) {
	fromDate := "NOW() - INTERVAL '30 days'"
	toDate := "NOW()"
	args := []any{projectID}
	if from != nil {
		fromDate = "$2"
		args = append(args, *from)
	}
	if to != nil {
		toDate = fmt.Sprintf("$%d", len(args)+1)
		args = append(args, *to)
	}

	// Fetch distinct statuses for project
	statusRows, err := r.db.Query(ctx, `
		SELECT DISTINCT ws.id, ws.name FROM workflow_statuses ws
		JOIN workflows w ON w.id = ws.workflow_id
		JOIN projects p ON p.workflow_id = w.id
		WHERE p.id = $1
		ORDER BY ws.name`, projectID)
	if err != nil {
		return nil, err
	}
	defer statusRows.Close()

	type statusInfo struct{ id, name string }
	var statuses []statusInfo
	for statusRows.Next() {
		var s statusInfo
		if err := statusRows.Scan(&s.id, &s.name); err != nil {
			return nil, err
		}
		statuses = append(statuses, s)
	}

	// Daily counts per status
	query := fmt.Sprintf(`
		SELECT to_char(gs, 'YYYY-MM-DD') AS day, i.status_id, COUNT(*) AS cnt
		FROM generate_series((%s)::date, (%s)::date, '1 day'::interval) gs
		JOIN issues i ON i.project_id = $1 AND i.deleted_at IS NULL
			AND i.created_at::date <= gs::date
		GROUP BY gs, i.status_id
		ORDER BY gs`, fromDate, toDate)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("sprintRepo.GetCFD: %w", err)
	}
	defer rows.Close()

	dayMap := map[string]map[string]int{}
	for rows.Next() {
		var day, statusID string
		var cnt int
		if err := rows.Scan(&day, &statusID, &cnt); err != nil {
			return nil, err
		}
		if dayMap[day] == nil {
			dayMap[day] = map[string]int{}
		}
		dayMap[day][statusID] = cnt
	}

	statusNames := make([]string, len(statuses))
	for i, s := range statuses {
		statusNames[i] = s.name
	}

	var points []entity.CFDPoint
	for day, counts := range dayMap {
		t := entity.CFDPoint{Counts: counts}
		if parsed, err2 := time.Parse("2006-01-02", day); err2 == nil {
			t.Date = parsed
		}
		points = append(points, t)
	}

	return &entity.CFDChart{
		ProjectID: projectID,
		Statuses:  statusNames,
		Points:    points,
	}, nil
}

func (r *sprintRepo) GetVelocity(ctx context.Context, projectID string, limit int) (*entity.VelocityReport, error) {
	if limit <= 0 {
		limit = 10
	}
	rows, err := r.db.Query(ctx, `
		SELECT s.id, s.name,
			COALESCE(SUM(i.story_points), 0)                                     AS committed,
			COALESCE(SUM(i.story_points) FILTER (WHERE ws.category = 'done'), 0) AS completed
		FROM sprints s
		LEFT JOIN issues i ON i.sprint_id = s.id AND i.deleted_at IS NULL
		LEFT JOIN workflow_statuses ws ON ws.id = i.status_id
		WHERE s.project_id = $1 AND s.status = 'completed' AND s.deleted_at IS NULL
		GROUP BY s.id, s.name, s.completed_at
		ORDER BY s.completed_at DESC
		LIMIT $2`, projectID, limit,
	)
	if err != nil {
		return nil, fmt.Errorf("sprintRepo.GetVelocity: %w", err)
	}
	defer rows.Close()

	var sprints []entity.VelocityPoint
	var totalCompleted float64
	for rows.Next() {
		var vp entity.VelocityPoint
		if err := rows.Scan(&vp.SprintID, &vp.SprintName, &vp.Committed, &vp.Completed); err != nil {
			return nil, err
		}
		totalCompleted += float64(vp.Completed)
		sprints = append(sprints, vp)
	}

	avg := 0.0
	if len(sprints) > 0 {
		avg = totalCompleted / float64(len(sprints))
	}

	return &entity.VelocityReport{
		ProjectID: projectID,
		Sprints:   sprints,
		Average:   avg,
	}, rows.Err()
}
