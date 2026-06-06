package helper

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	defaultPage  = 1
	defaultLimit = 20
	maxLimit     = 100
)

type Pagination struct {
	Page  int `json:"page"  form:"page"`
	Limit int `json:"limit" form:"limit"`
	Total int `json:"total"`
}

func (p *Pagination) Offset() int {
	if p.Page < 1 {
		return 0
	}
	return (p.Page - 1) * p.Limit
}

func (p *Pagination) TotalPages() int {
	if p.Limit == 0 {
		return 0
	}
	pages := p.Total / p.Limit
	if p.Total%p.Limit != 0 {
		pages++
	}
	return pages
}

// ParsePagination reads page/limit query params from the Gin context.
func ParsePagination(c *gin.Context) Pagination {
	page := parseIntQuery(c, "page", defaultPage)
	limit := parseIntQuery(c, "limit", defaultLimit)

	if page < 1 {
		page = defaultPage
	}
	if limit < 1 {
		limit = defaultLimit
	}
	if limit > maxLimit {
		limit = maxLimit
	}

	return Pagination{Page: page, Limit: limit}
}

func parseIntQuery(c *gin.Context, key string, def int) int {
	v := c.Query(key)
	if v == "" {
		return def
	}
	n, err := strconv.Atoi(v)
	if err != nil || n < 0 {
		return def
	}
	return n
}
