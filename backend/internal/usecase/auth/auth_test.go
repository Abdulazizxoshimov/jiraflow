package auth_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
	apperr "github.com/jira-backend/jiraflow-backend/internal/pkg/errors"
	"github.com/jira-backend/jiraflow-backend/internal/testutil"
	"github.com/jira-backend/jiraflow-backend/internal/usecase/auth"
)

func newUC(
	userRepo *testutil.UserRepoMock,
	authRepo *testutil.AuthRepoMock,
	tokens *testutil.TokenMakerMock,
	h *testutil.HasherMock,
) auth.UseCase {
	return auth.New(
		userRepo, authRepo, tokens, h,
		15*time.Minute, &testutil.EmailSenderMock{}, "http://localhost:3000",
		testutil.NopLogger{},
	)
}

// ─── Login ───────────────────────────────────────────────────────────────────

func TestLogin_Success(t *testing.T) {
	ctx := context.Background()

	userRepo := &testutil.UserRepoMock{
		GetByEmailFn: func(_ context.Context, email string) (*entity.User, error) {
			return &entity.User{ID: "user-1", Email: email, PasswordHash: "hashed:secret", IsActive: true, Role: "member"}, nil
		},
	}
	tokens := &testutil.TokenMakerMock{}
	h := &testutil.HasherMock{}

	uc := newUC(userRepo, &testutil.AuthRepoMock{}, tokens, h)
	pair, err := uc.Login(ctx, &entity.LoginReq{Email: "a@b.com", Password: "secret"}, "", "")

	require.NoError(t, err)
	assert.Equal(t, "access-token", pair.AccessToken)
	assert.Equal(t, "refresh-token", pair.RefreshToken)
}

func TestLogin_WrongPassword(t *testing.T) {
	ctx := context.Background()

	userRepo := &testutil.UserRepoMock{
		GetByEmailFn: func(_ context.Context, email string) (*entity.User, error) {
			return &entity.User{ID: "user-1", Email: email, PasswordHash: "hashed:correct", IsActive: true}, nil
		},
	}
	h := &testutil.HasherMock{}

	uc := newUC(userRepo, &testutil.AuthRepoMock{}, &testutil.TokenMakerMock{}, h)
	_, err := uc.Login(ctx, &entity.LoginReq{Email: "a@b.com", Password: "wrong"}, "", "")

	require.Error(t, err)
	var appErr *apperr.AppError
	require.True(t, errors.As(err, &appErr))
	assert.Equal(t, apperr.CodeUnauthorized, appErr.Code)
}

func TestLogin_UserNotFound(t *testing.T) {
	ctx := context.Background()

	userRepo := &testutil.UserRepoMock{
		GetByEmailFn: func(_ context.Context, _ string) (*entity.User, error) {
			return nil, apperr.NotFound("user not found")
		},
	}

	uc := newUC(userRepo, &testutil.AuthRepoMock{}, &testutil.TokenMakerMock{}, &testutil.HasherMock{})
	_, err := uc.Login(ctx, &entity.LoginReq{Email: "x@y.com", Password: "pw"}, "", "")

	require.Error(t, err)
	var appErr *apperr.AppError
	require.True(t, errors.As(err, &appErr))
	assert.Equal(t, apperr.CodeUnauthorized, appErr.Code)
}

func TestLogin_InactiveUser(t *testing.T) {
	ctx := context.Background()

	userRepo := &testutil.UserRepoMock{
		GetByEmailFn: func(_ context.Context, email string) (*entity.User, error) {
			return &entity.User{ID: "user-1", Email: email, PasswordHash: "hashed:pw", IsActive: false}, nil
		},
	}

	uc := newUC(userRepo, &testutil.AuthRepoMock{}, &testutil.TokenMakerMock{}, &testutil.HasherMock{})
	_, err := uc.Login(ctx, &entity.LoginReq{Email: "a@b.com", Password: "pw"}, "", "")

	require.Error(t, err)
	var appErr *apperr.AppError
	require.True(t, errors.As(err, &appErr))
	assert.Equal(t, apperr.CodeForbidden, appErr.Code)
}

// ─── Register ─────────────────────────────────────────────────────────────────

func TestRegister_Success(t *testing.T) {
	ctx := context.Background()

	var created *entity.User
	userRepo := &testutil.UserRepoMock{
		CreateFn: func(_ context.Context, u *entity.User) error {
			created = u
			return nil
		},
	}

	uc := newUC(userRepo, &testutil.AuthRepoMock{}, &testutil.TokenMakerMock{}, &testutil.HasherMock{})
	pair, err := uc.Register(ctx, &entity.RegisterReq{
		Email:    "new@user.com",
		Password: "password123",
		FullName: "New User",
	}, "127.0.0.1", "test-agent")

	require.NoError(t, err)
	assert.NotEmpty(t, pair.AccessToken)
	require.NotNil(t, created)
	assert.Equal(t, "new@user.com", created.Email)
	assert.Equal(t, "member", created.Role)
	assert.True(t, created.IsActive)
}

func TestRegister_DuplicateEmail(t *testing.T) {
	ctx := context.Background()

	userRepo := &testutil.UserRepoMock{
		CreateFn: func(_ context.Context, _ *entity.User) error {
			return apperr.Conflict("email already exists")
		},
	}

	uc := newUC(userRepo, &testutil.AuthRepoMock{}, &testutil.TokenMakerMock{}, &testutil.HasherMock{})
	_, err := uc.Register(ctx, &entity.RegisterReq{
		Email: "dup@user.com", Password: "pw", FullName: "Dup",
	}, "", "")

	require.Error(t, err)
}

// ─── Refresh ─────────────────────────────────────────────────────────────────

func TestRefresh_Success(t *testing.T) {
	ctx := context.Background()

	tokens := &testutil.TokenMakerMock{
		RotateFn: func(_ context.Context, _ string) (string, string, error) {
			return "new-access", "new-refresh", nil
		},
	}

	uc := newUC(&testutil.UserRepoMock{}, &testutil.AuthRepoMock{}, tokens, &testutil.HasherMock{})
	pair, err := uc.Refresh(ctx, &entity.RefreshReq{RefreshToken: "old-refresh"})

	require.NoError(t, err)
	assert.Equal(t, "new-access", pair.AccessToken)
	assert.Equal(t, "new-refresh", pair.RefreshToken)
}

func TestRefresh_InvalidToken(t *testing.T) {
	ctx := context.Background()

	tokens := &testutil.TokenMakerMock{
		RotateFn: func(_ context.Context, _ string) (string, string, error) {
			return "", "", errors.New("invalid token")
		},
	}

	uc := newUC(&testutil.UserRepoMock{}, &testutil.AuthRepoMock{}, tokens, &testutil.HasherMock{})
	_, err := uc.Refresh(ctx, &entity.RefreshReq{RefreshToken: "bad-token"})

	require.Error(t, err)
	var appErr *apperr.AppError
	require.True(t, errors.As(err, &appErr))
	assert.Equal(t, apperr.CodeUnauthorized, appErr.Code)
}

// ─── Logout ──────────────────────────────────────────────────────────────────

func TestLogout_Success(t *testing.T) {
	ctx := context.Background()

	var revokedID string
	authRepo := &testutil.AuthRepoMock{
		GetRefreshTokenByHashFn: func(_ context.Context, _ string) (*entity.RefreshToken, error) {
			return &entity.RefreshToken{ID: "rt-42", UserID: "user-1"}, nil
		},
		RevokeRefreshTokenFn: func(_ context.Context, id string) error {
			revokedID = id
			return nil
		},
	}

	uc := newUC(&testutil.UserRepoMock{}, authRepo, &testutil.TokenMakerMock{}, &testutil.HasherMock{})
	err := uc.Logout(ctx, &entity.LogoutReq{RefreshToken: "some-token"})

	require.NoError(t, err)
	assert.Equal(t, "rt-42", revokedID)
}

func TestLogout_UnknownToken(t *testing.T) {
	ctx := context.Background()

	authRepo := &testutil.AuthRepoMock{
		GetRefreshTokenByHashFn: func(_ context.Context, _ string) (*entity.RefreshToken, error) {
			return nil, apperr.NotFound("not found")
		},
	}

	uc := newUC(&testutil.UserRepoMock{}, authRepo, &testutil.TokenMakerMock{}, &testutil.HasherMock{})
	// Logout with unknown token should be silent (no error)
	err := uc.Logout(ctx, &entity.LogoutReq{RefreshToken: "unknown"})
	require.NoError(t, err)
}
