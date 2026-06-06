package entity

import "time"

type SpaceCategory struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Color     *string   `json:"color,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateSpaceCategoryReq struct {
	Name  string  `json:"name"  validate:"required,min=1,max=100"`
	Color *string `json:"color" validate:"omitempty,max=20"`
}

type UpdateSpaceCategoryReq struct {
	Name  string  `json:"name"  validate:"omitempty,min=1,max=100"`
	Color *string `json:"color" validate:"omitempty,max=20"`
}
