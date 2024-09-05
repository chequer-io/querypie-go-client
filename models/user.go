package models

import "qpc/rest"

type User struct {
	CreatedAt   string     `json:"createdAt"`
	Deleted     bool       `json:"deleted"`
	Email       string     `json:"email"`
	Expired     bool       `json:"expired"`
	LastLoginAt string     `json:"lastLoginAt"`
	Locked      bool       `json:"locked"`
	LoginId     string     `json:"loginId"`
	Name        string     `json:"name"`
	UpdatedAt   string     `json:"updatedAt"`
	UserRoles   []UserRole `json:"userRoles"`
	Uuid        string     `json:"uuid"`
}

func (u *User) Status() string {
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

func (u *User) ShortCreatedAt() string {
	return rest.ShortDatetimeWithTZ(u.CreatedAt)
}

func (u *User) ShortUpdatedAt() string {
	return rest.ShortDatetimeWithTZ(u.UpdatedAt)
}

type PagedUserList struct {
	PagedList[User]
}
