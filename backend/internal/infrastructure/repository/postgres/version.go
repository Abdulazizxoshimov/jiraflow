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

type versionRepo struct {
	db      *pgxpool.Pool
	builder sq.StatementBuilderType
}

func NewVersionRepo(p *pg.Postgres) repository.VersionRepository {
	return &versionRepo{db: p.DB, builder: p.Builder}
}

const versionCols = "id, project_id, name, description, status, start_date, release_date, released_at, created_at, updated_at"

func scanVersion(row pgx.Row) (*entity.Version, error) {
	v := &entity.Version{}
	err := row.Scan(
		&v.ID, &v.ProjectID, &v.Name, &v.Description, &v.Status,
		&v.StartDate, &v.ReleaseDate, &v.ReleasedAt,
		&v.CreatedAt, &v.UpdatedAt,
	)
	return v, err
}

func (r *versionRepo) Create(ctx context.Context, v *entity.Version) error {
	sql, args, err := r.builder.
		Insert("project_versions").
		Columns("id", "project_id", "name", "description", "status", "start_date", "release_date", "created_at", "updated_at").
		Values(v.ID, v.ProjectID, v.Name, v.Description, v.Status, v.StartDate, v.ReleaseDate, v.CreatedAt, v.UpdatedAt).
		ToSql()
	if err != nil {
		return fmt.Errorf("versionRepo.Create: %w", err)
	}
	_, err = r.db.Exec(ctx, sql, args...)
	return err
}

func (r *versionRepo) GetByID(ctx context.Context, id string) (*entity.Version, error) {
	sql, args, err := r.builder.
		Select(versionCols).From("project_versions").
		Where(sq.And{sq.Eq{"id": id}, sq.Eq{"deleted_at": nil}}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("versionRepo.GetByID: %w", err)
	}
	v, err := scanVersion(r.db.QueryRow(ctx, sql, args...))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, apperr.NotFound("version")
	}
	return v, err
}

func (r *versionRepo) List(ctx context.Context, projectID string) ([]*entity.Version, error) {
	sql, args, err := r.builder.
		Select(versionCols).From("project_versions").
		Where(sq.And{sq.Eq{"project_id": projectID}, sq.Eq{"deleted_at": nil}}).
		OrderBy("created_at DESC").ToSql()
	if err != nil {
		return nil, fmt.Errorf("versionRepo.List: %w", err)
	}
	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("versionRepo.List query: %w", err)
	}
	defer rows.Close()

	var result []*entity.Version
	for rows.Next() {
		v, err := scanVersion(rows)
		if err != nil {
			return nil, err
		}
		result = append(result, v)
	}
	return result, rows.Err()
}

func (r *versionRepo) Update(ctx context.Context, v *entity.Version) error {
	sql, args, err := r.builder.
		Update("project_versions").
		Set("name", v.Name).
		Set("description", v.Description).
		Set("start_date", v.StartDate).
		Set("release_date", v.ReleaseDate).
		Set("updated_at", sq.Expr("NOW()")).
		Where(sq.And{sq.Eq{"id": v.ID}, sq.Eq{"deleted_at": nil}}).ToSql()
	if err != nil {
		return fmt.Errorf("versionRepo.Update: %w", err)
	}
	_, err = r.db.Exec(ctx, sql, args...)
	return err
}

func (r *versionRepo) Release(ctx context.Context, id string, releasedAt time.Time) error {
	_, err := r.db.Exec(ctx,
		`UPDATE project_versions SET status='released', released_at=$2, updated_at=NOW() WHERE id=$1 AND deleted_at IS NULL`,
		id, releasedAt,
	)
	return err
}

func (r *versionRepo) Archive(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx,
		`UPDATE project_versions SET status='archived', updated_at=NOW() WHERE id=$1 AND deleted_at IS NULL`, id)
	return err
}

func (r *versionRepo) Delete(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx,
		`UPDATE project_versions SET deleted_at=NOW() WHERE id=$1 AND deleted_at IS NULL`, id)
	return err
}

func (r *versionRepo) setIssueVersionsByType(ctx context.Context, issueID string, versionIDs []string, vType string) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, `DELETE FROM issue_versions WHERE issue_id=$1 AND version_type=$2`, issueID, vType); err != nil {
		return err
	}
	if len(versionIDs) > 0 {
		b := &pgx.Batch{}
		for _, vid := range versionIDs {
			b.Queue(`INSERT INTO issue_versions(issue_id,version_id,version_type) VALUES($1,$2,$3) ON CONFLICT DO NOTHING`, issueID, vid, vType)
		}
		br := tx.SendBatch(ctx, b)
		for range versionIDs {
			if _, err := br.Exec(); err != nil {
				br.Close()
				return err
			}
		}
		br.Close()
	}
	return tx.Commit(ctx)
}

func (r *versionRepo) SetIssueVersions(ctx context.Context, issueID string, versionIDs []string) error {
	return r.setIssueVersionsByType(ctx, issueID, versionIDs, "fix")
}

func (r *versionRepo) SetIssueAffectsVersions(ctx context.Context, issueID string, versionIDs []string) error {
	return r.setIssueVersionsByType(ctx, issueID, versionIDs, "affects")
}

func (r *versionRepo) getIssueVersionsByType(ctx context.Context, issueID, vType string) ([]*entity.Version, error) {
	query := `
		SELECT v.id, v.project_id, v.name, v.description, v.status,
		       v.start_date, v.release_date, v.released_at, v.created_at, v.updated_at
		FROM project_versions v
		JOIN issue_versions iv ON iv.version_id = v.id
		WHERE iv.issue_id = $1 AND iv.version_type = $2 AND v.deleted_at IS NULL`
	rows, err := r.db.Query(ctx, query, issueID, vType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*entity.Version
	for rows.Next() {
		v, err := scanVersion(rows)
		if err != nil {
			return nil, err
		}
		result = append(result, v)
	}
	return result, rows.Err()
}

func (r *versionRepo) GetIssueVersions(ctx context.Context, issueID string) ([]*entity.Version, error) {
	return r.getIssueVersionsByType(ctx, issueID, "fix")
}

func (r *versionRepo) GetIssueAffectsVersions(ctx context.Context, issueID string) ([]*entity.Version, error) {
	return r.getIssueVersionsByType(ctx, issueID, "affects")
}

func (r *versionRepo) GetProgress(ctx context.Context, versionID string) (int, int, error) {
	var total, done int
	err := r.db.QueryRow(ctx, `
		SELECT
			COUNT(*)                                                        AS total,
			COUNT(*) FILTER (WHERE ws.category = 'done')                   AS done
		FROM issue_versions iv
		JOIN issues i ON i.id = iv.issue_id AND i.deleted_at IS NULL
		JOIN workflow_statuses ws ON ws.id = i.status_id
		WHERE iv.version_id = $1`, versionID,
	).Scan(&total, &done)
	return total, done, err
}
