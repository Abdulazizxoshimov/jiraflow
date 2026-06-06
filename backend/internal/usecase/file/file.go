package file

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/google/uuid"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
	"github.com/jira-backend/jiraflow-backend/internal/infrastructura/minio"
	"github.com/jira-backend/jiraflow-backend/internal/pkg/logger"
)

const presignTTL = 1 * time.Hour

type useCase struct {
	minioClient minio.Client
	log         logger.Logger
}

func New(minioClient minio.Client, log logger.Logger) UseCase {
	return &useCase{minioClient: minioClient, log: log}
}

func (uc *useCase) Upload(ctx context.Context, uploaderID string, name string, size int64, mimeType string, r io.Reader) (*entity.Attachment, error) {
	id := uuid.NewString()
	objectName := fmt.Sprintf("files/%s/%s", id, name)

	storagePath, err := uc.minioClient.Upload(ctx, objectName, mimeType, r, size)
	if err != nil {
		uc.log.Error(ctx, "file.Upload: minio error", logger.String("uploader_id", uploaderID), logger.SafeString("err", err.Error()))
		return nil, fmt.Errorf("file.Upload: %w", err)
	}

	url, _ := uc.minioClient.PresignedURL(ctx, storagePath, presignTTL)

	a := &entity.Attachment{
		ID:          id,
		ParentType:  "file",
		ParentID:    uploaderID,
		FileName:    name,
		FileSize:    size,
		MimeType:    mimeType,
		StoragePath: storagePath,
		StorageType: "s3",
		UploadedBy:  uploaderID,
		CreatedAt:   time.Now().UTC(),
		DownloadURL: url,
	}
	uc.log.Info(ctx, "file uploaded", logger.String("id", id), logger.String("uploader_id", uploaderID))
	return a, nil
}

func (uc *useCase) GetPresignedURL(ctx context.Context, storagePath string) (string, error) {
	return uc.minioClient.PresignedURL(ctx, storagePath, presignTTL)
}
