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

type oauthRepo struct {
	db      *pgxpool.Pool
	builder sq.StatementBuilderType
}

func NewOAuthRepo(p *pg.Postgres) repository.OAuthRepository {
	return &oauthRepo{db: p.DB, builder: p.Builder}
}

func (r *oauthRepo) SaveState(ctx context.Context, s *entity.OAuthState) error {
	sql, args, err := r.builder.
		Insert("oauth_states").
		Columns("state", "redirect_url", "created_at", "expires_at").
		Values(s.State, s.RedirectURL, s.CreatedAt, s.ExpiresAt).
		ToSql()
	if err != nil {
		return fmt.Errorf("oauthRepo.SaveState build: %w", err)
	}
	_, err = r.db.Exec(ctx, sql, args...)
	return err
}

func (r *oauthRepo) GetState(ctx context.Context, state string) (*entity.OAuthState, error) {
	sql, args, err := r.builder.
		Select("state", "redirect_url", "created_at", "expires_at").
		From("oauth_states").
		Where(sq.Eq{"state": state}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("oauthRepo.GetState build: %w", err)
	}
	var s entity.OAuthState
	err = r.db.QueryRow(ctx, sql, args...).Scan(&s.State, &s.RedirectURL, &s.CreatedAt, &s.ExpiresAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, apperr.NotFound("oauth state not found")
	}
	if err != nil {
		return nil, fmt.Errorf("oauthRepo.GetState scan: %w", err)
	}
	if time.Now().UTC().After(s.ExpiresAt) {
		_ = r.DeleteState(ctx, state)
		return nil, apperr.NotFound("oauth state expired")
	}
	return &s, nil
}

func (r *oauthRepo) DeleteState(ctx context.Context, state string) error {
	sql, args, err := r.builder.Delete("oauth_states").Where(sq.Eq{"state": state}).ToSql()
	if err != nil {
		return fmt.Errorf("oauthRepo.DeleteState build: %w", err)
	}
	_, err = r.db.Exec(ctx, sql, args...)
	return err
}

func (r *oauthRepo) UpsertAccount(ctx context.Context, acc *entity.OAuthAccount) error {
	sql, args, err := r.builder.
		Insert("oauth_accounts").
		Columns("id", "user_id", "provider", "provider_user_id", "email", "name", "avatar_url", "created_at", "updated_at").
		Values(acc.ID, acc.UserID, acc.Provider, acc.ProviderUserID, acc.Email, acc.Name, acc.AvatarURL, acc.CreatedAt, acc.UpdatedAt).
		Suffix(`ON CONFLICT (provider, provider_user_id) DO UPDATE
			SET email = EXCLUDED.email, name = EXCLUDED.name, avatar_url = EXCLUDED.avatar_url,
			    updated_at = EXCLUDED.updated_at`).
		ToSql()
	if err != nil {
		return fmt.Errorf("oauthRepo.UpsertAccount build: %w", err)
	}
	_, err = r.db.Exec(ctx, sql, args...)
	return err
}

func (r *oauthRepo) GetAccountByProvider(ctx context.Context, provider, providerUserID string) (*entity.OAuthAccount, error) {
	sql, args, err := r.builder.
		Select("id", "user_id", "provider", "provider_user_id", "email", "name", "avatar_url", "created_at", "updated_at").
		From("oauth_accounts").
		Where(sq.Eq{"provider": provider, "provider_user_id": providerUserID}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("oauthRepo.GetAccountByProvider build: %w", err)
	}
	var acc entity.OAuthAccount
	err = r.db.QueryRow(ctx, sql, args...).Scan(
		&acc.ID, &acc.UserID, &acc.Provider, &acc.ProviderUserID,
		&acc.Email, &acc.Name, &acc.AvatarURL, &acc.CreatedAt, &acc.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, apperr.NotFound("oauth account not found")
	}
	if err != nil {
		return nil, fmt.Errorf("oauthRepo.GetAccountByProvider scan: %w", err)
	}
	return &acc, nil
}

func (r *oauthRepo) ListByUser(ctx context.Context, userID string) ([]*entity.OAuthAccount, error) {
	sql, args, err := r.builder.
		Select("id", "user_id", "provider", "provider_user_id", "email", "name", "avatar_url", "created_at", "updated_at").
		From("oauth_accounts").
		Where(sq.Eq{"user_id": userID}).
		OrderBy("created_at ASC").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("oauthRepo.ListByUser build: %w", err)
	}
	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("oauthRepo.ListByUser query: %w", err)
	}
	defer rows.Close()
	var list []*entity.OAuthAccount
	for rows.Next() {
		var acc entity.OAuthAccount
		if err := rows.Scan(
			&acc.ID, &acc.UserID, &acc.Provider, &acc.ProviderUserID,
			&acc.Email, &acc.Name, &acc.AvatarURL, &acc.CreatedAt, &acc.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("oauthRepo.ListByUser scan: %w", err)
		}
		list = append(list, &acc)
	}
	return list, nil
}

func (r *oauthRepo) DeleteAccount(ctx context.Context, userID, provider string) error {
	sql, args, err := r.builder.
		Delete("oauth_accounts").
		Where(sq.Eq{"user_id": userID, "provider": provider}).
		ToSql()
	if err != nil {
		return fmt.Errorf("oauthRepo.DeleteAccount build: %w", err)
	}
	_, err = r.db.Exec(ctx, sql, args...)
	return err
}
