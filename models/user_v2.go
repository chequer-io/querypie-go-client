package models

import (
	"fmt"
	"qpc/rest"
)

type UserV2 struct {
	Uuid            string      `json:"uuid"`
	LoginId         string      `json:"loginId"`
	Email           string      `json:"email"`
	Name            string      `json:"name"`
	AdminRoles      []AdminRole `json:"adminRoles"`
	Status          string      `json:"status"`
	Factor          Factor      `json:"factor"`
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
	return rest.ShortDatetimeWithTZ(u.CreatedAt)
}

func (u UserV2) ShortUpdatedAt() string {
	return rest.ShortDatetimeWithTZ(u.UpdatedAt)
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
