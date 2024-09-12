package model

import (
	"fmt"
	"qpc/utils"
)

type UserV1 struct {
	Uuid        string     `json:"uuid" gorm:"primaryKey"`
	LoginId     string     `json:"loginId"`
	Email       string     `json:"email"`
	Name        string     `json:"name"`
	UserRoles   []UserRole `json:"userRoles" gorm:"foreignKey:UserV1Uuid"`
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
	return utils.ShortDatetimeWithTZ(u.CreatedAt)
}

func (u UserV1) ShortUpdatedAt() string {
	return utils.ShortDatetimeWithTZ(u.UpdatedAt)
}

func (u UserV1) ShortID() string {
	return fmt.Sprintf(
		"{ Uuid=%s, LoginId=%s }",
		u.Uuid, u.LoginId,
	)
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
	List []UserV1 `json:"list"`
	Page Page     `json:"page"`
}

func (pul PagedUserV1List) GetPage() Page {
	return pul.Page
}

func (pul PagedUserV1List) GetList() []UserV1 {
	return pul.List
}

type UserRole struct {
	Uuid       string `json:"uuid" gorm:"primaryKey"`
	UserV1Uuid string `gorm:"index"`
	RoleUuid   string `json:"-"`
	Role       Role   `json:"role" gorm:"foreignKey:RoleUuid"`
	ObjectUuid string `json:"objectUuid"`
	ObjectType string `json:"objectType"`
}

func (r UserRole) String() string {
	return fmt.Sprintf(
		"{ Uuid=%s, Role=%s, "+
			"ObjectUuid=%s, ObjectType=%s }",
		r.Uuid, r.Role,
		r.ObjectUuid, r.ObjectType,
	)
}
