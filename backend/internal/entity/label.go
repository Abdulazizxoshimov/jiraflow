package entity

import "time"

type Label struct {
	ID        string    `json:"id"`
	ProjectID string    `json:"project_id"`
	Name      string    `json:"name"`
	Color     string    `json:"color"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateLabelReq struct {
	Name  string `json:"name"  validate:"required,min=1,max=64"`
	Color string `json:"color" validate:"omitempty,len=7"`
}

type UpdateLabelReq struct {
	Name  *string `json:"name"  validate:"omitempty,min=1,max=64"`
	Color *string `json:"color" validate:"omitempty,len=7"`
}
