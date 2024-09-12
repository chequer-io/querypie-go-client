package dac_privilege

import (
	"github.com/sirupsen/logrus"
	"qpc/config"
	"qpc/utils"
)

func (pl *PrivilegePagedList) Save() {
	for _, p := range pl.GetList() {
		p.Save()
	}
}

func (pl *PrivilegePagedList) FetchAllAndPrintAndSave() {
	utils.FetchPrintAndSave(
		"/api/external/v2/privileges",
		&PrivilegePagedList{},
		func(result *PrivilegePagedList, first bool, last bool) {
			result.Print()
		},
		func(result *PrivilegePagedList) bool {
			result.Save()
			return !result.Page.HasNext()
		},
	)
}

func (p *Privilege) Save() {
	p.CreatedByUuid = p.CreatedBy.Uuid
	p.UpdatedByUuid = p.UpdatedBy.Uuid

	// Attempt to update the privilege
	result := config.LocalDatabase.Model(&Privilege{}).Where("uuid = ?", p.Uuid).Updates(map[string]interface{}{
		"Name":             p.Name,
		"PrivilegeTypes":   p.PrivilegeTypes,
		"Description":      p.Description,
		"CanImport":        p.CanImport,
		"CanExport":        p.CanExport,
		"CanCopyClipboard": p.CanCopyClipboard,
		"PrivilegeVendor":  p.PrivilegeVendor,
		"Status":           p.Status,
		"CreatedAt":        p.CreatedAt,
		"CreatedByUuid":    p.CreatedByUuid,
		"UpdatedAt":        p.UpdatedAt,
		"UpdatedByUuid":    p.UpdatedByUuid,
	})

	// If no rows were affected, create a new privilege
	if result.RowsAffected == 0 {
		if err := config.LocalDatabase.Create(&p).Error; err != nil {
			logrus.Errorf("Failed to save privilege %s: %v", p.ShortID(), err)
		}
	} else if result.Error != nil {
		logrus.Errorf("Failed to update privilege %s: %v", p.ShortID(), result.Error)
	}
}

func RunAutoMigrate() {
	db := config.LocalDatabase
	err := db.AutoMigrate(
		&Privilege{},
	)
	if err != nil {
		logrus.Fatal(err)
	}
}
