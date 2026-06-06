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

type dataImportRepo struct {
	db      *pgxpool.Pool
	builder sq.StatementBuilderType
}

func NewDataImportRepo(p *pg.Postgres) repository.DataImportRepository {
	return &dataImportRepo{db: p.DB, builder: p.Builder}
}

func (r *dataImportRepo) Create(ctx context.Context, imp *entity.DataImport) error {
	sql, args, err := r.builder.
		Insert("data_imports").
		Columns("id", "user_id", "source", "status", "total_items", "processed_items", "created_at").
		Values(imp.ID, imp.UserID, imp.Source, imp.Status, imp.TotalItems, imp.ProcessedItems, imp.CreatedAt).
		ToSql()
	if err != nil {
		return fmt.Errorf("dataImportRepo.Create build: %w", err)
	}
	_, err = r.db.Exec(ctx, sql, args...)
	return err
}

func (r *dataImportRepo) GetByID(ctx context.Context, id string) (*entity.DataImport, error) {
	sql, args, err := r.builder.
		Select("id", "user_id", "source", "status", "total_items", "processed_items", "error_message", "created_at", "completed_at").
		From("data_imports").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("dataImportRepo.GetByID build: %w", err)
	}
	var imp entity.DataImport
	err = r.db.QueryRow(ctx, sql, args...).Scan(
		&imp.ID, &imp.UserID, &imp.Source, &imp.Status,
		&imp.TotalItems, &imp.ProcessedItems, &imp.ErrorMessage,
		&imp.CreatedAt, &imp.CompletedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, apperr.NotFound("import job not found")
	}
	if err != nil {
		return nil, fmt.Errorf("dataImportRepo.GetByID scan: %w", err)
	}
	return &imp, nil
}

func (r *dataImportRepo) UpdateStatus(ctx context.Context, id, status string, total, processed int, errMsg string) error {
	sql, args, err := r.builder.
		Update("data_imports").
		Set("status", status).
		Set("total_items", total).
		Set("processed_items", processed).
		Set("error_message", errMsg).
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return fmt.Errorf("dataImportRepo.UpdateStatus build: %w", err)
	}
	_, err = r.db.Exec(ctx, sql, args...)
	return err
}

func (r *dataImportRepo) MarkCompleted(ctx context.Context, id string) error {
	now := time.Now().UTC()
	sql, args, err := r.builder.
		Update("data_imports").
		Set("status", "done").
		Set("completed_at", now).
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return fmt.Errorf("dataImportRepo.MarkCompleted build: %w", err)
	}
	_, err = r.db.Exec(ctx, sql, args...)
	return err
}
