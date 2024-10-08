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

var configQuerypieCmd = &cobra.Command{
	Use:     "config <key-name>",
	Short:   "Show detailed configuration for a specific key",
	Example: `  qpc config querypie     # List querypie servers in config`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			_ = cmd.Help()
		} else if args[0] == "querypie" {
			listQuerypieServers()
		} else {
			_ = cmd.Help()
			os.Exit(1)
		}
	},
}

func listQuerypieServers() {
	const format = "%-30s  %-36s  %-38s  %-4s\n"
	fmt.Printf(format, // Two spaces as delimiter
		"NAME",
		"BASE_URL",
		"ACCESS_TOKEN",
		"STATUS",
	)
	// Iterate over querypieServerConfigs and print each server's configuration
	for _, server := range utils.QuerypieServerConfigs {
		defaultFlag := ""
		status := "FAIL"
		if checkEndpoint(server, "/api/external/v2/security") {
			status = "OK"
		}
		if server.Default {
			defaultFlag = "[*]"
		}
		fmt.Printf(format,
			server.Name+defaultFlag,
			server.BaseURL,
			utils.MaskAccessToken(server.AccessToken),
			status,
		)
	}
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
	if err := viper.UnmarshalKey("querypie-servers", &utils.QuerypieServerConfigs); err != nil {
		fmt.Println("Unable to decode into struct:", err)
		os.Exit(1)
	}

	for i := range utils.QuerypieServerConfigs {
		server := &utils.QuerypieServerConfigs[i]
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
			if utils.DefaultQuerypieServer == (utils.QueryPieServerConfig{}) {
				utils.DefaultQuerypieServer = *server
			} else {
				logrus.Fatalf("Configuration error: Multiple default querypie-server configurations found. Name: %s, URL: %s\n", server.Name, server.BaseURL)
				os.Exit(1)
			}
		}
	}

	if utils.DefaultQuerypieServer == (utils.QueryPieServerConfig{}) {
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
