package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"qpc/models"
	"qpc/rest"
)

var fetchUserCmdV1 = &cobra.Command{
	Use:   "user-v1",
	Short: "Fetch users from QueryPie API v0.9",
	Run: func(cmd *cobra.Command, args []string) {
		page := 0
		size := 40 // Set the desired page size

		for {
			users, err := fetchUsersV1FromQueryPie(defaultQuerypieServer, size, page)
			if err != nil {
				logrus.Fatalf("Failed to fetch user data: %v", err)
			}
			printUserListV1(*users, page == 0, !users.Page.HasNext())

			if !users.Page.HasNext() {
				break
			}
			page++
		}
	},
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

func fetchUsersV1FromQueryPie(querypie QueryPieServerConfig, size int, page int) (*models.PagedUserV1List, error) {
	uri := fmt.Sprintf("/api/external/users?pageSize=%d&pageNumber=%d", size, page)
	client := rest.NewAPIClient(querypie.BaseURL, querypie.AccessToken)

	// Call the GetData method
	result, err := client.GetData(uri)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user data: %v", err)
	}

	// Convert result to []byte
	resultBytes, err := json.Marshal(result)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal result: %v", err)
	}

	// Unmarshal the JSON data into a PagedUserV1List struct
	var users models.PagedUserV1List
	if err := json.Unmarshal(resultBytes, &users); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON data: %v", err)
	}

	return &users, nil
}

var fetchUserCmdV2 = &cobra.Command{
	Use:     "user",
	Aliases: []string{"user-v2"},
	Short:   "Fetch users from QueryPie API v2",
	Run: func(cmd *cobra.Command, args []string) {
		page := 0
		size := 40 // Set the desired page size

		for {
			users, err := fetchUsersV2FromQueryPie(defaultQuerypieServer, size, page)
			if err != nil {
				logrus.Fatalf("Failed to fetch user data: %v", err)
			}
			printUserListV2(*users, page == 0, !users.Page.HasNext())

			if !users.Page.HasNext() {
				break
			}
			page++
		}
	},
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

func fetchUsersV2FromQueryPie(querypie QueryPieServerConfig, size int, page int) (*models.PagedUserV2List, error) {
	uri := fmt.Sprintf("/api/external/v2/users?pageSize=%d&pageNumber=%d", size, page)
	client := rest.NewAPIClient(querypie.BaseURL, querypie.AccessToken)

	// Call the GetData method
	result, err := client.GetData(uri)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user data: %v", err)
	}

	// Convert result to []byte
	resultBytes, err := json.Marshal(result)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal result: %v", err)
	}

	// Unmarshal the JSON data into a PagedUserV1List struct
	var users models.PagedUserV2List
	if err := json.Unmarshal(resultBytes, &users); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON data: %v", err)
	}

	return &users, nil
}

func init() {
	// Add fetch subcommands to fetchCmd
	fetchCmd.AddCommand(fetchUserCmdV1)
	fetchCmd.AddCommand(fetchUserCmdV2)
}
