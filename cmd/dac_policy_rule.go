package cmd

import (
	"github.com/spf13/cobra"
	"qpc/entity/dac_policy"
)

var dacSensitiveDataRuleCmd = &cobra.Command{
	Use:   "sensitive-data-rule <command> [...] [flags]",
	Short: "Manage Sensitive Data Rules",
}

var dacSensitiveDataRuleListCmd = &cobra.Command{
	Use:     "ls [flags]",
	Aliases: []string{"list"},
	Short:   "List sensitive data rules in local sqlite database",
	Example: `  ls # List all sensitive data rules in local sqlite database
  ls --connection=<name> # List sensitive data rules with connection of the name`,
	PreRun: func(cmd *cobra.Command, args []string) {
		dac_policy.RunAutoMigrate()
	},
	Run: func(cmd *cobra.Command, args []string) {
		const silent = false
		connection, _ := cmd.Flags().GetString("connection")

		if len(connection) == 0 {
			listOrFetchAllSensitiveDataRulesInYaml(false, silent)
		} else {
			listOrFetchSensitiveDataRulesInYaml(connection, dac_policy.DataLevel, false, silent)
		}
	},
}

func addFlagsForSensitiveDataRuleList(cmd *cobra.Command) {
	cmd.Flags().SortFlags = false
	cmd.Flags().String("connection", "", "Connection name of the policy")
}

var dacSensitiveDataRuleFetchCmd = &cobra.Command{
	Use:   "fetch [flags]",
	Short: "Fetch sensitive-data-rule from QueryPie API v0.9",
	Example: `  fetch # Fetch all policies from QueryPie API v2
  fetch --policy-type=<policy-type> # DATA_LEVEL, DATA_ACCESS, DATA_MASKING, or NOTIFICATION`,
	PreRun: func(cmd *cobra.Command, args []string) {
		dac_policy.RunAutoMigrate()
	},
	Run: func(cmd *cobra.Command, args []string) {
		silent, _ := cmd.Flags().GetBool("silent")

		listOrFetchAllSensitiveDataRulesInYaml(true, silent)
	},
}

func addFlagsForSensitiveDataRuleFetch(cmd *cobra.Command) {
	cmd.Flags().SortFlags = false
	cmd.Flags().Bool("silent", false, "Silent or quiet mode. Do not print outputs")
}

func listOrFetchAllSensitiveDataRulesInYaml(fetch bool, silent bool) {
	var r dac_policy.SensitiveDataRule
	var p dac_policy.Policy
	var count = 0
	r.PrintYamlHeader(silent)
	if fetch {
		p.FindAllOfPolicyTypeAndForEach(dac_policy.DataLevel, func(policy *dac_policy.Policy) bool {
			r.FetchAllOfPolicyAndForEach(policy, func(fetched *dac_policy.SensitiveDataRule) bool {
				fetched.SaveAndLoad().PrintYaml(silent)
				count++
				return true // OK to continue finding
			})
			return true // OK to continue finding
		})
	} else {
		p.FindAllOfPolicyTypeAndForEach(dac_policy.DataLevel, func(policy *dac_policy.Policy) bool {
			r.FindAllOfPolicyAndForEach(policy, func(found *dac_policy.SensitiveDataRule) bool {
				found.PrintYaml(silent)
				count++
				return true // OK to continue finding
			})
			return true // OK to continue finding
		})
	}
	r.PrintYamlFooter(silent, count)
}

func listOrFetchSensitiveDataRulesInYaml(
	connection string, policyType dac_policy.PolicyType,
	fetch bool, silent bool,
) {
	var r dac_policy.SensitiveDataRule
	var p dac_policy.Policy
	var count = 0
	var policies []dac_policy.Policy
	p.PrintYamlHeader(silent)
	p.FindByConnectionAndPolicyType(connection, policyType, &policies)
	for _, policy := range policies {
		if fetch {
			r.FetchAllOfPolicyAndForEach(&policy, func(fetched *dac_policy.SensitiveDataRule) bool {
				fetched.SaveAndLoad().PrintYaml(silent)
				count++
				return true // OK to continue finding
			})
		} else {
			r.FindAllOfPolicyAndForEach(&policy, func(found *dac_policy.SensitiveDataRule) bool {
				found.PrintYaml(silent)
				count++
				return true // OK to continue finding
			})
		}
	}
	count += len(policies)
	p.PrintYamlFooter(silent, count)
}

func init() {
	addFlagsForSensitiveDataRuleList(dacSensitiveDataRuleListCmd)
	addFlagsForSensitiveDataRuleFetch(dacSensitiveDataRuleFetchCmd)

	dacSensitiveDataRuleCmd.AddCommand(dacSensitiveDataRuleListCmd)
	dacSensitiveDataRuleCmd.AddCommand(dacSensitiveDataRuleFetchCmd)
	// dacSensitiveDataRuleCmd is added dacCmd in init() of dac.go

}
