package dac_policy

import (
	"fmt"
	"qpc/model"
)

type Policy struct {
	Uuid             string                        `json:"uuid" gorm:"primaryKey" yaml:"uuid"`
	ClusterGroupUuid string                        `json:"clusterGroupUuid" yaml:"clusterGroupUuid"`
	Connection       SummarizedConnectionForPolicy `json:"-" gorm:"foreignKey:ClusterGroupUuid" yaml:"connection"`
	Title            string                        `json:"title" yaml:"title"`
	PolicyType       PolicyType                    `json:"policyType" yaml:"policyType"`
	NumberOfRules    int                           `json:"numberOfRules" yaml:"numberOfRules"`
	Enabled          bool                          `json:"enabled" yaml:"enabled"`
	CreatedAt        string                        `json:"createdAt" yaml:"createdAt"`
	CreatedByUuid    string                        `json:"-" yaml:"-"`
	CreatedBy        model.Modifier                `json:"createdUser" gorm:"foreignKey:CreatedByUuid" yaml:"createdBy"`
	UpdatedAt        string                        `json:"updatedAt" yaml:"updatedAt"`
	UpdatedByUuid    string                        `json:"-" yaml:"-"`
	UpdatedBy        model.Modifier                `json:"updatedUser" gorm:"foreignKey:UpdatedByUuid" yaml:"updatedBy"`

	model.WithHttpResponse `json:"-" gorm:"-" yaml:"-"`
}

func (p *Policy) ShortID() string {
	return fmt.Sprintf(
		"{ Uuid=%s, Title=%s, PolicyType=%s }",
		p.Uuid, p.Title, p.PolicyType,
	)
}

func (p *Policy) String() string {
	return fmt.Sprintf(
		"{ Uuid=%s, ClusterGroupUuid=%s, Title=%s, PolicyType=%s"+
			"NumberOfRules=%d, Enabled=%t, "+
			"CreatedAt=%s, CreatedBy=%s, UpdatedAt=%s, UpdatedBy=%s }",
		p.Uuid, p.ClusterGroupUuid, p.Title, p.PolicyType,
		p.NumberOfRules, p.Enabled,
		p.CreatedAt, p.CreatedBy, p.UpdatedAt, p.UpdatedBy,
	)
}

type PolicyType string // PolicyType is a string type alias.
const (
	UnknownPolicyType PolicyType = ""
	DataLevel         PolicyType = "DATA_LEVEL"
	DataAccess        PolicyType = "DATA_ACCESS"
	DataMasking       PolicyType = "DATA_MASKING"
	Notification      PolicyType = "NOTIFICATION"
	Ledger            PolicyType = "LEDGER"
)

func (pt PolicyType) IsValid() bool {
	if pt == DataLevel || pt == DataAccess || pt == DataMasking || pt == Notification {
		return true
	}
	return false
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
