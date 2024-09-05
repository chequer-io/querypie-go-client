package models

import "fmt"

type UserRole struct {
	Uuid string `json:"uuid"`
	User struct {
		Uuid         string `json:"uuid"`
		Id           int    `json:"id"`
		LoginId      string `json:"loginId"`
		Email        string `json:"email"`
		Name         string `json:"name"`
		Type         string `json:"type"`
		Active       bool   `json:"active"`
		AuthProvider string `json:"authProvider"`
	} `json:"user"`
	Role struct {
		Uuid string `json:"uuid"`
		Name string `json:"name"`
	} `json:"role"`
	ObjectUuid string `json:"objectUuid"`
	ObjectType string `json:"objectType"`
}

func (r UserRole) String() string {
	return fmt.Sprintf(
		"{ Uuid=%s, "+
			"User={ Uuid=%s, Id=%d, LoginId=%s, Email=%s, "+
			"Name=%s, Type=%s, Active=%t, AuthProvider=%s }, "+
			"Role={ Uuid=%s, Name=%s }, "+
			"ObjectUuid=%s, ObjectType=%s }",
		r.Uuid,
		r.User.Uuid, r.User.Id, r.User.LoginId, r.User.Email,
		r.User.Name, r.User.Type, r.User.Active, r.User.AuthProvider,
		r.Role.Uuid, r.Role.Name,
		r.ObjectUuid, r.ObjectType,
	)
}
