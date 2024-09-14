package user_v1

import (
	"github.com/sirupsen/logrus"
	"qpc/config"
	"qpc/model"
	"qpc/utils"
)

func (u *UserV1) FetchAllAndForEach(
	forEachFunc func(fetched *UserV1) bool,
) {
	utils.FetchPagedListAndForEach(
		"/api/external/users",
		&PagedUserV1List{},
		func(page *PagedUserV1List) bool {
			for _, it := range page.List {
				forEachFunc(&it)
			}
			return true
		},
	)
}

func (u *UserV1) FindAllAndForEach(
	forEachFunc func(sc *UserV1) bool,
) {
	var connections []UserV1
	var total, fetched int64 = 0, 0
	result := config.LocalDatabase.Model(&UserV1{}).Count(&total)
	if result.Error != nil {
		logrus.Fatalf("Failed to count items: %v", result.Error)
	}
	logrus.Debugf("Found %d items", total)

	result = config.LocalDatabase.Find(&connections)
	if result.Error != nil {
		logrus.Fatalf("Failed to select data from local database: %v", result.Error)
		return
	}
	fetched = int64(len(connections))
	for _, sc := range connections {
		forEachFunc(&sc)
	}
	if fetched != total {
		logrus.Errorf("Selected %d, whereas total count was %d, difference: %d",
			fetched, total, total-fetched)
	}
}

func (pul *PagedUserV1List) FindAllAsPagedList(currentPage, pageSize, totalElements int) (PagedUserV1List, error) {
	var list PagedUserV1List

	var page model.Page
	page.CurrentPage = currentPage
	page.PageSize = pageSize
	page.TotalElements = totalElements
	page.TotalPages = (totalElements + pageSize - 1) / pageSize

	var items []UserV1
	offset := currentPage * pageSize
	result := config.LocalDatabase.
		Offset(offset).
		Limit(pageSize).
		Find(&items)
	if result.Error != nil {
		return PagedUserV1List{}, result.Error
	}
	list.List = items
	list.Page = page
	return list, nil
}

func (u *UserV1) Save() *UserV1 {
	// Attempt to update the user
	result := config.LocalDatabase.Model(&UserV1{}).Where("uuid = ?", u.Uuid).Updates(&u)

	// If no rows were affected, create a new user
	if result.RowsAffected == 0 {
		if err := config.LocalDatabase.Create(&u).Error; err != nil {
			logrus.Fatalf("Failed to create user %s: %v", u.ShortID(), err)
		}
	} else if result.Error != nil {
		logrus.Fatalf("Failed to update user %s: %v", u.ShortID(), result.Error)
	}
	return u
}

func RunAutoMigrate() {
	db := config.LocalDatabase
	err := db.AutoMigrate(
		&UserV1{},
		&UserRole{},
		&model.Role{},
	)
	if err != nil {
		logrus.Fatal(err)
	}
}
