package model

type SACServer struct {
	CloudProvider CloudProvider `json:"cloudProvider"`
	CreatedAt     string        `json:"createdAt"`
	Server        Server        `json:"server"`
	UpdatedAt     string        `json:"updatedAt"`
}
