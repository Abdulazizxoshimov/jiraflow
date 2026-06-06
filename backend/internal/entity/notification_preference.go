package entity

import "time"

type NotificationPreference struct {
	UserID         string    `json:"user_id"`
	EmailAssigned  bool      `json:"email_assigned"`
	EmailMentioned bool      `json:"email_mentioned"`
	EmailCommented bool      `json:"email_commented"`
	EmailStatus    bool      `json:"email_status"`
	EmailWatcher   bool      `json:"email_watcher"`
	DailyDigest    bool      `json:"daily_digest"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type UpdateNotificationPreferenceReq struct {
	EmailAssigned  *bool `json:"email_assigned"`
	EmailMentioned *bool `json:"email_mentioned"`
	EmailCommented *bool `json:"email_commented"`
	EmailStatus    *bool `json:"email_status"`
	EmailWatcher   *bool `json:"email_watcher"`
	DailyDigest    *bool `json:"daily_digest"`
}
