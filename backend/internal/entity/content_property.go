package entity

import "time"

type ContentProperty struct {
	ID         string         `json:"id"`
	EntityType string         `json:"entity_type"` // page | issue | space
	EntityID   string         `json:"entity_id"`
	Key        string         `json:"key"`
	Value      map[string]any `json:"value"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
}

type SetContentPropertyReq struct {
	Value map[string]any `json:"value" validate:"required"`
}
