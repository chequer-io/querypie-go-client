package user_v1

import (
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
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
	forEachFunc func(found *UserV1) bool,
) {
	utils.FindAllAndForEach(
		func(tx *gorm.DB, total *int64) *gorm.DB {
			return tx.Model(&UserV1{}).Count(total)
		},
		func(tx *gorm.DB, items *[]UserV1) *gorm.DB {
			return tx.Model(&UserV1{}).Find(items)
		},
		forEachFunc,
	)
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
