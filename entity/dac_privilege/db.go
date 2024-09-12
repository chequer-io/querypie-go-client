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
	p.PrivilegeTypesStr = utils.JsonFromStringArray(p.PrivilegeTypes)
	p.CreatedByUuid = p.CreatedBy.Uuid
	p.UpdatedByUuid = p.UpdatedBy.Uuid

	// Attempt to update the user
	result := config.LocalDatabase.Model(&Privilege{}).Where("uuid = ?", p.Uuid).Updates(&p)

	// If no rows were affected, create a new user
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
