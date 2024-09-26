package cmd

import (
	"github.com/spf13/cobra"
	"qpc/config"
	"qpc/entity/dac_connection"
)

var dacConnectionCmd = &cobra.Command{
	Use:   "connection <command> [...] [flags]",
	Short: "Manage DAC Connections",
}

var dacConnectionListCmd = &cobra.Command{
	Use:   "ls [flags]",
	Short: "List connections in local sqlite database",
	Example: `  ls # List all connections in local sqlite database
  ls --summarized # List all connection as summarized output`,
	PreRun: func(cmd *cobra.Command, args []string) {
		m := config.LocalDatabase.Migrator()
		if !m.HasTable(&dac_connection.SummarizedConnectionV2{}) {
			dac_connection.RunAutoMigrate()
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		silent := false // Silent mode of listing does not make sense.
		summarized, _ := cmd.Flags().GetBool("summarized")

		if summarized {
			var sc dac_connection.SummarizedConnectionV2
			var count = 0
			sc.PrintYamlHeader(silent)
			sc.FindAllAndForEach(func(found *dac_connection.SummarizedConnectionV2) bool {
				found.PrintYaml(silent)
				count++
				return true // OK to continue finding
			})
			sc.PrintYamlFooter(silent, count)
		} else {
			var c dac_connection.ConnectionV2
			var count = 0
			c.PrintYamlHeader(silent)
			c.FindAllAndForEach(func(found *dac_connection.ConnectionV2) bool {
				found.PrintYaml(silent)
				count++
				return true
			})
			c.PrintYamlFooter(silent, count)
		}
	},
}

var dacConnectionFetchCmd = &cobra.Command{
	Use:   "fetch [flags]",
	Short: "Fetch connections from QueryPie API v2 and save them to local sqlite database",
	Example: `  fetch # Fetch all detailed connections from QueryPie API v2
  fetch --summarized # Fetch summarized connections from QueryPie API v2`,
	PreRun: func(cmd *cobra.Command, args []string) {
		m := config.LocalDatabase.Migrator()
		if !m.HasTable(&dac_connection.SummarizedConnectionV2{}) {
			dac_connection.RunAutoMigrate()
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		summarized, _ := cmd.Flags().GetBool("summarized")
		silent, _ := cmd.Flags().GetBool("silent")

		if summarized {
			var sc dac_connection.SummarizedConnectionV2
			var count = 0
			sc.PrintYamlHeader(silent)
			sc.FetchAllAndForEach(func(fetched *dac_connection.SummarizedConnectionV2) bool {
				fetched.PrintYaml(silent).Save()
				count++
				return true // OK to continue fetching
			})
			sc.PrintYamlFooter(silent, count)
		} else {
			var c dac_connection.ConnectionV2
			var count = 0
			c.PrintYamlHeader(silent)
			(&dac_connection.SummarizedConnectionV2{}).
				FetchAllAndForEach(func(fetched *dac_connection.SummarizedConnectionV2) bool {
					fetched.Save()
					count++
					conn := (&dac_connection.ConnectionV2{}).FetchByUuid(fetched.Uuid)
					conn.PrintYaml(silent).SaveAlsoForServerError()
					return true // OK to continue fetching
				})
			c.PrintYamlFooter(silent, count)
		}
	},
}

func addFlagsForConnectionList(cmd *cobra.Command) {
	cmd.Flags().Bool("summarized", false, "Use summarized connection")
}

func addFlagsForConnectionFetch(cmd *cobra.Command) {
	cmd.Flags().Bool("summarized", false, "Use summarized connection")
	cmd.Flags().Bool("silent", false, "Silent or quiet mode. Do not print outputs")
}

func init() {
	addFlagsForConnectionList(dacConnectionListCmd)
	addFlagsForConnectionFetch(dacConnectionFetchCmd)

	dacConnectionCmd.AddCommand(dacConnectionListCmd)
	dacConnectionCmd.AddCommand(dacConnectionFetchCmd)
	// dacConnectionCmd is added dacCmd in init() of dac.go
}
