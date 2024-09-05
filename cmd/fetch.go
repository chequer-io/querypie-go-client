package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var fetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "Fetch various resources from QueryPie server",
}

var fetchDbCmd = &cobra.Command{
	Use:   "db",
	Short: "Fetch database connections from QueryPie server",
	Run: func(cmd *cobra.Command, args []string) {
		// Implement logic to fetch database connections
		fmt.Println("TODO(JK): Fetching database connections...")
	},
}

var fetchServerCmd = &cobra.Command{
	Use:   "server",
	Short: "Fetch server connections from QueryPie server",
	Run: func(cmd *cobra.Command, args []string) {
		// Implement logic to fetch server connections
		fmt.Println("TODO(JK): Fetching server connections...")
	},
}

func init() {
	// Add fetch subcommands to fetchCmd
	fetchCmd.AddCommand(fetchDbCmd)
	fetchCmd.AddCommand(fetchServerCmd)

	// Add fetchCmd to the root command
	rootCmd.AddCommand(fetchCmd)
}
