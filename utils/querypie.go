package utils

type QueryPieServerConfig struct {
	Name        string `mapstructure:"name"`
	BaseURL     string `mapstructure:"url"`
	AccessToken string `mapstructure:"token"`
	Default     bool   `mapstructure:"default"`
}
