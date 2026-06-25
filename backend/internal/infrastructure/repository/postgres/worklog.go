package postgres

import (
	"context"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
	"github.com/jira-backend/jiraflow-backend/internal/infrastructure/repository"
	apperr "github.com/jira-backend/jiraflow-backend/internal/pkg/errors"
	pg "github.com/jira-backend/jiraflow-backend/internal/pkg/postgres"
)

type worklogRepo struct {
	db      *pgxpool.Pool
	builder sq.StatementBuilderType
}

func NewWorklogRepo(p *pg.Postgres) repository.WorklogRepository {
	return &worklogRepo{db: p.DB, builder: p.Builder}
}

const worklogCols = `
	w.id, w.issue_id, w.user_id, w.time_spent, w.started_at, w.description, w.created_at, w.updated_at,
	u.id, u.full_name, u.email, u.avatar_url, u.color`

func scanWorklog(row pgx.Row) (*entity.Worklog, error) {
	w := &entity.Worklog{User: &entity.UserShort{}}
	err := row.Scan(
		&w.ID, &w.IssueID, &w.UserID, &w.TimeSpent, &w.StartedAt, &w.Description, &w.CreatedAt, &w.UpdatedAt,
		&w.User.ID, &w.User.FullName, &w.User.Email, &w.User.AvatarURL, &w.User.Color,
	)
	if err != nil {
		return nil, err
	}
	return w, nil
}

func (r *worklogRepo) Create(ctx context.Context, w *entity.Worklog) error {
	sql, args, err := r.builder.
		Insert("issue_worklogs").
		Columns("id", "issue_id", "user_id", "time_spent", "started_at", "description", "created_at", "updated_at").
		Values(w.ID, w.IssueID, w.UserID, w.TimeSpent, w.StartedAt, w.Description, w.CreatedAt, w.UpdatedAt).
		ToSql()
	if err != nil {
		return fmt.Errorf("worklogRepo.Create: %w", err)
	}
	_, err = r.db.Exec(ctx, sql, args...)
	return err
}

func (r *worklogRepo) GetByID(ctx context.Context, id string) (*entity.Worklog, error) {
	query := `
		SELECT ` + worklogCols + `
		FROM issue_worklogs w
		JOIN users u ON u.id = w.user_id
		WHERE w.id = $1`
	w, err := scanWorklog(r.db.QueryRow(ctx, query, id))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, apperr.NotFound("worklog")
	}
	return w, err
}

func (r *worklogRepo) List(ctx context.Context, filter *entity.WorklogFilter) ([]*entity.Worklog, int, error) {
	cond := sq.And{}
	if filter.IssueID != "" {
		cond = append(cond, sq.Eq{"w.issue_id": filter.IssueID})
	}
	if filter.UserID != "" {
		cond = append(cond, sq.Eq{"w.user_id": filter.UserID})
	}

	countSQL, countArgs, err := r.builder.
		Select("COUNT(*)").From("issue_worklogs w").Where(cond).ToSql()
	if err != nil {
		return nil, 0, fmt.Errorf("worklogRepo.List count build: %w", err)
	}
	var total int
	if err := r.db.QueryRow(ctx, countSQL, countArgs...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("worklogRepo.List count: %w", err)
	}

	dataSQL, dataArgs, err := r.builder.
		Select(worklogCols).
		From("issue_worklogs w").
		Join("users u ON u.id = w.user_id").
		Where(cond).
		OrderBy("w.started_at DESC").
		Limit(uint64(filter.GetLimit())).
		Offset(uint64(filter.Offset())).
		ToSql()
	if err != nil {
		return nil, 0, fmt.Errorf("worklogRepo.List query build: %w", err)
	}

	rows, err := r.db.Query(ctx, dataSQL, dataArgs...)
	if err != nil {
		return nil, 0, fmt.Errorf("worklogRepo.List query: %w", err)
	}
	defer rows.Close()

	var result []*entity.Worklog
	for rows.Next() {
		w, err := scanWorklog(rows)
		if err != nil {
			return nil, 0, err
		}
		result = append(result, w)
	}
	return result, total, rows.Err()
}

func (r *worklogRepo) Update(ctx context.Context, w *entity.Worklog) error {
	sql, args, err := r.builder.
		Update("issue_worklogs").
		Set("time_spent", w.TimeSpent).
		Set("started_at", w.StartedAt).
		Set("description", w.Description).
		Set("updated_at", sq.Expr("NOW()")).
		Where(sq.Eq{"id": w.ID}).
		ToSql()
	if err != nil {
		return fmt.Errorf("worklogRepo.Update: %w", err)
	}
	_, err = r.db.Exec(ctx, sql, args...)
	return err
}

func (r *worklogRepo) Delete(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM issue_worklogs WHERE id = $1`, id)
	return err
}

func (r *worklogRepo) SumByIssue(ctx context.Context, issueID string) (int, error) {
	var total int
	err := r.db.QueryRow(ctx,
		`SELECT COALESCE(SUM(time_spent), 0) FROM issue_worklogs WHERE issue_id = $1`,
		issueID,
	).Scan(&total)
	return total, err
}
