package model

const (
	defaultPageLimit = 20
	maxPageLimit     = 100
)

type Page[T any] struct {
	Data   []T `json:"data"`
	Total  int `json:"total"`
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

type PageQuery struct {
	Limit  int `query:"limit"`
	Offset int `query:"offset"`
}

func (q PageQuery) Normalize() (limit, offset int) {
	limit = q.Limit
	if limit <= 0 {
		limit = defaultPageLimit
	}
	limit = min(limit, maxPageLimit)
	offset = max(q.Offset, 0)
	return limit, offset
}
