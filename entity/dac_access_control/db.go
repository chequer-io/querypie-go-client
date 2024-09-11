package dac_access_control

import (
	"github.com/sirupsen/logrus"
	"qpc/config"
)

func (acl *SummarizedAccessControlPagedList) Save() {
	for _, sac := range acl.GetList() {
		sac.Save()
	}
}

func (sac *SummarizedAccessControl) Save() {
	sac.PopulateMemberStr()

	// Attempt to update the user
	result := config.LocalDatabase.Model(&SummarizedAccessControl{}).Where("uuid = ?", sac.Uuid).Updates(&sac)

	// If no rows were affected, create a new user
	if result.RowsAffected == 0 {
		if err := config.LocalDatabase.Create(&sac).Error; err != nil {
			logrus.Errorf("Failed to save access control %s: %v", sac.ShortID(), err)
		}
	} else if result.Error != nil {
		logrus.Errorf("Failed to update access control %s: %v", sac.ShortID(), result.Error)
	}
}

func RunAutoMigrate() {
	db := config.LocalDatabase
	err := db.AutoMigrate(
		&SummarizedAccessControl{},
	)
	if err != nil {
		logrus.Fatal(err)
	}
}
