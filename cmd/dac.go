package cmd

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/tidwall/pretty"
	"os"
	"qpc/config"
	"qpc/entity/dac_access_control"
	"qpc/entity/dac_connection"
	"qpc/entity/dac_privilege"
	"qpc/utils"
	"strconv"
)

var dacCmd = &cobra.Command{
	Use:   "dac",
	Short: "Manage DAC resources",
}

var dacFetchAllCmd = &cobra.Command{
	Use:   "fetch-all <resource>",
	Short: "Fetch all DAC resources from QueryPie server and save them to local sqlite database",
	Example: `  fetch-all connections # from QueryPie API v2
  fetch-all detailed-connections # from QueryPie API v2
  fetch-all access-controls # from QueryPie API v2
  fetch-all privileges # from QueryPie API v2`,
	Args: cobra.ExactArgs(1),
	PreRun: func(cmd *cobra.Command, args []string) {
		m := config.LocalDatabase.Migrator()
		if !m.HasTable(&dac_connection.SummarizedConnectionV2{}) {
			dac_connection.RunAutoMigrate()
		}
		if !m.HasTable(&dac_access_control.SummarizedAccessControl{}) {
			dac_access_control.RunAutoMigrate()
		}
		if !m.HasTable(&dac_privilege.Privilege{}) {
			dac_privilege.RunAutoMigrate()
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		resource := args[0]
		switch resource {
		case "connections":
			var sc dac_connection.SummarizedConnectionV2
			sc.PrintHeader()
			sc.FetchAllAndForEach(func(fetched *dac_connection.SummarizedConnectionV2) bool {
				fetched.Print().Save()
				return true // OK to continue fetching
			})
		case "detailed-connections":
			dac_connection.PrintHeaderOfDetailedConnection()
			(&dac_connection.SummarizedConnectionV2{}).
				FetchAllAndForEach(func(fetched *dac_connection.SummarizedConnectionV2) bool {
					fetched.Print().Save()
					return fetched.FetchDetailedConnectionAndPrintAndSave()
				})
		case "access-controls":
			var sac dac_access_control.SummarizedAccessControl
			sac.PrintHeader()
			sac.FetchAllAndForEach(func(fetched *dac_access_control.SummarizedAccessControl) bool {
				fetched.Print().Save()
				return true // OK to continue fetching
			})
		case "privileges":
			var p dac_privilege.Privilege
			p.PrintHeader()
			p.FetchAllAndForEach(func(fetched *dac_privilege.Privilege) bool {
				fetched.Print().Save()
				return true // OK to continue fetching
			})
		default:
			logrus.Fatalf("Unknown resource: %s", resource)
		}
	},
}

var dacListCmd = &cobra.Command{
	Use:   "ls",
	Short: "List DAC connections in local sqlite database",
	Example: `  ls connections # from local sqlite database
  ls detailed-connections # from local sqlite database
  ls access-controls # from local sqlite database
  ls privileges # from QueryPie API v2
`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			_ = cmd.Help()
			return
		}
		resource := args[0]
		switch resource {
		case "connections":
			var sc dac_connection.SummarizedConnectionV2
			sc.PrintHeader()
			sc.FindAllAndForEach(func(found *dac_connection.SummarizedConnectionV2) bool {
				found.Print()
				return true // OK to continue finding
			})
		case "detailed-connections":
			dac_connection.PrintHeaderOfDetailedConnection()
			(&dac_connection.SummarizedConnectionV2{}).
				FindAllAndForEach(func(found *dac_connection.SummarizedConnectionV2) bool {
					found.Print()
					return found.FirstDetailedConnectionAndPrint()
				})
		case "access-controls":
			var sac dac_access_control.SummarizedAccessControl
			sac.PrintHeader()
			sac.FindAllAndForEach(func(found *dac_access_control.SummarizedAccessControl) bool {
				found.Print()
				return true // OK to continue finding
			})
		case "privileges":
			var p dac_privilege.Privilege
			p.PrintHeader()
			p.FindAllAndForEach(func(found *dac_privilege.Privilege) bool {
				found.Print()
				return true // OK to continue finding
			})
		case "clusters":
			var c dac_connection.Cluster
			c.PrintHeaderWithConnection()
			c.FindAllAndForEach(func(found *dac_connection.Cluster) bool {
				logrus.Debug(found)
				found.PrintWithConnection()
				return true // OK to continue finding
			})
		default:
			logrus.Fatalf("Unknown resource: %s", resource)
		}
	},
}

var grantCmd = &cobra.Command{
	Use:   "grant <user> <privilege> <cluster|connection>",
	Short: "Grants a <privilege> to <user> for accessing a <cluster> in a DAC connection",
	Example: `  <user>       - login_id, email, or uuid
  <privilege>  - name, or uuid
  <cluster>    - host:port, cloud_identifier, or uuid
  <connection> - name, or uuid`,
	Args: cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		force, _ := cmd.Flags().GetBool("force")
		dryRun, _ := cmd.Flags().GetBool("dry-run")
		var req = dac_access_control.DraftGrantRequest{
			UserQuery:      args[0],
			PrivilegeQuery: args[1],
			ClusterQuery:   args[2],
			Force:          force,
		}
		req.LookUpEntities().Print().Validate(
			func() {
				fmt.Print("VALIDATION: success\n\n")
			},
			func(reason string) {
				fmt.Printf("VALIDATION: failure - %s\n\n", reason)
				// Exit with error code 4 for validation failure
				os.Exit(4)
			},
		)
		r := req.ToGrantRequest()
		if dryRun {
			logrus.Infof("Dry-run mode is enabled. No actual grant is performed.")
		} else {
			r.Post(utils.DefaultQuerypieServer).Print()
		}
	},
}

func addFlagsForGrant(cmd *cobra.Command) {
	cmd.Flags().Bool("dry-run", false, "Dry-run mode to verify the request")
	cmd.Flags().Bool("force", false, "Force to replace the existing privilege (default: false)")
}

var dacFetchByUuidCmd = &cobra.Command{
	Use:     "fetch-by-uuid <resource> <uuid>",
	Short:   "[Debug] Fetch a DAC resource specified as UUID, and save it to local sqlite database",
	Example: `  fetch-by-uuid connection <connection-uuid> # from QueryPie API v2`,
	Args:    cobra.ExactArgs(2),
	PreRun: func(cmd *cobra.Command, args []string) {
		m := config.LocalDatabase.Migrator()
		resource := args[0]
		switch resource {
		case "connection":
			if !m.HasTable(&dac_connection.SummarizedConnectionV2{}) ||
				!m.HasTable(&dac_connection.ConnectionV2{}) {
				dac_connection.RunAutoMigrate()
			}
		default:
			logrus.Fatalf("Unknown resource: %s", resource)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		resource := args[0]
		uuid := args[1]
		switch resource {
		case "connection":
			var c = (&dac_connection.ConnectionV2{}).FetchByUuid(uuid)
			utils.PrintHttpRequestLineAndResponseStatus(c.HttpResponse)
			if c.HttpResponse.IsSuccess() {
				c.PrintJson().Save()
			} else if utils.IsServerError(c.HttpResponse) {
				fmt.Printf("%s\n", pretty.Pretty(c.HttpResponse.Body()))
				// When the API returns a server error,
				// it is necessary to save an empty object with the Uuid,
				// so that following Save() can save a non-nil object,
				// as a workaround to prevent error due to missing entities in local database.
				c.SaveAlsoForServerError()
			} else if utils.IsClientError(c.HttpResponse) {
				fmt.Printf("%s\n", pretty.Pretty(c.HttpResponse.Body()))
				// Do not save for client errors such as 404 Not Found
			} else {
				logrus.Fatalf("Unexpected error: %s", c.HttpResponse.Status())
			}
		default:
			logrus.Fatalf("Unknown resource: %s", resource)
		}
	},
}

var dacFindByUuidCmd = &cobra.Command{
	Use:     "find-by-uuid <resource> <uuid>",
	Short:   "[Debug] Find a DAC resource specified as UUID from local sqlite database",
	Example: `  find-by-uuid connection <connection-uuid>`,
	Args:    cobra.ExactArgs(2),
	PreRun: func(cmd *cobra.Command, args []string) {
		if !config.LocalDatabase.Migrator().HasTable(&dac_connection.ConnectionV2{}) {
			dac_connection.RunAutoMigrate()
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		resource := args[0]
		uuid := args[1]
		switch resource {
		case "connection":
			var c dac_connection.ConnectionV2
			c.FirstByUuid(uuid).PrintJson()
		default:
			logrus.Fatalf("Unknown resource: %s", resource)
		}
	},
}

var grantByUuidCmd = &cobra.Command{
	Use:   "grant-by-uuid <user-uuid> <connection-uuid> <privilege-uuid> [<force>]",
	Short: "[Debug] Grant access to a DAC connection using UUIDs as argument",
	Example: `  grant-by-uuid <uuid> <uuid> <uuid>
  grant-by-uuid <uuid> <uuid> <uuid> false
  grant-by-uuid <uuid> <uuid> <uuid> true`,
	Args: cobra.RangeArgs(3, 4),
	Run: func(cmd *cobra.Command, args []string) {
		userUuid := args[0]
		clusterUuid := args[1]
		privilegeUuid := args[2]
		force := false
		if len(args) > 3 {
			force, _ = strconv.ParseBool(args[3])
		}
		(&dac_access_control.GrantRequest{
			UserUuid:      userUuid,
			ClusterUuid:   clusterUuid,
			PrivilegeUuid: privilegeUuid,
			Force:         force,
		}).
			Post(utils.DefaultQuerypieServer).
			Print()
	},
}

func init() {

	addFlagsForGrant(grantCmd)

	dacCmd.AddCommand(dacConnectionCmd)
	dacCmd.AddCommand(dacListCmd)
	dacCmd.AddCommand(dacFetchAllCmd)
	dacCmd.AddCommand(grantCmd)
	dacCmd.AddCommand(dacPolicyCmd)
	dacCmd.AddCommand(dacSensitiveDataRuleCmd)

	dacCmd.AddCommand(dacFetchByUuidCmd)
	dacCmd.AddCommand(dacFindByUuidCmd)
	dacCmd.AddCommand(grantByUuidCmd)

	// dacPolicyCmd is added rootCmd in init() of root.go
}
