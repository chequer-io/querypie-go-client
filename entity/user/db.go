package user

import (
	"github.com/sirupsen/logrus"
	"qpc/config"
	"qpc/utils"
)

// Save @Deprecated
func (pul *PagedUserList) Save() {
	for _, user := range pul.GetList() {
		user.Save()
	}
}

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

func (u *User) Save() *User {
	// Attempt to update the user
	result := config.LocalDatabase.Model(&User{}).Where("uuid = ?", u.Uuid).Updates(&u)

	// If no rows were affected, create a new user
	if result.RowsAffected == 0 {
		if err := config.LocalDatabase.Create(&u).Error; err != nil {
			logrus.Fatalf("Failed to save user %s: %v", u.ShortID(), err)
		}
	} else if result.Error != nil {
		logrus.Fatalf("Failed to update user %s: %v", u.ShortID(), result.Error)
	}
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
