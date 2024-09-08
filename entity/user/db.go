package user

import (
	"github.com/sirupsen/logrus"
	"qpc/config"
)

func (pul PagedUserList) Save() {
	for _, user := range pul.List {
		user.Save()
	}
}

func (u User) Save() {
	// Attempt to update the user
	result := config.LocalDatabase.Model(&User{}).Where("uuid = ?", u.Uuid).Updates(&u)

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
		&User{},
		&AdminRole{},
	)
	if err != nil {
		logrus.Fatal(err)
	}
}
