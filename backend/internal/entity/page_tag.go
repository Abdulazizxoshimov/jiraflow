package entity

import "time"

type PageTag struct {
	ID        string    `json:"id"`
	SpaceID   string    `json:"space_id"`
	Name      string    `json:"name"`
	Color     string    `json:"color"`
	CreatedAt time.Time `json:"created_at"`
}

type CreatePageTagReq struct {
	Name  string `json:"name"  binding:"required,min=1,max=64"`
	Color string `json:"color"`
}

type UpdatePageTagReq struct {
	Name  *string `json:"name"`
	Color *string `json:"color"`
}

type SetPageTagsReq struct {
	TagIDs []string `json:"tag_ids" binding:"required"`
}
