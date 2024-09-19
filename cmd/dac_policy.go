package cmd

import (
	"github.com/spf13/cobra"
	"qpc/config"
	"qpc/entity/dac_policy"
)

var dacPoliciesCmd = &cobra.Command{
	Use:   "policies [flags]",
	Short: "List policies or fetch them from QueryPie API v0.9",
	Example: `  policies # List all policies in local sqlite database
  policies --fetch # Fetch all detailed policies from QueryPie API v2`,
	PreRun: func(cmd *cobra.Command, args []string) {
		m := config.LocalDatabase.Migrator()
		if !m.HasTable(&dac_policy.Policy{}) {
			dac_policy.RunAutoMigrate()
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		fetch, _ := cmd.Flags().GetBool("fetch")
		silent, _ := cmd.Flags().GetBool("silent")

		if fetch {
			var p dac_policy.Policy
			var count = 0
			p.PrintYamlHeader(silent)
			p.FetchAllAndForEach(func(fetched *dac_policy.Policy) bool {
				fetched.SaveAndLoad().PrintYaml(silent)
				count++
				return true // OK to continue fetching
			})
			p.PrintYamlFooter(silent, count)
		} else {
			var p dac_policy.Policy
			var count = 0
			p.PrintYamlHeader(silent)
			p.FindAllAndForEach(func(found *dac_policy.Policy) bool {
				found.PrintYaml(silent)
				count++
				return true // OK to continue finding
			})
			p.PrintYamlFooter(silent, count)
		}
	},
}

func addFlagsForPolicies(cmd *cobra.Command) {
	cmd.Flags().Bool("fetch", false, "Fetch from QueryPie API v0.9 and save to local sqlite database")
	cmd.Flags().Bool("silent", false, "Silent or quiet mode. Do not print outputs")
}

func init() {
	addFlagsForPolicies(dacPoliciesCmd)
	// dacCmd is added rootCmd in init() of root.go
}
