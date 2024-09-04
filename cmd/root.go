package cmd

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"

	"github.com/spf13/cobra"
)

var (
	configFile string
	logLevel   string
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
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&configFile, "config", "", "config file (default is ./.querypie-client.yaml)")
	rootCmd.PersistentFlags().StringVar(&logLevel, "log-level", "info", "Set the logging level (debug, info, warn, error, fatal, panic)")
	// Add global flags or subcommands here
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(serverCmd)
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

	initConfigForServer(v)
}
