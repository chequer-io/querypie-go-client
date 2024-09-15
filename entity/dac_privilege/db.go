package dac_privilege

import (
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"qpc/config"
	"qpc/utils"
)

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
			logrus.Fatalf("Failed to create privilege %s: %v", p.ShortID(), err)
		}
	} else if result.Error != nil {
		logrus.Fatalf("Failed to update privilege %s: %v", p.ShortID(), result.Error)
	}
	return p
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
