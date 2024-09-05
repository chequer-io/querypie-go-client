package models

type Zone struct {
	CreatedAt string `json:"createdAt"`
	IpBands   string `json:"ipBands"`
	Name      string `json:"name"`
	UpdatedAt string `json:"updatedAt"`
	Uuid      string `json:"uuid"`
}
