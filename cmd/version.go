package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version of QueryPie Client for Operation",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("qpc v0.1.0")
	},
}
