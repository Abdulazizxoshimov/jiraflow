package entity

import "time"

type PageLock struct {
	PageID    string    `json:"page_id"`
	UserID    string    `json:"user_id"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

type AcquireLockReq struct {
	TTLSeconds int `json:"ttl_seconds"`
}
