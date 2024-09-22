package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"qpc/config"
	"qpc/entity/dac_policy"
	"qpc/utils"
)

var dacPolicyCmd = &cobra.Command{
	Use:   "policy <command> [...] [flags]",
	Short: "Manage DAC Policies",
}

var dacPolicyListCmd = &cobra.Command{
	Use:     "ls [<name|uuid> ...] [flags]",
	Aliases: []string{"list"},
	Short:   "List policies in local sqlite database",
	Example: `  ls # List all policies in local sqlite database
  ls pattern% # List policies with a prefix 'pattern' in name
  ls %pattern # List policies with a postfix 'pattern' in name
  ls <name|uuid> # Show a policy with given name or uuid`,
	PreRun: func(cmd *cobra.Command, args []string) {
		m := config.LocalDatabase.Migrator()
		if !m.HasTable(&dac_policy.Policy{}) {
			dac_policy.RunAutoMigrate()
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		const silent = false

		if len(args) > 0 {
			listOrFetchPoliciesInYaml(args, false, silent)
		} else {
			listOrFetchAllPoliciesInYaml(false, silent)
		}
	},
}

var dacPolicyFetchCmd = &cobra.Command{
	Use:   "fetch [<name|uuid> ...] [flags]",
	Short: "Fetch policies from QueryPie API v2",
	Example: `  fetch # Fetch all policies from QueryPie API v2
  fetch pattern% # Fetch policies with a prefix 'pattern' in name
  fetch %pattern # Fetch policies with a postfix 'pattern' in name
  fetch <name|uuid> # Fetch a policy from QueryPie API v2`,
	PreRun: func(cmd *cobra.Command, args []string) {
		m := config.LocalDatabase.Migrator()
		if !m.HasTable(&dac_policy.Policy{}) {
			dac_policy.RunAutoMigrate()
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		silent, _ := cmd.Flags().GetBool("silent")

		if len(args) > 0 {
			listOrFetchPoliciesInYaml(args, true, silent)
		} else {
			listOrFetchAllPoliciesInYaml(true, silent)
		}
	},
}

func addFlagsForPolicyFetch(cmd *cobra.Command) {
	cmd.Flags().SortFlags = false
	cmd.Flags().Bool("silent", false, "Silent or quiet mode. Do not print outputs")
}

func listOrFetchAllPoliciesInYaml(fetch bool, silent bool) {
	var p dac_policy.Policy
	var count = 0
	p.PrintYamlHeader(silent)
	if fetch {
		p.FetchAllAndForEach(func(fetched *dac_policy.Policy) bool {
			fetched.SaveAndLoad().PrintYaml(silent)
			count++
			return true // OK to continue fetching
		})
	} else {
		p.FindAllAndForEach(func(found *dac_policy.Policy) bool {
			found.PrintYaml(silent)
			count++
			return true // OK to continue finding
		})
	}
	p.PrintYamlFooter(silent, count)
}

func listOrFetchPoliciesInYaml(args []string, fetch bool, silent bool) {
	var p dac_policy.Policy
	var count = 0
	p.PrintYamlHeader(silent)
	for _, pattern := range args {
		var list []dac_policy.Policy
		p.FindByNameOrUuid(pattern, &list)
		for _, found := range list {
			if fetch {
				fetched := found.FetchByUuid(found.Uuid)
				fetched.SaveAndLoad().PrintYaml(silent)
			} else {
				found.PrintYaml(silent)
			}
		}
		count += len(list)
	}
	p.PrintYamlFooter(silent, count)
}

var dacPolicyUpsertCmd = &cobra.Command{
	Use:   "upsert <connection> <policy-type> <name> [flags]",
	Short: "Create or update a policy",
	Example: `  <connection>  - Name, or uuid of a DAC connection
  <policy-type> - DATA_LEVEL, DATA_ACCESS, DATA_MASKING, NOTIFICATION, or LEDGER
  <name>        - Name, or title of the policy

  upsert My-Connection DATA_LEVEL My-Policy # Upsert a policy`,

	Args: cobra.ExactArgs(3),
	PreRun: func(cmd *cobra.Command, args []string) {
		m := config.LocalDatabase.Migrator()
		if !m.HasTable(&dac_policy.Policy{}) {
			dac_policy.RunAutoMigrate()
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		silent, _ := cmd.Flags().GetBool("silent")

		upsertPolicy(args[0], args[1], args[2], silent)
	},
}

func upsertPolicy(connection string, policyType string, name string, silent bool) {
	policy := dac_policy.GeneratePolicyRequest(
		connection,
		dac_policy.PolicyType(policyType),
		name,
	).Validate().PrintYaml(silent).UnlessValidated(func() {
		logrus.Warn("Validation failed")
		os.Exit(4) // Exit code 4 means input error.
	}).
		PolicyRequest.
		UpdateOrCreateRemotely(utils.DefaultQuerypieServer)

	policy.HandleResponse(func() {
		policy.PrintHttpRequestLineAndResponseStatus(silent)

		// NOTE(JK): Response from API v0.9 does not give the policy type.
		// So, we need to set it manually.
		policy.PolicyType = dac_policy.PolicyType(policyType)

		policy.SaveAndLoad().PrintYaml(silent)
	}, func() {
		policy.PrintHttpRequestLineAndResponseStatus(silent).
			PrintRawBody(silent)
		os.Exit(4) // Exit code 4 means input error.
	}, func() {
		policy.PrintHttpRequestLineAndResponseStatus(silent).
			PrintRawBody(silent)
		os.Exit(5) // Exit code 5 means remote error.
	})
}

func addFlagsForPolicyUpsert(cmd *cobra.Command) {
	cmd.Flags().SortFlags = false
	cmd.Flags().Bool("silent", false, "Silent or quiet mode. Do not print outputs")
}

var dacPolicyDeleteCmd = &cobra.Command{
	Use:   "delete [<connection> <policy-type> <name>] [--uuid=<uuid>] [flags]",
	Short: "Delete a policy",
	Example: `  <connection>  - Name, or uuid of a DAC connection
  <policy-type> - DATA_LEVEL, DATA_ACCESS, DATA_MASKING, NOTIFICATION, or LEDGER
  <name>        - Name, or title of the policy
  <uuid>        - Optional. Uuid of the policy

  delete My-Connection DATA_LEVEL My-Policy # Delete a policy
  delete --uuid <uuid> # Delete a policy by uuid`,

	Args: cobra.RangeArgs(0, 3),
	PreRun: func(cmd *cobra.Command, args []string) {
		m := config.LocalDatabase.Migrator()
		if !m.HasTable(&dac_policy.Policy{}) {
			dac_policy.RunAutoMigrate()
		}

		uuid, _ := cmd.Flags().GetString("uuid")
		if len(args) == 0 {
			if len(uuid) == 0 {
				logrus.Warn("No arguments given")
				os.Exit(4) // Exit code 4 means input error.
			}
		} else if len(args) != 3 {
			logrus.Warn("Invalid number of arguments")
			os.Exit(4) // Exit code 4 means input error.
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		var connection, policyType, name string
		if len(args) == 3 {
			connection, policyType, name = args[0], args[1], args[2]
		}
		silent, _ := cmd.Flags().GetBool("silent")
		uuid, _ := cmd.Flags().GetString("uuid")

		deletePolicy(connection, policyType, name, uuid, silent)
	},
}

func deletePolicy(connection string, policyType string, name string, uuid string, silent bool) {
	var policy *dac_policy.Policy

	if len(uuid) > 0 {
		policy = (&dac_policy.PolicyRequest{PolicyUuid: uuid}).
			DeleteRemotely(utils.DefaultQuerypieServer)
	} else {
		policy = dac_policy.GeneratePolicyRequest(
			connection,
			dac_policy.PolicyType(policyType),
			name,
		).ValidateForDelete().PrintYaml(silent).UnlessValidated(func() {
			logrus.Warn("Validation failed")
			os.Exit(4) // Exit code 4 means input error.
		}).
			PolicyRequest.
			DeleteRemotely(utils.DefaultQuerypieServer)
	}

	policy.HandleResponse(func() {
		policy.PrintHttpRequestLineAndResponseStatus(silent)
		policy.Delete()
	}, func() {
		policy.PrintHttpRequestLineAndResponseStatus(silent).
			PrintRawBody(silent)
		os.Exit(4) // Exit code 4 means input error.
	}, func() {
		policy.PrintHttpRequestLineAndResponseStatus(silent).
			PrintRawBody(silent)
		os.Exit(5) // Exit code 5 means remote error.
	})
}

func addFlagsForPolicyDelete(cmd *cobra.Command) {
	cmd.Flags().SortFlags = false
	cmd.Flags().Bool("silent", false, "Silent or quiet mode. Do not print outputs")
	cmd.Flags().String("uuid", "", "Uuid of the policy")
}

func init() {
	addFlagsForPolicyFetch(dacPolicyFetchCmd)
	addFlagsForPolicyUpsert(dacPolicyUpsertCmd)
	addFlagsForPolicyDelete(dacPolicyDeleteCmd)

	dacPolicyCmd.AddCommand(dacPolicyListCmd)
	dacPolicyCmd.AddCommand(dacPolicyFetchCmd)
	dacPolicyCmd.AddCommand(dacPolicyUpsertCmd)
	dacPolicyCmd.AddCommand(dacPolicyDeleteCmd)
	// dacCmd is added rootCmd in init() of root.go
}
