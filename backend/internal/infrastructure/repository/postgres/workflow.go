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
	pgpkg "github.com/jira-backend/jiraflow-backend/internal/pkg/postgres"
)

type workflowRepo struct {
	db      *pgxpool.Pool
	builder sq.StatementBuilderType
}

func NewWorkflowRepo(p *pgpkg.Postgres) repository.WorkflowRepository {
	return &workflowRepo{db: p.DB, builder: p.Builder}
}

const wfCols = "id, name, description, is_default, created_by, created_at, updated_at, deleted_at"

func scanWorkflow(row pgx.Row) (*entity.Workflow, error) {
	wf := &entity.Workflow{}
	err := row.Scan(&wf.ID, &wf.Name, &wf.Description, &wf.IsDefault,
		&wf.CreatedBy, &wf.CreatedAt, &wf.UpdatedAt, &wf.DeletedAt)
	return wf, err
}

const wsCols = "id, workflow_id, name, category, color, position, is_initial, created_at, updated_at"

func scanStatus(row pgx.Row) (*entity.WorkflowStatus, error) {
	s := &entity.WorkflowStatus{}
	err := row.Scan(&s.ID, &s.WorkflowID, &s.Name, &s.Category, &s.Color, &s.Position, &s.IsInitial, &s.CreatedAt, &s.UpdatedAt)
	return s, err
}

const wtCols = "id, workflow_id, from_status_id, to_status_id, name, created_at"

func scanTransition(row pgx.Row) (*entity.WorkflowTransition, error) {
	t := &entity.WorkflowTransition{}
	err := row.Scan(&t.ID, &t.WorkflowID, &t.FromStatusID, &t.ToStatusID, &t.Name, &t.CreatedAt)
	return t, err
}

func (r *workflowRepo) Create(ctx context.Context, wf *entity.Workflow) error {
	sql, args, err := r.builder.
		Insert("workflows").
		Columns("id", "name", "description", "is_default", "created_by", "created_at", "updated_at").
		Values(wf.ID, wf.Name, wf.Description, wf.IsDefault, wf.CreatedBy, wf.CreatedAt, wf.UpdatedAt).
		ToSql()
	if err != nil {
		return fmt.Errorf("workflowRepo.Create: %w", err)
	}
	_, err = r.db.Exec(ctx, sql, args...)
	return err
}

func (r *workflowRepo) GetByID(ctx context.Context, id string) (*entity.Workflow, error) {
	sql, args, err := r.builder.
		Select(wfCols).From("workflows").
		Where(sq.And{sq.Eq{"id": id}, sq.Eq{"deleted_at": nil}}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("workflowRepo.GetByID: %w", err)
	}
	wf, err := scanWorkflow(r.db.QueryRow(ctx, sql, args...))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, apperr.NotFound("workflow")
	}
	return wf, err
}

func (r *workflowRepo) GetWithDetails(ctx context.Context, id string) (*entity.Workflow, error) {
	wf, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	statuses, err := r.ListStatuses(ctx, id)
	if err != nil {
		return nil, err
	}
	for _, s := range statuses {
		wf.Statuses = append(wf.Statuses, *s)
	}
	transitions, err := r.ListTransitions(ctx, id)
	if err != nil {
		return nil, err
	}
	for _, t := range transitions {
		wf.Transitions = append(wf.Transitions, *t)
	}
	return wf, nil
}

func (r *workflowRepo) List(ctx context.Context, filter *entity.Filter) ([]*entity.Workflow, int, error) {
	where := sq.And{sq.Eq{"deleted_at": nil}}
	if filter.Search != "" {
		where = append(where, sq.ILike{"name": "%" + filter.Search + "%"})
	}

	var total int
	cntSQL, cntArgs, _ := r.builder.Select("COUNT(*)").From("workflows").Where(where).ToSql()
	if err := r.db.QueryRow(ctx, cntSQL, cntArgs...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("workflowRepo.List count: %w", err)
	}

	dataSQL, dataArgs, _ := r.builder.
		Select(wfCols).From("workflows").Where(where).
		OrderBy("created_at DESC").
		Limit(uint64(filter.GetLimit())).Offset(uint64(filter.Offset())).ToSql()

	rows, err := r.db.Query(ctx, dataSQL, dataArgs...)
	if err != nil {
		return nil, 0, fmt.Errorf("workflowRepo.List query: %w", err)
	}
	defer rows.Close()

	var wfs []*entity.Workflow
	for rows.Next() {
		wf, err := scanWorkflow(rows)
		if err != nil {
			return nil, 0, err
		}
		wfs = append(wfs, wf)
	}
	return wfs, total, rows.Err()
}

func (r *workflowRepo) Update(ctx context.Context, wf *entity.Workflow) error {
	sql, args, err := r.builder.
		Update("workflows").
		Set("name", wf.Name).Set("description", wf.Description).
		Set("updated_at", sq.Expr("NOW()")).
		Where(sq.And{sq.Eq{"id": wf.ID}, sq.Eq{"deleted_at": nil}}).ToSql()
	if err != nil {
		return fmt.Errorf("workflowRepo.Update: %w", err)
	}
	_, err = r.db.Exec(ctx, sql, args...)
	return err
}

func (r *workflowRepo) SoftDelete(ctx context.Context, id string) error {
	sql, args, err := r.builder.
		Update("workflows").Set("deleted_at", sq.Expr("NOW()")).
		Where(sq.And{sq.Eq{"id": id}, sq.Eq{"deleted_at": nil}}).ToSql()
	if err != nil {
		return fmt.Errorf("workflowRepo.SoftDelete: %w", err)
	}
	_, err = r.db.Exec(ctx, sql, args...)
	return err
}

func (r *workflowRepo) SetDefault(ctx context.Context, id string) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("workflowRepo.SetDefault: begin: %w", err)
	}
	defer tx.Rollback(ctx)

	if _, err = tx.Exec(ctx, `UPDATE workflows SET is_default=FALSE, updated_at=NOW() WHERE is_default=TRUE AND deleted_at IS NULL`); err != nil {
		return fmt.Errorf("workflowRepo.SetDefault: clear: %w", err)
	}
	if _, err = tx.Exec(ctx, `UPDATE workflows SET is_default=TRUE, updated_at=NOW() WHERE id=$1`, id); err != nil {
		return fmt.Errorf("workflowRepo.SetDefault: set: %w", err)
	}
	return tx.Commit(ctx)
}

func (r *workflowRepo) GetDefault(ctx context.Context) (*entity.Workflow, error) {
	sql, args, err := r.builder.
		Select(wfCols).From("workflows").
		Where(sq.And{sq.Eq{"is_default": true}, sq.Eq{"deleted_at": nil}}).
		Limit(1).ToSql()
	if err != nil {
		return nil, fmt.Errorf("workflowRepo.GetDefault: %w", err)
	}
	wf, err := scanWorkflow(r.db.QueryRow(ctx, sql, args...))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, apperr.NotFound("default workflow")
	}
	return wf, err
}

// ─── Statuses ─────────────────────────────────────────────────────────────────

func (r *workflowRepo) CreateStatus(ctx context.Context, s *entity.WorkflowStatus) error {
	sql, args, err := r.builder.
		Insert("workflow_statuses").
		Columns("id", "workflow_id", "name", "category", "color", "position", "is_initial", "created_at", "updated_at").
		Values(s.ID, s.WorkflowID, s.Name, s.Category, s.Color, s.Position, s.IsInitial, s.CreatedAt, s.UpdatedAt).
		ToSql()
	if err != nil {
		return fmt.Errorf("workflowRepo.CreateStatus: %w", err)
	}
	if _, err = r.db.Exec(ctx, sql, args...); err != nil {
		if isUniqueViolation(err) {
			return apperr.Conflict("status name already exists in this workflow")
		}
		return fmt.Errorf("workflowRepo.CreateStatus: %w", err)
	}
	return nil
}

func (r *workflowRepo) GetStatusByID(ctx context.Context, id string) (*entity.WorkflowStatus, error) {
	sql, args, err := r.builder.
		Select(wsCols).From("workflow_statuses").Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("workflowRepo.GetStatusByID: %w", err)
	}
	s, err := scanStatus(r.db.QueryRow(ctx, sql, args...))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, apperr.NotFound("workflow status")
	}
	return s, err
}

func (r *workflowRepo) ListStatuses(ctx context.Context, workflowID string) ([]*entity.WorkflowStatus, error) {
	sql, args, err := r.builder.
		Select(wsCols).From("workflow_statuses").
		Where(sq.Eq{"workflow_id": workflowID}).
		OrderBy("position ASC").ToSql()
	if err != nil {
		return nil, fmt.Errorf("workflowRepo.ListStatuses: %w", err)
	}
	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("workflowRepo.ListStatuses query: %w", err)
	}
	defer rows.Close()

	var statuses []*entity.WorkflowStatus
	for rows.Next() {
		s, err := scanStatus(rows)
		if err != nil {
			return nil, err
		}
		statuses = append(statuses, s)
	}
	return statuses, rows.Err()
}

func (r *workflowRepo) UpdateStatus(ctx context.Context, s *entity.WorkflowStatus) error {
	sql, args, err := r.builder.
		Update("workflow_statuses").
		Set("name", s.Name).Set("category", s.Category).
		Set("color", s.Color).Set("position", s.Position).
		Set("updated_at", sq.Expr("NOW()")).
		Where(sq.Eq{"id": s.ID}).ToSql()
	if err != nil {
		return fmt.Errorf("workflowRepo.UpdateStatus: %w", err)
	}
	_, err = r.db.Exec(ctx, sql, args...)
	return err
}

func (r *workflowRepo) DeleteStatus(ctx context.Context, id string) error {
	sql, args, err := r.builder.Delete("workflow_statuses").Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		return fmt.Errorf("workflowRepo.DeleteStatus: %w", err)
	}
	_, err = r.db.Exec(ctx, sql, args...)
	return err
}

// ─── Transitions ──────────────────────────────────────────────────────────────

func (r *workflowRepo) CreateTransition(ctx context.Context, t *entity.WorkflowTransition) error {
	sql, args, err := r.builder.
		Insert("workflow_transitions").
		Columns("id", "workflow_id", "from_status_id", "to_status_id", "name", "created_at").
		Values(t.ID, t.WorkflowID, t.FromStatusID, t.ToStatusID, t.Name, t.CreatedAt).
		ToSql()
	if err != nil {
		return fmt.Errorf("workflowRepo.CreateTransition: %w", err)
	}
	if _, err = r.db.Exec(ctx, sql, args...); err != nil {
		if isUniqueViolation(err) {
			return apperr.Conflict("transition already exists")
		}
		return fmt.Errorf("workflowRepo.CreateTransition: %w", err)
	}
	return nil
}

func (r *workflowRepo) GetTransitionByID(ctx context.Context, id string) (*entity.WorkflowTransition, error) {
	sql, args, err := r.builder.
		Select(wtCols).From("workflow_transitions").Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("workflowRepo.GetTransitionByID: %w", err)
	}
	t, err := scanTransition(r.db.QueryRow(ctx, sql, args...))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, apperr.NotFound("workflow transition")
	}
	return t, err
}

func (r *workflowRepo) ListTransitions(ctx context.Context, workflowID string) ([]*entity.WorkflowTransition, error) {
	sql, args, err := r.builder.
		Select(wtCols).From("workflow_transitions").
		Where(sq.Eq{"workflow_id": workflowID}).
		OrderBy("created_at ASC").ToSql()
	if err != nil {
		return nil, fmt.Errorf("workflowRepo.ListTransitions: %w", err)
	}
	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("workflowRepo.ListTransitions query: %w", err)
	}
	defer rows.Close()

	var transitions []*entity.WorkflowTransition
	for rows.Next() {
		t, err := scanTransition(rows)
		if err != nil {
			return nil, err
		}
		transitions = append(transitions, t)
	}
	return transitions, rows.Err()
}

func (r *workflowRepo) DeleteTransition(ctx context.Context, id string) error {
	sql, args, err := r.builder.Delete("workflow_transitions").Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		return fmt.Errorf("workflowRepo.DeleteTransition: %w", err)
	}
	_, err = r.db.Exec(ctx, sql, args...)
	return err
}

func (r *workflowRepo) IsTransitionAllowed(ctx context.Context, workflowID, fromStatusID, toStatusID string) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx, `
		SELECT EXISTS(
			SELECT 1 FROM workflow_transitions
			WHERE workflow_id=$1 AND to_status_id=$3
			AND (from_status_id=$2 OR from_status_id IS NULL)
		)`, workflowID, fromStatusID, toStatusID,
	).Scan(&exists)
	return exists, err
}
