package utils

var QuerypieServerConfigs []QueryPieServerConfig
var DefaultQuerypieServer QueryPieServerConfig

type QueryPieServerConfig struct {
	Name        string `mapstructure:"name"`
	BaseURL     string `mapstructure:"url"`
	AccessToken string `mapstructure:"token"`
	Default     bool   `mapstructure:"default"`
}
