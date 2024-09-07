package cmd

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"qpc/local_db"

	"github.com/spf13/cobra"
)

var (
	configFile string
	logLevel   string
)

var rootCmd = &cobra.Command{
	Use:   "qpc",
	Short: "QueryPie Client for Operation",
	Long:  `QueryPie Client for Operation - you can manage resources, access control policies`,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
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
	cobra.OnInitialize(initConfig)
	cobra.EnableCommandSorting = false // Do not sort commands alphabetically
	rootCmd.PersistentFlags().StringVar(&configFile,
		"config", ".querypie-client.yaml",
		"yaml file for configuration")
	rootCmd.PersistentFlags().StringVar(&logLevel,
		"log-level", "warn",
		"Set the logging level (debug, info, warn, error, fatal, panic)")
	// Add global flags or subcommands here
	rootCmd.AddCommand(dacCmd)
	rootCmd.AddCommand(userCmd)
	rootCmd.AddCommand(userV1Cmd)
	rootCmd.AddCommand(fetchAllCmd)
	rootCmd.AddCommand(configQuerypieCmd)
	rootCmd.AddCommand(versionCmd)
}

func initConfig() {
	v := viper.New()
	if configFile != "" {
		v.SetConfigFile(configFile)
	} else {
		v.AddConfigPath(".")
		v.SetConfigFile(".querypie-client.yaml")
	}
	v.AutomaticEnv()
	if err := v.ReadInConfig(); err != nil {
		fmt.Println("Can't read config:", err)
		os.Exit(1)
	}

	// Parse and set log level
	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		fmt.Printf("Invalid log level: %s\n", logLevel)
		os.Exit(1)
	}
	logrus.SetLevel(level)

	initConfigForQueryPieServer(v)
	local_db.InitConfigForResource(v)
}
