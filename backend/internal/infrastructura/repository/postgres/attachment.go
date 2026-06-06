package postgres

import (
	"context"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
	apperr "github.com/jira-backend/jiraflow-backend/internal/pkg/errors"
	pg "github.com/jira-backend/jiraflow-backend/internal/pkg/postgres"
	"github.com/jira-backend/jiraflow-backend/internal/infrastructura/repository"
)

type attachmentRepo struct {
	db      *pgxpool.Pool
	builder sq.StatementBuilderType
}

func NewAttachmentRepo(p *pg.Postgres) repository.AttachmentRepository {
	return &attachmentRepo{db: p.DB, builder: p.Builder}
}

const attachmentCols = "id, parent_type, parent_id, file_name, file_size, mime_type, storage_path, storage_type, checksum, uploaded_by, created_at, deleted_at"

func scanAttachment(row pgx.Row) (*entity.Attachment, error) {
	a := &entity.Attachment{}
	err := row.Scan(
		&a.ID, &a.ParentType, &a.ParentID,
		&a.FileName, &a.FileSize, &a.MimeType,
		&a.StoragePath, &a.StorageType, &a.Checksum,
		&a.UploadedBy, &a.CreatedAt, &a.DeletedAt,
	)
	return a, err
}

func (r *attachmentRepo) Create(ctx context.Context, a *entity.Attachment) error {
	sql, args, err := r.builder.
		Insert("attachments").
		Columns("id", "parent_type", "parent_id", "file_name", "file_size", "mime_type",
			"storage_path", "storage_type", "checksum", "uploaded_by", "created_at").
		Values(a.ID, a.ParentType, a.ParentID, a.FileName, a.FileSize, a.MimeType,
			a.StoragePath, a.StorageType, a.Checksum, a.UploadedBy, a.CreatedAt).
		ToSql()
	if err != nil {
		return fmt.Errorf("attachmentRepo.Create: %w", err)
	}
	_, err = r.db.Exec(ctx, sql, args...)
	return err
}

func (r *attachmentRepo) GetByID(ctx context.Context, id string) (*entity.Attachment, error) {
	sql, args, err := r.builder.
		Select(attachmentCols).From("attachments").
		Where(sq.And{sq.Eq{"id": id}, sq.Eq{"deleted_at": nil}}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("attachmentRepo.GetByID: %w", err)
	}
	a, err := scanAttachment(r.db.QueryRow(ctx, sql, args...))
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, apperr.NotFound("attachment")
	}
	return a, err
}

func (r *attachmentRepo) ListByParent(ctx context.Context, parentType, parentID string) ([]*entity.Attachment, error) {
	sql, args, err := r.builder.
		Select(attachmentCols).From("attachments").
		Where(sq.And{sq.Eq{"parent_type": parentType}, sq.Eq{"parent_id": parentID}, sq.Eq{"deleted_at": nil}}).
		OrderBy("created_at ASC").ToSql()
	if err != nil {
		return nil, fmt.Errorf("attachmentRepo.ListByParent: %w", err)
	}
	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("attachmentRepo.ListByParent query: %w", err)
	}
	defer rows.Close()

	var list []*entity.Attachment
	for rows.Next() {
		a, err := scanAttachment(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, a)
	}
	return list, rows.Err()
}

func (r *attachmentRepo) SoftDelete(ctx context.Context, id string) error {
	sql, args, err := r.builder.
		Update("attachments").
		Set("deleted_at", sq.Expr("NOW()")).
		Where(sq.And{sq.Eq{"id": id}, sq.Eq{"deleted_at": nil}}).
		ToSql()
	if err != nil {
		return fmt.Errorf("attachmentRepo.SoftDelete: %w", err)
	}
	_, err = r.db.Exec(ctx, sql, args...)
	return err
}

func (r *attachmentRepo) GetTotalSizeByParent(ctx context.Context, parentType, parentID string) (int64, error) {
	var total int64
	err := r.db.QueryRow(ctx,
		`SELECT COALESCE(SUM(file_size), 0) FROM attachments WHERE parent_type=$1 AND parent_id=$2 AND deleted_at IS NULL`,
		parentType, parentID,
	).Scan(&total)
	return total, err
}
