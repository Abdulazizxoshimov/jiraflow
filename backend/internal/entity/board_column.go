package entity

import "time"

type BoardColumn struct {
	ID       string  `json:"id"`
	BoardID  string  `json:"board_id"`
	Name     string  `json:"name"`
	Position int     `json:"position"`
	WIPLimit *int    `json:"wip_limit,omitempty"`
	CreatedAt time.Time `json:"created_at"`

	StatusIDs []string `json:"status_ids,omitempty"`
}

type CreateBoardColumnReq struct {
	Name      string   `json:"name"      validate:"required,min=1,max=100"`
	Position  int      `json:"position"`
	WIPLimit  *int     `json:"wip_limit" validate:"omitempty,gt=0"`
	StatusIDs []string `json:"status_ids"`
}

type UpdateBoardColumnReq struct {
	Name      *string  `json:"name"      validate:"omitempty,min=1,max=100"`
	Position  *int     `json:"position"`
	WIPLimit  *int     `json:"wip_limit" validate:"omitempty,gt=0"`
	StatusIDs []string `json:"status_ids"`
}
