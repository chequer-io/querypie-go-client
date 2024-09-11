package dac_connection

import (
	"github.com/sirupsen/logrus"
	"qpc/config"
)

func (cl *PagedConnectionV2List) Save() {
	for _, sac := range cl.GetList() {
		sac.Save()
	}
}

func (sc *SummarizedConnectionV2) Save() {

	// Attempt to update the user
	result := config.LocalDatabase.Model(&SummarizedConnectionV2{}).Where("uuid = ?", sc.Uuid).Updates(&sc)

	// If no rows were affected, create a new user
	if result.RowsAffected == 0 {
		if err := config.LocalDatabase.Create(&sc).Error; err != nil {
			logrus.Errorf("Failed to save connection %s: %v", sc.ShortID(), err)
		}
	} else if result.Error != nil {
		logrus.Errorf("Failed to update connection %s: %v", sc.ShortID(), result.Error)
	}
}

func RunAutoMigrate() {
	db := config.LocalDatabase
	err := db.AutoMigrate(
		&SummarizedConnectionV2{},
	)
	if err != nil {
		logrus.Fatal(err)
	}
}
