package entity

import "time"

type Space struct {
	ID         string     `json:"id"`
	Key        string     `json:"key"`
	Name       string     `json:"name"`
	Description *string   `json:"description,omitempty"`
	IconURL    *string    `json:"icon_url,omitempty"`
	Type       string     `json:"type"` // team | personal | project
	LeadID     string     `json:"lead_id"`
	ProjectID  *string    `json:"project_id,omitempty"`
	IsArchived bool       `json:"is_archived"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	DeletedAt  *time.Time `json:"-"`

	Lead *UserShort `json:"lead,omitempty"`
}

type CreateSpaceReq struct {
	Key        string  `json:"key"         validate:"required,min=2,max=10"`
	Name       string  `json:"name"        validate:"required,min=2,max=255"`
	Description *string `json:"description" validate:"omitempty,max=2000"`
	IconURL    *string `json:"icon_url"    validate:"omitempty,url"`
	Type       string  `json:"type"        validate:"required,oneof=team personal project"`
	ProjectID  *string `json:"project_id"  validate:"omitempty,uuid4"`
}

type UpdateSpaceReq struct {
	Name       *string `json:"name"        validate:"omitempty,min=2,max=255"`
	Description *string `json:"description" validate:"omitempty,max=2000"`
	IconURL    *string `json:"icon_url"    validate:"omitempty,url"`
	IsArchived *bool   `json:"is_archived"`
}

type SpaceStatistics struct {
	TotalPages     int           `json:"total_pages"`
	PublishedPages int           `json:"published_pages"`
	DraftPages     int           `json:"draft_pages"`
	TotalBlogPosts int           `json:"total_blog_posts"`
	TotalMembers   int           `json:"total_members"`
	TotalViews     int           `json:"total_views"`
	RecentActivity int           `json:"recent_activity_7d"`
	TopContributors []Contributor `json:"top_contributors"`
}

type Contributor struct {
	UserID    string `json:"user_id"`
	FullName  string `json:"full_name"`
	AvatarURL string `json:"avatar_url"`
	PageCount int    `json:"page_count"`
}
