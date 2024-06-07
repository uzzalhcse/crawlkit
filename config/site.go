package config

import (
	"github.com/spf13/viper"
)

type SiteConfig struct {
	UserAgent string
	SiteEnv   string
	Proxy     string
}

func loadSiteConfig() SiteConfig {
	return SiteConfig{
		UserAgent: viper.GetString("SITE_USER_AGENT"),
		SiteEnv:   viper.GetString("SITE_ENV"),
		Proxy:     viper.GetString("SITE_PROXY"),
	}
}
