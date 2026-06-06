package entity

import "time"

type SavedFilter struct {
	ID          string         `json:"id"`
	UserID      string         `json:"user_id"`
	Name        string         `json:"name"`
	Description *string        `json:"description,omitempty"`
	FilterType  string         `json:"filter_type"` // issue | page
	Filters     map[string]any `json:"filters"`
	IsShared    bool           `json:"is_shared"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

type CreateSavedFilterReq struct {
	Name        string         `json:"name"        validate:"required,min=1,max=255"`
	Description *string        `json:"description"`
	FilterType  string         `json:"filter_type" validate:"required,oneof=issue page"`
	Filters     map[string]any `json:"filters"     validate:"required"`
	IsShared    bool           `json:"is_shared"`
}

type UpdateSavedFilterReq struct {
	Name        *string        `json:"name"        validate:"omitempty,min=1,max=255"`
	Description *string        `json:"description"`
	Filters     map[string]any `json:"filters"`
	IsShared    *bool          `json:"is_shared"`
}
