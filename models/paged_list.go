package models

type Page struct {
	CurrentPage   int `json:"currentPage"`
	PageSize      int `json:"pageSize"`
	TotalElements int `json:"totalElements"`
	TotalPages    int `json:"totalPages"`
}

type PagedList[T any] struct {
	List []T  `json:"list"`
	Page Page `json:"page"`
}
