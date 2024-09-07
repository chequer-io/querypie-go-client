package cmd

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"net/url"
	"os"
	"qpc/utils"
)

var querypieServerConfigs []utils.QueryPieServerConfig
var defaultQuerypieServer utils.QueryPieServerConfig

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
			if checkEndpoint(server, "/api/external/v2/security") {
				status = "OK"
			}
			if server.Default {
				defaultFlag = "[*]"
			}
			fmt.Printf("%-30s %-40s %-40s %-5s\n",
				server.Name+defaultFlag,
				server.BaseURL,
				utils.MaskAccessToken(server.AccessToken),
				status,
			)
		}
	},
}

func checkEndpoint(server utils.QueryPieServerConfig, uri string) bool {
	client := resty.New()

	resp, err := client.R().
		SetHeader("Authorization", "Bearer "+server.AccessToken).
		SetHeader("Content-Type", "application/json").
		Get(server.BaseURL + uri)

	if err != nil {
		logrus.Errorf("Failed to check endpoint: %v", err)
		return false
	}

	if resp.StatusCode() != 200 {
		logrus.Errorf("Received non-200 status code: %d", resp.StatusCode())
		return false
	}

	return true
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
			if defaultQuerypieServer == (utils.QueryPieServerConfig{}) {
				defaultQuerypieServer = *server
			} else {
				logrus.Fatalf("Configuration error: Multiple default querypie-server configurations found. Name: %s, URL: %s\n", server.Name, server.BaseURL)
				os.Exit(1)
			}
		}
	}

	if defaultQuerypieServer == (utils.QueryPieServerConfig{}) {
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
