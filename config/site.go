package config

import (
	"github.com/spf13/viper"
	"log/slog"
	"net/url"
)

type SiteConfig struct {
	Name      string
	Url       string
	BaseUrl   string
	UserAgent string
	SiteEnv   string
}

func loadSiteConfig() SiteConfig {
	return SiteConfig{
		Name:      viper.GetString("SITE_NAME"),
		Url:       viper.GetString("SITE_URL"),
		BaseUrl:   GetBaseUrl(),
		UserAgent: viper.GetString("SITE_USER_AGENT"),
		SiteEnv:   viper.GetString("SITE_ENV"),
	}
}
func GetBaseUrl() string {
	urlString := viper.GetString("SITE_URL")
	parsedURL, err := url.Parse(urlString)
	if err != nil {
		slog.Error("failed to parse Url:", "Error", err)
		return ""
	}

	baseURL := parsedURL.Scheme + "://" + parsedURL.Host
	return baseURL
}
