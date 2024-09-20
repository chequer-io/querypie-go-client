package dac_policy

import (
	"fmt"
	"qpc/model"
)

type Policy struct {
	Uuid             string                        `json:"uuid" gorm:"primaryKey" yaml:"uuid"`
	ClusterGroupUuid string                        `json:"clusterGroupUuid" yaml:"-"`
	Connection       SummarizedConnectionForPolicy `json:"-" gorm:"foreignKey:ClusterGroupUuid" yaml:"connection"`
	Name             string                        `json:"title" yaml:"name"`
	PolicyType       string                        `json:"policyType" yaml:"policyType"`
	NumberOfRules    int                           `json:"numberOfRules" yaml:"numberOfRules"`
	Enabled          bool                          `json:"enabled" yaml:"enabled"`
	CreatedAt        string                        `json:"createdAt" yaml:"createdAt"`
	CreatedByUuid    string                        `json:"-" yaml:"-"`
	CreatedBy        model.Modifier                `json:"createdUser" gorm:"foreignKey:CreatedByUuid" yaml:"createdBy"`
	UpdatedAt        string                        `json:"updatedAt" yaml:"updatedAt"`
	UpdatedByUuid    string                        `json:"-" yaml:"-"`
	UpdatedBy        model.Modifier                `json:"updatedUser" gorm:"foreignKey:UpdatedByUuid" yaml:"updatedBy"`
}

func (p *Policy) ShortID() string {
	return fmt.Sprintf(
		"{ Uuid=%s, Name=%s, NumberOfRules=%d }",
		p.Uuid, p.Name, p.NumberOfRules,
	)
}

func (p *Policy) String() string {
	return fmt.Sprintf(
		"{ Uuid=%s, ClusterGroupUuid=%s, Name=%s, PolicyType=%s"+
			"NumberOfRules=%d, Enabled=%t, "+
			"CreatedAt=%s, CreatedBy=%s, UpdatedAt=%s, UpdatedBy=%s }",
		p.Uuid, p.ClusterGroupUuid, p.ClusterGroupUuid, p.PolicyType,
		p.NumberOfRules, p.Enabled,
		p.CreatedAt, p.CreatedBy, p.UpdatedAt, p.UpdatedBy,
	)
}

type PolicyPagedList struct {
	List []Policy   `json:"list"`
	Page model.Page `json:"page"`
}

func (pl *PolicyPagedList) GetPage() model.Page {
	return pl.Page
}

func (pl *PolicyPagedList) GetList() []Policy {
	return pl.List
}

type SummarizedConnectionForPolicy struct {
	Uuid         string `json:"uuid" gorm:"primaryKey" yaml:"uuid"`
	DatabaseType string `json:"databaseType" yaml:"databaseType"`
	Name         string `json:"name" yaml:"name"`
}

func (s SummarizedConnectionForPolicy) String() string {
	return fmt.Sprintf(
		"{ Uuid=%s, DatabaseType=%s, Name=%s }",
		s.Uuid, s.DatabaseType, s.Name,
	)
}

func (s SummarizedConnectionForPolicy) TableName() string {
	return "connection_v2"
}
