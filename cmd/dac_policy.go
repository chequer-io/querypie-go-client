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

func getTargetPolicyTypes(name string) []dac_policy.PolicyType {
	var policyType = dac_policy.PolicyType(name)
	if policyType.IsValid() {
		// valid policy type
		return []dac_policy.PolicyType{policyType}
	} else if policyType == dac_policy.UnknownPolicyType {
		// On empty input, it returns all valid policy types
		return []dac_policy.PolicyType{
			dac_policy.DataLevel,
			dac_policy.DataAccess,
			dac_policy.DataMasking,
			dac_policy.Notification,
		}
	}
	logrus.Warnf("Invalid policy type: %s", name)
	os.Exit(4) // Exit code 4 means input error.
	return nil // Unreachable
}

var dacPolicyListCmd = &cobra.Command{
	Use:     "ls [flags]",
	Aliases: []string{"list"},
	Short:   "List policies in local sqlite database",
	Example: `  ls # List all policies in local sqlite database
  ls --connection=<name> # List policies with connection of the name
  ls --policy-type=<policy-type> # DATA_LEVEL, DATA_ACCESS, DATA_MASKING, or NOTIFICATION`,
	PreRun: func(cmd *cobra.Command, args []string) {
		m := config.LocalDatabase.Migrator()
		if !m.HasTable(&dac_policy.Policy{}) {
			dac_policy.RunAutoMigrate()
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		const silent = false
		connection, _ := cmd.Flags().GetString("connection")
		policyType, _ := cmd.Flags().GetString("policy-type")

		policyTypes := getTargetPolicyTypes(policyType)
		if len(connection) == 0 {
			listOrFetchAllPoliciesInYaml(policyTypes, false, silent)
		} else {
			listOrFetchPoliciesInYaml(connection, policyTypes, false, silent)
		}
	},
}

func addFlagsForPolicyList(cmd *cobra.Command) {
	cmd.Flags().SortFlags = false
	cmd.Flags().String("connection", "", "Connection name of the policy")
	cmd.Flags().String("policy-type", "", "DATA_LEVEL, DATA_ACCESS, DATA_MASKING, or NOTIFICATION")
	cmd.Flags().String("title", "", "Policy title to search")
	cmd.Flags().String("uuid", "", "Policy uuid to search")
}

var dacPolicyFetchCmd = &cobra.Command{
	Use:   "fetch [flags]",
	Short: "Fetch policies from QueryPie API v2",
	Example: `  fetch # Fetch all policies from QueryPie API v2
  fetch --policy-type=<policy-type> # DATA_LEVEL, DATA_ACCESS, DATA_MASKING, or NOTIFICATION`,
	PreRun: func(cmd *cobra.Command, args []string) {
		m := config.LocalDatabase.Migrator()
		if !m.HasTable(&dac_policy.Policy{}) {
			dac_policy.RunAutoMigrate()
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		policyType, _ := cmd.Flags().GetString("policy-type")
		silent, _ := cmd.Flags().GetBool("silent")

		policyTypes := getTargetPolicyTypes(policyType)
		listOrFetchAllPoliciesInYaml(policyTypes, true, silent)
	},
}

func addFlagsForPolicyFetch(cmd *cobra.Command) {
	cmd.Flags().SortFlags = false
	cmd.Flags().String("policy-type", "", "DATA_LEVEL, DATA_ACCESS, DATA_MASKING, or NOTIFICATION")
	cmd.Flags().Bool("silent", false, "Silent or quiet mode. Do not print outputs")
}

func listOrFetchAllPoliciesInYaml(policyTypes []dac_policy.PolicyType, fetch bool, silent bool) {
	var p dac_policy.Policy
	var count = 0
	p.PrintYamlHeader(silent)
	if fetch {
		for _, policyType := range policyTypes {
			p.FetchAllOfPolicyTypeAndForEach(policyType,
				func(fetched *dac_policy.Policy) bool {
					// NOTE(JK): Response from API v0.9 does not give the policy type.
					// So, we need to set it manually.
					fetched.PolicyType = policyType
					fetched.SaveAndLoad().PrintYaml(silent)
					count++
					return true // OK to continue fetching
				})
		}
	} else {
		for _, policyType := range policyTypes {
			p.FindAllOfPolicyTypeAndForEach(policyType, func(found *dac_policy.Policy) bool {
				found.PrintYaml(silent)
				count++
				return true // OK to continue finding
			})
		}
	}
	p.PrintYamlFooter(silent, count)
}

func listOrFetchPoliciesInYaml(
	connection string, policyTypes []dac_policy.PolicyType,
	fetch bool, silent bool,
) {
	var p dac_policy.Policy
	var count = 0
	var policies []dac_policy.Policy
	p.PrintYamlHeader(silent)
	for _, policyType := range policyTypes {
		p.FindByConnectionAndPolicyType(connection, policyType, &policies)
		for _, found := range policies {
			if fetch {
				fetched := found.FetchByUuid(found.Uuid)
				fetched.SaveAndLoad().PrintYaml(silent)
			} else {
				found.PrintYaml(silent)
			}
		}
	}
	count += len(policies)
	p.PrintYamlFooter(silent, count)
}

var dacPolicyUpsertCmd = &cobra.Command{
	Use:   "upsert <connection> <policy-type> <title> [flags]",
	Short: "Create or update a policy",
	Example: `  <connection>  - Name, or uuid of a DAC connection
  <policy-type> - DATA_LEVEL, DATA_ACCESS, DATA_MASKING, or NOTIFICATION
  <title>       - Title of the policy

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

func upsertPolicy(connection string, policyType string, title string, silent bool) {
	policy := dac_policy.GeneratePolicyRequest(
		connection,
		dac_policy.PolicyType(policyType),
		title,
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
	Use:   "delete [<connection> <policy-type>] [--uuid=<uuid>] [flags]",
	Short: "Delete a policy",
	Example: `  <connection>  - Name, or uuid of a DAC connection
  <policy-type> - DATA_LEVEL, DATA_ACCESS, DATA_MASKING, or NOTIFICATION
  --uuid=<uuid> - Optional. Uuid of the policy

  delete My-Connection DATA_ACCESS # Delete a policy
  delete --uuid <uuid> # Delete a policy by uuid`,

	Args: cobra.RangeArgs(0, 2),
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
		} else if len(args) != 2 {
			logrus.Warn("Invalid number of arguments")
			os.Exit(4) // Exit code 4 means input error.
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		var connection, policyType string
		if len(args) == 2 {
			connection, policyType = args[0], args[1]
		}
		silent, _ := cmd.Flags().GetBool("silent")
		uuid, _ := cmd.Flags().GetString("uuid")

		deletePolicy(connection, policyType, uuid, silent)
	},
}

func deletePolicy(connection string, policyType string, uuid string, silent bool) {
	var policy *dac_policy.Policy

	if len(uuid) > 0 {
		policy = (&dac_policy.PolicyRequest{PolicyUuid: uuid}).
			DeleteRemotely(utils.DefaultQuerypieServer)
	} else {
		policy = dac_policy.GeneratePolicyRequest(
			connection,
			dac_policy.PolicyType(policyType),
			"",
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
	addFlagsForPolicyList(dacPolicyListCmd)
	addFlagsForPolicyFetch(dacPolicyFetchCmd)
	addFlagsForPolicyUpsert(dacPolicyUpsertCmd)
	addFlagsForPolicyDelete(dacPolicyDeleteCmd)

	dacPolicyCmd.AddCommand(dacPolicyListCmd)
	dacPolicyCmd.AddCommand(dacPolicyFetchCmd)
	dacPolicyCmd.AddCommand(dacPolicyUpsertCmd)
	dacPolicyCmd.AddCommand(dacPolicyDeleteCmd)
	// dacCmd is added rootCmd in init() of root.go
}
