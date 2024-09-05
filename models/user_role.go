package models

type UserRole struct {
	ObjectType string `json:"objectType"`
	ObjectUuid string `json:"objectUuid"`
	Role       struct {
		Name string `json:"name"`
		Uuid string `json:"uuid"`
	} `json:"role"`
	User struct {
		Active       bool   `json:"active"`
		AuthProvider string `json:"authProvider"`
		Email        string `json:"email"`
		Id           int    `json:"id"`
		LoginId      string `json:"loginId"`
		Name         string `json:"name"`
		Type         string `json:"type"`
		Uuid         string `json:"uuid"`
	} `json:"user"`
	Uuid string `json:"uuid"`
}
