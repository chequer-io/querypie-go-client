package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "qpc",
	Short: "QueryPie Client for Operation",
	Long:  `QueryPie Client for Operation is a CLI client for managing QueryPie operations.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Default action when no subcommands are provided
		fmt.Println("Hello from QueryPie Client for Operation!")
	},
}

// Execute executes the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// Add global flags or subcommands here
	rootCmd.AddCommand(versionCmd)
}
