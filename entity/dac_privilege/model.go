package dac_privilege

import (
	"fmt"
	"qpc/model"
)

type Privilege struct {
	Uuid             string           `json:"uuid" gorm:"primaryKey"`
	Name             string           `json:"name"`
	PrivilegeTypes   model.StringList `json:"privilegeTypes"`
	Description      string           `json:"description"`
	CanImport        bool             `json:"canImport"`
	CanExport        bool             `json:"canExport"`
	CanCopyClipboard bool             `json:"canCopyClipboard"`
	PrivilegeVendor  string           `json:"privilegeVendor"`
	Status           string           `json:"status"`
	CreatedAt        string           `json:"createdAt"`
	CreatedByUuid    string           `json:"-"`
	CreatedBy        model.Modifier   `json:"createdBy" gorm:"foreignKey:CreatedByUuid"`
	UpdatedAt        string           `json:"updatedAt"`
	UpdatedByUuid    string           `json:"-"`
	UpdatedBy        model.Modifier   `json:"updatedBy" gorm:"foreignKey:UpdatedByUuid"`
}

func (p *Privilege) ShortID() string {
	return fmt.Sprintf(
		"{ Uuid=%s, Name=%s, PrivilegeTypes=%s }",
		p.Uuid, p.Name, p.PrivilegeTypes,
	)
}

func (p *Privilege) String() string {
	return fmt.Sprintf(
		"{ Uuid=%s, Name=%s, PrivilegeTypes=%s, Description=%s, CanImport=%t, "+
			"CanExport=%t, CanCopyClipboard=%t, PrivilegeVendor=%s, Status=%s, "+
			"CreatedAt=%s, CreatedBy=%v, UpdatedAt=%s, UpdatedBy=%v }",
		p.Uuid, p.Name, p.PrivilegeTypes, p.Description, p.CanImport, p.CanExport,
		p.CanCopyClipboard, p.PrivilegeVendor, p.Status, p.CreatedAt, p.CreatedBy,
		p.UpdatedAt, p.UpdatedBy,
	)
}

type PrivilegePagedList struct {
	List []Privilege `json:"list"`
	Page model.Page  `json:"page"`
}

func (pl *PrivilegePagedList) GetPage() model.Page {
	return pl.Page
}

func (pl *PrivilegePagedList) GetList() []Privilege {
	return pl.List
}
