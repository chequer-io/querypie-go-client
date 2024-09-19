package cmd

import (
	"github.com/spf13/cobra"
	"qpc/config"
	"qpc/entity/dac_connection"
)

var dacConnectionsCmd = &cobra.Command{
	Use:   "connections [flags]",
	Short: "List connections or fetch them from QueryPie API v2",
	Example: `  connections # List all connections in local sqlite database
  connections --summarized # List all summarized connections in local sqlite database
  connections --fetch # Fetch all detailed connections from QueryPie API v2
  connections --fetch --summarized # Fetch summarized connections from QueryPie API v2`,
	PreRun: func(cmd *cobra.Command, args []string) {
		m := config.LocalDatabase.Migrator()
		if !m.HasTable(&dac_connection.SummarizedConnectionV2{}) {
			dac_connection.RunAutoMigrate()
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		fetch, _ := cmd.Flags().GetBool("fetch")
		summarized, _ := cmd.Flags().GetBool("summarized")
		silent, _ := cmd.Flags().GetBool("silent")

		if fetch {
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
		} else {
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
		}
	},
}

func addFlagsForConnections(cmd *cobra.Command) {
	cmd.Flags().Bool("fetch", false, "Fetch from QueryPie API v2 and save to local sqlite database")
	cmd.Flags().Bool("summarized", false, "Use summarized connection")
	cmd.Flags().Bool("silent", false, "Silent or quiet mode. Do not print outputs")
}

func init() {
	addFlagsForConnections(dacConnectionsCmd)
	// dacCmd is added rootCmd in init() of root.go
}
