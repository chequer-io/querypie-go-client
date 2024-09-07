package cmd

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"qpc/local_db"
	"qpc/models"
)

// @Deprecated - Use fetchAllCmd instead
var fetchUserCmdV1 = &cobra.Command{
	Use:   "fetch-all-user-v1",
	Short: "Fetch users from QueryPie API v0.9",
	Run: func(cmd *cobra.Command, args []string) {
		page := 0
		size := 40 // Set the desired page size

		for {
			var restClient = resty.New()
			var list models.PagedUserV1List
			resp, err := restClient.R().
				SetQueryParams(
					map[string]string{
						"pageSize":   fmt.Sprintf("%d", size),
						"pageNumber": fmt.Sprintf("%d", page),
					},
				).
				SetHeader("Accept", "application/json").
				SetAuthToken(defaultQuerypieServer.AccessToken).
				SetResult(&list).
				Get(defaultQuerypieServer.BaseURL + "/api/external/users")
			logrus.Debugf("Response: %v", resp)
			if err != nil {
				logrus.Fatalf("Failed to fetch user data: %v", err)
			}
			printUserListV1(list, page == 0, !list.Page.HasNext())
			saveUserListV1(list.List)

			if !list.Page.HasNext() {
				break
			}
			page++
		}
	},
}

func saveUserListV1(list []models.UserV1) {
	for _, user := range list {
		// Attempt to update the user
		result := local_db.LocalDatabase.Model(&models.UserV1{}).Where("uuid = ?", user.Uuid).Updates(&user)

		// If no rows were affected, create a new user
		if result.RowsAffected == 0 {
			if err := local_db.LocalDatabase.Create(&user).Error; err != nil {
				logrus.Errorf("Failed to save user %s: %v", user.ShortID(), err)
			}
		} else if result.Error != nil {
			logrus.Errorf("Failed to update user %s: %v", user.ShortID(), result.Error)
		}
	}
}

func printUserListV1(list models.PagedUserV1List, first bool, last bool) {
	format := "%-36s  %-24s  %-24s  %-20s  %-8s  %-16s  %-16s\n"
	if first {
		logrus.Debugf("Page: %v", list.Page)
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
	for _, user := range list.List {
		logrus.Debug(user)
		fmt.Printf(format,
			user.Uuid,
			user.LoginId,
			user.Email,
			user.Name,
			user.Status(),
			user.ShortCreatedAt(),
			user.ShortUpdatedAt(),
		)
	}
	if last {
		logrus.Infof("TotalElements: %v", list.Page.TotalElements)
	}
}

// @Deprecated - Use fetchAllCmd instead
var fetchUserCmdV2 = &cobra.Command{
	Use:   "fetch-all-user",
	Short: "Fetch users from QueryPie API v2",
	Run: func(cmd *cobra.Command, args []string) {
		page := 0
		size := 40 // Set the desired page size

		for {
			var restClient = resty.New()
			var list models.PagedUserV2List
			resp, err := restClient.R().
				SetQueryParams(
					map[string]string{
						"pageSize":   fmt.Sprintf("%d", size),
						"pageNumber": fmt.Sprintf("%d", page),
					},
				).
				SetHeader("Accept", "application/json").
				SetAuthToken(defaultQuerypieServer.AccessToken).
				SetResult(&list).
				Get(defaultQuerypieServer.BaseURL + "/api/external/v2/users")
			logrus.Debugf("Response: %v", resp)
			if err != nil {
				logrus.Fatalf("Failed to fetch user data: %v", err)
			}
			printUserListV2(list, page == 0, !list.Page.HasNext())
			saveUserListV2(list.List)

			if !list.Page.HasNext() {
				break
			}
			page++
		}
	},
}

func saveUserListV2(list []models.UserV2) {
	for _, user := range list {
		// Attempt to update the user
		result := local_db.LocalDatabase.Model(&models.UserV2{}).Where("uuid = ?", user.Uuid).Updates(&user)

		// If no rows were affected, create a new user
		if result.RowsAffected == 0 {
			if err := local_db.LocalDatabase.Create(&user).Error; err != nil {
				logrus.Errorf("Failed to save user %s: %v", user.ShortID(), err)
			}
		} else if result.Error != nil {
			logrus.Errorf("Failed to update user %s: %v", user.ShortID(), result.Error)
		}
	}
}

func printUserListV2(list models.PagedUserV2List, first bool, last bool) {
	format := "%-36s  %-22s  %-22s  %-18s  %-8s  %-8s %-16s  %-16s\n"
	if first {
		logrus.Debugf("Page: %v", list.Page)
		fmt.Printf(format,
			"UUID",
			"LOGIN_ID",
			"EMAIL",
			"NAME",
			"STATUS",
			"MORE",
			"CREATED",
			"UPDATED",
		)

	}
	for _, user := range list.List {
		logrus.Debug(user)
		fmt.Printf(format,
			user.Uuid,
			user.LoginId,
			user.Email,
			user.Name,
			user.Status,
			user.StatusMore(),
			user.ShortCreatedAt(),
			user.ShortUpdatedAt(),
		)
	}
	if last {
		logrus.Infof("TotalElements: %v", list.Page.TotalElements)
	}
}

func init() {
	// Add fetch subcommands to fetchCmd
}
