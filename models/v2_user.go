package models

type V2User struct {
	AdminRoles      []AdminRole `json:"adminRoles"`
	CreatedAt       string      `json:"createdAt"`
	Deleted         bool        `json:"deleted"`
	Email           string      `json:"email"`
	Expired         bool        `json:"expired"`
	Factor          Factor      `json:"factor"`
	Locked          bool        `json:"locked"`
	LoginId         string      `json:"loginId"`
	Name            string      `json:"name"`
	PasswordExpired bool        `json:"passwordExpired"`
	Status          string      `json:"status"`
	UpdatedAt       string      `json:"updatedAt"`
	Uuid            string      `json:"uuid"`
}
