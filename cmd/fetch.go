package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"qpc/entity/dac_connection"
	"qpc/entity/user"
	"qpc/entity/user_v1"
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
			var pcl dac_connection.PagedConnectionV2List
			pcl.FetchAllAndPrintAndSave()
		case "users":
			var pul user.PagedUserList
			pul.FetchAllAndPrintAndSave()
		case "users-v1":
			var pul user_v1.PagedUserV1List
			pul.FetchAllAndPrintAndSave()
		default:
			logrus.Fatalf("Unknown resource: %s", resource)
		}
	},
}

func init() {
	// fetchAllCmd is added rootCmd in init() of root.go
}
