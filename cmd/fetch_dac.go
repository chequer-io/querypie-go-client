package cmd

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"qpc/config"
	"qpc/models"
)

// @Deprecated - Use fetchAllCmd instead
var fetchDacCmd = &cobra.Command{
	Use:   "dac",
	Short: "Fetch DAC resources from QueryPie API v2",
	Run: func(cmd *cobra.Command, args []string) {
		page := 0
		size := 40 // Set the desired page size

		for {
			var restClient = resty.New()
			var list models.PagedConnectionV2List
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
				Get(defaultQuerypieServer.BaseURL + "/api/external/v2/dac/connections")
			logrus.Debugf("Response: %v", resp)
			if err != nil {
				logrus.Fatalf("Failed to fetch dac connections: %v", err)
			}
			printConnectionV2List(list, page == 0, !list.Page.HasNext())
			saveConnectionV2List(list.List)

			if !list.Page.HasNext() {
				break
			}
			page++
		}
	},
}

func saveConnectionV2List(list []models.SummarizedConnectionV2) {
	for _, conn := range list {
		// Attempt to update the user
		result := config.LocalDatabase.Model(&models.SummarizedConnectionV2{}).Where("uuid = ?", conn.Uuid).Updates(&conn)

		// If no rows were affected, create a new summarized connection
		if result.RowsAffected == 0 {
			if err := config.LocalDatabase.Create(&conn).Error; err != nil {
				logrus.Errorf("Failed to save connection %s: %v", conn.ShortID(), err)
			}
		} else if result.Error != nil {
			logrus.Errorf("Failed to update connection %s: %v", conn.ShortID(), result.Error)
		}
	}
}

func printConnectionV2List(list models.PagedConnectionV2List, first bool, last bool) {
	format := "%-36s  %-10s  %-5s  %-36s  %-8s  %-16s  %-16s\n"
	if first {
		logrus.Debugf("Page of the first: %v", list.Page)
		fmt.Printf(format,
			"UUID",
			"DB_TYPE",
			"CLOUD",
			"NAME",
			"STATUS",
			"CREATED",
			"UPDATED",
		)

	}
	for _, conn := range list.List {
		logrus.Debug(conn)
		cloudProviderType := conn.CloudProviderType
		if cloudProviderType == "" {
			cloudProviderType = "-"
		}
		fmt.Printf(format,
			conn.Uuid,
			conn.DatabaseType,
			cloudProviderType,
			conn.Name,
			conn.Status(),
			conn.ShortCreatedAt(),
			conn.ShortUpdatedAt(),
		)
	}
	if last {
		logrus.Infof("TotalElements: %v", list.Page.TotalElements)
	}
}

func init() {
	// Add fetch subcommands to fetchCmd
	// rootCmd.AddCommand(fetchDacCmd) // @Deprecated - Use fetchAllCmd instead
}
