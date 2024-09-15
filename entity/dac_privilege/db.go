package dac_privilege

import (
	"github.com/sirupsen/logrus"
	"qpc/config"
	"qpc/model"
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
	var privileges []Privilege
	var total, fetched int64 = 0, 0
	result := config.LocalDatabase.Model(&Privilege{}).Count(&total)
	if result.Error != nil {
		logrus.Fatalf("Failed to count items: %v", result.Error)
	}
	logrus.Debugf("Found %d items", total)

	result = config.LocalDatabase.Find(&privileges)
	if result.Error != nil {
		logrus.Fatalf("Failed to select data from local database: %v", result.Error)
		return
	}
	fetched = int64(len(privileges))
	for _, sc := range privileges {
		forEachFunc(&sc)
	}
	if fetched != total {
		logrus.Errorf("Selected %d, whereas total count was %d, difference: %d",
			fetched, total, total-fetched)
	}
}

func (pl *PrivilegePagedList) FindAllAsPagedList(currentPage, pageSize, totalElements int) (PrivilegePagedList, error) {
	var list PrivilegePagedList

	var page model.Page
	page.CurrentPage = currentPage
	page.PageSize = pageSize
	page.TotalElements = totalElements
	page.TotalPages = (totalElements + pageSize - 1) / pageSize

	var items []Privilege
	offset := currentPage * pageSize
	result := config.LocalDatabase.
		Offset(offset).
		Limit(pageSize).
		Find(&items)
	if result.Error != nil {
		return PrivilegePagedList{}, result.Error
	}
	list.List = items
	list.Page = page
	return list, nil
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
