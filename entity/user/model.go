package user

import (
	"fmt"
	"qpc/models"
	"qpc/utils"
)

type User struct {
	Uuid            string        `json:"uuid" gorm:"primaryKey"`
	LoginId         string        `json:"loginId"`
	Email           string        `json:"email"`
	Name            string        `json:"name"`
	AdminRoles      []AdminRole   `json:"adminRoles" gorm:"foreignKey:UserV2Uuid"`
	Status          string        `json:"status"`
	Factor          models.Factor `json:"factor" gorm:"-"`
	PasswordExpired bool          `json:"passwordExpired"`
	Locked          bool          `json:"locked"`
	Expired         bool          `json:"expired"`
	Deleted         bool          `json:"deleted"`
	CreatedAt       string        `json:"createdAt"`
	UpdatedAt       string        `json:"updatedAt"`
}

func (u User) StatusMore() string {
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

func (u User) ShortCreatedAt() string {
	return utils.ShortDatetimeWithTZ(u.CreatedAt)
}

func (u User) ShortUpdatedAt() string {
	return utils.ShortDatetimeWithTZ(u.UpdatedAt)
}

func (u User) ShortID() string {
	return fmt.Sprintf(
		"{ Uuid=%s, LoginId=%s }",
		u.Uuid, u.LoginId,
	)
}

func (u User) String() string {
	return fmt.Sprintf(
		"{ Uuid=%s, LoginId=%s, Email=%s, Name=%s, AdminRoles=%v, "+
			"Status=%s, Factor=%v, PasswordExpired=%t, Locked=%t, "+
			"Expired=%t, Deleted=%t, CreatedAt=%s, UpdatedAt=%s }",
		u.Uuid, u.LoginId, u.Email, u.Name, u.AdminRoles,
		u.Status, u.Factor, u.PasswordExpired, u.Locked,
		u.Expired, u.Deleted, u.CreatedAt, u.UpdatedAt)
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

type PagedUserList struct {
	models.PagedList[User]
}
