package dac_policy

import (
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"qpc/config"
	"qpc/model"
	"qpc/utils"
)

func (p *Policy) FindByNameOrUuid(query string, like bool, policies *[]Policy) {
	where := "name = ? OR uuid = ?"
	if like {
		where = "name LIKE ? OR uuid LIKE ?"
	}
	utils.FindMultiple(policies, func(db *gorm.DB) *gorm.DB {
		return db.
			Model(&Policy{}).
			// Note: Column names are snake_case in the database.
			Where(where, query, query).
			Preload("UpdatedBy").
			Preload("CreatedBy").
			Preload("Connection").
			Find(policies)
	})
}

func (p *Policy) FetchAllAndForEach(
	forEachFunc func(fetched *Policy) bool,
) {
	utils.FetchPagedListAndForEach(
		"/api/external/policies",
		&PolicyPagedList{},
		func(page *PolicyPagedList) bool {
			for _, it := range page.List {
				forEachFunc(&it)
			}
			return true
		},
	)
}

func (p *Policy) FindAllAndForEach(
	forEachFunc func(found *Policy) bool,
) {
	utils.FindAllAndForEach(
		func(tx *gorm.DB, total *int64) *gorm.DB {
			return tx.Model(&Policy{}).Count(total)
		},
		func(tx *gorm.DB, items *[]Policy) *gorm.DB {
			return tx.Model(&Policy{}).
				Preload("UpdatedBy").
				Preload("CreatedBy").
				Preload("Connection").
				Find(items)
		},
		forEachFunc,
	)
}

func (p *Policy) Save() *Policy {
	// NOTE Don’t use Save with Model, it’s an Undefined Behavior.
	// https://gorm.io/docs/update.html#Save
	db := config.LocalDatabase.Save(p)
	logrus.Debugf("Saved it, RowsAffected: %d", db.RowsAffected)
	return p
}

func (p *Policy) SaveAndLoad() *Policy {
	db := config.LocalDatabase.Save(p)
	logrus.Debugf("Saved it, RowsAffected: %d", db.RowsAffected)
	db.
		Preload("UpdatedBy").
		Preload("CreatedBy").
		Preload("Connection").
		First(p, "uuid = ?", p.Uuid)
	return p
}

func RunAutoMigrate() {
	db := config.LocalDatabase
	err := db.AutoMigrate(
		&model.Modifier{},
		&Policy{},
	)
	if err != nil {
		logrus.Fatal(err)
	}
}
