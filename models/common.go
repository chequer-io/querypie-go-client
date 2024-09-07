package models

import "fmt"

type CloudProvider struct {
	Uuid string `json:"uuid"`
	Type string `json:"type"`
	Name string `json:"name"`
}

type Modifier struct {
	Uuid    string `json:"uuid" gorm:"primaryKey"`
	LoginId string `json:"loginId"`
	Email   string `json:"email"`
	Name    string `json:"name"`
}

func (m Modifier) String() string {
	return fmt.Sprint(m.LoginId)
}

type SummarizedZone struct {
	Uuid string `json:"uuid"`
	Name string `json:"name"`
}

type Zone struct {
	Uuid      string   `json:"uuid"`
	Name      string   `json:"name"`
	IpBands   []string `json:"ipBands"`
	CreatedAt string   `json:"createdAt"`
	UpdatedAt string   `json:"updatedAt"`
}

type Tag struct {
	Key          string `json:"key"`
	Value        string `json:"value"`
	Synchronized bool   `json:"synchronized"`
}

type Role struct {
	Uuid string `json:"uuid" gorm:"primaryKey"`
	Name string `json:"name"`
}

func (r Role) String() string {
	return fmt.Sprintf(
		"{ Uuid=%s, Name=%s }",
		r.Uuid, r.Name,
	)
}

type SecretStore struct {
	Uuid    string `json:"uuid"`
	Name    string `json:"name"`
	Account string `json:"account"`
}
