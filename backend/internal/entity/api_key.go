package entity

import "time"

type APIKey struct {
	ID          string     `json:"id"`
	UserID      string     `json:"user_id"`
	Name        string     `json:"name"`
	KeyPrefix   string     `json:"key_prefix"`
	Scopes      []string   `json:"scopes"`
	LastUsedAt  *time.Time `json:"last_used_at"`
	ExpiresAt   *time.Time `json:"expires_at"`
	CreatedAt   time.Time  `json:"created_at"`
	RevokedAt   *time.Time `json:"revoked_at,omitempty"`
}

type CreateAPIKeyReq struct {
	Name      string     `json:"name"       validate:"required,max=100"`
	Scopes    []string   `json:"scopes"`
	ExpiresAt *time.Time `json:"expires_at"`
}

// CreateAPIKeyResp is returned only on creation — PlainKey is shown once.
type CreateAPIKeyResp struct {
	*APIKey
	PlainKey string `json:"key"`
}
