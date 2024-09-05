package models

import (
	"fmt"
	"qpc/rest"
)

type UserV1 struct {
	Uuid        string     `json:"uuid"`
	LoginId     string     `json:"loginId"`
	Email       string     `json:"email"`
	Name        string     `json:"name"`
	UserRoles   []UserRole `json:"userRoles"`
	LastLoginAt string     `json:"lastLoginAt"`
	Locked      bool       `json:"locked"`
	Expired     bool       `json:"expired"`
	Deleted     bool       `json:"deleted"`
	CreatedAt   string     `json:"createdAt"`
	UpdatedAt   string     `json:"updatedAt"`
}

func (u UserV1) Status() string {
	if u.Deleted {
		return "deleted"
	}
	if u.Expired {
		return "expired"
	}
	if u.Locked {
		return "locked"
	}
	return "active"
}

func (u UserV1) ShortCreatedAt() string {
	return rest.ShortDatetimeWithTZ(u.CreatedAt)
}

func (u UserV1) ShortUpdatedAt() string {
	return rest.ShortDatetimeWithTZ(u.UpdatedAt)
}

func (u UserV1) String() string {
	return fmt.Sprintf(
		"{ Uuid=%s, LoginId=%s, Email=%s, Name=%s, UserRoles=%v, "+
			"LastLoginAt=%s, Locked=%t, Expired=%t, Deleted=%t, "+
			"CreatedAt=%s, UpdatedAt=%s }",
		u.Uuid, u.LoginId, u.Email, u.Name, u.UserRoles,
		u.LastLoginAt, u.Locked, u.Expired, u.Deleted,
		u.CreatedAt, u.UpdatedAt)
}

type PagedUserV1List struct {
	PagedList[UserV1]
}
