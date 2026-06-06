package user

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
	"github.com/jira-backend/jiraflow-backend/internal/infrastructure/repository"
	apperr "github.com/jira-backend/jiraflow-backend/internal/pkg/errors"
	"github.com/jira-backend/jiraflow-backend/internal/pkg/hasher"
	"github.com/jira-backend/jiraflow-backend/internal/pkg/logger"
)

type useCase struct {
	repo      repository.UserRepository
	spaceRepo repository.SpaceRepository
	hasher    hasher.Hasher
	log       logger.Logger
}

func New(repo repository.UserRepository, spaceRepo repository.SpaceRepository, h hasher.Hasher, log logger.Logger) UseCase {
	return &useCase{repo: repo, spaceRepo: spaceRepo, hasher: h, log: log}
}

func (uc *useCase) Create(ctx context.Context, req *entity.CreateUserReq) (*entity.User, error) {
	hashed, err := uc.hasher.Hash(req.Password)
	if err != nil {
		uc.log.Error(ctx, "user.Create: hash failed", logger.SafeString("err", err.Error()))
		return nil, fmt.Errorf("user.Create hash: %w", err)
	}
	role := req.Role
	if role == "" {
		role = "member"
	}
	u := &entity.User{
		ID:           uuid.NewString(),
		Email:        req.Email,
		PasswordHash: hashed,
		FullName:     req.FullName,
		Color:        "#6366F1",
		Role:         role,
		Timezone:     "UTC",
		Language:     "en",
		IsActive:     true,
	}
	if err := uc.repo.Create(ctx, u); err != nil {
		uc.log.Error(ctx, "user.Create: db error", logger.SafeEmail("email", req.Email), logger.SafeString("err", err.Error()))
		return nil, err
	}

	go uc.createPersonalSpace(ctx, u)

	uc.log.Info(ctx, "user created", logger.String("id", u.ID))
	return u, nil
}

func (uc *useCase) createPersonalSpace(ctx context.Context, u *entity.User) {
	username := strings.ToLower(strings.ReplaceAll(u.FullName, " ", "."))
	key := "~" + username
	now := time.Now().UTC()
	s := &entity.Space{
		ID:        uuid.NewString(),
		Key:       key,
		Name:      u.FullName + "'s Space",
		Type:      "personal",
		LeadID:    u.ID,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := uc.spaceRepo.Create(ctx, s); err != nil {
		uc.log.Warn(ctx, "createPersonalSpace: failed", logger.String("user_id", u.ID), logger.SafeString("err", err.Error()))
		return
	}
	member := &entity.SpaceMember{
		SpaceID:   s.ID,
		UserID:    u.ID,
		Role:      "admin",
		CreatedAt: now,
	}
	_ = uc.spaceRepo.AddMember(ctx, member)
}

func (uc *useCase) GetByID(ctx context.Context, id string) (*entity.User, error) {
	u, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		uc.log.Debug(ctx, "user.GetByID: not found", logger.String("id", id))
		return nil, err
	}
	return u, nil
}

func (uc *useCase) List(ctx context.Context, filter *entity.UserFilter) ([]*entity.User, int, error) {
	return uc.repo.List(ctx, filter)
}

func (uc *useCase) Update(ctx context.Context, id string, req *entity.UpdateUserReq) (*entity.User, error) {
	u, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if req.FullName != nil {
		u.FullName = *req.FullName
	}
	if req.AvatarURL != nil {
		u.AvatarURL = req.AvatarURL
	}
	if req.Color != nil {
		u.Color = *req.Color
	}
	if req.Timezone != nil {
		u.Timezone = *req.Timezone
	}
	if req.Language != nil {
		u.Language = *req.Language
	}
	if err := uc.repo.Update(ctx, u); err != nil {
		uc.log.Error(ctx, "user.Update: db error", logger.String("id", id), logger.SafeString("err", err.Error()))
		return nil, err
	}
	uc.log.Info(ctx, "user updated", logger.String("id", id))
	return u, nil
}

func (uc *useCase) ChangePassword(ctx context.Context, id string, req *entity.ChangePasswordReq) error {
	u, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if !uc.hasher.Check(req.CurrentPassword, u.PasswordHash) {
		uc.log.Warn(ctx, "user.ChangePassword: wrong current password", logger.String("id", id))
		return apperr.BadRequest("current password is incorrect")
	}
	hashed, err := uc.hasher.Hash(req.NewPassword)
	if err != nil {
		uc.log.Error(ctx, "user.ChangePassword: hash failed", logger.SafeString("err", err.Error()))
		return fmt.Errorf("user.ChangePassword hash: %w", err)
	}
	if err := uc.repo.UpdatePassword(ctx, id, hashed); err != nil {
		uc.log.Error(ctx, "user.ChangePassword: db error", logger.String("id", id), logger.SafeString("err", err.Error()))
		return err
	}
	uc.log.Info(ctx, "user password changed", logger.String("id", id))
	return nil
}

func (uc *useCase) Deactivate(ctx context.Context, id string) error {
	u, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	u.IsActive = false
	if err := uc.repo.Update(ctx, u); err != nil {
		uc.log.Error(ctx, "user.Deactivate: db error", logger.String("id", id), logger.SafeString("err", err.Error()))
		return err
	}
	uc.log.Info(ctx, "user deactivated", logger.String("id", id))
	return nil
}

func (uc *useCase) Activate(ctx context.Context, id string) error {
	u, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	u.IsActive = true
	if err := uc.repo.Update(ctx, u); err != nil {
		uc.log.Error(ctx, "user.Activate: db error", logger.String("id", id), logger.SafeString("err", err.Error()))
		return err
	}
	uc.log.Info(ctx, "user activated", logger.String("id", id))
	return nil
}
