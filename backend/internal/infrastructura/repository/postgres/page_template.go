package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
	"github.com/jira-backend/jiraflow-backend/internal/infrastructura/repository"
	apperr "github.com/jira-backend/jiraflow-backend/internal/pkg/errors"
	pg "github.com/jira-backend/jiraflow-backend/internal/pkg/postgres"
)

type pageTemplateRepo struct {
	db      *pgxpool.Pool
	builder sq.StatementBuilderType
}

func NewPageTemplateRepo(p *pg.Postgres) repository.PageTemplateRepository {
	return &pageTemplateRepo{db: p.DB, builder: p.Builder}
}

func (r *pageTemplateRepo) Create(ctx context.Context, t *entity.PageTemplate) error {
	if t.ID == "" {
		t.ID = uuid.NewString()
	}
	now := time.Now().UTC()
	t.CreatedAt = now
	t.UpdatedAt = now

	contentJSON, err := json.Marshal(t.Content)
	if err != nil {
		return fmt.Errorf("pageTemplateRepo.Create marshal: %w", err)
	}

	_, err = r.db.Exec(ctx,
		`INSERT INTO page_templates(id, space_id, name, description, category, content, content_text, icon, created_by, is_global, created_at, updated_at)
		 VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)`,
		t.ID, t.SpaceID, t.Name, t.Description, t.Category,
		contentJSON, t.ContentText, t.Icon, t.CreatedBy, t.IsGlobal,
		t.CreatedAt, t.UpdatedAt,
	)
	return err
}

func (r *pageTemplateRepo) GetByID(ctx context.Context, id string) (*entity.PageTemplate, error) {
	row := r.db.QueryRow(ctx, `
		SELECT id, space_id, name, description, category, content, content_text, icon, created_by, is_global, created_at, updated_at
		FROM page_templates WHERE id=$1
	`, id)

	t := &entity.PageTemplate{}
	var contentJSON []byte
	err := row.Scan(
		&t.ID, &t.SpaceID, &t.Name, &t.Description, &t.Category,
		&contentJSON, &t.ContentText, &t.Icon, &t.CreatedBy, &t.IsGlobal,
		&t.CreatedAt, &t.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, apperr.NotFound("page template")
	}
	if err != nil {
		return nil, fmt.Errorf("pageTemplateRepo.GetByID: %w", err)
	}
	if len(contentJSON) > 0 {
		_ = json.Unmarshal(contentJSON, &t.Content)
	}
	return t, nil
}

func (r *pageTemplateRepo) List(ctx context.Context, filter *entity.PageTemplateFilter) ([]*entity.PageTemplate, int, error) {
	where := sq.Or{}
	if filter.SpaceID != "" {
		where = append(where, sq.Eq{"space_id": filter.SpaceID})
	}
	where = append(where, sq.Eq{"is_global": true})

	if filter.Category != "" {
		// wrap current where in AND with category filter
		combinedWhere := sq.And{where, sq.Eq{"category": filter.Category}}
		return r.listWithWhere(ctx, combinedWhere, filter)
	}
	return r.listWithWhere(ctx, where, filter)
}

func (r *pageTemplateRepo) listWithWhere(ctx context.Context, where sq.Sqlizer, filter *entity.PageTemplateFilter) ([]*entity.PageTemplate, int, error) {
	var total int
	cntSQL, cntArgs, _ := r.builder.Select("COUNT(*)").From("page_templates").Where(where).ToSql()
	if err := r.db.QueryRow(ctx, cntSQL, cntArgs...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("pageTemplateRepo.List count: %w", err)
	}

	dataSQL, dataArgs, _ := r.builder.
		Select("id", "space_id", "name", "description", "category", "content", "content_text", "icon", "created_by", "is_global", "created_at", "updated_at").
		From("page_templates").Where(where).
		OrderBy("is_global DESC, name ASC").
		Limit(uint64(filter.GetLimit())).Offset(uint64(filter.Offset())).ToSql()

	rows, err := r.db.Query(ctx, dataSQL, dataArgs...)
	if err != nil {
		return nil, 0, fmt.Errorf("pageTemplateRepo.List query: %w", err)
	}
	defer rows.Close()

	var templates []*entity.PageTemplate
	for rows.Next() {
		t := &entity.PageTemplate{}
		var contentJSON []byte
		if err := rows.Scan(
			&t.ID, &t.SpaceID, &t.Name, &t.Description, &t.Category,
			&contentJSON, &t.ContentText, &t.Icon, &t.CreatedBy, &t.IsGlobal,
			&t.CreatedAt, &t.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		if len(contentJSON) > 0 {
			_ = json.Unmarshal(contentJSON, &t.Content)
		}
		templates = append(templates, t)
	}
	return templates, total, rows.Err()
}

func (r *pageTemplateRepo) Update(ctx context.Context, t *entity.PageTemplate) error {
	contentJSON, err := json.Marshal(t.Content)
	if err != nil {
		return fmt.Errorf("pageTemplateRepo.Update marshal: %w", err)
	}
	_, err = r.db.Exec(ctx,
		`UPDATE page_templates SET name=$1, description=$2, category=$3, content=$4, content_text=$5, icon=$6, updated_at=NOW()
		 WHERE id=$7`,
		t.Name, t.Description, t.Category, contentJSON, t.ContentText, t.Icon, t.ID,
	)
	return err
}

func (r *pageTemplateRepo) Delete(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM page_templates WHERE id=$1`, id)
	return err
}
