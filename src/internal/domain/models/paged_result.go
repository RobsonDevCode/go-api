package domain

type PagedResult[T any] struct {
	Items T   `json:"items"`
	Page  int `json:"page"`
	Size  int `json:"size"`
	Total int `json:"total"`
}
