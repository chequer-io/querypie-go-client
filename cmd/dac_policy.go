package cmd

import (
	"github.com/spf13/cobra"
	"qpc/config"
	"qpc/entity/dac_policy"
)

var dacPolicyCmd = &cobra.Command{
	Use:   "policy [<name|uuid> ...] [flags]",
	Short: "Manage DAC Policies",
	Example: `  policy # List all policies in local sqlite database
  policy --fetch # Fetch all policies from QueryPie API v2
  policy <name|uuid> # Show a policy in local sqlite database
  policy <name|uuid> --fetch # Fetch a policy from QueryPie API v0.9
  (TODO) policy --upsert # Update or create a policy to QueryPie API v2`,
	PreRun: func(cmd *cobra.Command, args []string) {
		m := config.LocalDatabase.Migrator()
		if !m.HasTable(&dac_policy.Policy{}) {
			dac_policy.RunAutoMigrate()
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		fetch, _ := cmd.Flags().GetBool("fetch")
		like, _ := cmd.Flags().GetBool("like")
		silent, _ := cmd.Flags().GetBool("silent")

		if len(args) == 0 {
			listPoliciesInYaml(fetch, silent)
			return
		}

		var p dac_policy.Policy
		var count = 0
		p.PrintYamlHeader(silent)
		for _, arg := range args {
			var list []dac_policy.Policy
			p.FindByNameOrUuid(arg, like, &list)
			for _, found := range list {
				found.PrintYaml(silent)
			}
			count += len(list)
		}
		p.PrintYamlFooter(silent, count)
	},
}

func addFlagsForPolicy(cmd *cobra.Command) {
	cmd.Flags().SortFlags = false
	cmd.Flags().Bool("fetch", false, "Fetch from QueryPie API v0.9 and save to local sqlite database")
	cmd.Flags().Bool("like", false, "Use LIKE instead of = in SQL queries")
	cmd.Flags().Bool("silent", false, "Silent or quiet mode. Do not print outputs")
}

func listPoliciesInYaml(fetch bool, silent bool) {
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
}

func init() {
	addFlagsForPolicy(dacPolicyCmd)
	// dacCmd is added rootCmd in init() of root.go
}
