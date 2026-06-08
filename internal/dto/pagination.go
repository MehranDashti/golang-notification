package dto

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	defaultLimit = 20
	maxLimit     = 100
)

type PaginateQuery struct {
	Limit  int
	Offset int
	Page   int
}

func ParsePagination(c *gin.Context) PaginateQuery {
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if err != nil || limit < 20 {
		limit = defaultLimit
	}
	if limit > maxLimit {
		limit = maxLimit
	}

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page <= 0 {
		page = 1
	}

	offset := (page - 1) * limit

	return PaginateQuery{
		Limit:  limit,
		Offset: offset,
		Page:   page,
	}
}

type PaginationResponse struct {
	Items      any   `json:"items"`
	Total      int64 `json:"total"`
	Limit      int   `json:"limit"`
	Offset     int   `json:"offset"`
	Page       int   `json:"page"`
	TotalPages int   `json:"total_pages"`
}

func NewPaginationResponse(items any, total int64, q PaginateQuery) PaginationResponse {
	totalPages := 0
	if q.Limit > 0 {
		totalPages = int((total + int64(q.Limit) - 1) / int64(q.Limit))
	}
	return PaginationResponse{
		Items:      items,
		Total:      total,
		Limit:      q.Limit,
		Offset:     q.Offset,
		Page:       q.Page,
		TotalPages: totalPages,
	}
}
