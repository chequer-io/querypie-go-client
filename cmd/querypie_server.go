package cmd

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"net/url"
	"os"
	"qpc/rest"
	"unicode"
)

type QueryPieServerConfig struct {
	Name        string `mapstructure:"name"`
	BaseURL     string `mapstructure:"url"`
	AccessToken string `mapstructure:"token"`
}

var querypieServerConfigs []QueryPieServerConfig

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
			status := "FAIL"
			if checkEndpoint(server, "/api/external/users?pageSize=3") {
				status = "OK"
			}
			fmt.Printf("%-30s %-40s %-40s %-5s\n",
				server.Name,
				server.BaseURL,
				maskAccessToken(server.AccessToken),
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

func maskAccessToken(token string) string {
	if len(token) <= 11 {
		return token
	}
	masked := []rune(token[:11])
	for _, r := range token[11:] {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			masked = append(masked, '*')
		} else {
			masked = append(masked, r)
		}
	}
	return string(masked)
}

func initConfigForQueryPieServer(viper *viper.Viper) {
	if err := viper.UnmarshalKey("querypie-servers", &querypieServerConfigs); err != nil {
		fmt.Println("Unable to decode into struct:", err)
		os.Exit(1)
	}

	for i, server := range querypieServerConfigs {
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
		querypieServerConfigs[i].BaseURL = baseURL
	}
}

func isValidURL(u string) bool {
	parsedURL, err := url.Parse(u)
	if err != nil {
		return false
	}
	return parsedURL.Scheme == "http" || parsedURL.Scheme == "https"
}
