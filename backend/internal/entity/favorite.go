package entity

import "time"

type Favorite struct {
	ID         string    `json:"id"`
	UserID     string    `json:"user_id"`
	EntityType string    `json:"entity_type"` // page | space
	EntityID   string    `json:"entity_id"`
	CreatedAt  time.Time `json:"created_at"`

	Page  *Page  `json:"page,omitempty"`
	Space *Space `json:"space,omitempty"`
}

type AddFavoriteReq struct {
	EntityType string `json:"entity_type" binding:"required,oneof=page space"`
	EntityID   string `json:"entity_id"   binding:"required"`
}

type FavoriteFilter struct {
	EntityType string
	Filter
}
