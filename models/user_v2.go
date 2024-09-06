package models

import (
	"fmt"
	"qpc/utils"
)

type UserV2 struct {
	Uuid            string      `json:"uuid" gorm:"primaryKey"`
	LoginId         string      `json:"loginId"`
	Email           string      `json:"email"`
	Name            string      `json:"name"`
	AdminRoles      []AdminRole `json:"adminRoles" gorm:"foreignKey:UserV2Uuid"`
	Status          string      `json:"status"`
	Factor          Factor      `json:"factor" gorm:"-"`
	PasswordExpired bool        `json:"passwordExpired"`
	Locked          bool        `json:"locked"`
	Expired         bool        `json:"expired"`
	Deleted         bool        `json:"deleted"`
	CreatedAt       string      `json:"createdAt"`
	UpdatedAt       string      `json:"updatedAt"`
}

func (u UserV2) StatusMore() string {
	if u.Deleted {
		return "deleted"
	}
	if u.Expired {
		return "expired"
	}
	if u.Locked {
		return "locked"
	}
	if u.PasswordExpired {
		return "pwd-exp"
	}
	return "-"
}

func (u UserV2) ShortCreatedAt() string {
	return utils.ShortDatetimeWithTZ(u.CreatedAt)
}

func (u UserV2) ShortUpdatedAt() string {
	return utils.ShortDatetimeWithTZ(u.UpdatedAt)
}

func (u UserV2) ShortID() string {
	return fmt.Sprintf(
		"{ Uuid=%s, LoginId=%s }",
		u.Uuid, u.LoginId,
	)
}

func (u UserV2) String() string {
	return fmt.Sprintf(
		"{ Uuid=%s, LoginId=%s, Email=%s, Name=%s, AdminRoles=%v, "+
			"Status=%s, Factor=%v, PasswordExpired=%t, Locked=%t, "+
			"Expired=%t, Deleted=%t, CreatedAt=%s, UpdatedAt=%s }",
		u.Uuid, u.LoginId, u.Email, u.Name, u.AdminRoles,
		u.Status, u.Factor, u.PasswordExpired, u.Locked,
		u.Expired, u.Deleted, u.CreatedAt, u.UpdatedAt)
}

type PagedUserV2List struct {
	PagedList[UserV2]
}

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
