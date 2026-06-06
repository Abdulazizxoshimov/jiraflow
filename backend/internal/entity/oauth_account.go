package entity

import "time"

type OAuthAccount struct {
	ID             string     `json:"id"`
	UserID         string     `json:"user_id"`
	Provider       string     `json:"provider"`
	ProviderUserID string     `json:"provider_user_id"`
	Email          string     `json:"email"`
	Name           string     `json:"name"`
	AvatarURL      string     `json:"avatar_url"`
	RefreshToken   *string    `json:"refresh_token,omitempty"`
	TokenExpiry    *time.Time `json:"token_expiry,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

type OAuthState struct {
	State       string    `json:"state"`
	RedirectURL string    `json:"redirect_url"`
	CreatedAt   time.Time `json:"created_at"`
	ExpiresAt   time.Time `json:"expires_at"`
}

// GoogleUserInfo is the response from Google's userinfo endpoint.
type GoogleUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
}

type OAuthCallbackResp struct {
	Tokens    *TokenPair `json:"tokens"`
	IsNewUser bool       `json:"is_new_user"`
}
