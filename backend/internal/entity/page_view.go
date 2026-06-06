package entity

import "time"

type PageView struct {
	ID        string     `json:"id"`
	PageID    string     `json:"page_id"`
	UserID    *string    `json:"user_id,omitempty"`
	IPAddress *string    `json:"ip_address,omitempty"`
	ViewedAt  time.Time  `json:"viewed_at"`
}

type PageAnalytics struct {
	PageID        string `json:"page_id"`
	TotalViews    int    `json:"total_views"`
	UniqueVisitors int   `json:"unique_visitors"`
	ViewsToday    int    `json:"views_today"`
	ViewsThisWeek int    `json:"views_this_week"`
}

type RecentPage struct {
	PageID    string    `json:"page_id"`
	Title     string    `json:"title"`
	SpaceID   string    `json:"space_id"`
	SpaceName string    `json:"space_name"`
	ViewedAt  time.Time `json:"viewed_at"`
}
