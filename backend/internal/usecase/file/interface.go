package file

import (
	"context"
	"io"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

type UseCase interface {
	Upload(ctx context.Context, uploaderID string, name string, size int64, mimeType string, r io.Reader) (*entity.Attachment, error)
	GetPresignedURL(ctx context.Context, storagePath string) (string, error)
}
