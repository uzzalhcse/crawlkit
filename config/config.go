package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Site     SiteConfig
	Database DatabaseConfig
}

// NewConfig loads the common configuration from environment variables
func NewConfig() *Config {
	viper.SetConfigName(".env") // allow directly reading from .env file
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("/")
	viper.AllowEmptyEnv(true)
	viper.AutomaticEnv()
	_ = viper.ReadInConfig()

	return &Config{
		Site:     loadSiteConfig(),
		Database: loadDBConfig(),
	}
}
