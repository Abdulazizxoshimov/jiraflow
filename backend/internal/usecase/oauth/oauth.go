package oauth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/jira-backend/jiraflow-backend/internal/entity"
	"github.com/jira-backend/jiraflow-backend/internal/infrastructura/repository"
	apperr "github.com/jira-backend/jiraflow-backend/internal/pkg/errors"
	"github.com/jira-backend/jiraflow-backend/internal/pkg/token"
)

const (
	googleAuthURL    = "https://accounts.google.com/o/oauth2/v2/auth"
	googleTokenURL   = "https://oauth2.googleapis.com/token"
	googleUserURL    = "https://www.googleapis.com/oauth2/v2/userinfo"
	googleProvider   = "google"
)

type useCase struct {
	oauthRepo  repository.OAuthRepository
	userRepo   repository.UserRepository
	authRepo   repository.AuthRepository
	tokens     token.Maker
	clientID   string
	clientSecret string
	redirectURL  string
	httpClient *http.Client
}

func New(
	oauthRepo repository.OAuthRepository,
	userRepo repository.UserRepository,
	authRepo repository.AuthRepository,
	tokens token.Maker,
	clientID, clientSecret, redirectURL string,
) UseCase {
	return &useCase{
		oauthRepo:    oauthRepo,
		userRepo:     userRepo,
		authRepo:     authRepo,
		tokens:       tokens,
		clientID:     clientID,
		clientSecret: clientSecret,
		redirectURL:  redirectURL,
		httpClient:   &http.Client{Timeout: 10 * time.Second},
	}
}

func (uc *useCase) GenerateAuthURL(ctx context.Context, redirectURL string) (string, error) {
	stateBytes := make([]byte, 16)
	if _, err := rand.Read(stateBytes); err != nil {
		return "", fmt.Errorf("oauth.GenerateAuthURL rand: %w", err)
	}
	state := hex.EncodeToString(stateBytes)

	now := time.Now().UTC()
	oauthState := &entity.OAuthState{
		State:       state,
		RedirectURL: redirectURL,
		CreatedAt:   now,
		ExpiresAt:   now.Add(10 * time.Minute),
	}
	if err := uc.oauthRepo.SaveState(ctx, oauthState); err != nil {
		return "", fmt.Errorf("oauth.GenerateAuthURL save state: %w", err)
	}

	params := url.Values{
		"client_id":     {uc.clientID},
		"redirect_uri":  {uc.redirectURL},
		"response_type": {"code"},
		"scope":         {"openid email profile"},
		"state":         {state},
		"access_type":   {"offline"},
	}
	return googleAuthURL + "?" + params.Encode(), nil
}

func (uc *useCase) HandleCallback(ctx context.Context, state, code string) (*entity.OAuthCallbackResp, error) {
	if _, err := uc.oauthRepo.GetState(ctx, state); err != nil {
		return nil, apperr.BadRequest("invalid or expired oauth state")
	}
	_ = uc.oauthRepo.DeleteState(ctx, state)

	result, err := uc.exchangeAndGetUserInfo(code)
	if err != nil {
		return nil, err
	}
	ui := result.UserInfo

	isNewUser := false
	var userID string

	acc, err := uc.oauthRepo.GetAccountByProvider(ctx, googleProvider, ui.ID)
	if err != nil {
		// New user — register them
		existingUser, _ := uc.userRepo.GetByEmail(ctx, ui.Email)
		if existingUser != nil {
			userID = existingUser.ID
		} else {
			pic := ui.Picture
			newUser := &entity.User{
				ID:           uuid.NewString(),
				Email:        ui.Email,
				FullName:     ui.Name,
				PasswordHash: "",
				Role:         "member",
				IsActive:     true,
				AvatarURL:    &pic,
				CreatedAt:    time.Now().UTC(),
				UpdatedAt:    time.Now().UTC(),
			}
			if err := uc.userRepo.Create(ctx, newUser); err != nil {
				return nil, fmt.Errorf("oauth.HandleCallback create user: %w", err)
			}
			userID = newUser.ID
			isNewUser = true
		}

		now := time.Now().UTC()
		newAcc := &entity.OAuthAccount{
			ID:             uuid.NewString(),
			UserID:         userID,
			Provider:       googleProvider,
			ProviderUserID: ui.ID,
			Email:          ui.Email,
			Name:           ui.Name,
			AvatarURL:      ui.Picture,
			RefreshToken:   result.RefreshToken,
			TokenExpiry:    result.TokenExpiry,
			CreatedAt:      now,
			UpdatedAt:      now,
		}
		if err := uc.oauthRepo.UpsertAccount(ctx, newAcc); err != nil {
			return nil, fmt.Errorf("oauth.HandleCallback upsert account: %w", err)
		}
	} else {
		userID = acc.UserID
		now := time.Now().UTC()
		acc.Email = ui.Email
		acc.Name = ui.Name
		acc.AvatarURL = ui.Picture
		if result.RefreshToken != nil {
			acc.RefreshToken = result.RefreshToken
		}
		acc.TokenExpiry = result.TokenExpiry
		acc.UpdatedAt = now
		_ = uc.oauthRepo.UpsertAccount(ctx, acc)
	}

	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("oauth.HandleCallback get user: %w", err)
	}
	if !user.IsActive {
		return nil, apperr.Forbidden("account is deactivated")
	}

	sessionID := uuid.NewString()
	access, refresh, err := uc.tokens.Generate(ctx, user.ID, sessionID, user.Role)
	if err != nil {
		return nil, fmt.Errorf("oauth.HandleCallback generate tokens: %w", err)
	}

	return &entity.OAuthCallbackResp{
		Tokens:    &entity.TokenPair{AccessToken: access, RefreshToken: refresh},
		IsNewUser: isNewUser,
	}, nil
}

func (uc *useCase) LinkAccount(ctx context.Context, userID, state, code string) error {
	if _, err := uc.oauthRepo.GetState(ctx, state); err != nil {
		return apperr.BadRequest("invalid or expired oauth state")
	}
	_ = uc.oauthRepo.DeleteState(ctx, state)

	result, err := uc.exchangeAndGetUserInfo(code)
	if err != nil {
		return err
	}
	ui := result.UserInfo

	now := time.Now().UTC()
	acc := &entity.OAuthAccount{
		ID:             uuid.NewString(),
		UserID:         userID,
		Provider:       googleProvider,
		ProviderUserID: ui.ID,
		Email:          ui.Email,
		Name:           ui.Name,
		AvatarURL:      ui.Picture,
		RefreshToken:   result.RefreshToken,
		TokenExpiry:    result.TokenExpiry,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
	return uc.oauthRepo.UpsertAccount(ctx, acc)
}

func (uc *useCase) UnlinkAccount(ctx context.Context, userID, provider string) error {
	return uc.oauthRepo.DeleteAccount(ctx, userID, provider)
}

func (uc *useCase) ListLinkedAccounts(ctx context.Context, userID string) ([]*entity.OAuthAccount, error) {
	return uc.oauthRepo.ListByUser(ctx, userID)
}

type googleTokenResp struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"` // seconds
}

type exchangeResult struct {
	UserInfo     *entity.GoogleUserInfo
	RefreshToken *string
	TokenExpiry  *time.Time
}

// exchangeAndGetUserInfo exchanges the authorization code for tokens and fetches the user profile.
func (uc *useCase) exchangeAndGetUserInfo(code string) (*exchangeResult, error) {
	body := url.Values{
		"code":          {code},
		"client_id":     {uc.clientID},
		"client_secret": {uc.clientSecret},
		"redirect_uri":  {uc.redirectURL},
		"grant_type":    {"authorization_code"},
	}
	resp, err := uc.httpClient.Post(googleTokenURL, "application/x-www-form-urlencoded", strings.NewReader(body.Encode()))
	if err != nil {
		return nil, fmt.Errorf("oauth: token exchange request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		raw, _ := io.ReadAll(resp.Body)
		return nil, apperr.BadRequest(fmt.Sprintf("google token exchange failed: %s", string(raw)))
	}

	var tokenResp googleTokenResp
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, fmt.Errorf("oauth: decode token response: %w", err)
	}

	req, _ := http.NewRequest(http.MethodGet, googleUserURL, nil)
	req.Header.Set("Authorization", "Bearer "+tokenResp.AccessToken)
	infoResp, err := uc.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("oauth: user info request: %w", err)
	}
	defer infoResp.Body.Close()

	var userInfo entity.GoogleUserInfo
	if err := json.NewDecoder(infoResp.Body).Decode(&userInfo); err != nil {
		return nil, fmt.Errorf("oauth: decode user info: %w", err)
	}
	if userInfo.ID == "" {
		return nil, apperr.BadRequest("google returned empty user info")
	}

	result := &exchangeResult{UserInfo: &userInfo}
	if tokenResp.RefreshToken != "" {
		result.RefreshToken = &tokenResp.RefreshToken
	}
	if tokenResp.ExpiresIn > 0 {
		expiry := time.Now().UTC().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)
		result.TokenExpiry = &expiry
	}
	return result, nil
}
