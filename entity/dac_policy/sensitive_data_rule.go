package dac_policy

import (
	"fmt"
	"qpc/model"
)

type SensitiveDataRule struct {
	Uuid       string   `json:"uuid" gorm:"primaryKey" yaml:"uuid"`
	PolicyUuid string   `json:"policyUuid" yaml:"policyUuid"`
	Policy     Policy   `json:"-" gorm:"foreignKey:PolicyUuid" yaml:"policy"`
	ObjectType string   `json:"objectType" yaml:"objectType"`
	ObjectPath []string `json:"objectPath" gorm:"json" yaml:"objectPath"`
	Level      int      `json:"level"`

	CreatedAt string `json:"createdAt" yaml:"createdAt"`
	UpdatedAt string `json:"updatedAt" yaml:"updatedAt"`

	model.WithHttpResponse `json:"-" gorm:"-" yaml:"-"`
}

func (r SensitiveDataRule) String() string {
	return fmt.Sprintf(
		"{ Uuid=%s, PolicyUuid=%s, ObjectType=%s, ObjectPath=%s, Level=%d, "+
			"CreatedAt=%s, UpdatedAt=%s }",
		r.Uuid, r.PolicyUuid, r.ObjectType, r.ObjectPath, r.Level,
		r.CreatedAt, r.UpdatedAt,
	)
}
