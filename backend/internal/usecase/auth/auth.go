package auth

import (
	"context"
	"crypto/sha256"
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

type useCase struct {
	userRepo repository.UserRepository
	authRepo repository.AuthRepository
	tokens   token.Maker
	hasher   hasher.Hasher
	resetTTL time.Duration
	log      logger.Logger
}

func New(
	userRepo repository.UserRepository,
	authRepo repository.AuthRepository,
	tokens token.Maker,
	h hasher.Hasher,
	resetTTL time.Duration,
	log logger.Logger,
) UseCase {
	return &useCase{
		userRepo: userRepo,
		authRepo: authRepo,
		tokens:   tokens,
		hasher:   h,
		resetTTL: resetTTL,
		log:      log,
	}
}

func (uc *useCase) Register(ctx context.Context, req *entity.RegisterReq, ip, userAgent string) (*entity.TokenPair, error) {
	hashed, err := uc.hasher.Hash(req.Password)
	if err != nil {
		return nil, fmt.Errorf("auth.Register hash: %w", err)
	}
	now := time.Now().UTC()
	u := &entity.User{
		ID:           uuid.NewString(),
		Email:        req.Email,
		PasswordHash: hashed,
		FullName:     req.FullName,
		Color:        "#6366F1",
		Role:         "member",
		Timezone:     "UTC",
		Language:     "en",
		IsActive:     true,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	if err := uc.userRepo.Create(ctx, u); err != nil {
		return nil, err
	}

	sessionID := uuid.NewString()
	access, refresh, err := uc.tokens.Generate(ctx, u.ID, sessionID, u.Role)
	if err != nil {
		return nil, fmt.Errorf("auth.Register generate tokens: %w", err)
	}

	rt := &entity.RefreshToken{
		ID:        uuid.NewString(),
		UserID:    u.ID,
		TokenHash: hashToken(refresh),
		ExpiresAt: now.Add(30 * 24 * time.Hour),
		CreatedAt: now,
	}
	if ip != "" {
		rt.IPAddress = &ip
	}
	if userAgent != "" {
		rt.UserAgent = &userAgent
	}
	if err := uc.authRepo.CreateRefreshToken(ctx, rt); err != nil {
		return nil, fmt.Errorf("auth.Register store token: %w", err)
	}
	if err := uc.tokens.StoreSession(ctx, sessionID, u.ID, 30*24*time.Hour); err != nil {
		return nil, fmt.Errorf("auth.Register store session: %w", err)
	}

	uc.log.Info(ctx, "auth.Register: success", logger.String("user_id", u.ID))
	return &entity.TokenPair{AccessToken: access, RefreshToken: refresh}, nil
}

func (uc *useCase) Login(ctx context.Context, req *entity.LoginReq, ip, userAgent string) (*entity.TokenPair, error) {
	user, err := uc.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		uc.log.Warn(ctx, "auth.Login: user not found", logger.SafeEmail("email", req.Email))
		return nil, apperr.Unauthorized("invalid credentials")
	}
	if !user.IsActive {
		uc.log.Warn(ctx, "auth.Login: account deactivated", logger.String("user_id", user.ID))
		return nil, apperr.Forbidden("account is deactivated")
	}
	if !uc.hasher.Check(req.Password, user.PasswordHash) {
		uc.log.Warn(ctx, "auth.Login: wrong password", logger.String("user_id", user.ID))
		return nil, apperr.Unauthorized("invalid credentials")
	}

	sessionID := uuid.NewString()
	access, refresh, err := uc.tokens.Generate(ctx, user.ID, sessionID, user.Role)
	if err != nil {
		uc.log.Error(ctx, "auth.Login: generate tokens failed", logger.String("user_id", user.ID), logger.SafeString("err", err.Error()))
		return nil, fmt.Errorf("auth.Login generate tokens: %w", err)
	}

	now := time.Now().UTC()
	rt := &entity.RefreshToken{
		ID:        uuid.NewString(),
		UserID:    user.ID,
		TokenHash: hashToken(refresh),
		ExpiresAt: now.Add(30 * 24 * time.Hour),
		CreatedAt: now,
	}
	if ip != "" {
		rt.IPAddress = &ip
	}
	if userAgent != "" {
		rt.UserAgent = &userAgent
	}
	if err := uc.authRepo.CreateRefreshToken(ctx, rt); err != nil {
		uc.log.Error(ctx, "auth.Login: store refresh token failed", logger.String("user_id", user.ID), logger.SafeString("err", err.Error()))
		return nil, fmt.Errorf("auth.Login store token: %w", err)
	}
	if err := uc.tokens.StoreSession(ctx, sessionID, user.ID, 30*24*time.Hour); err != nil {
		uc.log.Error(ctx, "auth.Login: store session failed", logger.String("user_id", user.ID), logger.SafeString("err", err.Error()))
		return nil, fmt.Errorf("auth.Login store session: %w", err)
	}

	uc.log.Info(ctx, "auth.Login: success", logger.String("user_id", user.ID))
	return &entity.TokenPair{AccessToken: access, RefreshToken: refresh}, nil
}

func (uc *useCase) Refresh(ctx context.Context, req *entity.RefreshReq) (*entity.TokenPair, error) {
	access, refresh, err := uc.tokens.Rotate(ctx, req.RefreshToken)
	if err != nil {
		uc.log.Warn(ctx, "auth.Refresh: invalid token")
		return nil, apperr.Unauthorized("invalid or expired refresh token")
	}
	return &entity.TokenPair{AccessToken: access, RefreshToken: refresh}, nil
}

func (uc *useCase) Logout(ctx context.Context, req *entity.LogoutReq) error {
	h := hashToken(req.RefreshToken)
	rt, err := uc.authRepo.GetRefreshTokenByHash(ctx, h)
	if err != nil {
		return nil
	}
	_ = uc.authRepo.RevokeRefreshToken(ctx, rt.ID)
	_ = uc.tokens.RevokeSession(ctx, rt.UserID)
	uc.log.Info(ctx, "auth.Logout: success", logger.String("user_id", rt.UserID))
	return nil
}

func (uc *useCase) ForgotPassword(ctx context.Context, req *entity.ForgotPasswordReq) error {
	user, err := uc.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil
	}

	rawToken := uuid.NewString()
	h := hashToken(rawToken)
	now := time.Now().UTC()
	pr := &entity.PasswordReset{
		ID:        uuid.NewString(),
		UserID:    user.ID,
		TokenHash: h,
		ExpiresAt: now.Add(uc.resetTTL),
		CreatedAt: now,
	}
	if err := uc.authRepo.CreatePasswordReset(ctx, pr); err != nil {
		uc.log.Error(ctx, "auth.ForgotPassword: store reset token failed", logger.String("user_id", user.ID), logger.SafeString("err", err.Error()))
		return fmt.Errorf("auth.ForgotPassword: %w", err)
	}
	uc.log.Info(ctx, "auth.ForgotPassword: reset token created", logger.String("user_id", user.ID))
	_ = rawToken
	return nil
}

func (uc *useCase) ResetPassword(ctx context.Context, req *entity.ResetPasswordReq) error {
	h := hashToken(req.Token)
	pr, err := uc.authRepo.GetPasswordResetByHash(ctx, h)
	if err != nil {
		return apperr.BadRequest("invalid or expired reset token")
	}

	hashed, err := uc.hasher.Hash(req.NewPassword)
	if err != nil {
		uc.log.Error(ctx, "auth.ResetPassword: hash failed", logger.SafeString("err", err.Error()))
		return fmt.Errorf("auth.ResetPassword hash: %w", err)
	}
	if err := uc.userRepo.UpdatePassword(ctx, pr.UserID, hashed); err != nil {
		uc.log.Error(ctx, "auth.ResetPassword: update failed", logger.String("user_id", pr.UserID), logger.SafeString("err", err.Error()))
		return fmt.Errorf("auth.ResetPassword update: %w", err)
	}
	uc.log.Info(ctx, "auth.ResetPassword: success", logger.String("user_id", pr.UserID))
	return uc.authRepo.MarkPasswordResetUsed(ctx, pr.ID)
}

func hashToken(t string) string {
	sum := sha256.Sum256([]byte(t))
	return fmt.Sprintf("%x", sum)
}
