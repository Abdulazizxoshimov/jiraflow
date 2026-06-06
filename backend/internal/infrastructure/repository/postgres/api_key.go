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

type apiKeyRepo struct {
	db      *pgxpool.Pool
	builder sq.StatementBuilderType
}

func NewAPIKeyRepo(p *pg.Postgres) repository.APIKeyRepository {
	return &apiKeyRepo{db: p.DB, builder: p.Builder}
}

func (r *apiKeyRepo) Create(ctx context.Context, key *entity.APIKey, keyHash string) error {
	sql, args, err := r.builder.
		Insert("api_keys").
		Columns("id", "user_id", "name", "key_prefix", "key_hash", "scopes", "expires_at", "created_at").
		Values(key.ID, key.UserID, key.Name, key.KeyPrefix, keyHash, key.Scopes, key.ExpiresAt, key.CreatedAt).
		ToSql()
	if err != nil {
		return fmt.Errorf("apiKeyRepo.Create build: %w", err)
	}
	_, err = r.db.Exec(ctx, sql, args...)
	return err
}

func (r *apiKeyRepo) ListByUser(ctx context.Context, userID string) ([]*entity.APIKey, error) {
	sql, args, err := r.builder.
		Select("id", "user_id", "name", "key_prefix", "scopes", "last_used_at", "expires_at", "created_at", "revoked_at").
		From("api_keys").
		Where(sq.Eq{"user_id": userID}).
		Where("revoked_at IS NULL").
		OrderBy("created_at DESC").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("apiKeyRepo.ListByUser build: %w", err)
	}
	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("apiKeyRepo.ListByUser query: %w", err)
	}
	defer rows.Close()
	var list []*entity.APIKey
	for rows.Next() {
		k := &entity.APIKey{}
		if err := rows.Scan(
			&k.ID, &k.UserID, &k.Name, &k.KeyPrefix, &k.Scopes,
			&k.LastUsedAt, &k.ExpiresAt, &k.CreatedAt, &k.RevokedAt,
		); err != nil {
			return nil, fmt.Errorf("apiKeyRepo.ListByUser scan: %w", err)
		}
		list = append(list, k)
	}
	return list, nil
}

func (r *apiKeyRepo) GetByID(ctx context.Context, id string) (*entity.APIKey, error) {
	sql, args, err := r.builder.
		Select("id", "user_id", "name", "key_prefix", "scopes", "last_used_at", "expires_at", "created_at", "revoked_at").
		From("api_keys").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("apiKeyRepo.GetByID build: %w", err)
	}
	k := &entity.APIKey{}
	err = r.db.QueryRow(ctx, sql, args...).Scan(
		&k.ID, &k.UserID, &k.Name, &k.KeyPrefix, &k.Scopes,
		&k.LastUsedAt, &k.ExpiresAt, &k.CreatedAt, &k.RevokedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, apperr.NotFound("api key not found")
	}
	if err != nil {
		return nil, fmt.Errorf("apiKeyRepo.GetByID scan: %w", err)
	}
	return k, nil
}

func (r *apiKeyRepo) GetByHash(ctx context.Context, keyHash string) (*entity.APIKey, error) {
	sql, args, err := r.builder.
		Select("id", "user_id", "name", "key_prefix", "scopes", "last_used_at", "expires_at", "created_at", "revoked_at").
		From("api_keys").
		Where(sq.Eq{"key_hash": keyHash}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("apiKeyRepo.GetByHash build: %w", err)
	}
	k := &entity.APIKey{}
	err = r.db.QueryRow(ctx, sql, args...).Scan(
		&k.ID, &k.UserID, &k.Name, &k.KeyPrefix, &k.Scopes,
		&k.LastUsedAt, &k.ExpiresAt, &k.CreatedAt, &k.RevokedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, apperr.NotFound("api key not found")
	}
	if err != nil {
		return nil, fmt.Errorf("apiKeyRepo.GetByHash scan: %w", err)
	}
	return k, nil
}

func (r *apiKeyRepo) Revoke(ctx context.Context, id, userID string) error {
	now := time.Now().UTC()
	sql, args, err := r.builder.
		Update("api_keys").
		Set("revoked_at", now).
		Where(sq.Eq{"id": id, "user_id": userID}).
		ToSql()
	if err != nil {
		return fmt.Errorf("apiKeyRepo.Revoke build: %w", err)
	}
	tag, err := r.db.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("apiKeyRepo.Revoke exec: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return apperr.NotFound("api key not found")
	}
	return nil
}

func (r *apiKeyRepo) UpdateLastUsed(ctx context.Context, id string) error {
	now := time.Now().UTC()
	sql, args, err := r.builder.
		Update("api_keys").
		Set("last_used_at", now).
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return fmt.Errorf("apiKeyRepo.UpdateLastUsed build: %w", err)
	}
	_, err = r.db.Exec(ctx, sql, args...)
	return err
}
