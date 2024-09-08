package user

import "fmt"

type AdminRole struct {
	UserV2Uuid string `gorm:"primaryKey"`
	RoleUuid   string `json:"roleUuid" gorm:"primaryKey"`
	RoleName   string `json:"roleName"`
}

func (r AdminRole) String() string {
	return fmt.Sprintf(
		"{ RoleUuid=%s, RoleName=%s }",
		r.RoleUuid, r.RoleName,
	)
}
