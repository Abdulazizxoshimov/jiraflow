package entity

import "time"

type BlogPost struct {
	ID          string     `json:"id"`
	SpaceID     string     `json:"space_id"`
	Title       string     `json:"title"`
	Body        string     `json:"body"`
	AuthorID    string     `json:"author_id"`
	IsPublished bool       `json:"is_published"`
	PublishedAt *time.Time `json:"published_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`

	Author *UserShort `json:"author,omitempty"`
}

type CreateBlogPostReq struct {
	Title string `json:"title" validate:"required,min=1,max=500"`
	Body  string `json:"body"`
}

type UpdateBlogPostReq struct {
	Title *string `json:"title" validate:"omitempty,min=1,max=500"`
	Body  *string `json:"body"`
}

type ListBlogPostsFilter struct {
	SpaceID       string
	AuthorID      string
	OnlyPublished bool
	Limit         int
	Offset        int
}
