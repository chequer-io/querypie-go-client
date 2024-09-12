package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"qpc/config"
	"qpc/entity/dac_access_control"
	"qpc/entity/dac_connection"
	"qpc/model"
	"qpc/utils"
	"strconv"
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
			var pcl dac_connection.PagedConnectionV2List
			pcl.FetchAllAndPrintAndSave()
		case "access-controls":
			var acl dac_access_control.SummarizedAccessControlPagedList
			acl.FetchAllAndPrintAndSave()
		default:
			logrus.Fatalf("Unknown resource: %s", resource)
		}
	},
}

var dacListCmd = &cobra.Command{
	Use:   "ls",
	Short: "List DAC connections in local sqlite database",
	Example: `  ls connections # from local sqlite database
  ls access-controls # from local sqlite database`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			_ = cmd.Help()
			return
		}
		resource := args[0]
		switch resource {
		case "connections":
			var pcl dac_connection.PagedConnectionV2List
			pcl.FindAllAndPrint()
		case "access-controls":
			selectFromDatabaseAndPrintSummarizedAccessControlPagedList()
		default:
			logrus.Fatalf("Unknown resource: %s", resource)
		}
	},
}

func selectFromDatabaseAndPrintSummarizedAccessControlPagedList() {
	var total, fetched int64 = 0, 0
	result := config.LocalDatabase.Model(&dac_access_control.SummarizedAccessControl{}).Count(&total)
	if result.Error != nil {
		logrus.Fatalf("Failed to count dac connections: %v", result.Error)
	}
	logrus.Debugf("Found %d dac connections", total)

	page := 0
	size := 30 // Set the desired page size

	for {
		list, err := selectSummarizedAccessControlPagedList(page, size, int(total))
		if err != nil {
			logrus.Fatalf("Failed to select data from local database: %v", err)
		}
		logrus.Debugf("Selected %d, page %d, size %d, total %d",
			len(list.List), page, size, total)
		fetched += int64(len(list.List))
		list.Print()
		if !list.Page.HasNext() {
			break
		}
		page++
	}
	logrus.Debugf("Selected %d, whereas total count was %d, difference: %d",
		fetched, total, total-fetched)

}

func selectSummarizedAccessControlPagedList(
	currentPage, pageSize, totalElements int,
) (dac_access_control.SummarizedAccessControlPagedList, error) {
	var acl dac_access_control.SummarizedAccessControlPagedList
	var page model.Page
	var sac []dac_access_control.SummarizedAccessControl
	offset := currentPage * pageSize
	result := config.LocalDatabase.
		Offset(offset).
		Limit(pageSize).
		Find(&sac)
	if result.Error != nil {
		return dac_access_control.SummarizedAccessControlPagedList{}, result.Error
	}
	for i := range sac {
		sac[i].PopulateMembers()
		logrus.Debugf("Populated Members[%d]: %v from MembersStr: %v", i, sac[i].Members, sac[i].MembersStr)
	}

	page.CurrentPage = currentPage
	page.PageSize = pageSize
	page.TotalElements = totalElements
	page.TotalPages = (totalElements + pageSize - 1) / pageSize
	acl.List = sac
	acl.Page = page
	return acl, nil
}

var grantByUuidCmd = &cobra.Command{
	Use:   "grant-by-uuid <user-uuid> <connection-uuid> <privilege-uuid> [<force>]",
	Short: "Grant access to a DAC connection using UUIDs as argument",
	Example: `  grant-by-uuid <uuid> <uuid> <uuid>
  grant-by-uuid <uuid> <uuid> <uuid> false
  grant-by-uuid <uuid> <uuid> <uuid> true`,
	Args: cobra.RangeArgs(3, 4),
	Run: func(cmd *cobra.Command, args []string) {
		userUuid := args[0]
		clusterUuid := args[1]
		privilegeUuid := args[2]
		force := false
		if len(args) > 3 {
			force, _ = strconv.ParseBool(args[3])
		}
		(&dac_access_control.GrantRequest{
			UserUuid:      userUuid,
			ClusterUuid:   clusterUuid,
			PrivilegeUuid: privilegeUuid,
			Force:         force,
		}).
			Post(utils.DefaultQuerypieServer).
			Print()
	},
}

func init() {
	// Add dacListCmd subcommands to dacCmd
	dacCmd.AddCommand(dacListCmd)
	dacCmd.AddCommand(dacFetchAllConnectionsCmd)
	dacCmd.AddCommand(grantByUuidCmd)

	// dacCmd is added rootCmd in init() of root.go
}
