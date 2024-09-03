package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"net/url"
	"os"
)

type ServerConfig struct {
	Name string `mapstructure:"name"`
	URL  string `mapstructure:"url"`
}

var serverConfigs []ServerConfig

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "List all servers",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%-30s %-40s\n", "Name", "URL")
		// Iterate over serverConfigs and print each server's configuration
		for _, server := range serverConfigs {
			fmt.Printf("%-30s %-40s\n", server.Name, server.URL)
		}
	},
}

func initConfigForServer(viper *viper.Viper) {
	if err := viper.UnmarshalKey("servers", &serverConfigs); err != nil {
		fmt.Println("Unable to decode into struct:", err)
		os.Exit(1)
	}

	for i, server := range serverConfigs {
		if !isValidURL(server.URL) {
			fmt.Printf("Invalid URL for server %s: %s\n", server.Name, server.URL)
			os.Exit(1)
		}

		// Extract base URL
		parsedURL, err := url.Parse(server.URL)
		if err != nil {
			fmt.Printf("Error parsing URL for server %s: %s\n", server.Name, server.URL)
			os.Exit(1)
		}
		baseURL := fmt.Sprintf("%s://%s", parsedURL.Scheme, parsedURL.Host)
		serverConfigs[i].URL = baseURL
	}
}

func isValidURL(u string) bool {
	parsedURL, err := url.Parse(u)
	if err != nil {
		return false
	}
	return parsedURL.Scheme == "http" || parsedURL.Scheme == "https"
}
