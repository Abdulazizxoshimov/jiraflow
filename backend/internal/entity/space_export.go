package entity

import "time"

type SpaceExport struct {
	ID          string     `json:"id"`
	SpaceID     string     `json:"space_id"`
	RequestedBy string     `json:"requested_by"`
	Status      string     `json:"status"` // pending | processing | done | failed
	FileURL     *string    `json:"file_url,omitempty"`
	ErrorMsg    *string    `json:"error_msg,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}
