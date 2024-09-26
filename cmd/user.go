package cmd

import (
	"github.com/spf13/cobra"
	"qpc/config"
	"qpc/entity/user"
	"qpc/entity/user_v1"
)

var userCmd = &cobra.Command{
	Use:   "user",
	Short: "Manage user accounts with QueryPie API v2",
}

var userFetchAllCmd = &cobra.Command{
	Use:   "fetch",
	Short: "Fetch all users from QueryPie server and save them to local sqlite database",
	PreRun: func(cmd *cobra.Command, args []string) {
		if !config.LocalDatabase.Migrator().HasTable(&user.User{}) {
			user.RunAutoMigrate()
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		var u user.User
		u.PrintHeader()
		u.FetchAllAndForEach(func(fetched *user.User) bool {
			fetched.Print().Save()
			return true // OK to continue fetching
		})
	},
}

var userListCmd = &cobra.Command{
	Use:   "ls",
	Short: "List all users in local sqlite database",
	PreRun: func(cmd *cobra.Command, args []string) {
		if !config.LocalDatabase.Migrator().HasTable(&user.User{}) {
			user.RunAutoMigrate()
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		var u user.User
		u.PrintHeader()
		u.FetchAllAndForEach(func(fetched *user.User) bool {
			fetched.Print()
			return true // OK to continue finding
		})
	},
}

var userV1Cmd = &cobra.Command{
	Use:   "user-v1",
	Short: "Manage user accounts with QueryPie API v0.9",
}

var userV1FetchAllCmd = &cobra.Command{
	Use:   "fetch",
	Short: "Fetch all users from QueryPie server and save them to local sqlite database",
	PreRun: func(cmd *cobra.Command, args []string) {
		if !config.LocalDatabase.Migrator().HasTable(&user_v1.UserV1{}) {
			user_v1.RunAutoMigrate()
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		var u user_v1.UserV1
		u.PrintHeader()
		u.FetchAllAndForEach(func(fetched *user_v1.UserV1) bool {
			fetched.Print().Save()
			return true // OK to continue fetching
		})
	},
}

var userV1ListCmd = &cobra.Command{
	Use:   "ls",
	Short: "List all users in local sqlite database",
	PreRun: func(cmd *cobra.Command, args []string) {
		if !config.LocalDatabase.Migrator().HasTable(&user_v1.UserV1{}) {
			user_v1.RunAutoMigrate()
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		var u user_v1.UserV1
		u.PrintHeader()
		u.FetchAllAndForEach(func(fetched *user_v1.UserV1) bool {
			fetched.Print()
			return true // OK to continue finding
		})
	},
}

func init() {
	userCmd.AddCommand(userFetchAllCmd)
	userCmd.AddCommand(userListCmd)
	userV1Cmd.AddCommand(userV1FetchAllCmd)
	userV1Cmd.AddCommand(userV1ListCmd)

	// userCmd, userV1Cmd are added rootCmd in init() of root.go
}
