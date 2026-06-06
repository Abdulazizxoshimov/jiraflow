package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
	"github.com/jira-backend/jiraflow-backend/internal/infrastructura/repository"
	apperr "github.com/jira-backend/jiraflow-backend/internal/pkg/errors"
	pg "github.com/jira-backend/jiraflow-backend/internal/pkg/postgres"
)

type contentPropertyRepo struct {
	db      *pgxpool.Pool
	builder sq.StatementBuilderType
}

func NewContentPropertyRepo(p *pg.Postgres) repository.ContentPropertyRepository {
	return &contentPropertyRepo{db: p.DB, builder: p.Builder}
}

func scanContentProperty(row pgx.Row) (*entity.ContentProperty, error) {
	p := &entity.ContentProperty{}
	var valueJSON []byte
	if err := row.Scan(&p.ID, &p.EntityType, &p.EntityID, &p.Key, &valueJSON, &p.CreatedAt, &p.UpdatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperr.NotFound("content property not found")
		}
		return nil, fmt.Errorf("contentPropertyRepo.scan: %w", err)
	}
	if len(valueJSON) > 0 {
		_ = json.Unmarshal(valueJSON, &p.Value)
	}
	return p, nil
}

func (r *contentPropertyRepo) Set(ctx context.Context, p *entity.ContentProperty) error {
	valueJSON, err := json.Marshal(p.Value)
	if err != nil {
		return fmt.Errorf("contentPropertyRepo.Set marshal: %w", err)
	}
	if p.ID == "" {
		p.ID = uuid.NewString()
	}
	_, err = r.db.Exec(ctx, `
		INSERT INTO content_properties (id, entity_type, entity_id, key, value)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (entity_type, entity_id, key)
		DO UPDATE SET value = EXCLUDED.value, updated_at = NOW()
	`, p.ID, p.EntityType, p.EntityID, p.Key, valueJSON)
	return err
}

func (r *contentPropertyRepo) Get(ctx context.Context, entityType, entityID, key string) (*entity.ContentProperty, error) {
	row := r.db.QueryRow(ctx,
		`SELECT id, entity_type, entity_id, key, value, created_at, updated_at FROM content_properties WHERE entity_type=$1 AND entity_id=$2 AND key=$3`,
		entityType, entityID, key)
	return scanContentProperty(row)
}

func (r *contentPropertyRepo) List(ctx context.Context, entityType, entityID string) ([]*entity.ContentProperty, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, entity_type, entity_id, key, value, created_at, updated_at FROM content_properties WHERE entity_type=$1 AND entity_id=$2 ORDER BY key ASC`,
		entityType, entityID)
	if err != nil {
		return nil, fmt.Errorf("contentPropertyRepo.List: %w", err)
	}
	defer rows.Close()
	var list []*entity.ContentProperty
	for rows.Next() {
		p, err := scanContentProperty(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, p)
	}
	return list, rows.Err()
}

func (r *contentPropertyRepo) Delete(ctx context.Context, entityType, entityID, key string) error {
	_, err := r.db.Exec(ctx,
		`DELETE FROM content_properties WHERE entity_type=$1 AND entity_id=$2 AND key=$3`,
		entityType, entityID, key)
	return err
}
