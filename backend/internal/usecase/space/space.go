package space

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
	"github.com/jira-backend/jiraflow-backend/internal/infrastructure/repository"
	apperr "github.com/jira-backend/jiraflow-backend/internal/pkg/errors"
	"github.com/jira-backend/jiraflow-backend/internal/pkg/logger"
)

type useCase struct {
	spaceRepo repository.SpaceRepository
	log       logger.Logger
}

func New(spaceRepo repository.SpaceRepository, log logger.Logger) UseCase {
	return &useCase{spaceRepo: spaceRepo, log: log}
}

func (uc *useCase) Create(ctx context.Context, s *entity.Space, leadID string) (*entity.Space, error) {
	exists, err := uc.spaceRepo.ExistsByKey(ctx, s.Key)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, apperr.Conflict("space key already exists")
	}

	now := time.Now().UTC()
	s.ID = uuid.NewString()
	s.LeadID = leadID
	s.CreatedAt = now
	s.UpdatedAt = now

	if err := uc.spaceRepo.Create(ctx, s); err != nil {
		uc.log.Error(ctx, "space.Create: db error", logger.SafeString("err", err.Error()))
		return nil, err
	}

	member := &entity.SpaceMember{
		SpaceID:   s.ID,
		UserID:    leadID,
		Role:      "admin",
		CreatedAt: now,
	}
	_ = uc.spaceRepo.AddMember(ctx, member)

	uc.log.Info(ctx, "space created", logger.String("id", s.ID), logger.String("lead_id", leadID))
	return s, nil
}

func (uc *useCase) GetByID(ctx context.Context, id string) (*entity.Space, error) {
	return uc.spaceRepo.GetByID(ctx, id)
}

func (uc *useCase) GetByKey(ctx context.Context, key string) (*entity.Space, error) {
	return uc.spaceRepo.GetByKey(ctx, key)
}

func (uc *useCase) List(ctx context.Context, filter *entity.Filter) ([]*entity.Space, int, error) {
	return uc.spaceRepo.List(ctx, filter)
}

func (uc *useCase) Update(ctx context.Context, id string, s *entity.Space, actorID string) (*entity.Space, error) {
	existing, err := uc.spaceRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if s.Name != "" {
		existing.Name = s.Name
	}
	if s.Description != nil {
		existing.Description = s.Description
	}
	if s.IconURL != nil {
		existing.IconURL = s.IconURL
	}
	existing.IsArchived = s.IsArchived
	existing.UpdatedAt = time.Now().UTC()

	if err := uc.spaceRepo.Update(ctx, existing); err != nil {
		uc.log.Error(ctx, "space.Update: db error", logger.String("id", id), logger.SafeString("err", err.Error()))
		return nil, err
	}
	uc.log.Info(ctx, "space updated", logger.String("id", id))
	return existing, nil
}

func (uc *useCase) Delete(ctx context.Context, id, actorID string) error {
	s, err := uc.spaceRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if s.LeadID != actorID {
		uc.log.Warn(ctx, "space.Delete: forbidden", logger.String("id", id), logger.String("actor_id", actorID))
		return apperr.Forbidden("only the space lead can delete this space")
	}
	if err := uc.spaceRepo.SoftDelete(ctx, id); err != nil {
		uc.log.Error(ctx, "space.Delete: db error", logger.String("id", id), logger.SafeString("err", err.Error()))
		return err
	}
	uc.log.Info(ctx, "space deleted", logger.String("id", id))
	return nil
}

func (uc *useCase) AddMember(ctx context.Context, spaceID string, m *entity.SpaceMember) error {
	isMember, err := uc.spaceRepo.IsMember(ctx, spaceID, m.UserID)
	if err != nil {
		return err
	}
	if isMember {
		return apperr.Conflict("user is already a member")
	}
	m.SpaceID = spaceID
	if m.CreatedAt.IsZero() {
		m.CreatedAt = time.Now().UTC()
	}
	if err := uc.spaceRepo.AddMember(ctx, m); err != nil {
		uc.log.Error(ctx, "space.AddMember: db error", logger.String("space_id", spaceID), logger.String("user_id", m.UserID), logger.SafeString("err", err.Error()))
		return err
	}
	uc.log.Info(ctx, "space member added", logger.String("space_id", spaceID), logger.String("user_id", m.UserID))
	return nil
}

func (uc *useCase) UpdateMemberRole(ctx context.Context, spaceID, userID, role string) error {
	if err := uc.spaceRepo.UpdateMemberRole(ctx, spaceID, userID, role); err != nil {
		uc.log.Error(ctx, "space.UpdateMemberRole: db error", logger.String("space_id", spaceID), logger.String("user_id", userID), logger.SafeString("err", err.Error()))
		return err
	}
	uc.log.Info(ctx, "space member role updated", logger.String("space_id", spaceID), logger.String("user_id", userID), logger.String("role", role))
	return nil
}

func (uc *useCase) RemoveMember(ctx context.Context, spaceID, userID string) error {
	if err := uc.spaceRepo.RemoveMember(ctx, spaceID, userID); err != nil {
		uc.log.Error(ctx, "space.RemoveMember: db error", logger.String("space_id", spaceID), logger.String("user_id", userID), logger.SafeString("err", err.Error()))
		return err
	}
	uc.log.Info(ctx, "space member removed", logger.String("space_id", spaceID), logger.String("user_id", userID))
	return nil
}

func (uc *useCase) ListMembers(ctx context.Context, spaceID string, filter *entity.Filter) ([]*entity.SpaceMember, int, error) {
	return uc.spaceRepo.ListMembers(ctx, spaceID, filter)
}

func (uc *useCase) Archive(ctx context.Context, id, actorID string) error {
	s, err := uc.spaceRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if s.IsArchived {
		return apperr.BadRequest("space is already archived")
	}
	if s.LeadID != actorID {
		member, merr := uc.spaceRepo.GetMember(ctx, id, actorID)
		if merr != nil || member.Role != "admin" {
			return apperr.Forbidden("only space admin can archive this space")
		}
	}
	if err := uc.spaceRepo.Archive(ctx, id); err != nil {
		uc.log.Error(ctx, "space.Archive: db error", logger.String("id", id), logger.SafeString("err", err.Error()))
		return err
	}
	uc.log.Info(ctx, "space archived", logger.String("id", id), logger.String("actor_id", actorID))
	return nil
}

func (uc *useCase) Restore(ctx context.Context, id, actorID string) error {
	s, err := uc.spaceRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if !s.IsArchived {
		return apperr.BadRequest("space is not archived")
	}
	if s.LeadID != actorID {
		member, merr := uc.spaceRepo.GetMember(ctx, id, actorID)
		if merr != nil || member.Role != "admin" {
			return apperr.Forbidden("only space admin can restore this space")
		}
	}
	if err := uc.spaceRepo.Restore(ctx, id); err != nil {
		uc.log.Error(ctx, "space.Restore: db error", logger.String("id", id), logger.SafeString("err", err.Error()))
		return err
	}
	uc.log.Info(ctx, "space restored", logger.String("id", id), logger.String("actor_id", actorID))
	return nil
}

func (uc *useCase) GetStatistics(ctx context.Context, id string) (*entity.SpaceStatistics, error) {
	if _, err := uc.spaceRepo.GetByID(ctx, id); err != nil {
		return nil, err
	}
	stats, err := uc.spaceRepo.GetStatistics(ctx, id)
	if err != nil {
		uc.log.Error(ctx, "space.GetStatistics: db error", logger.String("id", id), logger.SafeString("err", err.Error()))
		return nil, err
	}
	return stats, nil
}
