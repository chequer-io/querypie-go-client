package cmd

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"qpc/models"
)

var fetchAllCmd = &cobra.Command{
	Use:   "fetch-all <resource>",
	Short: "Fetch various resources from QueryPie server, and save them to local sqlite database",
	Example: `  fetch-all dac       # DAC resources from QueryPie API v2
  fetch-all sac       # N/A Yet - SAC resources from QueryPie API v2
  fetch-all user      # Users from QueryPie API v2
  fetch-all user-v1   # Users from QueryPie API v0.9`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			_ = cmd.Help()
			return
		}
		resource := validate(args[0])
		uri := getUri(resource)
		logrus.Debugf("URI: %s", uri)
		fetchPrintSave(resource, uri)
	},
}

func validate(resource string) string {
	switch resource {
	case "dac", "sac", "user", "user-v1":
		return resource
	default:
		logrus.Fatalf("Unknown resource: %s", resource)
		return ""
	}
}

func getUri(resource string) string {
	switch resource {
	case "dac":
		return "/api/external/v2/dac/connections"
	case "sac":
		return "/api/external/v2/sac/servers"
	case "user":
		return "/api/external/v2/users"
	case "user-v1":
		return "/api/external/users"
	default:
		logrus.Fatalf("Unknown resource: %s", resource)
		return ""
	}
}

func fetchPrintSave(resource string, uri string) {
	page := 0
	size := 40 // Set the desired page size
	restClient := resty.New()

	var result interface{}
	switch resource {
	case "dac":
		result = &models.PagedConnectionV2List{}
	case "sac":
		result = nil
	case "user":
		result = &models.PagedUserV2List{}
	case "user-v1":
		result = &models.PagedUserV1List{}
	default:
		logrus.Fatalf("Unknown resource: %s", resource)
		return
	}
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

		switch v := result.(type) {
		case *models.PagedConnectionV2List:
			printConnectionV2List(*v, page == 0, !v.Page.HasNext())
		case *models.PagedUserV2List:
			printUserListV2(*v, page == 0, !v.Page.HasNext())
		case *models.PagedUserV1List:
			printUserListV1(*v, page == 0, !v.Page.HasNext())
		default:
			logrus.Fatalf("printPagedList() Unknown type: %T", v)
		}

		shouldBreak := false
		switch v := result.(type) {
		case *models.PagedConnectionV2List:
			saveConnectionV2List(v.List)
			if !v.Page.HasNext() {
				shouldBreak = true
			}
		case *models.PagedUserV2List:
			saveUserListV2(v.List)
			if !v.Page.HasNext() {
				shouldBreak = true
			}
		case *models.PagedUserV1List:
			saveUserListV1(v.List)
			if !v.Page.HasNext() {
				shouldBreak = true
			}
		default:
			logrus.Fatalf("savePagedList() Unknown type: %T", v)
		}

		if shouldBreak {
			break
		}

		page++
	}
}

func init() {
	// fetchAllCmd is added rootCmd in init() of root.go
}
