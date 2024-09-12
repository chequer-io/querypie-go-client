package cmd

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"qpc/config"
	"qpc/model"
)

func saveUserListV1(list []model.UserV1) {
	for _, userV1 := range list {
		// Attempt to update the user
		result := config.LocalDatabase.Model(&model.UserV1{}).Where("uuid = ?", userV1.Uuid).Updates(&userV1)

		// If no rows were affected, create a new user
		if result.RowsAffected == 0 {
			if err := config.LocalDatabase.Create(&userV1).Error; err != nil {
				logrus.Errorf("Failed to save user %s: %v", userV1.ShortID(), err)
			}
		} else if result.Error != nil {
			logrus.Errorf("Failed to update user %s: %v", userV1.ShortID(), result.Error)
		}
	}
}

func printUserListV1(list model.PagedUserV1List, first bool, last bool) {
	format := "%-36s  %-24s  %-24s  %-20s  %-8s  %-16s  %-16s\n"
	if first {
		logrus.Debugf("Page: %v", list.GetPage())
		fmt.Printf(format,
			"UUID",
			"LOGIN_ID",
			"EMAIL",
			"NAME",
			"STATUS",
			"CREATED",
			"UPDATED",
		)

	}
	for _, u := range list.GetList() {
		logrus.Debug(u)
		fmt.Printf(format,
			u.Uuid,
			u.LoginId,
			u.Email,
			u.Name,
			u.Status(),
			u.ShortCreatedAt(),
			u.ShortUpdatedAt(),
		)
	}
	if last {
		logrus.Infof("TotalElements: %v", list.GetPage().TotalElements)
	}
}

func init() {
	// Add fetch subcommands to fetchCmd
}
