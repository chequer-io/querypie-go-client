package user

import (
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"qpc/config"
	"qpc/utils"
)

func (u *User) FetchAllAndForEach(
	forEachFunc func(fetched *User) bool,
) {
	utils.FetchPagedListAndForEach(
		"/api/external/v2/users",
		&PagedUserList{},
		func(page *PagedUserList) bool {
			for _, it := range page.List {
				forEachFunc(&it)
			}
			return true
		},
	)
}

func (u *User) FindByLoginIdOrEmailOrUuid(query string, users *[]User) {
	utils.FindMultiple(users, func(db *gorm.DB) *gorm.DB {
		return db.
			Model(&User{}).
			// Note: Column names are snake_case in the database.
			Where("login_id = ? OR email = ? OR uuid = ?", query, query, query).
			Find(users)
	})
}

func (u *User) FindAllAndForEach(
	forEachFunc func(found *User) bool,
) {
	utils.FindAllAndForEach(
		func(tx *gorm.DB, total *int64) *gorm.DB {
			return tx.Model(&User{}).Count(total)
		},
		func(tx *gorm.DB, items *[]User) *gorm.DB {
			return tx.Model(&User{}).Find(items)
		},
		forEachFunc,
	)
}

func (u *User) Save() *User {
	// NOTE Don’t use Save with Model, it’s an Undefined Behavior.
	// https://gorm.io/docs/update.html#Save
	db := config.LocalDatabase.Save(u)
	logrus.Debugf("Saved it, RowsAffected: %d", db.RowsAffected)
	return u
}

func RunAutoMigrate() {
	db := config.LocalDatabase
	err := db.AutoMigrate(
		&User{},
		&AdminRole{},
	)
	if err != nil {
		logrus.Fatal(err)
	}
}
