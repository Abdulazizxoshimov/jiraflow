package postgres

import (
	"context"
	"encoding/json"
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

type pageVersionRepo struct {
	db      *pgxpool.Pool
	builder sq.StatementBuilderType
}

func NewPageVersionRepo(p *pg.Postgres) repository.PageVersionRepository {
	return &pageVersionRepo{db: p.DB, builder: p.Builder}
}

const pvCols = "id, page_id, version, title, content, content_text, author_id, change_note, created_at"

func scanPageVersion(row pgx.Row) (*entity.PageVersion, error) {
	v := &entity.PageVersion{}
	var contentJSON []byte
	err := row.Scan(
		&v.ID, &v.PageID, &v.Version, &v.Title,
		&contentJSON, &v.ContentText, &v.AuthorID, &v.ChangeNote, &v.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	if len(contentJSON) > 0 {
		_ = json.Unmarshal(contentJSON, &v.Content)
	}
	return v, nil
}

func (r *pageVersionRepo) Create(ctx context.Context, v *entity.PageVersion) error {
	contentJSON, err := json.Marshal(v.Content)
	if err != nil {
		return fmt.Errorf("pageVersionRepo.Create marshal content: %w", err)
	}
	sql, args, err := r.builder.
		Insert("page_versions").
		Columns("id", "page_id", "version", "title", "content", "content_text", "author_id", "change_note", "created_at").
		Values(v.ID, v.PageID, v.Version, v.Title, contentJSON, v.ContentText, v.AuthorID, v.ChangeNote, v.CreatedAt).
		ToSql()
	if err != nil {
		return fmt.Errorf("pageVersionRepo.Create: %w", err)
	}
	_, err = r.db.Exec(ctx, sql, args...)
	return err
}

func (r *pageVersionRepo) GetByID(ctx context.Context, id string) (*entity.PageVersion, error) {
	sql, args, err := r.builder.
		Select(pvCols).From("page_versions").Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("pageVersionRepo.GetByID: %w", err)
	}
	v, err := scanPageVersion(r.db.QueryRow(ctx, sql, args...))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, apperr.NotFound("page version")
	}
	return v, err
}

func (r *pageVersionRepo) GetByVersion(ctx context.Context, pageID string, version int) (*entity.PageVersion, error) {
	sql, args, err := r.builder.
		Select(pvCols).From("page_versions").
		Where(sq.And{sq.Eq{"page_id": pageID}, sq.Eq{"version": version}}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("pageVersionRepo.GetByVersion: %w", err)
	}
	v, err := scanPageVersion(r.db.QueryRow(ctx, sql, args...))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, apperr.NotFound("page version")
	}
	return v, err
}

func (r *pageVersionRepo) ListByPage(ctx context.Context, pageID string, filter *entity.Filter) ([]*entity.PageVersion, int, error) {
	where := sq.Eq{"page_id": pageID}

	var total int
	cntSQL, cntArgs, _ := r.builder.Select("COUNT(*)").From("page_versions").Where(where).ToSql()
	if err := r.db.QueryRow(ctx, cntSQL, cntArgs...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("pageVersionRepo.ListByPage count: %w", err)
	}

	dataSQL, dataArgs, _ := r.builder.
		Select(pvCols).From("page_versions").Where(where).
		OrderBy("version DESC").
		Limit(uint64(filter.GetLimit())).Offset(uint64(filter.Offset())).ToSql()

	rows, err := r.db.Query(ctx, dataSQL, dataArgs...)
	if err != nil {
		return nil, 0, fmt.Errorf("pageVersionRepo.ListByPage query: %w", err)
	}
	defer rows.Close()

	var versions []*entity.PageVersion
	for rows.Next() {
		v, err := scanPageVersion(rows)
		if err != nil {
			return nil, 0, err
		}
		versions = append(versions, v)
	}
	return versions, total, rows.Err()
}

func (r *pageVersionRepo) GetLatestVersion(ctx context.Context, pageID string) (int, error) {
	var version int
	err := r.db.QueryRow(ctx,
		`SELECT COALESCE(MAX(version), 0) FROM page_versions WHERE page_id=$1`, pageID,
	).Scan(&version)
	return version, err
}
