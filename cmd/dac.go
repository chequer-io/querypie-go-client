package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"qpc/config"
	"qpc/entity/dac_access_control"
	"qpc/models"
)

var dacCmd = &cobra.Command{
	Use:   "dac",
	Short: "Manage DAC resources",
}

var dacFetchAllConnectionsCmd = &cobra.Command{
	Use:   "fetch-all <resource>",
	Short: "Fetch all DAC connections from QueryPie server and save them to local sqlite database",
	Example: `  fetch-all connections # from QueryPie API v2
  fetch-all access-controls # from QueryPie API v2`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			_ = cmd.Help()
			return
		}
		resource := args[0]
		switch resource {
		case "connections":
			fetchDACPrintAndSave()
		case "access-controls":
			fetchDACAccessControlPrintAndSave()
		default:
			logrus.Fatalf("Unknown resource: %s", resource)
		}
	},
}

var dacListCmd = &cobra.Command{
	Use:   "ls",
	Short: "List DAC connections in local sqlite database",
	Run: func(cmd *cobra.Command, args []string) {
		var total, fetched int64 = 0, 0
		result := config.LocalDatabase.Model(&models.SummarizedConnectionV2{}).Count(&total)
		if result.Error != nil {
			logrus.Fatalf("Failed to count dac connections: %v", result.Error)
		}
		logrus.Debugf("Found %d dac connections", total)

		page := 0
		size := 30 // Set the desired page size

		for {
			list, err := selectPagedConnectionV2List(page, size, int(total))
			if err != nil {
				logrus.Fatalf("Failed to fetch dac connections: %v", err)
			}
			logrus.Debugf("Fetched %d, page %d, size %d, total %d",
				len(list.List), page, size, total)
			fetched += int64(len(list.List))
			printConnectionV2List(list, page == 0, !list.Page.HasNext())

			if !list.Page.HasNext() {
				break
			}
			page++
		}
		logrus.Debugf("Fetched %d, whereas total count was %d, difference: %d",
			fetched, total, total-fetched)
	},
}

func fetchDACPrintAndSave() {
	fetchPrintAndSave(
		"/api/external/v2/dac/connections",
		&models.PagedConnectionV2List{},
		func(result *models.PagedConnectionV2List, first bool, last bool) {
			printConnectionV2List(*result, first, last)
		},
		func(result *models.PagedConnectionV2List) bool {
			saveConnectionV2List(result.List)
			return !result.Page.HasNext()
		},
	)
}

func fetchDACAccessControlPrintAndSave() {
	fetchPrintAndSave(
		"/api/external/v2/dac/access-controls",
		&dac_access_control.SummarizedAccessControlPagedList{},
		func(result *dac_access_control.SummarizedAccessControlPagedList, first bool, last bool) {
			result.Print()
		},
		func(result *dac_access_control.SummarizedAccessControlPagedList) bool {
			// result.Save()
			return !result.Page.HasNext()
		},
	)
}

func selectPagedConnectionV2List(currentPage, pageSize, totalElements int) (models.PagedConnectionV2List, error) {
	var pagedConnectionV2List models.PagedConnectionV2List
	var page models.Page
	var connections []models.SummarizedConnectionV2
	offset := currentPage * pageSize
	result := config.LocalDatabase.
		Offset(offset).
		Limit(pageSize).
		Find(&connections)
	if result.Error != nil {
		return models.PagedConnectionV2List{}, result.Error
	}
	page.CurrentPage = currentPage
	page.PageSize = pageSize
	page.TotalElements = totalElements
	page.TotalPages = (totalElements + pageSize - 1) / pageSize
	pagedConnectionV2List.List = connections
	pagedConnectionV2List.Page = page
	return pagedConnectionV2List, nil
}

func init() {
	// Add dacListCmd subcommands to dacCmd
	dacCmd.AddCommand(dacListCmd)
	dacCmd.AddCommand(dacFetchAllConnectionsCmd)

	// dacCmd is added rootCmd in init() of root.go
}
