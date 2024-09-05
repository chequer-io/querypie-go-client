package cmd

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"net/url"
	"os"
	"qpc/rest"
)

type QueryPieServerConfig struct {
	Name        string `mapstructure:"name"`
	BaseURL     string `mapstructure:"url"`
	AccessToken string `mapstructure:"token"`
	Default     bool   `mapstructure:"default"`
}

var querypieServerConfigs []QueryPieServerConfig
var defaultQuerypieServer QueryPieServerConfig

var querypieServerCmd = &cobra.Command{
	Use:   "querypie-servers",
	Short: "List all querypie servers",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%-30s %-40s %-40s %-5s\n",
			"NAME",
			"BASE_URL",
			"ACCESS_TOKEN",
			"STATUS",
		)
		// Iterate over querypieServerConfigs and print each server's configuration
		for _, server := range querypieServerConfigs {
			defaultFlag := ""
			status := "FAIL"
			if checkEndpoint(server, "/api/external/users?pageSize=3") {
				status = "OK"
			}
			if server.Default {
				defaultFlag = "[*]"
			}
			fmt.Printf("%-30s %-40s %-40s %-5s\n",
				server.Name+defaultFlag,
				server.BaseURL,
				rest.MaskAccessToken(server.AccessToken),
				status,
			)
		}
	},
}

func checkEndpoint(server QueryPieServerConfig, uri string) bool {
	client := rest.NewAPIClient(server.BaseURL, server.AccessToken)
	// Call the GetData method
	result, err := client.GetData(uri)
	logrus.Debugf("Result: %v", result)
	if err == nil {
		return true
	} else {
		return false
	}
}

func initConfigForQueryPieServer(viper *viper.Viper) {
	if err := viper.UnmarshalKey("querypie-servers", &querypieServerConfigs); err != nil {
		fmt.Println("Unable to decode into struct:", err)
		os.Exit(1)
	}

	for i := range querypieServerConfigs {
		server := &querypieServerConfigs[i]
		if !isValidURL(server.BaseURL) {
			logrus.Fatalf("Invalid URL for server %s: %s\n", server.Name, server.BaseURL)
			os.Exit(1)
		}

		// Extract base URL
		parsedURL, err := url.Parse(server.BaseURL)
		if err != nil {
			logrus.Fatalf("Error parsing URL for server %s: %s\n", server.Name, server.BaseURL)
			os.Exit(1)
		}
		baseURL := fmt.Sprintf("%s://%s", parsedURL.Scheme, parsedURL.Host)
		server.BaseURL = baseURL

		// Check if the server is the default server
		if server.Default {
			if defaultQuerypieServer == (QueryPieServerConfig{}) {
				defaultQuerypieServer = *server
			} else {
				logrus.Fatalf("Configuration error: Multiple default querypie-server configurations found. Name: %s, URL: %s\n", server.Name, server.BaseURL)
				os.Exit(1)
			}
		}
	}

	if defaultQuerypieServer == (QueryPieServerConfig{}) {
		logrus.Fatalf("Configuration error: No default querypie-server configuration found.\n")
		os.Exit(1)
	}
}

func isValidURL(u string) bool {
	parsedURL, err := url.Parse(u)
	if err != nil {
		return false
	}
	return parsedURL.Scheme == "http" || parsedURL.Scheme == "https"
}
