package models

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

type PagedUserList struct {
	PagedList[User]
}
