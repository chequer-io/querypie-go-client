package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"qpc/models"
	"qpc/rest"
)

var fetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "Fetch various resources from QueryPie server",
}

var fetchUserCmd = &cobra.Command{
	Use:   "user",
	Short: "Fetch users from QueryPie server",
	Run: func(cmd *cobra.Command, args []string) {
		page := 0
		size := 40 // Set the desired page size

		for {
			users, err := fetchUsersFromQueryPie(defaultQuerypieServer, size, page)
			if err != nil {
				logrus.Fatalf("Failed to fetch user data: %v", err)
			}
			printUserList(*users, page == 0, !users.Page.HasNext())

			if !users.Page.HasNext() {
				break
			}
			page++
		}
	},
}

func printUserList(list models.PagedUserList, first bool, last bool) {
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

func fetchUsersFromQueryPie(querypie QueryPieServerConfig, size int, page int) (*models.PagedUserList, error) {
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

	// Unmarshal the JSON data into a PagedUserList struct
	var users models.PagedUserList
	if err := json.Unmarshal(resultBytes, &users); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON data: %v", err)
	}

	return &users, nil
}

var fetchDbCmd = &cobra.Command{
	Use:   "db",
	Short: "Fetch database connections from QueryPie server",
	Run: func(cmd *cobra.Command, args []string) {
		// Implement logic to fetch database connections
		fmt.Println("TODO(JK): Fetching database connections...")
	},
}

var fetchServerCmd = &cobra.Command{
	Use:   "server",
	Short: "Fetch server connections from QueryPie server",
	Run: func(cmd *cobra.Command, args []string) {
		// Implement logic to fetch server connections
		fmt.Println("TODO(JK): Fetching server connections...")
	},
}

func init() {
	// Add fetch subcommands to fetchCmd
	fetchCmd.AddCommand(fetchUserCmd)
	fetchCmd.AddCommand(fetchDbCmd)
	fetchCmd.AddCommand(fetchServerCmd)

	// Add fetchCmd to the root command
	rootCmd.AddCommand(fetchCmd)
}
