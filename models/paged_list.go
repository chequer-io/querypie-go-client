package models

import "fmt"

type Page struct {
	CurrentPage   int `json:"currentPage"`
	PageSize      int `json:"pageSize"`
	TotalElements int `json:"totalElements"`
	TotalPages    int `json:"totalPages"`
}

func (p Page) HasNext() bool {
	return p.CurrentPage+1 < p.TotalPages
}

func (p Page) NextPageQuery() string {
	if p.HasNext() {
		return fmt.Sprintf("page=%d&pageSize=%d", p.CurrentPage+1, p.PageSize)
	} else {
		return ""
	}
}

type PagedList[T any] interface {
	GetPage() Page
	GetList() []T
}
