package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"qpc/local_db"
	"qpc/models"
)

var dacCmd = &cobra.Command{
	Use:   "dac",
	Short: "Manage DAC resources",
}

var dacListCmd = &cobra.Command{
	Use:   "ls",
	Short: "List DAC connections in local sqlite database",
	Run: func(cmd *cobra.Command, args []string) {
		var total, fetched int64 = 0, 0
		result := local_db.LocalDatabase.Model(&models.SummarizedConnectionV2{}).Count(&total)
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

func selectPagedConnectionV2List(currentPage, pageSize, totalElements int) (models.PagedConnectionV2List, error) {
	var pagedConnectionV2List models.PagedConnectionV2List
	var page models.Page
	var connections []models.SummarizedConnectionV2
	offset := currentPage * pageSize
	result := local_db.LocalDatabase.
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

	// Add dacCmd to the root command
	rootCmd.AddCommand(dacCmd)
}
