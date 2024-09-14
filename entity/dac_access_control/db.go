package dac_access_control

import (
	"github.com/sirupsen/logrus"
	"qpc/config"
	"qpc/utils"
)

func (sac *SummarizedAccessControl) FetchAllAndForEach(
	forEachFunc func(fetched *SummarizedAccessControl) bool,
) {
	utils.FetchPagedListAndForEach(
		"/api/external/v2/dac/access-controls",
		&SummarizedAccessControlPagedList{},
		func(page *SummarizedAccessControlPagedList) bool {
			for _, it := range page.List {
				forEachFunc(&it)
			}
			return true
		},
	)
}

func (sac *SummarizedAccessControl) Save() *SummarizedAccessControl {
	sac.PopulateMemberStr()

	// Attempt to update the user
	result := config.LocalDatabase.Model(&SummarizedAccessControl{}).Where("uuid = ?", sac.Uuid).Updates(&sac)

	// If no rows were affected, create a new user
	if result.RowsAffected == 0 {
		if err := config.LocalDatabase.Create(&sac).Error; err != nil {
			logrus.Fatalf("Failed to save access control %s: %v", sac.ShortID(), err)
		}
	} else if result.Error != nil {
		logrus.Fatalf("Failed to update access control %s: %v", sac.ShortID(), result.Error)
	}
	return sac
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
