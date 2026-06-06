package notification_scheme

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
	"github.com/jira-backend/jiraflow-backend/internal/infrastructure/repository"
)

type UseCase interface {
	Create(ctx context.Context, req *entity.CreateNotificationSchemeReq) (*entity.NotificationScheme, error)
	GetByID(ctx context.Context, id string) (*entity.NotificationScheme, error)
	List(ctx context.Context) ([]*entity.NotificationScheme, error)
	Delete(ctx context.Context, id string) error
}

type useCase struct {
	repo repository.NotificationSchemeRepository
}

func New(repo repository.NotificationSchemeRepository) UseCase {
	return &useCase{repo: repo}
}

func (uc *useCase) Create(ctx context.Context, req *entity.CreateNotificationSchemeReq) (*entity.NotificationScheme, error) {
	s := &entity.NotificationScheme{
		ID:          uuid.NewString(),
		Name:        req.Name,
		Description: req.Description,
	}
	for _, r := range req.Rules {
		s.Rules = append(s.Rules, &entity.NotificationSchemeRule{
			ID:            uuid.NewString(),
			SchemeID:      s.ID,
			EventType:     r.EventType,
			RecipientType: r.RecipientType,
			RecipientID:   r.RecipientID,
		})
	}
	if err := uc.repo.Create(ctx, s); err != nil {
		return nil, fmt.Errorf("notification_scheme.Create: %w", err)
	}
	return uc.repo.GetByID(ctx, s.ID)
}

func (uc *useCase) GetByID(ctx context.Context, id string) (*entity.NotificationScheme, error) {
	return uc.repo.GetByID(ctx, id)
}

func (uc *useCase) List(ctx context.Context) ([]*entity.NotificationScheme, error) {
	return uc.repo.List(ctx)
}

func (uc *useCase) Delete(ctx context.Context, id string) error {
	return uc.repo.Delete(ctx, id)
}
