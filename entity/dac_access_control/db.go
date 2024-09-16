package dac_access_control

import (
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
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

func (sac *SummarizedAccessControl) FindAllAndForEach(
	forEachFunc func(found *SummarizedAccessControl) bool,
) {
	utils.FindAllAndForEach(
		func(db *gorm.DB, total *int64) *gorm.DB {
			return db.Model(&SummarizedAccessControl{}).Count(total)
		},
		func(db *gorm.DB, items *[]SummarizedAccessControl) *gorm.DB {
			return db.Model(&SummarizedAccessControl{}).Find(items)
		},
		forEachFunc,
	)
}

func (sac *SummarizedAccessControl) Save() *SummarizedAccessControl {
	// NOTE Don’t use Save with Model, it’s an Undefined Behavior.
	// https://gorm.io/docs/update.html#Save
	db := config.LocalDatabase.Save(sac)
	logrus.Debugf("Saved it, RowsAffected: %d", db.RowsAffected)
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
