package models

import "fmt"

type AdminRole struct {
	RoleUuid string `json:"roleUuid"`
	RoleName string `json:"roleName"`
}

func (r AdminRole) String() string {
	return fmt.Sprintf(
		"{ RoleUuid=%s, RoleName=%s }",
		r.RoleUuid, r.RoleName,
	)
}
