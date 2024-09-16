package dac_privilege

import (
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"qpc/config"
	"qpc/model"
	"qpc/utils"
)

func (p *Privilege) FindByNameOrUuid(query string, privileges *[]Privilege) {
	config.LocalDatabase.
		// Note: Column names are snake_case in the database.
		Where("name = ? OR uuid = ?", query, query).
		Find(privileges)
}

func (p *Privilege) FetchAllAndForEach(
	forEachFunc func(fetched *Privilege) bool,
) {
	utils.FetchPagedListAndForEach(
		"/api/external/v2/privileges",
		&PrivilegePagedList{},
		func(page *PrivilegePagedList) bool {
			for _, it := range page.List {
				forEachFunc(&it)
			}
			return true
		},
	)
}

func (p *Privilege) FindAllAndForEach(
	forEachFunc func(found *Privilege) bool,
) {
	utils.FindAllAndForEach(
		func(tx *gorm.DB, total *int64) *gorm.DB {
			return tx.Model(&Privilege{}).Count(total)
		},
		func(tx *gorm.DB, items *[]Privilege) *gorm.DB {
			return tx.Model(&Privilege{}).Find(items)
		},
		forEachFunc,
	)
}

func (p *Privilege) Save() *Privilege {
	// NOTE Don’t use Save with Model, it’s an Undefined Behavior.
	// https://gorm.io/docs/update.html#Save
	db := config.LocalDatabase.Save(p)
	logrus.Debugf("Saved it, RowsAffected: %d", db.RowsAffected)
	return p
}

func RunAutoMigrate() {
	db := config.LocalDatabase
	err := db.AutoMigrate(
		&model.Modifier{},
		&Privilege{},
	)
	if err != nil {
		logrus.Fatal(err)
	}
}
