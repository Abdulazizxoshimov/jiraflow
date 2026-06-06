package attachment

import (
	"context"
	"io"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
)

type UseCase interface {
	Upload(ctx context.Context, parentType, parentID, uploaderID string, fileName string, size int64, mimeType string, r io.Reader) (*entity.Attachment, error)
	GetByID(ctx context.Context, id string) (*entity.Attachment, error)
	ListByParent(ctx context.Context, parentType, parentID string) ([]*entity.Attachment, error)
	Delete(ctx context.Context, id, actorID string) error
	PresignedURL(ctx context.Context, id string) (string, error)
}
