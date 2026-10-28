package model

const (
	DefaultPage  = 1
	DefaultLimit = 10
	MaxLimit     = 100
)

type Page[T any] struct {
	Data       []T `json:"data"`
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

func NewPage[T any](data []T, page, limit, total int) *Page[T] {
	totalPages := 0
	if limit > 0 {
		totalPages = (total + limit - 1) / limit // divide and round up
	}
	return &Page[T]{Data: data, Page: page, Limit: limit, Total: total, TotalPages: totalPages}
}