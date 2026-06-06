package attachment

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/google/uuid"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
	"github.com/jira-backend/jiraflow-backend/internal/infrastructura/minio"
	"github.com/jira-backend/jiraflow-backend/internal/infrastructura/repository"
	apperr "github.com/jira-backend/jiraflow-backend/internal/pkg/errors"
	"github.com/jira-backend/jiraflow-backend/internal/pkg/logger"
)

const presignTTL = 24 * time.Hour

type useCase struct {
	attachRepo  repository.AttachmentRepository
	minioClient minio.Client
	log         logger.Logger
}

func New(attachRepo repository.AttachmentRepository, minioClient minio.Client, log logger.Logger) UseCase {
	return &useCase{attachRepo: attachRepo, minioClient: minioClient, log: log}
}

func (uc *useCase) Upload(ctx context.Context, parentType, parentID, uploaderID string, fileName string, size int64, mimeType string, r io.Reader) (*entity.Attachment, error) {
	id := uuid.NewString()
	objectName := fmt.Sprintf("%s/%s/%s/%s", parentType, parentID, id, fileName)

	storagePath, err := uc.minioClient.Upload(ctx, objectName, mimeType, r, size)
	if err != nil {
		uc.log.Error(ctx, "attachment.Upload: minio error", logger.String("parent_id", parentID), logger.SafeString("err", err.Error()))
		return nil, fmt.Errorf("attachment.Upload: %w", err)
	}

	a := &entity.Attachment{
		ID:          id,
		ParentType:  parentType,
		ParentID:    parentID,
		FileName:    fileName,
		FileSize:    size,
		MimeType:    mimeType,
		StoragePath: storagePath,
		StorageType: "s3",
		UploadedBy:  uploaderID,
		CreatedAt:   time.Now().UTC(),
	}
	if err := uc.attachRepo.Create(ctx, a); err != nil {
		uc.log.Error(ctx, "attachment.Upload: db error", logger.String("parent_id", parentID), logger.SafeString("err", err.Error()))
		return nil, err
	}
	uc.log.Info(ctx, "attachment uploaded", logger.String("id", a.ID), logger.String("parent_id", parentID))
	return a, nil
}

func (uc *useCase) GetByID(ctx context.Context, id string) (*entity.Attachment, error) {
	a, err := uc.attachRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	url, _ := uc.minioClient.PresignedURL(ctx, a.StoragePath, presignTTL)
	a.DownloadURL = url
	return a, nil
}

func (uc *useCase) ListByParent(ctx context.Context, parentType, parentID string) ([]*entity.Attachment, error) {
	attachments, err := uc.attachRepo.ListByParent(ctx, parentType, parentID)
	if err != nil {
		return nil, err
	}
	for _, a := range attachments {
		url, _ := uc.minioClient.PresignedURL(ctx, a.StoragePath, presignTTL)
		a.DownloadURL = url
	}
	return attachments, nil
}

func (uc *useCase) Delete(ctx context.Context, id, actorID string) error {
	a, err := uc.attachRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if a.UploadedBy != actorID {
		uc.log.Warn(ctx, "attachment.Delete: forbidden", logger.String("id", id), logger.String("actor_id", actorID))
		return apperr.Forbidden("only the uploader can delete this attachment")
	}
	_ = uc.minioClient.Delete(ctx, a.StoragePath)
	if err := uc.attachRepo.SoftDelete(ctx, id); err != nil {
		uc.log.Error(ctx, "attachment.Delete: db error", logger.String("id", id), logger.SafeString("err", err.Error()))
		return err
	}
	uc.log.Info(ctx, "attachment deleted", logger.String("id", id))
	return nil
}

func (uc *useCase) PresignedURL(ctx context.Context, id string) (string, error) {
	a, err := uc.attachRepo.GetByID(ctx, id)
	if err != nil {
		return "", err
	}
	return uc.minioClient.PresignedURL(ctx, a.StoragePath, presignTTL)
}
