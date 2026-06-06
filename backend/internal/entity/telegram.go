package entity

import "time"

type TelegramConnection struct {
	ID               string     `json:"id"`
	UserID           string     `json:"user_id"`
	TelegramID       *int64     `json:"telegram_id,omitempty"`
	ChatID           *int64     `json:"chat_id,omitempty"`
	Username         *string    `json:"username,omitempty"`
	VerificationCode *string    `json:"-"`
	VerifiedAt       *time.Time `json:"verified_at,omitempty"`
	CreatedAt        time.Time  `json:"created_at"`
}

type TelegramUpdate struct {
	UpdateID int64            `json:"update_id"`
	Message  *TelegramMessage `json:"message,omitempty"`
}

type TelegramMessage struct {
	MessageID int64        `json:"message_id"`
	From      TelegramUser `json:"from"`
	Chat      TelegramChat `json:"chat"`
	Text      string       `json:"text"`
	Date      int64        `json:"date"`
}

type TelegramUser struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type TelegramChat struct {
	ID   int64  `json:"id"`
	Type string `json:"type"`
}

type TelegramStatusResp struct {
	Connected  bool       `json:"connected"`
	Username   *string    `json:"username,omitempty"`
	VerifiedAt *time.Time `json:"verified_at,omitempty"`
}
