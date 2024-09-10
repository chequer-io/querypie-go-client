package cmd

import (
	"github.com/spf13/cobra"
)

var userCmd = &cobra.Command{
	Use:   "user",
	Short: "Manage user accounts with QueryPie API v2",
}

var userFetchAllCmd = &cobra.Command{
	Use:   "fetch-all",
	Short: "Fetch all users from QueryPie server and save them to local sqlite database",
	Run: func(cmd *cobra.Command, args []string) {
		fetchUserPrintAndSave()
	},
}

var userV1Cmd = &cobra.Command{
	Use:   "user-v1",
	Short: "Manage user accounts with QueryPie API v0.9",
}

var userV1FetchAllCmd = &cobra.Command{
	Use:   "fetch-all",
	Short: "Fetch all users from QueryPie server and save them to local sqlite database",
	Run: func(cmd *cobra.Command, args []string) {
		fetchUserV1PrintAndSave()
	},
}

func init() {
	userCmd.AddCommand(userFetchAllCmd)
	userV1Cmd.AddCommand(userV1FetchAllCmd)

	// userCmd, userV1Cmd are added rootCmd in init() of root.go
}
