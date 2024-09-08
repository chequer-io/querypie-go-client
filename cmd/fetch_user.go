package cmd

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"qpc/config"
	"qpc/entity/user"
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
	for _, userV1 := range list {
		// Attempt to update the user
		result := config.LocalDatabase.Model(&models.UserV1{}).Where("uuid = ?", userV1.Uuid).Updates(&userV1)

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
	for _, u := range list.List {
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
			var list user.PagedUserList
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
			list.Print()
			list.Save()

			if !list.Page.HasNext() {
				break
			}
			page++
		}
	},
}

func init() {
	// Add fetch subcommands to fetchCmd
}
