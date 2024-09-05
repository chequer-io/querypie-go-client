package models

import "fmt"

type Factor struct {
	Type    string `json:"type"`
	Enabled bool   `json:"enabled"`
}

func (f Factor) String() string {
	return fmt.Sprintf(
		"{ Type=%s, Enabled=%t }",
		f.Type, f.Enabled,
	)
}
