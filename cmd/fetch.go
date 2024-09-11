package cmd

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"qpc/entity/user"
	"qpc/models"
)

var fetchAllCmd = &cobra.Command{
	Use:   "fetch-all <resource>",
	Short: "Fetch various resources from QueryPie server, and save them to local sqlite database",
	Example: `  fetch-all dac       # DAC resources from QueryPie API v2
  fetch-all sac       # N/A Yet - SAC resources from QueryPie API v2
  fetch-all users     # Users from QueryPie API v2
  fetch-all users-v1  # Users from QueryPie API v0.9`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			_ = cmd.Help()
			return
		}
		resource := args[0]
		switch resource {
		case "dac":
			fetchDACPrintAndSave()
		case "users":
			fetchUserPrintAndSave()
		case "users-v1":
			fetchUserV1PrintAndSave()
		default:
			logrus.Fatalf("Unknown resource: %s", resource)
		}
	},
}

func fetchUserPrintAndSave() {
	fetchPrintAndSave(
		"/api/external/v2/users",
		&user.PagedUserList{},
		func(result *user.PagedUserList, first bool, last bool) {
			result.Print()
		},
		func(result *user.PagedUserList) bool {
			result.Save()
			return !result.Page.HasNext()
		},
	)
}

func fetchUserV1PrintAndSave() {
	fetchPrintAndSave(
		"/api/external/users",
		&models.PagedUserV1List{},
		func(result *models.PagedUserV1List, first bool, last bool) {
			printUserListV1(*result, first, last)
		},
		func(result *models.PagedUserV1List) bool {
			saveUserListV1(result.GetList())
			return !result.Page.HasNext()
		},
	)
}

func fetchPrintAndSave[T any, P models.PagedList[T]](
	uri string,
	result P,
	printFunc func(object P, first bool, last bool),
	saveFunc func(object P) bool,
) {
	page := 0
	size := 40 // Set the desired page size
	restClient := resty.New()

	logrus.Debugf("Type of result: %T", result)

	for {
		resp, err := restClient.R().
			SetQueryParams(
				map[string]string{
					"pageSize":   fmt.Sprintf("%d", size),
					"pageNumber": fmt.Sprintf("%d", page),
				},
			).
			SetHeader("Accept", "application/json").
			SetAuthToken(defaultQuerypieServer.AccessToken).
			SetResult(&result).
			Get(defaultQuerypieServer.BaseURL + uri)
		logrus.Debugf("Response: %v", resp)
		if err != nil {
			logrus.Fatalf("Failed to fetch resources: %v", err)
		}

		printFunc(result, page == 0, !result.GetPage().HasNext())
		shouldBreak := saveFunc(result)
		if shouldBreak {
			break
		}

		page++
	}
}

func init() {
	// fetchAllCmd is added rootCmd in init() of root.go
}
