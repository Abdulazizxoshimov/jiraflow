package invite

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
	"github.com/jira-backend/jiraflow-backend/internal/infrastructura/repository"
	apperr "github.com/jira-backend/jiraflow-backend/internal/pkg/errors"
	"github.com/jira-backend/jiraflow-backend/internal/pkg/hasher"
	"github.com/jira-backend/jiraflow-backend/internal/pkg/logger"
	"github.com/jira-backend/jiraflow-backend/internal/pkg/token"
)

const inviteTTL = 7 * 24 * time.Hour

type useCase struct {
	inviteRepo repository.InviteRepository
	userRepo   repository.UserRepository
	tokens     token.Maker
	hasher     hasher.Hasher
	log        logger.Logger
}

func New(
	inviteRepo repository.InviteRepository,
	userRepo repository.UserRepository,
	tokens token.Maker,
	h hasher.Hasher,
	log logger.Logger,
) UseCase {
	return &useCase{
		inviteRepo: inviteRepo,
		userRepo:   userRepo,
		tokens:     tokens,
		hasher:     h,
		log:        log,
	}
}

func (uc *useCase) Create(ctx context.Context, req *entity.CreateInviteReq, invitedBy string) (*entity.Invite, error) {
	existing, err := uc.inviteRepo.GetPendingByEmail(ctx, req.Email)
	if err == nil && existing != nil {
		return nil, apperr.Conflict("pending invite already exists for this email")
	}

	rawToken, err := generateToken()
	if err != nil {
		uc.log.Error(ctx, "invite.Create: generate token failed", logger.SafeString("err", err.Error()))
		return nil, fmt.Errorf("invite.Create generate token: %w", err)
	}

	now := time.Now().UTC()
	invite := &entity.Invite{
		ID:        uuid.NewString(),
		Email:     req.Email,
		Role:      req.Role,
		TokenHash: hashToken(rawToken),
		InvitedBy: invitedBy,
		ExpiresAt: now.Add(inviteTTL),
		CreatedAt: now,
	}

	if err := uc.inviteRepo.Create(ctx, invite); err != nil {
		uc.log.Error(ctx, "invite.Create: db error", logger.SafeEmail("email", req.Email), logger.SafeString("err", err.Error()))
		return nil, err
	}
	uc.log.Info(ctx, "invite created", logger.String("id", invite.ID), logger.String("invited_by", invitedBy))

	_ = rawToken
	return invite, nil
}

func (uc *useCase) Accept(ctx context.Context, req *entity.AcceptInviteReq) (*entity.TokenPair, error) {
	h := hashToken(req.Token)
	invite, err := uc.inviteRepo.GetByTokenHash(ctx, h)
	if err != nil {
		return nil, apperr.BadRequest("invalid or expired invite token")
	}
	if time.Now().UTC().After(invite.ExpiresAt) {
		return nil, apperr.BadRequest("invite token has expired")
	}
	if invite.AcceptedAt != nil {
		return nil, apperr.BadRequest("invite already accepted")
	}

	passwordHash, err := uc.hasher.Hash(req.Password)
	if err != nil {
		uc.log.Error(ctx, "invite.Accept: hash password failed", logger.SafeString("err", err.Error()))
		return nil, fmt.Errorf("invite.Accept hash password: %w", err)
	}

	now := time.Now().UTC()
	user := &entity.User{
		ID:           uuid.NewString(),
		Email:        invite.Email,
		FullName:     req.FullName,
		PasswordHash: passwordHash,
		Role:         invite.Role,
		IsActive:     true,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	if err := uc.userRepo.Create(ctx, user); err != nil {
		uc.log.Error(ctx, "invite.Accept: create user failed", logger.SafeEmail("email", invite.Email), logger.SafeString("err", err.Error()))
		return nil, fmt.Errorf("invite.Accept create user: %w", err)
	}

	if err := uc.inviteRepo.MarkAccepted(ctx, invite.ID, user.ID); err != nil {
		uc.log.Error(ctx, "invite.Accept: mark accepted failed", logger.String("invite_id", invite.ID), logger.SafeString("err", err.Error()))
		return nil, fmt.Errorf("invite.Accept mark accepted: %w", err)
	}

	sessionID := uuid.NewString()
	access, refresh, err := uc.tokens.Generate(ctx, user.ID, sessionID, user.Role)
	if err != nil {
		uc.log.Error(ctx, "invite.Accept: generate tokens failed", logger.String("user_id", user.ID), logger.SafeString("err", err.Error()))
		return nil, fmt.Errorf("invite.Accept generate tokens: %w", err)
	}
	_ = uc.tokens.StoreSession(ctx, sessionID, user.ID, 30*24*time.Hour)

	uc.log.Info(ctx, "invite accepted", logger.String("user_id", user.ID), logger.String("invite_id", invite.ID))
	return &entity.TokenPair{AccessToken: access, RefreshToken: refresh}, nil
}

func (uc *useCase) ListPending(ctx context.Context, filter *entity.Filter) ([]*entity.Invite, int, error) {
	return uc.inviteRepo.List(ctx, filter)
}

func (uc *useCase) Revoke(ctx context.Context, id string) error {
	if err := uc.inviteRepo.Delete(ctx, id); err != nil {
		uc.log.Error(ctx, "invite.Revoke: db error", logger.String("id", id), logger.SafeString("err", err.Error()))
		return err
	}
	uc.log.Info(ctx, "invite revoked", logger.String("id", id))
	return nil
}

func generateToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func hashToken(t string) string {
	sum := sha256.Sum256([]byte(t))
	return fmt.Sprintf("%x", sum)
}
