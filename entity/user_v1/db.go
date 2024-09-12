package user_v1

import (
	"github.com/sirupsen/logrus"
	"qpc/config"
	"qpc/model"
	"qpc/utils"
)

func (pul *PagedUserV1List) Save() {
	for _, user := range pul.GetList() {
		user.Save()
	}
}

func (pul *PagedUserV1List) FetchAllAndPrintAndSave() {
	utils.FetchPrintAndSave(
		"/api/external/users",
		&PagedUserV1List{},
		func(result *PagedUserV1List, first bool, last bool) {
			result.Print()
		},
		func(result *PagedUserV1List) bool {
			result.Save()
			return !result.Page.HasNext()
		},
	)
}

func (u *UserV1) Save() {
	// Attempt to update the user
	result := config.LocalDatabase.Model(&UserV1{}).Where("uuid = ?", u.Uuid).Updates(&u)

	// If no rows were affected, create a new user
	if result.RowsAffected == 0 {
		if err := config.LocalDatabase.Create(&u).Error; err != nil {
			logrus.Errorf("Failed to save user %s: %v", u.ShortID(), err)
		}
	} else if result.Error != nil {
		logrus.Errorf("Failed to update user %s: %v", u.ShortID(), result.Error)
	}
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
