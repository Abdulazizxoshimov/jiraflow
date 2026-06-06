package postgres

import (
	"context"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
	"github.com/jira-backend/jiraflow-backend/internal/infrastructure/repository"
	apperr "github.com/jira-backend/jiraflow-backend/internal/pkg/errors"
	pg "github.com/jira-backend/jiraflow-backend/internal/pkg/postgres"
)

type notificationSchemeRepo struct {
	db      *pgxpool.Pool
	builder sq.StatementBuilderType
}

func NewNotificationSchemeRepo(p *pg.Postgres) repository.NotificationSchemeRepository {
	return &notificationSchemeRepo{db: p.DB, builder: p.Builder}
}

func (r *notificationSchemeRepo) Create(ctx context.Context, s *entity.NotificationScheme) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("notificationSchemeRepo.Create begin: %w", err)
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx,
		`INSERT INTO notification_schemes (id, name, description) VALUES ($1, $2, $3)`,
		s.ID, s.Name, s.Description)
	if err != nil {
		return fmt.Errorf("notificationSchemeRepo.Create scheme: %w", err)
	}

	for _, rule := range s.Rules {
		if rule.ID == "" {
			rule.ID = uuid.NewString()
		}
		rule.SchemeID = s.ID
		_, err = tx.Exec(ctx,
			`INSERT INTO notification_scheme_rules (id, scheme_id, event_type, recipient_type, recipient_id) VALUES ($1,$2,$3,$4,$5)`,
			rule.ID, rule.SchemeID, rule.EventType, rule.RecipientType, rule.RecipientID)
		if err != nil {
			return fmt.Errorf("notificationSchemeRepo.Create rule: %w", err)
		}
	}
	return tx.Commit(ctx)
}

func (r *notificationSchemeRepo) GetByID(ctx context.Context, id string) (*entity.NotificationScheme, error) {
	row := r.db.QueryRow(ctx,
		`SELECT id, name, description, created_at FROM notification_schemes WHERE id=$1`, id)
	s := &entity.NotificationScheme{}
	if err := row.Scan(&s.ID, &s.Name, &s.Description, &s.CreatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperr.NotFound("notification scheme not found")
		}
		return nil, fmt.Errorf("notificationSchemeRepo.GetByID: %w", err)
	}
	rows, err := r.db.Query(ctx,
		`SELECT id, scheme_id, event_type, recipient_type, recipient_id, created_at FROM notification_scheme_rules WHERE scheme_id=$1`, id)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			rule := &entity.NotificationSchemeRule{}
			if rows.Scan(&rule.ID, &rule.SchemeID, &rule.EventType, &rule.RecipientType, &rule.RecipientID, &rule.CreatedAt) == nil {
				s.Rules = append(s.Rules, rule)
			}
		}
	}
	return s, nil
}

func (r *notificationSchemeRepo) List(ctx context.Context) ([]*entity.NotificationScheme, error) {
	rows, err := r.db.Query(ctx, `SELECT id, name, description, created_at FROM notification_schemes ORDER BY created_at DESC`)
	if err != nil {
		return nil, fmt.Errorf("notificationSchemeRepo.List: %w", err)
	}
	defer rows.Close()
	var list []*entity.NotificationScheme
	for rows.Next() {
		s := &entity.NotificationScheme{}
		if err := rows.Scan(&s.ID, &s.Name, &s.Description, &s.CreatedAt); err != nil {
			return nil, fmt.Errorf("notificationSchemeRepo.List scan: %w", err)
		}
		list = append(list, s)
	}
	return list, rows.Err()
}

func (r *notificationSchemeRepo) Delete(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM notification_schemes WHERE id=$1`, id)
	return err
}
